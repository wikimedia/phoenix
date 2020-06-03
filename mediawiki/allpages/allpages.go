package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var titleFrom = flag.String("from", "", "Start iterating from closest matching title")
var serverName = flag.String("server-name", "simple.wikipedia.org", "Wiki server name")

const query = "/w/api.php?action=query&format=json&prop=revisions&rvprop=ids&generator=allpages&gapnamespace=0&gapfilterredir=nonredirects&gaplimit=100"

type revision struct {
	RevID    int `json:"revid"`
	ParentID int `json:"parentid"`
}

type page struct {
	PageID    int        `json:"pageid"`
	NS        int        `json:"ns"`
	Title     string     `json:"title"`
	Revisions []revision `json:"revisions"`
}

type continuation struct {
	GapContinue string `json:"gapcontinue"`
	Continue    string `json:"continue"`
}

type pages struct {
	BatchComplete string       `json:"batchcomplete"`
	Continue      continuation `json:"continue"`
	Query         struct {
		Pages map[string]page `json:"pages"`
	}
}

type allPages struct {
	From     string
	Continue continuation
}

func (a *allPages) hasMore() bool {
	if a.Continue.GapContinue == "" && a.Continue.Continue == "" {
		return false
	}
	return true
}

func (a *allPages) url() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "https://%s/%s", *serverName, query)

	if a.From != "" {
		fmt.Fprintf(&builder, "&gapfrom=%s", url.PathEscape(a.From))
	}
	if a.hasMore() {
		fmt.Fprintf(&builder, "&continue=%s&gapcontinue=%s", url.PathEscape(a.Continue.Continue), url.PathEscape(a.Continue.GapContinue))
	}

	return builder.String()
}

func (a *allPages) fetch() (*pages, error) {
	res, err := http.Get(a.url())
	if err != nil {
		return nil, fmt.Errorf("HTTP GET: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	obj := pages{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	a.Continue = obj.Continue
	return &obj, nil
}

func main() {
	flag.Parse()

	api := &allPages{From: *titleFrom}

	printPages := func() {
		// Fetch (another) resultset
		resp, err := api.fetch()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Iterate and print
		for _, v := range resp.Query.Pages {
			if len(v.Revisions) > 1 {
				panic(fmt.Sprintf("Too many revisions (%d != 1)!", len(v.Revisions)))
			}
			fmt.Println(v.Revisions[0].RevID, v.Title)
		}
	}

	printPages()

	for api.hasMore() {
		printPages()
	}

}
