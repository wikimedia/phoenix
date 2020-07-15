package common

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewThing(t *testing.T) {
	thing := NewThing()
	require.Equal(t, "https://schema.org", thing.Context)
	require.Equal(t, "Thing", thing.Type)
}

func TestJSONTime(t *testing.T) {
	now := time.Now()

	s := struct {
		JSONTime
	}{
		JSONTime(now),
	}

	b, err := json.Marshal(s)

	require.Nil(t, err)
	require.Equal(t, now.Format(time.RFC3339), strings.ReplaceAll(string(b), `"`, ""))

}
