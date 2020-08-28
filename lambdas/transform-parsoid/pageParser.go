package main

import (
	"fmt"
	"net/url"
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

func getPageUrl(head *goquery.Selection) (*url.URL, error) {
	var base *goquery.Selection
	var err error
	var exists bool
	var href string
	var urlObj *url.URL

	base = head.Find("base").First()

	if len(base.Nodes) == 0 {
		return nil, fmt.Errorf("No element named 'base' found")
	}

	if href, exists = base.Attr("href"); !exists {
		return nil, fmt.Errorf("'base' element contained no 'html' attribute")
	}

	if urlObj, err = url.Parse(href); err != nil {
		return nil, fmt.Errorf("Unable to parse URL: %w", err)
	}

	return urlObj, nil
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

func parseParsoidDocumentPage(document *goquery.Document) (*common.Page, error) {
	var head, html *goquery.Selection
	var page = &common.Page{}
	var pageURL *url.URL
	var err error

	page.HasPart = make([]string, 0)
	page.Source = common.Source{}

	head = document.Find("html>head")

	if page.Name, err = getPageName(head); err != nil {
		return nil, err
	}

	if pageURL, err = getPageUrl(head); err != nil {
		return nil, err
	}

	page.URL = pageURL.String()

	if page.DateModified, err = getPageModified(head); err != nil {
		return nil, err
	}

	if page.Source.ID, err = getPageSourceId(head); err != nil {
		return nil, err
	}

	if page.Source.TimeUUID, err = getPageSourceTimeUuid(head); err != nil {
		return nil, err
	}

	html = document.Find("html")

	if page.Source.Revision, err = getPageSourceRevision(html); err != nil {
		return nil, err
	}

	page.Source.Authority = pageURL.Hostname()

	return page, nil
}
