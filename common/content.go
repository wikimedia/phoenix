package common

import (
	"time"
)

// Page represents the root node of a document graph
type Page struct {
	// Globally unique identifier
	ID string `json:"identifier"`

	// Source data for this page
	Source Source `json:"_source"`

	// The Page name (corresponds with schema.org/Thing#name)
	Name string `json:"name"`

	// The Page url (corresponds with schema.org/Thing#url)
	URL string `json:"url"`

	// Date and time of last modification (corresponds with schema.org/CreativeWork#dateModified)
	DateModified time.Time `json:"dateModified"`

	// URLs of content which are a part of this one.  Loosely corresponds with schema.org/CreativeWork#hasPart,
	// but unlike its namesake, this attribute serves as an adjacency list of nodes in the document graph.
	HasPart []string `json:"hasPart"`

	// URLs of metadata associated with the topic of this page.  Loosely correponds with
	// schema.org/CreativeWork#about, though unlike its namesake, this attribute is an associative array of
	// metadata in an arbitrary set of vocabularies (keyed by the vocabulary).
	About map[string]string `json:"about"`
}

// Source represents information on the source of the document.
type Source struct {
	// Page ID according to MediaWiki
	ID int `json:"id"`

	// Revision ID according to MediaWiki
	Revision int `json:"revision"`

	// Type 1 UUID; Date and time the source document was rendered
	TimeUUID string `json:"tid"`

	// The wiki/project/hostname of source document
	Authority string `json:"authority"`
}

// Node represents a node in the document graph
type Node struct {
	// Globally unique identifier
	ID string `json:"id"`

	// Node name (corresponds with schema.org/Thing#name).  For a section, corresponds to the text of
	// the first header (Parsoid HTML output).
	Name string `json:"name,omitempty"`

	// URLs of content that this node is a part of.  Loosely corresponds with
	// schema.org/CreativeWork#isPartOf, yet unlike its namesake, this attribute serves as an adjacency
	// list of nodes in the document graph.
	IsPartOf []string `json:"isPartOf"`

	// Date and time of last modification (corresponds with schema.org/CreativeWork#dateModified)
	DateModified time.Time `json:"dateModified"`

	// The raw HTML context of the corresponding node.
	Unsafe string `json:"unsafe"`
}

type metadata struct {
	ID      string `json:"-"`
	Context string `json:"@context"`
	Type    string `json:"@type"`
}

// Thing corresponds to https://schema.org/Thing
type Thing struct {
	metadata
	AlternateName string `json:"alternateName,omitempty"`
	Description   string `json:"description,omitempty"`
	Image         string `json:"image,omitempty"`
	Name          string `json:"name,omitempty"`
	SameAs        string `json:"sameAs"`
}

// NewThing returns an initialized Thing
func NewThing() *Thing {
	return &Thing{metadata: metadata{Context: "https://schema.org", Type: "Thing"}}
}
