package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Index is an interface for indexing Phoenix documents
type Index interface {
	// Apply updates the index with new Phoenix document data
	Apply(update *Update) error

	// PageIDForName queries the index for page ID matching name
	PageIDForName(authority, name string) (string, error)

	// NodeIDForName queries the index for node ID matching name
	NodeIDForName(authority, pageName, name string) (string, error)
}

// ErrPageNotFound is an Error returned when a lookup by name fails
type ErrPageNotFound struct {
	Authority string
	Name      string
}

func (e *ErrPageNotFound) Error() string {
	return fmt.Sprintf(`("%s", "%s"): not found`, e.Authority, e.Name)
}

// ErrNodeNotFound is an Error returned when a lookup by name fails
type ErrNodeNotFound struct {
	Authority string
	PageName  string
	Name      string
}

func (e *ErrNodeNotFound) Error() string {
	return fmt.Sprintf(`("%s", "%s", "%s"): not found`, e.Authority, e.PageName, e.Name)
}

// MockIndex is a memory-backed Index used in testing
type MockIndex struct {
	pages map[string]string
	nodes map[string]string
}

// Apply updates the index with new Phoenix document data
func (i *MockIndex) Apply(update *Update) error {
	page := update.Page
	nodes := update.Nodes

	i.pages[fmt.Sprintf("%s:%s", page.Source.Authority, page.Name)] = page.ID

	for _, n := range nodes {
		i.nodes[fmt.Sprintf("%s:%s:%s", n.Source.Authority, page.Name, n.Name)] = n.ID
	}

	return nil
}

// PageIDForName queries the index for page ID matching name
func (i *MockIndex) PageIDForName(authority, name string) (string, error) {
	if v, ok := i.pages[fmt.Sprintf("%s:%s", authority, name)]; ok {
		return v, nil
	}

	return "", &ErrPageNotFound{Authority: authority, Name: name}
}

// NodeIDForName queries the index for node ID matching name
func (i *MockIndex) NodeIDForName(authority, pageName, name string) (string, error) {
	if v, ok := i.nodes[fmt.Sprintf("%s:%s:%s", authority, pageName, name)]; ok {
		return v, nil
	}

	return "", &ErrNodeNotFound{Authority: authority, PageName: pageName, Name: name}
}

// NewMockIndex creates a new MockIndex
func NewMockIndex() *MockIndex {
	return &MockIndex{make(map[string]string), make(map[string]string)}
}

// DynamoDBIndex is a Phoenix document indexer backed by DynamoDB
type DynamoDBIndex struct {
	Client      *dynamodb.DynamoDB
	TitlesTable string
	NamesTable  string
}

// Apply updates the index with new Phoenix document data
func (i *DynamoDBIndex) Apply(update *Update) error {
	var items []*dynamodb.TransactWriteItem
	var page = update.Page

	items = make([]*dynamodb.TransactWriteItem, 0)

	// Page titles
	items = append(items, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			Item: map[string]*dynamodb.AttributeValue{
				"Title":     {S: aws.String(page.Name)},
				"Authority": {S: aws.String(page.Source.Authority)},
				"ID":        {S: aws.String(page.ID)},
			},
			TableName: aws.String(i.TitlesTable),
		},
	})

	// Node names
	for _, n := range update.Nodes {
		items = append(items, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item: map[string]*dynamodb.AttributeValue{
					"Name":      {S: aws.String(encodeNodeName(page.Name, n.Name))},
					"Authority": {S: aws.String(n.Source.Authority)},
					"ID":        {S: aws.String(n.ID)},
				},
				TableName: aws.String(i.NamesTable),
			},
		})
	}

	_, err := i.Client.TransactWriteItems(&dynamodb.TransactWriteItemsInput{TransactItems: items})

	if err != nil {
		return err
	}

	return nil
}

// PageIDForName queries the index for page ID matching authority (wiki) and name
func (i *DynamoDBIndex) PageIDForName(authority, name string) (string, error) {
	result, err := i.Client.GetItem(
		&dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Title":     {S: aws.String(name)},
				"Authority": {S: aws.String(authority)},
			},
			TableName: aws.String(i.TitlesTable),
		})

	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", &ErrPageNotFound{Authority: authority, Name: name}
	}

	return *result.Item["ID"].S, nil
}

// NodeIDForName queries the index for page ID matching authority (wiki) and name
func (i *DynamoDBIndex) NodeIDForName(authority, pageName, name string) (string, error) {
	result, err := i.Client.GetItem(
		&dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Name":      {S: aws.String(encodeNodeName(pageName, name))},
				"Authority": {S: aws.String(authority)},
			},
			TableName: aws.String(i.NamesTable),
		})

	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", &ErrNodeNotFound{Authority: authority, PageName: pageName, Name: name}
	}

	return *result.Item["ID"].S, nil
}

func encodeNodeName(pageName, name string) string {
	return fmt.Sprintf("%s:%s", url.QueryEscape(pageName), url.QueryEscape(strings.ToLower(name)))
}

// ElasticsearchIndex is a Phoenix document indexer backed by Elasticsearch
type ElasticsearchIndex struct {
	Client *elasticsearch.Client
}

// Apply updates the index with new Phoenix document data
func (i *ElasticsearchIndex) Apply(update *Update) error {
	page := update.Page

	type document struct {
		ID string `json:"id"`
	}

	var b []byte
	var err error
	var res *esapi.Response

	if b, err = json.Marshal(document{ID: page.ID}); err != nil {
		return fmt.Errorf("unable to marshal json document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "page_name",
		DocumentID: url.PathEscape(fmt.Sprintf("%s:%s", page.Source.Authority, page.Name)),
		Body:       strings.NewReader(string(b)),
		Refresh:    "true",
	}

	if res, err = req.Do(context.Background(), i.Client); err != nil {
		return fmt.Errorf("Elasticsearch response error: %w", err)
	}

	defer res.Body.Close()

	// FIXME: Index node names too!

	if res.IsError() {
		return fmt.Errorf("error indexing %s (status=%s)", page.Name, res.Status())
	}

	return nil
}

// PageIDForName queries the index for page ID matching name
func (i *ElasticsearchIndex) PageIDForName(authority, name string) (string, error) {
	var err error
	var res *esapi.Response

	req := esapi.GetRequest{Index: "page_name", DocumentID: fmt.Sprintf("%s:%s", authority, name)}
	if res, err = req.Do(context.Background(), i.Client); err != nil {
		return "", fmt.Errorf("Elasticsearch response error: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return "", &ErrPageNotFound{Name: name}
		}
		return "", fmt.Errorf("unknown error retrieving %s (status=%s)", name, res.Status())
	}

	type response struct {
		Source struct {
			ID string `json:"id"`
		} `json:"_source"`
	}

	var r response
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("unable to decode JSON response: %w", err)
	}

	return r.Source.ID, nil
}

// NodeIDForName queries the index for page ID matching name
func (i *ElasticsearchIndex) NodeIDForName(authority, pageName, name string) (string, error) {
	// TODO: Do; Stubbed!
	return "", nil
}
