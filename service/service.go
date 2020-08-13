package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/rs/cors"
	"github.com/wikimedia/phoenix/common"
)

var mockPage = &common.Page{
	ID:           "abcdefghijklmn",
	Name:         "Foobar",
	URL:          "//en.wikipedia.org/wiki/Foobar",
	DateModified: time.Now(),
	HasPart: []string{
		"/page/385d6436a06b99d",
		"/page/644ed20cc75621c",
		"/page/42945840d44937c",
	},
	About: map[string]string{
		"//schema.org":        "/data/b6f7c05a-d367-11ea-af5c-2b020c033632",
		"//purl.org/dc/terms": "/data/b6f7c05a-d367-11ea-af5c-2b020c033632",
	},
}

// RootResolver is the top-level GraphQL resolver
type RootResolver struct{}

// Page returns a Page given its ID
func (r *RootResolver) Page(args struct{ ID graphql.ID }) (*PageResolver, error) {
	// TODO: This should retrieve the common.Page object from storage
	return &PageResolver{mockPage}, nil
}

// PageByName returns a Page given its Name
func (r *RootResolver) PageByName(args struct{ Name string }) (*PageResolver, error) {
	if args.Name == mockPage.Name {
		return &PageResolver{mockPage}, nil
	}
	return nil, nil
}

// PageResolver resolves a GraphQL page type
type PageResolver struct {
	p *common.Page
}

// ID resolves a page id attribute
func (r *PageResolver) ID() graphql.ID {
	return graphql.ID(r.p.ID)
}

// Name resolves a page name attribute
func (r *PageResolver) Name() string {
	return r.p.Name
}

// URL resolves a page url attribute
func (r *PageResolver) URL() string {
	return r.p.URL
}

// DateModified resolves a page dateModified attribute
func (r *PageResolver) DateModified() string {
	return r.p.DateModified.Format(time.RFC3339)
}

// HasPart resolves a page hasPart attribute
func (r *PageResolver) HasPart() []string {
	return r.p.HasPart
}

// About resolves a page about attribute
func (r *PageResolver) About() []*TupleResolver {
	res := make([]*TupleResolver, 0)
	for k, v := range r.p.About {
		res = append(res, &TupleResolver{key: k, val: v})
	}
	return res
}

// TupleResolver resolves a GraphQL Tuple type
type TupleResolver struct {
	key string
	val string
}

// Key resolves a tuple key attribute
func (r *TupleResolver) Key() string {
	return r.key
}

// Val resolves a tuple val attribute
func (r *TupleResolver) Val() string {
	return r.val
}

func main() {
	var b []byte
	var schema *graphql.Schema
	var err error

	// FIXME: Keeping the schema in a separate file could be convenient during early development (VS Code's
	// GraphQL editor addon has been pretty handy, for example), but at some point practical considerations
	// will probably dictate that we move it to a string variable in-code.
	if b, err = ioutil.ReadFile("schema.gql"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema.gql: %s", err)
		os.Exit(1)
	}

	if schema, err = graphql.ParseSchema(string(b), &RootResolver{}); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %s", err)
		os.Exit(1)
	}

	handler := cors.Default().Handler(&relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handler)))
}
