package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseParsoid(r io.Reader) (page *Page, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	//page := newPage()
	sections := doc.Find("html>body>section[data-mw-section-id]")
	for i := range sections.Nodes {

		section := sections.Eq(i)
		dataId, _ := section.Attr("data-mw-section-id")
		fmt.Printf("Section %d - %s:\n", i, dataId)
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
		} else {
			id, exists := name.Attr("id")
			if !exists {
				continue
			}
			fmt.Println(id)
			if id == "References" {
				var div = section.Find("div").First()
				data, exists := div.Attr("data-mw")
				if !exists {
					continue
				}

				references := strings.Split(data, "\\n\\n")
				for _, reference := range references {
					fmt.Println(reference)
				}

			} else {
				page.AddParagraph(newParagraph(id, content.Children))
			}
		}

	}
	return page, nil
}

func parseReferences(reference *goquery.Selection) {

}
