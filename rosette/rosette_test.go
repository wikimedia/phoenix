package rosette

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wikimedia/phoenix/common"
)

var (
	testData = `
		<section>
			<h2>Banana plant</h2>
			<p>The banana plant is the largest herbaceous flowering plant. Banana plants are often mistaken
			for trees. Bananas have a false stem (called pseudostem), which is made by the lower part of the
			leaves. This pseudostem can grow to be two to eight meters tall. Each pseudostem grows from a
			corm. A pseudostem is able to produce a single bunch of bananas. After fruiting, the pseudostem
			dies and is replaced. When most bananas are ripe, they turn yellow or, sometimes, red. Unripe
			bananas are green.</p>

			<p>Banana leaves grow in a spiral and may grow 2.7 metres (8.9 feet) long and 60 cm (2.0 ft)
			wide. They are easily torn by the wind, which results in a familiar, frayed look.</p>
		</section>`
	testNode = &common.Node{ID: "/node/abcdef0123456789", Unsafe: testData}
)

func TestExtractText(t *testing.T) {
	_, err := extractText(testData)
	require.Nil(t, err)
}

func TestRosetteTopics(t *testing.T) {
	var key string
	var log = common.NewLogger("DEBUG")
	var r Rosette

	if key = os.Getenv("ROSETTE_APIKEY"); key == "" {
		t.Skip("You must export env var ROSETTE_APIKEY")
	}

	r = Rosette{APIKey: key, Logger: log}

	topics, err := r.Topics(testNode)

	// Not a great test...
	require.Nil(t, err)
	require.NotNil(t, topics)
	require.NotEmpty(t, topics)
}
