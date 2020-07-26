package main

import (
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/wikimedia/phoenix/common"
)

var ignoredSections = map[string]bool{
	"html": true,
}

func getNameforSection(section *goquery.Selection) string {
	ret := section.Find("h2").First().Text()
	return ret
}

func getNameForPage(head *goquery.Selection) string {
	ret := head.Find("title").First().Text()
	return ret
}

func getUrlForPage(head *goquery.Selection) string {
	base := head.Find("base").First()
	ret := base.AttrOr("href", "")
	return ret
}

func getModifiedForPage(head *goquery.Selection) common.JSONTime {
	meta := head.Find("meta[property=\"dc:modified\"]").First()
	value := meta.AttrOr("content", "")
	date, _ := time.Parse(time.RFC3339, value)
	return common.JSONTime(date)
}

func getIdForSource(head *goquery.Selection) int {
	meta := head.Find("meta[property=\"mw:pageId\"]").First()
	value := meta.AttrOr("content", "")
	id, _ := strconv.Atoi(value)
	return id
}

func getTimeUuidForSource(head *goquery.Selection) string {
	meta := head.Find("meta[property=\"mw:TimeUuid\"]").First()
	value := meta.AttrOr("content", "")
	return value
}

func getRevisionForSource(head *goquery.Selection) string {
	meta := head.Find("meta[property=\"mw:revisionSHA1\"]").First()
	value := meta.AttrOr("content", "")
	return value
}

func parseParsoidDocument(document *goquery.Document) (page *common.Page, sections []*common.Section, err error) {

	page = &common.Page{}
	page.ID = uuid.Must(uuid.NewUUID()).String()
	page.HasPart = []string{}

	head := document.Find("html>head")
	page.Name = getNameForPage(head)
	page.URL = getUrlForPage(head)

	page.About = map[string]string{
		"@type": "Article",
		"name":  page.Name,
	}
	page.DateModified = getModifiedForPage(head)
	page.Source = common.Source{}
	page.Source.ID = getIdForSource(head)
	page.Source.TimeUUID = getTimeUuidForSource(head)
	page.Source.Revision = getRevisionForSource(head)

	sections = []*common.Section{}
	metaSelections := document.Find("html>head>meta")
	for i := range metaSelections.Nodes {
		metaNode := metaSelections.Eq(i)
		propertyName, _ := metaNode.Attr("property")
		value, _ := metaNode.Attr("content")
		if propertyName == "dc:modified" {
			date, _ := time.Parse(time.RFC3339, value)
			page.DateModified = common.JSONTime(date)
		}
		println(propertyName + " - " + value)
	}

	sectionsSelection := document.Find("html>body>section[data-mw-section-id]")
	for i := range sectionsSelection.Nodes {
		sectionSelection := sectionsSelection.Eq(i)
		htmlStr, _ := sectionSelection.Html()
		sec := common.Section{
			// Globally unique identifier
			ID: uuid.New().String(),

			// Section name (corresponds with schema.org/Thing#name).  Corresponds to the text of the first header
			// of a section in Parsoid HTML output.
			Name: getNameforSection(sectionSelection),

			// URLs of content that this section is a part of.  Loosely corresponds with
			// schema.org/CreativeWork#isPartOf, yet unlike its namesake, this attribute serves as an adjacency
			// list of nodes in the document graph.
			IsPartOf: []string{page.ID},

			// Date and time of last modification (corresponds with schema.org/CreativeWork#dateModified)
			DateModified: page.DateModified,

			// The raw HTML context of the corresponding section.
			Unsafe: htmlStr,
		}
		page.HasPart = append(page.HasPart, sec.ID)
		sections = append(sections, &sec)
	}
	return page, sections, nil
}

// func parseReferences(reference *goquery.Selection) {

// }
