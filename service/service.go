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

var (
	mockPage = &common.Page{
		ID:           "/page/abcdefghijklmn",
		Name:         "Foobar",
		URL:          "//en.wikipedia.org/wiki/Foobar",
		DateModified: time.Now(),
		HasPart: []string{
			"/node/365154aa-de4a-11ea-a27b-33aa6523fd57",
		},
		About: map[string]string{
			"//schema.org":        "/data/b6f7c05a-d367-11ea-af5c-2b020c033632",
			"//purl.org/dc/terms": "/data/b6f7c05a-d367-11ea-af5c-2b020c033632",
		},
	}
	mockNode = &common.Node{
		ID:           "/node/365154aa-de4a-11ea-a27b-33aa6523fd57",
		Name:         "",
		IsPartOf:     []string{"/page/abcdefghijklmn"},
		DateModified: time.Now(),
		Unsafe:       "<p>The rain in Spain falls mostly on the plains.</p>",
	}
)

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

// Node returns a Node given its ID
func (r *RootResolver) Node(args struct{ ID graphql.ID }) (*NodeResolver, error) {
	return &NodeResolver{mockNode}, nil
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

// NodeResolver resolves a GraphQL node type
type NodeResolver struct {
	n *common.Node
}

// ID resolves a node id attribute
func (r *NodeResolver) ID() graphql.ID {
	return graphql.ID(r.n.ID)
}

// Name resolves a node name attribute
func (r *NodeResolver) Name() string {
	return r.n.Name
}

// IsPartOf resolves a node isPartOf attribute
func (r *NodeResolver) IsPartOf() []graphql.ID {
	parents := make([]graphql.ID, len(r.n.IsPartOf))
	for _, id := range r.n.IsPartOf {
		parents = append(parents, graphql.ID(id))
	}
	return parents
}

// DateModified resolves a node dateModified attribute
func (r *NodeResolver) DateModified() string {
	return r.n.DateModified.Format(time.RFC3339)
}

// Unsafe resolves a node unsafe attribute
func (r *NodeResolver) Unsafe() string {
	return r.n.Unsafe
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

	if schema, err = graphql.ParseSchema(string(b), &RootResolver{}, graphql.UseFieldResolvers()); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %s", err)
		os.Exit(1)
	}

	handler := cors.Default().Handler(&relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handler)))
}
