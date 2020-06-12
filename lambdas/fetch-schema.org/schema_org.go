package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// URL for the Wikidata Query Service
const wdqsBase string = "https://query.wikidata.org/sparql"

// Thing represents an https://schema.org/Thing as JSON-LD
type Thing struct {
	Context       string `json:"@context"`
	Type          string `json:"@type"`
	AlternateName string `json:"alternateName,omitempty"`
	Description   string `json:"description,omitempty"`
	Image         string `json:"image,omitempty"`
	Name          string `json:"name,omitempty"`
	SameAs        string `json:"sameAs,omitempty"`
	// FIXME: What about these?
	// URL           string
	// Identifier    string
}

// NewThing returns a Thing struct with the @context and @type properties pre-filled.
func NewThing() *Thing {
	return &Thing{Context: "https://schema.org", Type: "Thing"}
}

func schemaOrgItem(item string) (*Thing, error) {
	// Construct the URL w/ query string
	var query strings.Builder
	fmt.Fprintf(&query, "%s?format=json&query=", wdqsBase)

	// Sparql query string
	sparql := fmt.Sprintf(`
		SELECT DISTINCT ?item ?itemLabel ?image ?itemDescription ?alias WHERE {
			BIND(wd:%s AS ?item)
			?item skos:altLabel ?alias. filter(lang(?alias)="en")
			OPTIONAL { ?item wdt:P18 ?image. }
			SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en". }
		}
	`, item)

	query.WriteString(url.QueryEscape(sparql))

	// Perform the HTTP request
	var body []byte
	var err error

	if body, err = request(query.String()); err != nil {
		return nil, fmt.Errorf("Error making HTTP request: %w", err)
	}

	// Unmarshal request body to JSON
	type value struct {
		Lang  string `json:"xml:lang"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	res := struct {
		Results struct {
			Bindings []struct {
				Item        value
				Alias       value
				Label       value `json:"itemLabel"`
				Description value `json:"itemDescription"`
				Image       value
			}
		}
	}{}

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON response: %w", err)
	}

	// This is our Thing
	thing := NewThing()

	// If there are no results, just return the empty struct
	if len(res.Results.Bindings) > 0 {
		binding := res.Results.Bindings[0]
		thing.AlternateName = binding.Alias.Value
		thing.Description = binding.Description.Value
		thing.Image = binding.Image.Value
		thing.Name = binding.Label.Value
		thing.SameAs = binding.Item.Value
	}

	return thing, nil
}
