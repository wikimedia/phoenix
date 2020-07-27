package main

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/wikimedia/phoenix/common"
)

func getSectionName(section *goquery.Selection) (string, error) {
	name := section.Find("h2").First().Text()
	return name, nil
}

func parseParsoidDocumentNodes(document *goquery.Document, modified time.Time) (nodes []common.Node, err error) {
	nodes = []common.Node{}

	sections := document.Find("html>body>section[data-mw-section-id]")
	for i := range sections.Nodes {
		section := sections.Eq(i)

		node := common.Node{}
		node.ID = uuid.New().String()
		if node.Name, err = getSectionName(section); err != nil {
			continue
		}
		if node.Unsafe, err = section.Html(); err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
