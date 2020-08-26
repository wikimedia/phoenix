package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/wikimedia/phoenix/common"
)

// Index is an interface for indexing Phoenix documents
type Index interface {
	// Apply updates the index with new Phoenix document data
	Apply(page *common.Page) error

	// PageIDForName queries the index for page ID matching name
	PageIDForName(authority, name string) (string, error)
}

// ErrNameNotFound is an Error returned when a lookup by name fails
type ErrNameNotFound struct {
	Name string
}

func (e *ErrNameNotFound) Error() string { return e.Name + ": not found" }

// MockIndex is a memory-backed Index used in testing
type MockIndex struct {
	pages map[string]string
}

// Apply updates the index with new Phoenix document data
func (i *MockIndex) Apply(page *common.Page) error {
	i.pages[fmt.Sprintf("%s:%s", page.Source.Authority, page.Name)] = page.ID
	return nil
}

// PageIDForName queries the index for page ID matching name
func (i *MockIndex) PageIDForName(authority, name string) (string, error) {
	if v, ok := i.pages[fmt.Sprintf("%s:%s", authority, name)]; ok {
		return v, nil
	}

	return "", &ErrNameNotFound{Name: name}
}

// NewMockIndex creates a new MockIndex
func NewMockIndex() *MockIndex {
	return &MockIndex{make(map[string]string)}
}

// ElasticsearchIndex is a Phoenix document indexer backed by Elasticsearch
type ElasticsearchIndex struct {
	Client *elasticsearch.Client
}

// Apply updates the index with new Phoenix document data
func (i *ElasticsearchIndex) Apply(page *common.Page) error {
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
			return "", &ErrNameNotFound{Name: name}
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
