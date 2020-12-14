
PHX_ACCOUNT_ID               = 113698225543
PHX_DEFAULT_REGION           = us-east-2
PHX_PREFIX                   = scpoc

######
# SNS resources
######

# Topic that receives change events originating from the Wikimedia
# Event Streams service.
PHX_SNS_EVENT_STREAMS_BRIDGE      = $(PHX_PREFIX)-event-streams-bridge

# Topic that receives events when new HTML is added to incoming (see
# PHX_S3_RAW_CONTENT_INCOMING).
PHX_SNS_RAW_CONTENT_INCOMING      = $(PHX_PREFIX)-sns-raw-content-incoming

# Topic that receives events when new linked-data (Wikidata) is added
# to the raw content store (see PHX_S3_RAW_CONTENT_WD_LINKED)
PHX_SNS_RAW_CONTENT_WD_LINKED     = $(PHX_PREFIX)-sns-raw-content-schemaorg

# Topic that receives events when new Node objects are added to the
# Structured Content Store
PHX_SNS_NODE_PUBLISHED            = $(PHX_PREFIX)-sns-node-published


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
# storage (See PHX_S3_STRUCTURED_CONTENT_BUCKET)
PHX_LAMBDA_TRANSFORM_PARSOID = $(PHX_PREFIX)-lambda-transform-parsoid

# Lambda subscribed to events signaling that a new Node object has been stored.
# Retrieves related topic information for the Node, and stores the result.
PHX_LAMBDA_RELATED_TOPICS = $(PHX_PREFIX)-lambda-related-topics


######
# DynamoDB resources
######

# Table used to index page titles
PHX_DYNAMODB_PAGE_TITLES = $(PHX_PREFIX)-dynamodb-page-titles

# Table used to index node names
PHX_DYNAMODB_NODE_NAMES  = $(PHX_PREFIX)-dynamodb-node-names


######
# Search index resources
######

# Use ../.config.mk to specify actual credentials
PHX_SEARCH_USERNAME = fauxuser
PHX_SEARCH_PASSWORD = fauxpass

# Elasticsearch endpoint URL
PHX_SEARCH_ENDPOINT = "https://search-scpoc-phoenix-zti4iohw623mbsmdabsrmhybm4.us-east-2.es.amazonaws.com"


# For internal use in ARN string formatting
_BASE_ARN = $(shell printf "arn:aws:%%s:%s:%s:%%s" "$(PHX_DEFAULT_REGION)" "$(PHX_ACCOUNT_ID)")

# Include user/developer overrides
include ../.config.mk
