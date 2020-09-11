package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikimedia/phoenix/common"
)

func TestIndex(t *testing.T) {
	index := GetTestIndex()
	source := common.Source{Authority: "fake.wikipedia.org"}

	err := index.Apply(
		&Update{
			Page: common.Page{
				ID:     "/page/a0a0a0a0a0a0a",
				Name:   "San Marcos",
				Source: source,
			},
			Nodes: []common.Node{
				{
					ID:     "/node/a0a0a0a0a0a0a",
					Source: source,
					Name:   "History",
				},
			},
		})
	require.Nil(t, err)

	id, err := index.PageIDForName("fake.wikipedia.org", "San Marcos")
	require.Nil(t, err)
	assert.Equal(t, "/page/a0a0a0a0a0a0a", id)

	id, err = index.PageIDForName("fake.wikipedia.org", "Bogus")
	require.NotNil(t, err)
	pageNotFound, ok := err.(*ErrPageNotFound)
	require.True(t, ok, "Expected an error of type ErrPageNotFound")
	assert.Equal(t, "Bogus", pageNotFound.Name)

	id, err = index.NodeIDForName("fake.wikipedia.org", "San Marcos", "History")
	require.Nil(t, err)
	assert.Equal(t, "/node/a0a0a0a0a0a0a", id)

	id, err = index.NodeIDForName("fake.wikipedia.org", "San Marcos", "Bogus")
	require.NotNil(t, err)
	nodeNotFound, ok := err.(*ErrNodeNotFound)
	require.True(t, ok, "Expected an error of type ErrNodeNotFound")
	assert.Equal(t, "Bogus", nodeNotFound.Name)
}
