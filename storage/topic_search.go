package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/wikimedia/phoenix/common"
)

const indexName = "topics"

// TopicSearch is an interface for topic search indexing.
type TopicSearch interface {
	// Search queries the index for nodes matching an ID
	Search(id string) ([]common.Node, error)

	// Update applies changes to the topic index
	Update(node *common.Node, topics []common.RelatedTopic) error
}

// ElasticTopicSearch is an Elasticsearch implementation of the TopicSearch interface.
type ElasticTopicSearch struct {
	Client *elasticsearch.Client
}

// Search queries the index for nodes matching an ID
func (t ElasticTopicSearch) Search(id string) ([]common.Node, error) {
	return nil, nil
}

// Update applies changes to the topic index
func (t ElasticTopicSearch) Update(node *common.Node, topics []common.RelatedTopic) error {
	var err error
	var indexer esutil.BulkIndexer
	var req esapi.Request
	var res *esapi.Response

	// Delete requests matching this node
	if req, err = deleteRequest(node); err != nil {
		return err
	}

	if res, err = req.Do(context.Background(), t.Client); err != nil {
		return fmt.Errorf("Elasticsearch response error: %w", err)
	}

	defer res.Body.Close()

	// TODO: Don't bail for status 404
	if res.IsError() {
		if res.StatusCode != 404 {
			return fmt.Errorf("error deleting entries for %s (status=%s)", node.ID, res.Status())
		}
	}

	// TODO: Issue bulk inserts
	if indexer, err = esutil.NewBulkIndexer(esutil.BulkIndexerConfig{Client: t.Client, Index: indexName}); err != nil {
		return fmt.Errorf("unable to create bulk indexer: %w", err)
	}

	type doc struct {
		NodeID   string  `json:"node_id"`
		ID       string  `json:"id"`
		Salience float32 `json:"salience"`
	}

	for _, i := range topics {
		var data []byte

		if data, err = json.Marshal(&doc{NodeID: node.ID, ID: i.ID, Salience: i.Salience}); err != nil {
			// TODO: Do
			continue
		}

		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   strings.NewReader(string(data)),
			},
		)
	}

	if err = indexer.Close(context.Background()); err != nil {
		// TODO: Do
	}

	// TODO: For debug; Remove this.
	fmt.Println(indexer.Stats().NumFailed)

	// TODO: Return stats in error(?)
	return nil
}

// Returns a new DeleteByQueryRequest for the supplied Node
func deleteRequest(node *common.Node) (esapi.Request, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"node_id": node.ID,
			},
		},
	}

	data, err := json.Marshal(query)

	if err != nil {
		return nil, fmt.Errorf("unable to marshal query to JSON: %w", err)
	}

	return esapi.DeleteByQueryRequest{Index: []string{indexName}, Body: strings.NewReader(string(data))}, nil
}
