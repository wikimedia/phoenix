terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "s3" {
    bucket = "--PUT S3 BUCKET HERE--"
    key    = "iac/state"
    region = "--PUT YOUR AWS REGION HERE--"
    dynamodb_table = "--PUT DYNAMODB TABLE HERE--"
    profile = "--PUT YOUR AWS REGION HERE--"
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

module "lambdas" {
  source = "./shared/lambda"
  count = length(var.lambdas)

  function_name = var.lambdas[count.index].name
  file_path = var.lambdas[count.index].path
  sns_subscription_topic = var.lambdas[count.index].sns_subscription_topic
  sns_publish_topics = var.lambdas[count.index].sns_publish_topics
  write_buckets = var.lambdas[count.index].write_buckets
  read_buckets = var.lambdas[count.index].read_buckets
  dynamodb_tables = var.lambdas[count.index].dynamodb_tables
}
