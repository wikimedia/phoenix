package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikimedia/phoenix/common"
)

// Test objects
var testAbout common.Thing
var testNode common.Node
var testPage common.Page
var testSource common.Source

// Initialization of test data...
func setup() {
	// About
	testAbout = *common.NewThing()
	testAbout.AlternateName = "Alamo City"
	testAbout.Description = "second-most populous city in Texas, United States of America"
	testAbout.Image = "https://commons.wikimedia.org/wiki/File:The_Alamo_2019_v2.jpg"
	testAbout.Name = "San Antonio"
	testAbout.SameAs = "https://www.wikidata.org/entity/Q975"

	var now time.Time
	var err error

	// Work around for the loss of resolution during the time.Time -> RFC3339 string -> time.Time round-trip
	// that would otherwise make equality comparisons difficult.
	if now, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)); err != nil {
		panic(err)
	}

	// Source
	testSource = common.Source{
		ID:        1,
		Revision:  1,
		TimeUUID:  uuid.Must(uuid.NewUUID()).String(),
		Authority: "fake.wikipedia.org",
	}

	// Node
	testNode = common.Node{
		Source:       testSource,
		Name:         "History",
		IsPartOf:     []string{},
		DateModified: now,
		Unsafe:       "<h1>History</h1><p>At the time of European encounter...</p>",
	}

	// Page
	testPage = common.Page{
		Source:       testSource,
		Name:         "San Antonio",
		URL:          "//fake.wikipedia.org/wiki/San_Antonio",
		DateModified: now,
		HasPart:      []string{testNode.ID},
		About:        map[string]string{"//schema.org": testAbout.ID},
	}

	testNode.IsPartOf = append(testNode.IsPartOf, testPage.ID)
}

// MockStore is a mock implementation of S3 storage
type MockStore struct {
	Pages  map[string]common.Page
	Nodes  map[string]common.Node
	Abouts map[string]common.Thing
}

// GetObject is a mock of s3.S3#GetObject
func (store *MockStore) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	var present bool
	var b []byte
	var err error

	switch {
	case strings.HasPrefix(*input.Key, "/page"):
		var page common.Page

		// Not found
		if page, present = store.Pages[*input.Key]; !present {
			return nil, awserr.New(s3.ErrCodeNoSuchKey, "Not found", nil)
		}

		if b, err = json.Marshal(&page); err != nil {
			return nil, fmt.Errorf("unabled to marshal Page to JSON: %w", err)
		}

	case strings.HasPrefix(*input.Key, "/node"):
		var node common.Node

		// Not found
		if node, present = store.Nodes[*input.Key]; !present {
			return nil, awserr.New(s3.ErrCodeNoSuchKey, "Not found", nil)
		}

		if b, err = json.Marshal(&node); err != nil {
			return nil, fmt.Errorf("unabled to marshal Node to JSON: %w", err)
		}

	case strings.HasPrefix(*input.Key, "/data"):
		var about common.Thing

		// Not found
		if about, present = store.Abouts[*input.Key]; !present {
			return nil, awserr.New(s3.ErrCodeNoSuchKey, "Not found", nil)
		}

		if b, err = json.Marshal(&about); err != nil {
			return nil, fmt.Errorf("unabled to marshal Node to JSON: %w", err)
		}

	default:
		return nil, fmt.Errorf("unrecognized key format (%s)", *input.Key)

	}

	return &s3.GetObjectOutput{Body: aws.ReadSeekCloser(bufio.NewReader(bytes.NewBuffer(b)))}, nil
}

// PutObject is a mock of s3.S3#PutObject
func (store *MockStore) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	var b []byte
	var err error

	// Copy the contents of Body
	if b, err = ioutil.ReadAll(input.Body); err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(*input.Key, "/page"):
		page := common.Page{}

		if err = json.Unmarshal(b, &page); err != nil {
			return nil, fmt.Errorf("unable to deserialize Page: %w", err)
		}

		store.Pages[*input.Key] = page

	case strings.HasPrefix(*input.Key, "/node"):
		node := common.Node{}

		if err = json.Unmarshal(b, &node); err != nil {
			return nil, fmt.Errorf("unable to deserialize Node: %w", err)
		}

		store.Nodes[*input.Key] = node

	case strings.HasPrefix(*input.Key, "/data"):
		about := common.Thing{}

		if err = json.Unmarshal(b, &about); err != nil {
			return nil, fmt.Errorf("unable to deserialize Thing: %w", err)
		}

		store.Abouts[*input.Key] = about

	default:
		return nil, fmt.Errorf("unrecognized key format (%s)", *input.Key)

	}

	return &s3.PutObjectOutput{}, nil
}

func (store *MockStore) DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	log.Println("DeleteObjects() not (yet) implemented")
	return &s3.DeleteObjectsOutput{}, nil
}

func NewMockStore() *MockStore {
	return &MockStore{
		Pages:  make(map[string]common.Page),
		Nodes:  make(map[string]common.Node),
		Abouts: make(map[string]common.Thing),
	}
}

// Get an environment variable if set, or a default otherwise.
func Getenv(envar string, def string) string {
	if val := os.Getenv(envar); val != "" {
		return val
	}
	return def
}

// Depending on the environment, return a Store to use in tests (either for S3 or a mock)
func GetTestStore() Store {
	useS3, _ := strconv.ParseBool(Getenv("TESTS_USE_S3", "false"))

	if useS3 {
		region := Getenv("AWS_REGION", "us-east-2")
		return s3.New(session.New(&aws.Config{Region: aws.String(region)}))
	}

	return NewMockStore()
}

// Depending on the environment, return an Index to use in tests (either for Elasticsearch or a mock)
func GetTestIndex() Index {
	useES, _ := strconv.ParseBool(Getenv("TESTS_USE_ELASTICSEARCH", "false"))
	useDynamoDB, _ := strconv.ParseBool(Getenv("TESTS_USE_DYNAMODB", "false"))

	if useES {
		client, _ := elasticsearch.NewDefaultClient()
		return &ElasticsearchIndex{Client: client}
	}

	if useDynamoDB {
		region := Getenv("AWS_REGION", "us-east-2")
		titles := Getenv("AWS_DYNAMODB_PAGE_TITLES_TABLE", "scpoc-dynamodb-page-titles")
		names := Getenv("AWS_DYNAMODB_NODE_NAMES_TABLE", "scpoc-dynamodb-node-names")

		return &DynamoDBIndex{
			Client:      dynamodb.New(session.New(&aws.Config{Region: aws.String(region)})),
			TitlesTable: titles,
			NamesTable:  names,
		}
	}

	return NewMockIndex()
}

// Tests
func TestRepository(t *testing.T) {
	repo := Repository{Store: GetTestStore(), Index: GetTestIndex(), Bucket: Getenv("AWS_BUCKET", "scpoc-structured-content-store")}

	// Page
	t.Run("PutPage", func(t *testing.T) {
		id, err := repo.PutPage(&testPage)
		require.Nil(t, err)
		testPage.ID = id
	})
	t.Run("GetPage", func(t *testing.T) {
		page, err := repo.GetPage(testPage.ID)
		require.Nil(t, err)
		assert.Equal(t, &testPage, page)
	})
	t.Run("GetPage (not found)", func(t *testing.T) {
		_, err := repo.GetPage("/page/bogus")
		require.NotNil(t, err)
		var s3err awserr.Error
		require.True(t, errors.As(err, &s3err))
		assert.Equal(t, s3.ErrCodeNoSuchKey, s3err.Code())
	})

	// Node
	t.Run("PutNode", func(t *testing.T) {
		id, err := repo.PutNode(&testNode)
		require.Nil(t, err)
		testNode.ID = id
	})
	t.Run("GetNode", func(t *testing.T) {
		node, err := repo.GetNode(testNode.ID)
		require.Nil(t, err)
		assert.Equal(t, &testNode, node)
	})

	// About
	t.Run("PutAbout", func(t *testing.T) {
		id, err := repo.PutAbout(&testAbout)
		require.Nil(t, err)
		testAbout.ID = id
	})
	t.Run("GetAbout", func(t *testing.T) {
		about, err := repo.GetAbout(testAbout.ID)
		require.Nil(t, err)
		assert.Equal(t, &testAbout, about)
	})

	// Function(s)
	t.Run("makePageID", func(t *testing.T) {
		id := makePageID(&testPage)
		assert.Equal(t, id, makePageID(&testPage))
	})
}

func TestRepositoryApply(t *testing.T) {
	repo := Repository{Store: GetTestStore(), Index: GetTestIndex(), Bucket: Getenv("AWS_BUCKET", "scpoc-structured-content-store")}

	t.Run("Apply", func(t *testing.T) {
		update := &Update{
			Page:   testPage,
			Nodes:  []common.Node{testNode},
			Abouts: map[string]common.Thing{"//schema.org": testAbout},
		}

		require.Nil(t, repo.Apply(update))

		// Retrieving a page by its ID
		page, err := repo.GetPage(testPage.ID)
		require.Nil(t, err)
		assert.Equal(t, testPage.DateModified, page.DateModified)
		assert.Equal(t, testPage.Name, page.Name)
		assert.Equal(t, testPage.Source, page.Source)
		assert.Equal(t, testPage.URL, page.URL)
		assert.Len(t, page.HasPart, 1)
		assert.Len(t, page.About, 1)

		// Retrieving a page by its name (exercises the index)
		page, err = repo.GetPageByName(testPage.Source.Authority, testPage.Name)
		require.Nil(t, err)
		assert.Equal(t, testPage.ID, page.ID, "By-name lookup of the Page failed (indexing)")

		// Retrieving a node by its ID
		node, err := repo.GetNode(page.HasPart[0])
		require.Nil(t, err)
		assert.Equal(t, testNode.Name, node.Name)
		assert.Equal(t, testNode.Unsafe, node.Unsafe)
		assert.Equal(t, testNode.DateModified, node.DateModified)
		assert.Equal(t, testPage.ID, node.IsPartOf[0])

		// Retrieving a node by its name (exercises the index)
		node, err = repo.GetNodeByName(testNode.Source.Authority, page.Name, node.Name)
		require.Nil(t, err)
		assert.Equal(t, testNode.ID, node.ID)

		// Retrieving a linked-data object
		about, err := repo.GetAbout(page.About["//schema.org"])
		require.Nil(t, err)
		assert.Equal(t, testAbout.AlternateName, about.AlternateName)
		assert.Equal(t, testAbout.Context, about.Context)
		assert.Equal(t, testAbout.Description, about.Description)
		assert.Equal(t, testAbout.Image, about.Image)
		assert.Equal(t, testAbout.Name, about.Name)
		assert.Equal(t, testAbout.SameAs, about.SameAs)
		assert.Equal(t, testAbout.Type, about.Type)
	})
}

func TestValidation(t *testing.T) {
	t.Run("Source", func(t *testing.T) {
		require.Nil(t, validateSource(&common.Source{ID: 1, Revision: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067", Authority: "s.wp.o"}))
		require.NotNil(t, validateSource(&common.Source{Revision: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067", Authority: "s.wp.o"}))
		require.NotNil(t, validateSource(&common.Source{ID: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067", Authority: "s.wp.o"}))
		require.NotNil(t, validateSource(&common.Source{ID: 1, Revision: 1, Authority: "s.wp.o"}))
		require.NotNil(t, validateSource(&common.Source{ID: 1, Revision: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067"}))
	})

	t.Run("Page", func(t *testing.T) {
		s := common.Source{ID: 1, Revision: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067", Authority: "s.wp.o"}
		require.Nil(t, validatePage(&common.Page{Name: "Foo", URL: "//s.wp.o", DateModified: time.Now(), HasPart: []string{"/n/a"}, Source: s}))
		require.NotNil(t, validatePage(&common.Page{URL: "//s.wp.o", DateModified: time.Now(), HasPart: []string{"/n/a"}, Source: s}))
		require.NotNil(t, validatePage(&common.Page{Name: "Foo", DateModified: time.Now(), HasPart: []string{"/n/a"}, Source: s}))
		require.NotNil(t, validatePage(&common.Page{Name: "Foo", URL: "//s.wp.o", HasPart: []string{"/n/a"}, Source: s}))
		require.NotNil(t, validatePage(&common.Page{Name: "Foo", URL: "//s.wp.o", DateModified: time.Now(), Source: s}))
		require.NotNil(t, validatePage(&common.Page{Name: "Foo", URL: "//s.wp.o", DateModified: time.Now(), HasPart: []string{"/n/a"}}))
	})

	t.Run("Node", func(t *testing.T) {
		s := common.Source{ID: 1, Revision: 1, TimeUUID: "61e16274-ed75-11ea-a791-9fba67228067", Authority: "s.wp.o"}
		require.Nil(t, validateNode(&common.Node{DateModified: time.Now(), Unsafe: "<p>...</p>", Source: s}))
		require.NotNil(t, validateNode(&common.Node{Unsafe: "<p>...</p>", Source: s}))
		require.NotNil(t, validateNode(&common.Node{DateModified: time.Now(), Source: s}))
		require.NotNil(t, validateNode(&common.Node{DateModified: time.Now(), Unsafe: "<p>...</p>"}))
	})
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}
