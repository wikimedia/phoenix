package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wikimedia/phoenix/common"
)

func getPageName(head *goquery.Selection) (string, error) {
	title := head.Find("title").First().Text()
	return title, nil
}

func getPageUrl(head *goquery.Selection) (string, error) {
	base := head.Find("base").First()
	url := base.AttrOr("href", "")
	return url, nil
}

func getPageModified(head *goquery.Selection) (time.Time, error) {
	meta := head.Find("meta[property=\"dc:modified\"]").First()
	modifiedDate := meta.AttrOr("content", "")
	return time.Parse(time.RFC3339, modifiedDate)
}

func getPageSourceId(head *goquery.Selection) (int, error) {
	meta := head.Find("meta[property=\"mw:pageId\"]").First()
	value := meta.AttrOr("content", "")
	return strconv.Atoi(value)
}

func getPageSourceTimeUuid(head *goquery.Selection) (string, error) {
	meta := head.Find("meta[property=\"mw:TimeUuid\"]").First()
	value := meta.AttrOr("content", "")
	return value, nil
}

func getPageSourceRevision(html *goquery.Selection) (int, error) {
	about := html.AttrOr("about", "")
	revision := about[strings.LastIndex(about, "/")+1:]
	return strconv.Atoi(revision)
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
