resource "aws_iam_role" "firehose_role" {
  name = "firehose_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "firehose.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "firehose_policy" {
  role = aws_iam_role.firehose_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "kinesis:GetRecords",
        "kinesis:GetShardIterator",
        "kinesis:DescribeStream",
        "kinesis:ListStreams",
      ],
      Resource = aws_kinesis_stream.visits.arn
      },
      {
        Effect = "Allow"
        Action = [
          "s3:AbortMultipartUpload",
          "s3:GetBucketLocation",
          "s3:ListBucket",
          "s3:PutObject",
        ],
        Resource = [
          aws_s3_bucket.logs.arn,
          "${aws_s3_bucket.logs.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:PutLogEvents",
        ],
        Resource = "*" # TODO:
      }

    ]
  })
}

resource "aws_kinesis_firehose_delivery_stream" "logs_delivery_stream" {
  name        = "cloudfront_raw_logs_firehose"
  destination = "extended_s3"

  kinesis_source_configuration {
    kinesis_stream_arn = aws_kinesis_stream.visits.arn
    role_arn           = aws_iam_role.firehose_role.arn
  }

  extended_s3_configuration {
    role_arn            = aws_iam_role.firehose_role.arn
    bucket_arn          = aws_s3_bucket.logs.arn
    buffering_size      = 64
    buffering_interval  = 300
    compression_format  = "GZIP"
    prefix              = "raw/year=!{timestamp:yyyy}/month=!{timestamp:MM}/day=!{timestamp:dd}/"
    error_output_prefix = "errors/"

    cloudwatch_logging_options {
      enabled         = true
      log_group_name  = "/aws/firehose/cloudfront-logs"
      log_stream_name = "s3-delivery"
    }
  }
}
