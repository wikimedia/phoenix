variable "file_path" {
  description = "Path to zip file with lambda"
  type = string
}

variable "function_name" {
  description = "Name of a lambda function"
  type = string
}

variable "write_buckets" {
  description = "List of buckets with write access"
  type = list(string)
}

variable "read_buckets" {
  description = "List of buckets with read access"
  type = list(string)
}

variable "sns_subscription_topic" {
  description = "SNS topic lambda is subscribed to"
  type = string
}

variable "sns_publish_topics" {
  description = "SNS topics lambda publishes to"
  type = list(string)
}

variable "dynamodb_tables" {
  description = "Dynamodb tables a lambda has access to"
  type = list(string)
}
