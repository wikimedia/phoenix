#!/bin/bash

set -e


query()
{
    cat <<EOF
{
  "query":
    "{
      node(name: { authority: \"$1\", pageName: \"$2\", name: \"$3\" } ) {
        dateModified
        name
        unsafe
        id
      }
    }"
}
EOF
}

echo "Query (JSON-encoded) -------------"
query "simple.wikipedia.org" "$1" "$2"

echo

echo "Response -------------------------"
query "simple.wikipedia.org" "$1" "$2" | curl -XPOST localhost:8080/query -d @- 2>/dev/null | json_pp
