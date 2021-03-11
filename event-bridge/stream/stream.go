package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/eevans/wikimedia/streams"
	"github.com/wikimedia/phoenix/common"
	"gopkg.in/yaml.v2"
)

// Event streams attributes
const namespace int = 0
const wiki string = "simplewiki"

var (
	awsAccount string
	awsRegion  string
	snsTopic   string

	allowFlag = flag.String("allow-list", "", "Yaml formatted allowlist")
	sinceFlag = flag.String("since", "", "Offset to begin stream from (ISO8601 timestamp or milliseconds since epoch)")
)

type access interface {
	allowed(evt streams.RecentChangeEvent) bool
}

// An access implementation that allows everything
type allowAny struct{}

func (a *allowAny) allowed(evt streams.RecentChangeEvent) bool {
	return true
}

// An access implementation read from a YAML-formatted file
type yamlAllow struct {
	Allowed []string `yaml:"allowed"`
}

func (a *yamlAllow) allowed(evt streams.RecentChangeEvent) bool {
	for _, title := range a.Allowed {
		if title == evt.Title {
			return true
		}
	}
	return false
}

func newYamlAllow(fname string) (*yamlAllow, error) {
	var data []byte
	var err error
	var result = &yamlAllow{}

	if data, err = ioutil.ReadFile(fname); err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", fname, err)
	}

	if err = yaml.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", fname, err)
	}

	return result, nil
}

func main() {
	var err error
	var events *streams.Client
	var list access = &allowAny{}

	flag.Parse()

	events = streams.NewClient().Match("namespace", namespace).Match("wiki", wiki).Match("type", "edit")

	if *sinceFlag != "" {
		// Validate that *since parses as RFC3339 (ISO8601)
		if _, err := time.Parse(time.RFC3339, *sinceFlag); err != nil {
			// If not RFC3339/ISO8601, is it at least a number?
			if _, err = strconv.Atoi(*sinceFlag); err != nil {
				fmt.Fprint(os.Stderr, "Invalid timestamp argument for `-since`!")
				os.Exit(1)
			}
		}

		events.Since = *sinceFlag
	}

	if *allowFlag != "" {
		if list, err = newYamlAllow(*allowFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading allow list: %s", err)
			os.Exit(1)
		}
	}

	client := common.NewChangeEventPublisher(awsAccount, awsRegion, snsTopic)

	// Loop in perpetuity (or until err is not nil)
	for err == nil {
		err = events.RecentChanges(func(event streams.RecentChangeEvent) {
			var err error
			var result *sns.PublishOutput

			fmt.Printf("Change event captured!\n")
			fmt.Printf("  Title ............: %s\n", event.Title)
			fmt.Printf("  Server name ......: %s\n", event.ServerName)
			fmt.Printf("  Wiki .............: %s\n", event.Wiki)
			fmt.Printf("  Namespace ........: %d\n", event.Namespace)
			fmt.Printf("  Type .............: %s\n", event.Type)
			fmt.Printf("  Revision .........: %d\n", event.Revision.New)
			fmt.Printf("  Timestamp ........: %s\n", event.Meta.Dt)

			// Only forward allowlisted events
			if !list.allowed(event) {
				fmt.Printf("  Status ...........: skipped\n")
				return
			}

			result, err = client.Send(
				&common.ChangeEvent{
					ServerName: event.ServerName,
					Title:      event.Title,
					Revision:   event.Revision.New,
				})
			if err != nil {
				fmt.Printf("Error enqueuing %s (%s)\n", event.Title, err)
				return
			}

			fmt.Printf("  Status ...........: queued as %s\n", *result.MessageId)
		})
	}

	fmt.Println()

	fmt.Printf("RecentChanges: %s\n", err)
	fmt.Printf("Last timestamp: %s\n", events.LastTimestamp())
}
