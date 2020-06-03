package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/wikimedia/phoenix/event-bridge/common"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <server> [<title> <revision>]\n\n", os.Args[0])
}

func main() {
	client, err := common.NewPublisher()
	if err != nil {
		fmt.Printf("Unable to create publisher: %s\n", err)
		os.Exit(1)
	}

	var serverName, title, revision string

	switch len(os.Args) {
	case 2:
		serverName = os.Args[1]

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			s := strings.SplitN(line, " ", 2)
			title = strings.TrimSpace(s[1])
			revision = strings.TrimSpace(s[0])

			revNum, err := strconv.Atoi(revision)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s not a valid revision, skipping...\n", revision)
				continue
			}

			result, err := client.Send(serverName, title, revNum)
			if err != nil {
				fmt.Printf("Error enqueuing %s (%s)\n", title, err)
				continue
			}

			fmt.Printf("Queued \"%s\" as %s\n", title, *result.MessageId)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

	case 4:
		serverName = os.Args[1]
		title = os.Args[2]
		revision = os.Args[3]

		revNum, err := strconv.Atoi(revision)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s not a valid revision\n", revision)
			printUsage()
			os.Exit(1)
		}

		result, err := client.Send(serverName, title, revNum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enqueuing %s (%s)\n", title, err)
			os.Exit(1)
		}

		fmt.Printf("Queued \"%s\" as %s\n", title, *result.MessageId)

	default:
		printUsage()
	}

}
