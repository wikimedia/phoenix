package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewThing(t *testing.T) {
	thing := NewThing()
	require.Equal(t, "https://schema.org", thing.Context)
	require.Equal(t, "Thing", thing.Type)
}
