package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// Content represents a piece of content on a page
type Content struct {
	Schema   string     `json:"schema,omitempty"`
	Rel      string     `json:"rel,omitempty"`
	TypeOf   string     `json:"typeOf,omitempty"`
	Text     string     `json:"text,omitempty"`
	DataMW   string     `json:"attributes,omitempty"`
	Parent   *Content   `json:"-"`
	Children []*Content `json:"content,omitempty"`
}

// AddChild to content
func (content *Content) AddChild(child *Content) {
	content.Children = append(content.Children, child)
}

var includedAttributes = map[string]bool{
	"rel":     true,
	"typeof":  true,
	"data-mw": true,
}

func getAttribtues(token html.Token) map[string]string {
	attributes := map[string]string{}
	for _, attr := range token.Attr {
		if !includedAttributes[attr.Key] {
			continue
		}
		// multiple attributes can have the same name so this might overwrite but not sure if it matters for our use case
		attributes[attr.Key] = attr.Val
	}
	return attributes
}

//AddAttribtues to content
func (content *Content) AddAttribtues(token html.Token) {
	attributes := getAttribtues(token)
	content.Rel = attributes["rel"]
	content.TypeOf = attributes["typeof"]
	content.DataMW = attributes["data-mw"]
}

func newContent(tagName string, parent *Content) *Content {
	new := Content{tagName, "", "", "", "", parent, []*Content{}}
	if parent != nil {
		parent.AddChild(&new)
	}
	return &new
}

// Use the existing endpoint for now for testing
func urlf(domain string, title string) string {
	return fmt.Sprintf("https://%s/api/rest_v1/page/html/%s", domain, url.PathEscape(title))
}

var ignoredTags = map[string]bool{
	"html":   true,
	"head":   true,
	"meta":   true,
	"link":   true,
	"body":   true,
	"script": true,
	"style":  true,
	"base":   true,
	"title":  true,
}

func isIgnored(token html.Token) bool {
	return ignoredTags[token.Data]
}

func requestParsoid(domain string, title string) (parsoidHTML string, err error) {
	res, err := http.Get(urlf(domain, title))
	if err != nil {
		return
	}
	defer res.Body.Close()
	z := html.NewTokenizer(res.Body)
	parent := newContent("page", nil)
	for {
		switch z.Next() {
		case html.ErrorToken:
			bytes, _ := json.Marshal(parent)
			fmt.Println(string(bytes))
			return "", z.Err()
		case html.TextToken:
			text := newContent("text", parent)
			text.Text = string(z.Text())
			continue
		case html.StartTagToken:
			token := z.Token()
			if isIgnored(token) {
				continue
			}
			parent = newContent(token.Data, parent)
			parent.AddAttribtues(token)
			continue
		case html.EndTagToken:
			token := z.Token()
			if isIgnored(token) {
				continue
			}
			parent = parent.Parent
			continue
		case html.SelfClosingTagToken:
			token := z.Token()
			if isIgnored(token) {
				continue
			}
			tag := newContent(token.Data, parent)
			tag.AddAttribtues(token)
			continue
		}
	}
}

func main() {
	requestParsoid("en.wikipedia.org", "United_States")
}
