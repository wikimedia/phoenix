package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var errNotFound = errors.New("wikibase_item not found")

// Returns the wikibase_item page property for an article
func wikibaseItemID(serverName string, title string) (string, error) {
	// Action API query parameters
	parameters := map[string]string{
		"action":        "query",
		"format":        "json",
		"prop":          "pageprops",
		"ppprop":        "wikibase_item",
		"formatversion": "2",
	}

	parameters["titles"] = url.PathEscape(title)

	// Build up the query string
	var query strings.Builder
	for k, v := range parameters {
		fmt.Fprintf(&query, "%s=%s&", k, v)
	}

	var body []byte
	var err error

	// Fetch the API results
	if body, err = request(fmt.Sprintf("https://%s/w/api.php?%s", serverName, query.String())); err != nil {
		return "", fmt.Errorf("Action API request error: %w", err)
	}

	// Just the minimal JSON
	res := struct {
		Query struct {
			Pages []struct {
				Pageprops struct {
					WikibaseItem string `json:"wikibase_item"`
				} `json:"pageprops"`
			} `json:"pages"`
		} `json:"query"`
	}{}

	if err = json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("Unable to unmarshal JSON response: %w", err)
	}

	if len(res.Query.Pages) < 1 {
		return "", errNotFound
	}

	return res.Query.Pages[0].Pageprops.WikibaseItem, nil
}
