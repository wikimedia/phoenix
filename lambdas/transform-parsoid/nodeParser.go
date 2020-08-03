package main

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
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

func parseParsoidDocumentNodes(document *goquery.Document, modified time.Time) (nodes []common.Node, err error) {
	nodes = []common.Node{}

	sections := document.Find("html>body>section[data-mw-section-id]")
	for i := range sections.Nodes {
		section := sections.Eq(i)

		node := common.Node{}
		node.ID = uuid.New().String()
		node.Name = getSectionName(section)

		if val, ok := ignoredNodes[node.Name]; !(val && ok) {
			continue
		}

		if node.Unsafe, err = section.Html(); err != nil {
			return []common.Node{}, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
