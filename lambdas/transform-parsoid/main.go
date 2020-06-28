package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type SectionType uint32

const (
	HeaderSection SectionType = iota
	ParagrpaphSection
	ReferencesSection
)

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

func sectionType(section *goquery.Selection) {

}

func parseReferences(reference *goquery.Selection) {

}

func processParsoid(domain string, title string) (page *Page, err error) {
	url := urlf(domain, title)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return parseParsoid(res.Body)
}

func tokenizeHtml(r io.Reader) (content *Content, err error) {
	z := html.NewTokenizer(r)
	parent := newContent("page", nil)
	for {
		switch z.Next() {
		case html.ErrorToken:
			//bytes, _ := json.Marshal(parent)
			//fmt.Println(z.Err())
			return parent, nil
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

func requestParsoid(domain string, title string) (content *Content, err error) {
	res, err := http.Get(urlf(domain, title))
	if err != nil {
		return
	}
	defer res.Body.Close()
	return tokenizeHtml(res.Body)

}

func main() {
	res, err := processParsoid("en.wikipedia.org", "Mars")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	//requestParsoid("en.wikipedia.org", "United_States")
}
