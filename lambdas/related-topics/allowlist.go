package main

import (
	"fmt"

	"github.com/wikimedia/phoenix/common"
)

// Hacky allowlist for processing of related topics.
func allowlist() []string {
	// Entries are strings in the format of "{authority}/{page ID}"
	return []string{
		"simple.wikipedia.org/2138",  // Albert Einstein
		"simple.wikipedia.org/39",    // Apple
		"simple.wikipedia.org/3715",  // Banana
		"simple.wikipedia.org/515",   // Mars
		"simple.wikipedia.org/17333", // Philadelphia
		"simple.wikipedia.org/31769", // San Antonio
	}
}

func allowed(node *common.Node) bool {
	for _, allowed := range allowlist() {
		if fmt.Sprintf("%s/%d", node.Source.Authority, node.Source.ID) == allowed {
			return true
		}
	}
	return false
}
