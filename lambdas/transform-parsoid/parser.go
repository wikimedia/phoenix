package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

func parseParsoidDocument(document *goquery.Document) (*common.Page, []common.Node, error) {
	var err error
	var page *common.Page
	var nodes []common.Node

	if page, err = parseParsoidDocumentPage(document); err != nil {
		return nil, nil, err
	}

	if nodes, err = parseParsoidDocumentNodes(document, page); err != nil {
		return nil, nil, err
	}

	return page, nodes, nil
}
