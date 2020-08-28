#!/bin/bash

set -e


query()
{
    cat <<EOF
{
  "query": "query Page(\$authority: String!, \$name: String!) { pageByName(authority: \$authority, name: \$name) { name dateModified hasPart about { key val } } }",
  "variables": { "authority": "simple.wikipedia.org", "name": "$1" }
}
EOF
}


query "$1" | curl -XPOST localhost:8080/query -d @- 2>/dev/null | json_pp
