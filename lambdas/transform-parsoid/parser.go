package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

func parseParsoidDocument(document *goquery.Document) (page *common.Page, nodes []common.Node, err error) {

	page, err = parseParsoidDocumentPage(document)
	if err != nil {
		return nil, nil, err
	}

	nodes, err = parseParsoidDocumentNodes(document, page.DateModified)
	if err != nil {
		return nil, nil, err
	}
	return page, nodes, nil
}
