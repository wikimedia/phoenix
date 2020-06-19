package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eevans/wikimedia/streams"
	"github.com/wikimedia/phoenix/common"
)

// Event streams attributes
const namespace int = 0
const wiki string = "simplewiki"

var (
	awsAccount string
	awsRegion  string
	snsTopic   string

	since = flag.String("since", "", "Offset to begin stream from (ISO8601 timestamp or milliseconds since epoch)")
)

func main() {
	flag.Parse()

	events := streams.NewClient().Match("namespace", namespace).Match("wiki", wiki).Match("type", "edit")

	if *since != "" {
		// Validate that *since parses as RFC3339 (ISO8601)
		if _, err := time.Parse(time.RFC3339, *since); err != nil {
			// If not RFC3339/ISO8601, is it at least a number?
			if _, err = strconv.Atoi(*since); err != nil {
				fmt.Fprint(os.Stderr, "Invalid timestamp argument for `-since`!")
				os.Exit(1)
			}
		}

		events.Since = *since
	}

	client := common.NewChangeEventPublisher(awsAccount, awsRegion, snsTopic)

	err := events.RecentChanges(func(event streams.RecentChangeEvent) {
		fmt.Printf("Change event captured!\n")
		fmt.Printf("  Title ............: %s\n", event.Title)
		fmt.Printf("  Server name ......: %s\n", event.ServerName)
		fmt.Printf("  Wiki .............: %s\n", event.Wiki)
		fmt.Printf("  Namespace ........: %d\n", event.Namespace)
		fmt.Printf("  Type .............: %s\n", event.Type)
		fmt.Printf("  Revision .........: %d\n", event.Revision.New)
		fmt.Printf("  Timestamp ........: %s\n", event.Meta.Dt)

		result, err := client.Send(
			&common.ChangeEvent{
				ServerName: event.ServerName,
				Title:      event.Title,
				Revision:   event.Revision.New,
			})
		if err != nil {
			fmt.Printf("Error enqueuing %s (%s)\n", event.Title, err)
			return
		}

		fmt.Printf("Queued \"%s\" as %s\n", event.Title, *result.MessageId)
	})

	fmt.Println()

	fmt.Printf("RecentChanges: %s\n", err)
	fmt.Printf("Last timestamp: %s\n", events.LastTimestamp())
}
