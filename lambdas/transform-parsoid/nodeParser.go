package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

var (
	ignoredNodes = map[string]bool{
		"References": true,
	}
)

func getSectionName(section *goquery.Selection) string {
	return section.Find("h2").First().Text()
}

func parseParsoidDocumentNodes(document *goquery.Document, page *common.Page) ([]common.Node, error) {
	var err error
	var modified = page.DateModified
	var nameCounts = make(map[string]int)
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

		// Since it is possible for a document to have more than one section with the same heading text, keep
		// track of the number of times we've assigned a name, and de-duplicate if necessary.
		nameCounts[node.Name]++

		if nameCounts[node.Name] > 1 {
			node.Name = fmt.Sprintf("%s_%d", node.Name, nameCounts[node.Name])
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
