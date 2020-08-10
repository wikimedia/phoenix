package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
	"github.com/wikimedia/phoenix/common"
)

// Store is a mockable interface corresponding to s3.S3.
type Store interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)
}

// Repository provides read/write access to the Phoenix Content Repository.
type Repository struct {
	Store  Store
	Index  Index
	Bucket string
}

// Helper method for downloading files from S3.
func (r *Repository) get(key string) (*json.Decoder, error) {
	var input *s3.GetObjectInput
	var output *s3.GetObjectOutput
	var err error

	input = &s3.GetObjectInput{Bucket: aws.String(r.Bucket), Key: aws.String(key)}
	if output, err = r.Store.GetObject(input); err != nil {
		return nil, err
	}

	return json.NewDecoder(output.Body), nil
}

// Helper method for uploading files to S3.
func (r *Repository) put(key string, data []byte, meta map[string]*string) error {
	_, err := r.Store.PutObject(
		&s3.PutObjectInput{
			Body:        aws.ReadSeekCloser(bytes.NewReader(data)),
			Bucket:      aws.String(r.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String("application/json"),
			Metadata:    meta,
		})

	return err
}

// Helper method for deleting files from S3.
func (r *Repository) delete(keys []string) error {
	var objects []*s3.ObjectIdentifier

	for _, key := range keys {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := r.Store.DeleteObjects(
		&s3.DeleteObjectsInput{
			Bucket: aws.String(r.Bucket),
			Delete: &s3.Delete{
				Objects: objects,
				Quiet:   aws.Bool(false),
			},
		})

	return err
}

// GetPage returns a Page by its ID
func (r *Repository) GetPage(id string) (*common.Page, error) {
	var data *json.Decoder
	var err error
	var page common.Page

	// Fetch
	if data, err = r.get(id); err != nil {
		return nil, fmt.Errorf("Error retrieving content: %w", err)
	}

	// Deserialize JSON
	if err = data.Decode(&page); err != nil {
		return nil, fmt.Errorf("Unable to deserialize JSON: %w", err)
	}

	return &page, nil
}

// GetNode returns a Node by its ID
func (r *Repository) GetNode(id string) (*common.Node, error) {
	var data *json.Decoder
	var err error
	var section common.Node

	// Fetch
	if data, err = r.get(id); err != nil {
		return nil, fmt.Errorf("Error retrieving content: %w", err)
	}

	// Deserialize JSON
	if err = data.Decode(&section); err != nil {
		return nil, fmt.Errorf("Unable to deserialize JSON: %w", err)
	}

	return &section, nil
}

// GetAbout returns an About by its ID
func (r *Repository) GetAbout(id string) (*common.Thing, error) {
	var data *json.Decoder
	var err error
	var about common.Thing

	// Fetch
	if data, err = r.get(id); err != nil {
		return nil, fmt.Errorf("Error retrieving content: %w", err)
	}

	// Deserialize JSON
	if err = data.Decode(&about); err != nil {
		return nil, fmt.Errorf("Unable to deserialize JSON: %w", err)
	}

	// Thing doesn't JSON serialize the ID
	about.ID = id

	return &about, nil
}

// PutPage stores a Page. This method generates a unique ID and returns it on success; NOTE: If
// you assign an ID it will be overwritten.
func (r *Repository) PutPage(page *common.Page) (string, error) {
	var data []byte
	var err error

	page.ID = pagef(makePageID(page))

	if data, err = encodeJSON(page); err != nil {
		return "", err
	}

	metadata := map[string]*string{"type": aws.String("common.Page")}

	if err = r.put(page.ID, data, metadata); err != nil {
		return "", err
	}

	return page.ID, nil
}

// PutNode stores a Node.  This method generates a unique ID and returns it on success; NOTE: If
// you assign an ID it will be overwritten.
func (r *Repository) PutNode(node *common.Node) (string, error) {
	var data []byte
	var err error

	node.ID = nodef(makeRandomID())

	if data, err = encodeJSON(node); err != nil {
		return "", err
	}

	metadata := map[string]*string{"type": aws.String("common.Node")}

	if err = r.put(node.ID, data, metadata); err != nil {
		return "", err
	}

	return node.ID, nil
}

// PutAbout stores a Thing.  This method generates a unique ID and returns it on success; NOTE: If
// you assign an ID it will be overwritten.
func (r *Repository) PutAbout(thing *common.Thing) (string, error) {
	var data []byte
	var err error

	thing.ID = aboutf(makeRandomID())

	if data, err = encodeJSON(thing); err != nil {
		return "", err
	}

	metadata := map[string]*string{"type": aws.String("common.Thing")}

	if err = r.put(thing.ID, data, metadata); err != nil {
		return "", err
	}

	return thing.ID, nil
}

// DeletePage removes a Page from storage by its ID
func (r *Repository) DeletePage(id string) {
	// TODO: Do.
}

// DeleteNode removes a Node from storage by its ID
func (r *Repository) DeleteNode(id string) {
	// TODO: Do.
}

// DeleteAbout removes a Thing from storage by its ID
func (r *Repository) DeleteAbout(id string) {
	// TODO: Do.
}

// Update encapsulates the parts of a document involved in an update of the content repository.
type Update struct {
	Page   common.Page
	Nodes  []common.Node
	Abouts map[string]common.Thing
}

// Apply updates a document in the content repository.
func (r *Repository) Apply(update *Update) error {
	var prevPage *common.Page
	var err error

	// Baby steps: An argument could be made for breaking down the steps here into events that
	// trigger the respective actions, but we're not going there just yet.  An argument could
	// also be made for handling some of these steps concurrently (we could easily parallelize
	// uploads of Node & Things, for example), but we're not going there yet either.

	// A zero-length ID might mean that this a first-time write, but if not, we'd miss an
	// opportunity to clean up referenced values.
	if update.Page.ID == "" {
		p := update.Page
		update.Page.ID = pagef(makePageID(&p))
	}

	if prevPage, err = r.GetPage(update.Page.ID); err != nil {
		// Continue for ErrCodeNoSuchKey (first write?), return any other error.
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				break
			default:
				return aerr
			}
		} else {
			return aerr
		}
	}

	update.Page.HasPart = make([]string, 0)

	// Upload node objects.  Remember: the ordering of HasPart matters (keep this in mind
	// when/if adding concurrency at a later date).
	for _, node := range update.Nodes {
		var id string
		var err error
		node.IsPartOf = []string{update.Page.ID}
		if id, err = r.PutNode(&node); err != nil {
			return fmt.Errorf("error storing node: %w", err)
		}
		update.Page.HasPart = append(update.Page.HasPart, id)
	}

	update.Page.About = make(map[string]string)

	// Upload linked data objects.
	for k, v := range update.Abouts {
		var id string
		var err error
		if id, err = r.PutAbout(&v); err != nil {
			return fmt.Errorf("error storing linked data object: %w", err)
		}
		update.Page.About[k] = id
	}

	// TODO: Overwrite the Page object
	if _, err = r.PutPage(&update.Page); err != nil {
		return err
	}

	// Delete previous Node and linked-data objects
	if prevPage != nil {
		for _, id := range prevPage.HasPart {
			// TODO: Do.
			r.DeleteNode(id)
		}
		for _, id := range prevPage.About {
			// TODO: Do.
			r.DeleteAbout(id)
		}
	}

	// Perform indexing
	return r.Index.Apply(&update.Page)
}

// Helpers are helpful.
func encodeJSON(v interface{}) ([]byte, error) {
	var buffer *bytes.Buffer
	var encoder *json.Encoder
	var err error

	buffer = bytes.NewBuffer(make([]byte, 0, 0))
	encoder = json.NewEncoder(buffer)

	// Don't escape HTML (we store raw HTML in attributes)
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), err
}

func makeRandomID() string {
	return uuid.New().String()
}

// To maintain page ID stability, we fake a globally unique generated ID using a hash of the
// underlying wiki and page ID.
func makePageID(page *common.Page) string {
	hasher := murmur3.New64()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", page.Source.Authority, page.Source.ID)))
	return fmt.Sprintf("%x", hasher.Sum64())
}

// Return formatted keys for page, node, and data objects.
func pagef(id string) string {
	return fmt.Sprintf("/page/%s", id)
}

func nodef(id string) string {
	return fmt.Sprintf("/node/%s", id)
}

func aboutf(id string) string {
	return fmt.Sprintf("/data/%s", id)
}
