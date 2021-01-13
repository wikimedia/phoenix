package storage

import (
	"io/ioutil"
	"os"
	"testing"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const cfgFile = ".config.yaml"

type TestConfig struct {
	ElasticsearchEndpoint string `yaml:"elasticsearch_endpoint"`
	ElasticsearchUsername string `yaml:"elasticsearch_username"`
	ElasticsearchPassword string `yaml:"elasticsearch_password"`
}

func TestTopicSearch(t *testing.T) {
	var cfg = TestConfig{}
	var data []byte
	var err error
	var topicSearch TopicSearch

	// If configuration does not exist, skip
	if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
		t.Skip("Elasticsearch tests not enabled")
	}

	// If configuration exists, but cannot be read, error out
	if data, err = ioutil.ReadFile(cfgFile); err != nil {
		t.Logf("Unable to read test configuration: %s (%+v)", cfgFile, err)
		t.FailNow()
	}

	// If configuration cannot be parsed, error out
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		t.Logf("Unable to parse %s as YAML (%+v)", cfgFile, err)
		t.FailNow()
	}

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchEndpoint},
		Username:  cfg.ElasticsearchUsername,
		Password:  cfg.ElasticsearchPassword,
	})

	require.Nil(t, err)

	topicSearch = ElasticTopicSearch{Client: esClient, IndexName: "topics_test"}

	t.Run("Update", func(t *testing.T) {
		err = topicSearch.Update(&testNode, testTopics)
		require.Nil(t, err)
	})

	// XXX: This can fail because the Update operation is race-y AF
	t.Run("Search", func(t *testing.T) {
		ids, err := topicSearch.Search("Q1")
		require.Nil(t, err)

		require.Len(t, ids, 1)
		assert.Equal(t, testNode.ID, ids[0])
	})
}
