package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type SectionType uint32

const (
	HeaderSection SectionType = iota
	ParagrpaphSection
	ReferencesSection
)

type Paragraph struct {
	Name    string     `json:"name"`
	Content []*Content `json:"content"`
}

func newParagraph(name string, content []*Content) *Paragraph {
	return &Paragraph{name, content}
}

type Header struct {
	Content []*Content `json:"content"`
}

type Reference struct {
}

type Page struct {
	Header     []*Content   `json:"header,omitempty"`
	Paragraphs []*Paragraph `json:"paragraphs,omitempty"`
	References []*Reference `json:"references,omitempty"`
}

func newPage() *Page {
	new := Page{nil, []*Paragraph{}, []*Reference{}}
	return &new
}

func (page *Page) AddParagraph(paragraph *Paragraph) {
	page.Paragraphs = append(page.Paragraphs, paragraph)
}

func (page *Page) SetHeader(header []*Content) {
	page.Header = header
}

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

func sectionType(section *goquery.Selection) {

}

func parseParsoid(r io.Reader) (page *Page, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	//page := newPage()
	sections := doc.Find("html body section[data-mw-section-id]")
	for i := range sections.Nodes {
		fmt.Printf("Review %d:\n", i)
		section := sections.Eq(i)
		html, err := section.Html()
		if err != nil {
			continue
		}

		content, err := tokenizeHtml(strings.NewReader(html))
		if err != nil {
			continue
		}
		page := newPage()
		name := section.Find("h2").First()
		if len(name.Nodes) == 0 {
			page.SetHeader(content.Children)
			fmt.Println("header")
		} else {
			id, exists := name.Attr("id")
			if !exists {
				continue
			}
			if id == "Reference" {

			} else {
				page.AddParagraph(newParagraph(name.Text(), content.Children))
			}
		}

	}
	return page, nil
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
