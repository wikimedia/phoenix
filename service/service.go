package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/handlers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/rs/cors"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	errorLogName  = flag.String("error-log", "error.log", "Path to the error log.")
	accessLogName = flag.String("access-log", "-", "Path to the access log.")

	// These values are passed in at build-time using -ldflags (see: Makefile)
	awsRegion string
	s3Bucket  string
)

// True if err is an awserr.Error, AND its code is s3.ErrCodeNoSuchKey, false otherwise.
func isS3NotFound(err error) bool {
	var s3err awserr.Error
	if errors.As(err, &s3err) {
		if s3err.Code() == s3.ErrCodeNoSuchKey {
			return true
		}
	}
	return false
}

// RootResolver is the top-level GraphQL resolver
type RootResolver struct {
	Repository *storage.Repository
	Logger     *common.Logger
}

// Page returns a Page given its ID
func (r *RootResolver) Page(args struct{ ID graphql.ID }) (*PageResolver, error) {
	var page *common.Page
	var err error

	if page, err = r.Repository.GetPage(string(args.ID)); err != nil {
		// If this was an error returned by S3 (it is an awserr.Error), and its code is s3.ErrCodeNoSuchKey
		// then the object was simply not found (read: this is not an error per say).
		if isS3NotFound(err) {
			return nil, nil
		}

		r.Logger.Error("Unable to retrieve Page (ID=%s): %s", string(args.ID), err)
		return nil, err
	}

	return &PageResolver{page}, nil
}

// PageByName returns a Page given its Name
func (r *RootResolver) PageByName(args struct {
	Authority string
	Name      string
}) (*PageResolver, error) {
	var page *common.Page
	var err error

	if page, err = r.Repository.GetPageByName(args.Authority, args.Name); err != nil {
		// If err is of type ErrNameNotFound, then this is not an error per say.
		if _, ok := err.(*storage.ErrNameNotFound); ok {
			return nil, nil
		}

		r.Logger.Error("Unable to retrieve Page (authority=%s, name=%s): %s", args.Authority, args.Name, err)
		return nil, err
	}

	return &PageResolver{page}, nil
}

// Node returns a Node given its ID
func (r *RootResolver) Node(args struct{ ID graphql.ID }) (*NodeResolver, error) {
	var node *common.Node
	var err error

	if node, err = r.Repository.GetNode(string(args.ID)); err != nil {
		// If this was an error returned by S3 (it is an awserr.Error) and its code is s3.ErrCodeNoSuchKey
		// then the object was simply not found (read: this is not an error per say).
		if isS3NotFound(err) {
			return nil, nil
		}

		r.Logger.Error("Unable to retrieve Node (ID=%s): %s", string(args.ID), err)
		return nil, err
	}

	return &NodeResolver{node}, nil
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

// Return configuration variables that are the union of defaults, and any values passed in the environment
func config() (region, bucket string) {
	// Retrieve environment variables
	env := func(name string, def string) string {
		if v := os.Getenv(name); v != "" {
			return v
		}
		return def
	}

	region = env("AWS_REGION", awsRegion)
	bucket = env("AWS_BUCKET", s3Bucket)

	return region, bucket
}

func main() {
	var accessLog io.Writer
	var b []byte
	var err error
	var logger *common.Logger
	var resolver *RootResolver
	var schema *graphql.Schema

	var region, bucket = config()
	var awsSession = session.New(&aws.Config{Region: aws.String(region)})

	flag.Parse()

	// Setup the error logger
	if logger, err = common.NewFileLogger("INFO", *errorLogName); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open error log: %s", err)
		os.Exit(1)
	}

	// Setup the access log
	if *accessLogName == "-" {
		accessLog = os.Stdout
	} else {
		if accessLog, err = os.OpenFile(*accessLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s for writing: %s\n", *accessLogName, err)
			os.Exit(1)
		}
	}

	// FIXME: Keeping the schema in a separate file could be convenient during early development (VS Code's
	// GraphQL editor addon has been pretty handy, for example), but at some point practical considerations
	// will probably dictate that we move it to a string variable in-code.
	if b, err = ioutil.ReadFile("schema.gql"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema.gql: %s", err)
		os.Exit(1)
	}

	resolver = &RootResolver{
		Repository: &storage.Repository{
			Store:  s3.New(awsSession),
			Index:  &storage.DynamoDBIndex{Client: dynamodb.New(awsSession)},
			Bucket: bucket,
		},
		Logger: logger,
	}

	if schema, err = graphql.ParseSchema(string(b), resolver, graphql.UseFieldResolvers()); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %s", err)
		os.Exit(1)
	}

	handler := cors.Default().Handler(&relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(accessLog, handler)))
}
