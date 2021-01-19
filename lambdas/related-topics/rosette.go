package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

const (
	rosetteEndpoint    = "https://api.rosette.com/rest/v1/topics?redirect=true"
	rosetteMinSalience = 0.1
	rosetteMaxTopics   = 10
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
	var client = &http.Client{}
	var err error
	var req *http.Request
	var reqData []byte
	var res *http.Response
	var resData []byte
	var topics topicsResponse

	// Serialize a requests body (JSON)
	content := &struct {
		Content string `json:"content"`
	}{
		text,
	}

	if reqData, err = json.Marshal(content); err != nil {
		return nil, fmt.Errorf("unable to serialize JSON request body: %w", err)
	}

	// Create the request
	if req, err = http.NewRequest("POST", rosetteEndpoint, bytes.NewBuffer(reqData)); err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("X-RosetteAPI-Key", rosetteAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	// Send the request
	if res, err = client.Do(req); err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code %d (expected %d)", res.StatusCode, http.StatusOK)
	}

	defer res.Body.Close()

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
