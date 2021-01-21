package main

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
	rosetteEndpoint    = "https://api.rosette.com/rest/v1/topics?redirect=true"
	rosetteMinSalience = 0.1
	rosetteMaxTopics   = 10
	rosetteRetries     = 10
)

var (
	// Regex for matching extraneous spaces
	xSpace = regexp.MustCompile(`\s+`)
)

// A helper for extracting text from an HTML snippet (very basic; Good Enough™ for now)
func extract(unsafe string) (string, error) {
	var doc *goquery.Document
	var err error

	if doc, err = goquery.NewDocumentFromReader(strings.NewReader(unsafe)); err != nil {
		return "", fmt.Errorf("failed to parse html string: %w", err)
	}

	return xSpace.ReplaceAllString(doc.Text(), " "), nil
}

type topicsResponse struct {
	Concepts []struct {
		ConceptID string  `json:"conceptId"`
		Phrase    string  `json:"phrase"`
		Salience  float32 `json:"salience"`
	} `json:"concepts"`
}

func requestTopics(text string) (*topicsResponse, error) {
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
			Content string `json:"content"`
		}{
			text,
		}

		if reqData, e = json.Marshal(content); e != nil {
			return nil, fmt.Errorf("unable to serialize JSON request body: %w", e)
		}

		// Create the request
		if req, e = http.NewRequest("POST", rosetteEndpoint, bytes.NewBuffer(reqData)); e != nil {
			return nil, fmt.Errorf("failed to create HTTP request: %w", e)
		}

		req.Header.Set("X-RosetteAPI-Key", rosetteAPIKey)
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
			log.Debug("Rosette returned status 429; reconnecting in %s (#%d/%d)", delay, i, rosetteRetries)
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

func rosetteTopics(node *common.Node) ([]common.RelatedTopic, error) {
	var content string
	var err error
	var res = make([]common.RelatedTopic, 0)
	var topics *topicsResponse

	if content, err = extract(node.Unsafe); err != nil {
		return nil, err
	}

	if topics, err = requestTopics(content); err != nil {
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
