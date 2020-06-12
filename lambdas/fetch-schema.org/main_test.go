package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikibaseItemID(t *testing.T) {
	val, err := wikibaseItemID("simple.wikipedia.org", "Banana")
	require.Nil(t, err)
	assert.Equal(t, "Q503", val)
}

func TestSchemaOrgItem(t *testing.T) {
	val, err := schemaOrgItem("Q60")
	require.Nil(t, err)
	assert.Equal(t, "New York City", val.Name)
}
