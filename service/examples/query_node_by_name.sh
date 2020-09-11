#!/bin/bash

set -e


query()
{
    cat <<EOF
{
  "query": "query Node(\$authority: String!, \$pageName: String!, \$name: String!) { nodeByName(authority: \$authority, pageName: \$pageName, name: \$name) { dateModified name unsafe } }",
  "variables": { "authority": "simple.wikipedia.org", "pageName": "$1", "name": "$2" }
}
EOF
}

query "$1" "$2" | curl -XPOST localhost:8080/query -d @- 2>/dev/null | json_pp
