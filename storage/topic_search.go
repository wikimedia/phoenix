package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/wikimedia/phoenix/common"
)

type UpdateStats esutil.BulkIndexerStats

// TopicSearch is an interface for topic search indexing.
type TopicSearch interface {
	// Search queries the index for nodes matching a Wikidata ID
	Search(qid string) ([]string, error)

	// Update applies changes to the topic index
	Update(node *common.Node, topics []common.RelatedTopic) (*UpdateStats, error)
}

// ElasticTopicSearch is an Elasticsearch implementation of the TopicSearch interface.
type ElasticTopicSearch struct {
	Client    *elasticsearch.Client
	IndexName string
}

// Search queries the index for nodes matching a Wikidata ID
func (t ElasticTopicSearch) Search(qid string) ([]string, error) {
	var err error
	var ids = make([]string, 0)
	var req esapi.Request
	var res *esapi.Response

	if req, err = t.searchRequest(qid); err != nil {
		return nil, err
	}

	if res, err = req.Do(context.Background(), t.Client); err != nil {
		return nil, fmt.Errorf("Elasticsearch response error: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed search for %s (status=%s)", qid, res.Status())
	}

	var r map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		ids = append(ids, hit.(map[string]interface{})["_source"].(map[string]interface{})["node_id"].(string))
	}

	return ids, nil
}

// Update applies changes to the topic index
func (t ElasticTopicSearch) Update(node *common.Node, topics []common.RelatedTopic) (*UpdateStats, error) {
	var err error
	var indexer esutil.BulkIndexer
	var req esapi.Request
	var res *esapi.Response
	var stats UpdateStats

	// Delete request matching this node
	if req, err = t.deleteRequest(node); err != nil {
		return nil, err
	}

	if res, err = req.Do(context.Background(), t.Client); err != nil {
		return nil, fmt.Errorf("Elasticsearch response error: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		// 404 is normal in the case of a first update
		if res.StatusCode != 404 {
			return nil, fmt.Errorf("error deleting entries for %s (status=%s)", node.ID, res.Status())
		}
	}

	// Bulk index
	if indexer, err = esutil.NewBulkIndexer(esutil.BulkIndexerConfig{Client: t.Client, Index: t.IndexName}); err != nil {
		return nil, fmt.Errorf("unable to create bulk indexer: %w", err)
	}

	type doc struct {
		NodeID   string  `json:"node_id"`
		ID       string  `json:"id"`
		Salience float32 `json:"salience"`
	}

	// TODO: Error handling could be improved here...

	for _, i := range topics {
		var data []byte

		if data, err = json.Marshal(&doc{NodeID: node.ID, ID: i.ID, Salience: i.Salience}); err != nil {
			return nil, fmt.Errorf("failed to marshal related topic document to JSON: %w", err)
		}

		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   strings.NewReader(string(data)),
				// Called for each successful operation
				// OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				// 	atomic.AddUint64(&countSuccessful, 1)
				// },
			},
		)
	}

	if err = indexer.Close(context.Background()); err != nil {
		return nil, fmt.Errorf("unexpected error encountered while closing the indexer %w", err)
	}

	stats = UpdateStats(indexer.Stats())

	return &stats, nil
}

// Convenience that returns a new DeleteByQueryRequest for the supplied Node
func (t ElasticTopicSearch) deleteRequest(node *common.Node) (esapi.Request, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"node_id": node.ID,
			},
		},
	}

	data, err := json.Marshal(query)

	if err != nil {
		return nil, fmt.Errorf("unable to marshal delete-by-query to JSON: %w", err)
	}

	return esapi.DeleteByQueryRequest{Index: []string{t.IndexName}, Body: strings.NewReader(string(data))}, nil
}

// Convenience that returns a new SearchRequest for the supplied Node
func (t ElasticTopicSearch) searchRequest(qid string) (esapi.Request, error) {
	query := map[string]interface{}{
		"sort": map[string]interface{}{
			"salience": map[string]interface{}{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": qid,
			},
		},
	}

	data, err := json.Marshal(query)

	if err != nil {
		return nil, fmt.Errorf("unable to marshal search query to JSON: %w", err)
	}

	return esapi.SearchRequest{Index: []string{t.IndexName}, Body: strings.NewReader(string(data))}, nil
}
