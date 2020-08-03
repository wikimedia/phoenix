package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

func getPageName(head *goquery.Selection) (string, error) {
	title := head.Find("title").First()
	if len(title.Nodes) == 0 {
		return "", fmt.Errorf("No `title` tag found")
	}
	return title.Text(), nil
}

func getPageUrl(head *goquery.Selection) (string, error) {
	base := head.Find("base").First()
	if len(base.Nodes) == 0 {
		return "", fmt.Errorf("No `base` tag found")
	}
	url, exists := base.Attr("href")
	if !exists {

	}
	return url, nil
}

func getPageModified(head *goquery.Selection) (time.Time, error) {
	meta := head.Find("meta[property=\"dc:modified\"]").First()
	modifiedDate := meta.AttrOr("content", "")
	return time.Parse(time.RFC3339, modifiedDate)
}

func getPageSourceId(head *goquery.Selection) (int, error) {
	meta := head.Find("meta[property=\"mw:pageId\"]").First()
	if len(meta.Nodes) == 0 {
		return 0, fmt.Errorf("No `meta` tag with property=\"mw:pageId\" found.")
	}
	value, exists := meta.Attr("content")
	if !exists {
		return 0, fmt.Errorf("No `content` for `meta` tag with property=\"mw:pageId\" found")
	}
	sourceId, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Unable to parse sourceId: %s", err)
	}
	return sourceId, nil

}

func getPageSourceTimeUuid(head *goquery.Selection) (string, error) {
	meta := head.Find("meta[property=\"mw:TimeUuid\"]").First()
	if len(meta.Nodes) == 0 {
		return "", fmt.Errorf("No `meta` tag property=\"mw:TimeUuid\" found")
	}
	value, exists := meta.Attr("content")
	if !exists {
		return "", fmt.Errorf("No `content` for `meta` tag property=\"mw:TimeUuid\" found")
	}
	return value, nil
}

func getPageSourceRevision(html *goquery.Selection) (int, error) {
	about, exits := html.Attr("about")
	if !exits {
		return 0, fmt.Errorf("No `about` html ")
	}
	revisionStr := about[strings.LastIndex(about, "/")+1:]
	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		return 0, fmt.Errorf("Unable to parse revision: %s", err)
	}
	return revision, nil
}

func parseParsoidDocumentPage(document *goquery.Document) (page *common.Page, err error) {
	page = &common.Page{}
	page.HasPart = []string{}

	head := document.Find("html>head")

	if page.Name, err = getPageName(head); err != nil {
		return nil, err
	}

	if page.URL, err = getPageUrl(head); err != nil {
		return nil, err
	}

	if page.DateModified, err = getPageModified(head); err != nil {
		return nil, err
	}

	page.Source = common.Source{}

	if page.Source.ID, err = getPageSourceId(head); err != nil {
		return nil, err
	}

	if page.Source.TimeUUID, err = getPageSourceTimeUuid(head); err != nil {
		return nil, err
	}

	html := document.Find("html")
	if page.Source.Revision, err = getPageSourceRevision(html); err != nil {
		return nil, err
	}
	return page, nil

}
