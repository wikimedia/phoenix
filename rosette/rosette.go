package rosette

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jpillora/backoff"
	"github.com/wikimedia/phoenix/common"
)

const (
	rosetteEndpoint        = "https://api.rosette.com/rest/v1/topics?redirect=true"
	rosetteMinSalience     = 0.1
	rosetteMaxContentSize  = 600 * 1024 // 600KB
	rosetteMaxContentChars = 50000
	rosetteMaxTopics       = 10
	rosetteRetries         = 10
)

var (
	// Regex for matching extraneous spaces
	xSpace = regexp.MustCompile(`\s+`)
)

// A helper for extracting text from an HTML snippet (very basic; Good Enough™ for now).  Note: If the resulting
// string would exceed Rosette's limits (thus triggering a 413), then it will be truncated accordingly.
func extractText(unsafe string) (string, error) {
	var doc *goquery.Document
	var err error
	var txtStr string
	var txtRunes []rune

	if doc, err = goquery.NewDocumentFromReader(strings.NewReader(unsafe)); err != nil {
		return "", fmt.Errorf("failed to parse html string: %w", err)
	}

	txtStr = xSpace.ReplaceAllString(doc.Text(), " ")
	txtRunes = []rune(txtStr)

	// Truncate the result, if necessary
	if len(txtStr) > rosetteMaxContentSize {
		txtStr = txtStr[0:rosetteMaxContentSize]
	}

	if len(txtRunes) > rosetteMaxContentChars {
		txtStr = string(txtRunes[0:rosetteMaxContentChars])
	}

	return txtStr, nil
}

type topicsResponse struct {
	Concepts []struct {
		ConceptID string  `json:"conceptId"`
		Phrase    string  `json:"phrase"`
		Salience  float32 `json:"salience"`
	} `json:"concepts"`
}

// Rosette is an abstraction for querying related topics from the Rosette service.
type Rosette struct {
	APIKey string
	Logger *common.Logger
}

func (obj *Rosette) requestTopics(text string) (*topicsResponse, error) {
	var b = &backoff.Backoff{Min: 500 * time.Millisecond, Max: 20 * time.Second, Jitter: true}
	var client = &http.Client{}
	var err error
	var res *http.Response
	var resData []byte
	var topics topicsResponse

	request := func(text string) (*http.Response, error) {
		var e error
		var req *http.Request
		var reqData []byte
		var r *http.Response

		// Serialize a requests body (JSON)
		content := &struct {
			Content  string `json:"content"`
			Language string `json:"language"`
		}{
			text,
			"eng",
		}

		if reqData, e = json.Marshal(content); e != nil {
			return nil, fmt.Errorf("unable to serialize JSON request body: %w", e)
		}

		// Create the request
		if req, e = http.NewRequest("POST", rosetteEndpoint, bytes.NewBuffer(reqData)); e != nil {
			return nil, fmt.Errorf("failed to create HTTP request: %w", e)
		}

		req.Header.Set("X-RosetteAPI-Key", obj.APIKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Cache-Control", "no-cache")

		// Send the request
		if r, e = client.Do(req); e != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", e)
		}

		return r, nil
	}

Loop:
	for i := 1; ; i++ {
		// Send the request
		if res, err = request(text); err != nil {
			return nil, err
		}

		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			break Loop
		case http.StatusTooManyRequests:
			if i >= rosetteRetries {
				return nil, fmt.Errorf("Reached max number of Rosette retries (%d)", rosetteRetries)
			}

			delay := b.Duration()
			obj.Logger.Debug("Rosette returned status 429; reconnecting in %s (#%d/%d)", delay, i, rosetteRetries)
			time.Sleep(delay)
			continue
		default:
			return nil, fmt.Errorf("unexpected HTTP status code %d (expected %d)", res.StatusCode, http.StatusOK)
		}
	}

	// Read the response
	if resData, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	// Deserialize results
	if err = json.Unmarshal(resData, &topics); err != nil {
		return nil, fmt.Errorf("failure to deserialize JSON results: %w", err)
	}

	return &topics, nil
}

// Topics returns topics correlated with a Node
func (obj *Rosette) Topics(node *common.Node) ([]common.RelatedTopic, error) {
	var content string
	var err error
	var res = make([]common.RelatedTopic, 0)
	var topics *topicsResponse

	if node.Unsafe == "" {
		return nil, fmt.Errorf("node.Unsafe is zero length")
	}

	if content, err = extractText(node.Unsafe); err != nil {
		return nil, err
	}

	if topics, err = obj.requestTopics(content); err != nil {
		return nil, fmt.Errorf("failure retrieving related topics: %w", err)
	}

	for i, topic := range topics.Concepts {
		// Limit our results to those greater than a defined minimum salience, and establish an upper bound on quantity.
		if topic.Salience > rosetteMinSalience && i < rosetteMaxTopics {
			res = append(res, common.RelatedTopic{ID: topic.ConceptID, Salience: topic.Salience})
		}
	}

	return res, nil
}
