
 # IAM role which dictates what other AWS services the Lambda function
 # may access.
resource "aws_iam_role" "lambda_exec" {
  name = "lambda_${var.function_name}_execution_role"
  assume_role_policy = "${file("./shared/iam/lamda_execution_role.json")}"
}

resource "aws_lambda_function" "lambda" {
  function_name = var.function_name
  filename = var.file_path
  handler = "main"
  runtime = "go1.x"
  source_code_hash = filebase64sha256(var.file_path)

  role = aws_iam_role.lambda_exec.arn
}

// Cloudwatch (logging)

resource "aws_cloudwatch_log_group" "lambda_log" {
  name        = "/aws/lambda/${var.function_name}"
  retention_in_days = 14
}

resource "aws_iam_policy" "lambda_logging_policy" {
  name        = "lambda_${var.function_name}_logging_policy"
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy = "${file("./shared/iam/lambda_logging_policy.json")}"
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_logging_policy.arn
}

// SNS subscription (lambda trigger)

resource "aws_sns_topic" "lambda_trigger" {
  name = var.sns_subscription_topic
}

resource "aws_sns_topic_subscription" "topic_lambda" {
  topic_arn = "${aws_sns_topic.lambda_trigger.arn}"
  protocol  = "lambda"
  endpoint  = "${aws_lambda_function.lambda.arn}"
}

resource "aws_lambda_permission" "with_sns" {
  statement_id = "AllowExecutionFromSNS"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda.arn}"
  principal = "sns.amazonaws.com"
  source_arn = "${aws_sns_topic.lambda_trigger.arn}"
}

// SNS publish

resource "aws_sns_topic" "publish_topic" {
  count = length(var.sns_publish_topics)
  name = "${var.sns_publish_topics[count.index]}"
}

data "template_file" "sns_publish_template" {
  template = "${file("./shared/iam/sns_publish.json.tpl")}"
  count = length(var.sns_publish_topics)
  vars = {
    resource = "${aws_sns_topic.publish_topic[count.index].arn}"
  }
}

resource "aws_iam_policy" "sns_publish_policy" {
  name = "lambda_${var.function_name}_sns_publish_policy"
  count = length(var.sns_publish_topics)
  policy = "${data.template_file.sns_publish_template[count.index].rendered}"
}

resource "aws_iam_role_policy_attachment" "sns_publich" {
  role  = aws_iam_role.lambda_exec.name
  count = length(var.sns_publish_topics)
  policy_arn = "${aws_iam_policy.sns_publish_policy[count.index].arn}"
}

// S3 write access

resource "aws_iam_policy" "s3_write_policy" {
 name = "lambda_${var.function_name}_s3_write_policy"
 count = length(var.write_buckets) > 0 ? 1 : 0
 policy = templatefile("./shared/iam/s3_write.json.tpl", {
   bucket_names = var.write_buckets
 })
}

resource "aws_iam_role_policy_attachment" "s3_write" {
  role  = aws_iam_role.lambda_exec.name
  count = length(var.write_buckets) > 0 ? 1 : 0
  policy_arn = "${aws_iam_policy.s3_write_policy[count.index].arn}"
}


// S3 read access

resource "aws_iam_policy" "s3_read_policy" {
  count = length(var.read_buckets) > 0 ? 1 : 0
  name = "lambda_${var.function_name}_s3_read_policy"
  policy = templatefile("./shared/iam/s3_read.json.tpl", {
    bucket_names = var.read_buckets
  })
}

resource "aws_iam_role_policy_attachment" "s3_read" {
  count = length(var.read_buckets) > 0 ? 1 : 0
  role  = aws_iam_role.lambda_exec.name
  policy_arn = "${aws_iam_policy.s3_read_policy[count.index].arn}"
}

// DynamoDB table access

data "aws_dynamodb_table" "dynamodb_table" {
 count = length(var.dynamodb_tables)
 name = "${var.dynamodb_tables[count.index]}"
}

data "template_file" "dynamodb_table_policy_template" {
 count = length(var.dynamodb_tables)
 template = "${file("./shared/iam/dynamodb_table_policy.tpl")}"
 vars = {
   resource = "${data.aws_dynamodb_table.dynamodb_table[count.index].arn}"
 }
}

resource "aws_iam_policy" "dynamodb_table_policy" {
 count = length(var.dynamodb_tables)
 name = "lambda_${var.function_name}_${var.dynamodb_tables[count.index]}_dynamodb_table_policy"
 policy = "${data.template_file.dynamodb_table_policy_template[count.index].rendered}"
}

resource "aws_iam_role_policy_attachment" "dynamodb_table" {
 count = length(var.dynamodb_tables)
 role  = aws_iam_role.lambda_exec.name
 policy_arn = "${aws_iam_policy.dynamodb_table_policy[count.index].arn}"
}
