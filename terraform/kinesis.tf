resource "aws_kinesis_stream" "visits" {
  name = "visits-stream"
  stream_mode_details {
    stream_mode = "ON_DEMAND"
  }
}

resource "aws_iam_role" "cloudfront_to_kinesis" {
  name = "cloudfront_to_kinesis"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "cloudfront.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "cloudfront_to_kinesis_policy" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "kinesis:DescribeStreamSummary",
        "kinesis:DescribeStream",
        "kinesis:PutRecord",
        "kinesis:PutRecords",
      ],
      Resource = aws_kinesis_stream.visits.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "cloudfront_to_kinesis_policy_attachment" {
  role       = aws_iam_role.cloudfront_to_kinesis.name
  policy_arn = aws_iam_policy.cloudfront_to_kinesis_policy.arn
}
