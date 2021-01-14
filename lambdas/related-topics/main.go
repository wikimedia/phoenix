package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	content storage.Repository
	log     *common.Logger

	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsAccount                string
	awsRegion                 string
	s3StructuredContentBucket string
)

type rosetteTopicsResponse struct {
	Concepts []struct {
		ConceptID string `json:"conceptId"`
		Phrase    string `json:"phrase"`
		Salience  int    `json:"salience"`
	} `json:"concepts"`
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	var minimumSalience = 0.1
	for _, record := range event.Records {
		var err error
		var msg = &common.NodeStoredEvent{}
		var node *common.Node
		var topics []common.RelatedTopic

		// Deserialize message
		if err = json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize message payload: %s", err)
			continue
		}

		log.Debug("Processing Node published event: %+v", msg)

		// Retrieve Node object from storage
		if node, err = content.GetNode(msg.ID); err != nil {
			log.Error("Failed to retreive S3 object for node (ID=%s): %s", msg.ID, err)
			continue
		}

		log.Debug("Processing Node.Unsafe='%.24s...'", node.Unsafe)

		// ?? Can we strip HTML? not a must, but recommended to send text-only
		//    since we have a char-limit on requests
		//    According to stackoverflow, "github.com/grokify/html-strip-tags-go"
		//    may be a suitable package if we trust the source
		//    (which we supposedly do since it's from Wikipedia parser)

		// Fetch topics from rosette.com
		topicsData, err := getRosetteTopics(node.Unsafe)
		if err != nil {
			log.Error("Unable to retrieve related topics for %s: %s", msg.ID, err)
			continue
		}

		// Unmarshall
		rosetteTopics := rosetteTopicsResponse{}
		if err := json.Unmarshal(topicsData, &rosetteTopics); err != nil {
			log.Error("Unable to deserialize topics data:", err)
			continue
		}

		// Sort topics by salience
		sort.SliceStable(rosetteTopics.Concepts, func(i, j int) bool {
			return rosetteTopics.Concepts[i].Salience > rosetteTopics.Concepts[j].Salience
		})

		// Collect and format to fit expected storage key/value
		topics := []common.RelatedTopic{}
		for _, conceptData := range rosetteTopics.Concepts {
			if (
				// Only add to the topics if the phrase is a wikidata item (starts with "Q")
				strings.HasPrefix(conceptData.Phrase,"Q") &&
				// Filter topics that are under the minimum salience
				conceptData.salience >= minimumSalience
			) {
				topics = append(topics, common.RelatedTopic{ID: conceptData.ID, Label: conceptData.Phrase, Salience: conceptData.Salience})
			}
		}

		// Store related topics
		if err = content.PutTopics(node, topics); err != nil {
			log.Error("Failed to store related-topics: %s", err)
		}
	}
}

// Fetch from Rosette's API
func getRosetteTopics(nodeUnsafe *string) ([]byte, error) {
	// TODO: Send headers with the rosette API key in a config file in .gitignore (?)
	// NOTE: From the Rosette API, this is how the request is sent via cURL
	// curl -s -X POST \
	// -H "X-RosetteAPI-Key: your_api_key" \
	// -H "Content-Type: application/json" \
	// -H "Accept: application/json" \
	// -H "Cache-Control: no-cache" \
	// -d '{"content": "Lily Collins is in talks to join Nicholas Hoult in Chernin Entertainment and Fox Searchlight\u0027s J.R.R. Tolkien biopic Tolkien. Anthony Boyle, known for playing Scorpius Malfoy in the British play Harry Potter and the Cursed Child, also has signed on for the film centered on the famed author. In Tolkien, Hoult will play the author of the Hobbit and Lord of the Rings book series that were later adapted into two Hollywood trilogies from Peter Jackson. Dome Karukoski is directing the project."}' \
	// "https://api.rosette.com/rest/v1/topics"

	client := &http.Client{}

	// Not sure how to take the content here while stripping quotes/making it safe for post params?
	var jsonStr = []byte(fmt.Sprintf(`{"content":"%s"}`, nodeUnsafe)
	req, err := http.NewRequest("POST", "https://api.rosette.com/rest/v1/topics?redirect=true", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %w", err)
	}

	req.Header.Set("X-RosetteAPI-Key", "xxxxx")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d (expected %d)", res.StatusCode, http.StatusOK)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading HTTP response body: %w", err)
	}

	return body, nil
}

func init() {
	// AWS S3 client obj
	content = storage.Repository{
		Store:  s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)})),
		Bucket: s3StructuredContentBucket,
	}

	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	log = common.NewLogger(level)
	log.Warn("%s LOGGING ENABLED (use LOG_LEVEL env var to configure)", common.LevelString(log.Level))

	log.Debug("AWS account ......................: %s", awsAccount)
	log.Debug("AWS region .......................: %s", awsRegion)
	log.Debug("S3 structured content bucket .....: %s", s3StructuredContentBucket)
}

func main() {
	lambda.Start(handleRequest)
}
