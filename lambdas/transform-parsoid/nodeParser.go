package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

var (
	ignoredNodes = map[string]bool{
		"References": true,
	}
)

func getSectionName(section *goquery.Selection) string {
	var id string
	var exists bool
	var header = section.Find("h2").First()

	if id, exists = header.Attr("id"); exists {
		return id
	}

	return header.Text()
}

func parseParsoidDocumentNodes(document *goquery.Document, page *common.Page) ([]common.Node, error) {
	var err error
	var modified = page.DateModified
	var nodes = make([]common.Node, 0)
	var sections = document.Find("html>body>section[data-mw-section-id]")

	for i := range sections.Nodes {
		var node = common.Node{}
		var section = sections.Eq(i)
		var unsafe string

		node.Source = page.Source

		// If this is the first section and the name is a zero length string, then we assign it
		// a constant to simplify lookups
		if i == 0 {
			if name := getSectionName(section); name == "" {
				node.Name = leadSectionName
			} else {
				node.Name = name
			}
		} else {
			node.Name = getSectionName(section)
		}

		node.DateModified = modified

		if val, ok := ignoredNodes[node.Name]; ok && val {
			continue
		}

		if unsafe, err = section.Html(); err != nil {
			return []common.Node{}, err
		}

		node.Unsafe = unsafe
		nodes = append(nodes, node)
	}

	return nodes, nil
}
