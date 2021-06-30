variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "--PUT YOUR AWS REGION HERE--"
}

variable "aws_profile" {
  description = "AWS profile name"
  type        = string
  default     = "--PUT YOUR PROFILE HERE--"
}

# "phoenix-" is a prefix for all required resources and this should match the one from ".config.mk" in the root directory
variable "lambdas" {
  type = list(object({
    name=string,
    path=string,
    sns_subscription_topic=string,
    sns_publish_topics=list(string)
    write_buckets=list(string)
    read_buckets=list(string)
    dynamodb_tables=list(string)
  }))
  default = [{
      name = "phoenix-fetch-changed",
      path = "../lambdas/fetch-changed/function.zip",
      sns_publish_topics=["phoenix-sns-raw-content-incoming"],
      sns_subscription_topic="phoenix-event-streams-bridge",
      write_buckets=["phoenix-raw-content"],
      read_buckets=[]
      dynamodb_tables=[]
    },
    {
      name = "phoenix-fetch-schemaorg",
      path = "../lambdas/fetch-schema.org/function.zip",
      sns_subscription_topic="phoenix-sns-raw-content-incoming",
      sns_publish_topics=["phoenix-sns-raw-content-schemaorg"],
      write_buckets=["phoenix-raw-content"],
      read_buckets=["phoenix-raw-content"]
      dynamodb_tables=[]
    },
    {
      name = "phoenix-merge-schemaorg",
      path = "../lambdas/merge-schema.org/function.zip",
      sns_subscription_topic="phoenix-sns-raw-content-schemaorg",
      write_buckets=["phoenix-raw-content"],
      sns_publish_topics=[]
      read_buckets=["phoenix-raw-content"]
      dynamodb_tables=[]
    },
    {
      name = "phoenix-transform-parsoid",
      path = "../lambdas/transform-parsoid/function.zip",
      sns_subscription_topic="phoenix-sns-raw-content-schemaorg",
      sns_publish_topics=["phoenix-sns-node-published"]
      write_buckets=["phoenix-structured-content", "phoenix-raw-content"],
      dynamodb_tables=["phoenix-node-names", "phoenix-page-titles"]
      read_buckets=["phoenix-raw-content", "phoenix-structured-content"]
    },
  ]
}
