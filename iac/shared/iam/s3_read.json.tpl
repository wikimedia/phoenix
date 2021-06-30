${jsonencode({
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ListObjectsInBucket",
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": [for name in bucket_names : "arn:aws:s3:::${name}"]
    },
    {
      "Sid": "GetObjectActions",
      "Effect": "Allow",
      "Action": ["s3:GetObject","s3:GetObjectVersion", "s3:GetObjectAcl"],
      "Resource": [for name in bucket_names : "arn:aws:s3:::${name}/*"]
    }
  ]
})}
