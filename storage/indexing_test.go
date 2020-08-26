package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikimedia/phoenix/common"
)

func TestIndex(t *testing.T) {
	index := GetTestIndex()

	err := index.Apply(
		&common.Page{
			ID:     "/page/a0a0a0a0a0a0a",
			Name:   "San Marcos",
			Source: common.Source{Authority: "simple.wikipedia.org"},
		})
	require.Nil(t, err)

	id, err := index.PageIDForName("simple.wikipedia.org", "San Marcos")
	require.Nil(t, err)
	assert.Equal(t, "/page/a0a0a0a0a0a0a", id)

	id, err = index.PageIDForName("simple.wikipedia.org", "Bogus")
	require.NotNil(t, err)
	notFound, ok := err.(*ErrNameNotFound)
	require.True(t, ok, "Expected an error of type ErrNameNotFound")
	assert.Equal(t, "Bogus", notFound.Name)

}
