package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/rosette"
	"github.com/wikimedia/phoenix/storage"
)

var (
	// Rosette gives us 80,000 requests per day, or one every 1.08 seconds
	maxRate = 1080 * time.Millisecond

	// Command-line flags
	limitFlag    = flag.Int("limit", -1, "number of items to process")
	resumeFlag   = flag.String("resume", "", "node ID to resume from")
	debugLogFlag = flag.String("debug-log", "/dev/null", "enable debug logging to file")

	// These are assigned during compilation using `-ldflags` (see: Makefile)
	awsRegion                 string
	esEndpoint                string
	esIndex                   string
	esUsername                string
	esPassword                string
	rosetteAPIKey             string
	s3StructuredContentBucket string
)

// Struct corresponding to the JSON output of an AWS CLI DynamoDB table scan (scpoc-dynamodb-node-names)
type tableScan struct {
	Items []struct {
		ID struct {
			S string `json:"S"`
		} `json:"ID"`
		Authority struct {
			S string `json:"S"`
		} `json:"Authority"`
		Name struct {
			S string `json:"S"`
		} `json:"Name"`
	} `json:"Items"`
	Count        int32 `json:"Count"`
	ScannedCount int32 `json:"ScannedCount"`
}

// Read JSON-serialized table scan data from stdin
func readTableScan() *tableScan {
	var b []byte
	var err error
	var table = &tableScan{}

	// Read from standard-in
	if b, err = ioutil.ReadAll(os.Stdin); err != nil {
		panic(fmt.Errorf("Unable to read input: %w", err))
	}

	// Deserialize JSON input
	if err = json.Unmarshal(b, table); err != nil {
		panic(fmt.Errorf("Failure to deserialize JSON input: %w", err))
	}

	return table
}

// Prefix output with a timestamp
func println(format string, a ...interface{}) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
}

func init() {
	if awsRegion == "" {
		panic("awsRegion is UNSET; not passed in ldflags during compilation!")
	}
	if esEndpoint == "" {
		panic("esEndpoint is UNSET; not passed in ldflags during compilation!")
	}
	if esIndex == "" {
		panic("esIndex is UNSET; not passed in ldflags during compilation!")
	}
	if esUsername == "" {
		panic("esUsername is UNSET; not passed in ldflags during compilation!")
	}
	if esPassword == "" {
		panic("esPassword is UNSET; not passed in ldflags during compilation!")
	}
	if rosetteAPIKey == "" {
		panic("rosetteAPIKey is UNSET; not passed in ldflags during compilation!")
	}
	if s3StructuredContentBucket == "" {
		panic("s3StructuredContentBucket is UNSET; not passed in ldflags during compilation!")
	}
}

func main() {
	var completed int
	var content *storage.Repository
	var data *tableScan
	var err error
	var limiter <-chan time.Time
	var log *common.Logger
	var paused bool
	var topicsIndex storage.TopicSearch
	var topicsService rosette.Rosette

	// Setup ---------

	flag.Parse()

	// Content repository
	content = &storage.Repository{
		Store:  s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)})),
		Bucket: s3StructuredContentBucket,
	}

	// Logging
	if log, err = common.NewFileLogger("DEBUG", *debugLogFlag); err != nil {
		panic(fmt.Errorf("Unable to create debug logger: %w", err))
	}
	log.SetFlags(common.Ltimestamp)

	// Table scan data
	log.Debug("Opening table scan data...")
	data = readTableScan()

	if *resumeFlag != "" {
		paused = true
	} else {
		paused = false
	}

	// Topics search indexer
	var esClient *elasticsearch.Client

	var esConfig = elasticsearch.Config{
		Addresses: []string{esEndpoint},
		Username:  esUsername,
		Password:  esPassword,
	}

	if esClient, err = elasticsearch.NewClient(esConfig); err == nil {
		topicsIndex = &storage.ElasticTopicSearch{Client: esClient, IndexName: esIndex}
	} else {
		panic(fmt.Errorf("Unable to create Elasticsearch client: %w", err))
	}

	// Related topics services
	topicsService = rosette.Rosette{APIKey: rosetteAPIKey, Logger: log}

	// Rate limiting
	limiter = time.Tick(maxRate)

	// Begin import -------------

	// Iterate over items
	for _, item := range data.Items {
		var id = item.ID.S
		var node *common.Node
		var topics []common.RelatedTopic

		log.Debug("Processing ID=%s", id)

		// Unpause as appropriate
		if paused && *resumeFlag == id {
			paused = false
			log.Debug("Resuming from ID=%s", id)
		}

		// Skip past this item
		if paused {
			log.Debug("Paused, skipping...")
			continue
		}

		if node, err = content.GetNode(id); err != nil {
			panic(fmt.Errorf("Failed to retrieve node from content repository: %w", err))
		}

		log.Debug("Retreived ID=%s (%s)", id, node.Name)

		if strings.Trim(node.Unsafe, " \n") == "" {
			log.Debug("Skipping %s (%s); node.Unsafe is zero length", node.ID, item.Name.S)
			continue
		}

		if topics, err = topicsService.Topics(node); err != nil {
			panic(fmt.Errorf("Failed to retreive topics for %s: %w", node.ID, err))
		}

		log.Debug("Retrieved related topics from Rosette service (ID=%s)", id)

		if err = content.PutTopics(node, topics); err != nil {
			panic(fmt.Errorf("Failed to store related topics to content repository: %w", err))
		} else {
			log.Debug("Stored topics object to content repository (ID=%s)", id)

			if err = topicsIndex.Update(node, topics); err != nil {
				panic(fmt.Errorf("Failed to update topics search index: %w", err))
			}

			log.Debug("Successfully updated related topics index (ID=%s)", id)
		}

		completed++

		println("%s complete (%s)", node.ID, item.Name.S)

		if *limitFlag > 0 && completed >= *limitFlag {
			break
		}

		// Wait maxRate ticks
		<-limiter
	}

	println("Completed %d items", completed)
}
