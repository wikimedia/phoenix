${jsonencode({
  "Version": "2012-10-17",
  "Statement": [{
      "Effect": "Allow",
      "Action": ["s3:PutObject","s3:PutObjectAcl"],
      "Resource": [for name in bucket_names : "arn:aws:s3:::${name}/*"]
  }]
})}
