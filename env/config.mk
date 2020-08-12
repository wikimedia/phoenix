
PHX_ACCOUNT_ID               = 195062990486
PHX_DEFAULT_REGION           = us-east-2
PHX_PREFIX                   = scpoc

######
# SNS resources
######

# Topic that receives change events originating from the Wikimedia
# Event Streams service.
PHX_SNS_EVENT_STREAMS_BRIDGE      = $(PHX_PREFIX)-event-streams-bridge
PHX_SNS_EVENT_STREAMS_BRIDGE_ARN  = $(shell printf "$(_BASE_ARN)" sns "$(PHX_SNS_EVENT_STREAMS_BRIDGE)")

# Topic that receives events when new HTML is added to incoming (see
# PHX_S3_RAW_CONTENT_INCOMING).
PHX_SNS_RAW_CONTENT_INCOMING      = $(PHX_PREFIX)-sns-raw-content-incoming
PHX_SNS_RAW_CONTENT_INCOMING_ARN  = $(shell printf "$(_BASE_ARN)" sns "$(PHX_SNS_RAW_CONTENT_INCOMING)")

# Topic that receives events when new linked-data (Wikidata) is added
# to the raw content store (see PHX_S3_RAW_CONTENT_WD_LINKED)
PHX_SNS_RAW_CONTENT_WD_LINKED     = $(PHX_PREFIX)-sns-raw-content-schemaorg
PHX_SNS_RAW_CONTENT_WD_LINKED_ARN = $(shell printf "$(_BASE_ARN)" sns "$(PHX_SNS_RAW_CONTENT_WD_LINKED)")

######
# S3 resources
######

# The "raw content" bucket; Corresponds with uses of "raw content
# store" in the architecture documents.
PHX_S3_RAW_CONTENT_BUCKET        = $(PHX_PREFIX)-raw-content-store

# Folder where HTML documents of a corresponding revision are
# downloaded to after a change event is received.
PHX_S3_RAW_CONTENT_INCOMING      = incoming

# Folder where linked data (in the schema.org vocabulary) is stored.
PHX_S3_RAW_CONTENT_WD_LINKED     = schema.org

# Folder where HTML augmented with linked data is stored.
PHX_S3_RAW_CONTENT_LINKED_HTML   = linked-html

# The "structured content" bucket, where parsed and transformed data are
# stored in canonical format
PHX_S3_STRUCTURED_CONTENT_BUCKET = $(PHX_PREFIX)-structured-content-store

######
# Lambda resources
######

# Lambda function subscribed to Wikimedia Event Stream change events.
# Downloads the corresponding HTML (revision) and writes it to S3 (see
# PHX_S3_RAW_CONTENT_INCOMING)
PHX_LAMBDA_FETCH_CHANGED   = $(PHX_PREFIX)-fetch-changed


# Function invoked when new content has been added to incoming
# (PHX_SNS_RAW_CONTENT_INCOMING).  Downloads corresponding Wikidata
# information, constructs linked data (JSON-LD) in the schema.org
# vocabulary, and uploads to S3 (PHX_S3_RAW_CONTENT_WD_LINKED)
PHX_LAMBDA_FETCH_SCHEMAORG = $(PHX_PREFIX)-lambda-fetch-schemaorg

# Lambda subscribed to events that signal the creation of new Wikidata
# linked data (PHX_SNS_RAW_CONTENT_WD_LINKED).  Transforms the HTML
# from incoming (PHX_S3_RAW_CONTENT_INCOMING) to include the linked
# data (as JSON-LD), and uploads the result
# (PHX_S3_RAW_CONTENT_LINKED_HTML)
PHX_LAMBDA_MERGE_SCHEMAORG = $(PHX_PREFIX)-lambda-merge-schemaorg

# Lambda subscribed to events that signal the saving raw content to S3 
# storage, transforms raw content into canonical tructure and save to S3 
# storage (See PHX_S3_CANONICAL_CONTENT_BUCKET)
PHX_LAMBDA_TRANSFORM_PARSOID = $(PHX_PREFIX)-transform-parsoid

# For internal use in ARN string formatting
_BASE_ARN = $(shell printf "arn:aws:%%s:%s:%s:%%s" "$(PHX_DEFAULT_REGION)" "$(PHX_ACCOUNT_ID)")

