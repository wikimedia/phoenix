Phoenix Terraform Module for Lambdas
=============

A Terraform module that allows to deliver all AWS lambda functions from `lambdas` directory.
The module allows to give each lambda all  permissions for accessing required resources 
such as S3 buckets, DynamoDB tables, SNS channels and subscriptions.

The module allows to pass a list of lambda configurations and iterate through the list applying all 
required permissions respectively to each lambda  

Module Input Variables
----------------------

- `file_path` -  Path to zip file with lambda.
- `function_name` - Name of a lambda function.
- `write_buckets` - A list of S3 buckets a lambda should have a read access to.
- `read_buckets` - A list of S3 buckets a lambda should have a write access to.
- `sns_subscription_topic` - An SNS trigger (subscription) for a lambda function.
- `sns_publish_topics` - A list of SNS topics a lambda should be able to publish to.
- `dynamodb_tables` - A list of Dynamodb tables a lambda has access to

Usage
-----

```hcl
# In your "main.tf" file:
# 1. Define 'lambdas' var

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
    name = "lambda-one",
    path = "../lambdas/lambda-one/function.zip",
    sns_publish_topics=["topic-one"],
    sns_subscription_topic="topic-two",
    write_buckets=["bucket1", "bucket2"],
    read_buckets=["bucket2", "bucket3"]
    dynamodb_tables=["table1", "table2"]
  }]
}

# 2. Initialize the module
module "lambdas" {
  source = 'path to the module'
  count = length(var.lambdas)
  function_name = var.lambdas[count.index].name
  file_path = var.lambdas[count.index].path
  sns_subscription_topic = var.lambdas[count.index].sns_subscription_topic
  sns_publish_topics = var.lambdas[count.index].sns_publish_topics
  write_buckets = var.lambdas[count.index].write_buckets
  read_buckets = var.lambdas[count.index].read_buckets
  dynamodb_tables = var.lambdas[count.index].dynamodb_tables
}
```
