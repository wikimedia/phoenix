package storage

import "github.com/wikimedia/phoenix/common"

// TopicSearch is an interface for topic search indexing.
type TopicSearch interface {
	// Search queries the index for nodes matching an ID
	Search(id string) ([]common.Node, error)

	// Update applies changes to the topic index
	Update(*common.Node, []common.RelatedTopic) error
}

// ElasticTopicSearch is an Elasticsearch implementation of the TopicSearch interface.
type ElasticTopicSearch struct {
}

// Search queries the index for nodes matching an ID
func (t *ElasticTopicSearch) Search(id string) ([]common.Node, error) {
	return nil, nil
}

// Update applies changes to the topic index
func (t *ElasticTopicSearch) Update(*common.Node, []common.RelatedTopic) error {
	return nil
}
