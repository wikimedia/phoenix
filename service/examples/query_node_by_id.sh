#!/bin/bash

set -e


query()
{
    cat <<EOF
{
  "query": "query Node(\$id: ID!) { node(id: \$id) { dateModified unsafe } }",
  "variables": { "id": "$1" }
}
EOF
}

query "$1" | curl -XPOST localhost:8080/query -d @- 2>/dev/null | json_pp
