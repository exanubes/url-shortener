# log procesor lambda

resource "aws_iam_role" "lambda_to_kinesis" {
  name = "lambda_to_kinesis"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "lambda_to_kinesis_policy" {
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
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_to_kinesis_policy_attachment" {
  role       = aws_iam_role.lambda_to_kinesis.name
  policy_arn = aws_iam_policy.lambda_to_kinesis_policy.arn
}

resource "aws_lambda_function" "cloudfront_rt_logs_processor" {
  function_name    = "cloudfront_rt_logs_processor"
  role             = aws_iam_role.lambda_to_kinesis.arn
  filename         = "../dist/logs/function.zip"
  source_code_hash = filebase64sha256("../dist/logs/function.zip")
  handler          = "bootstrap"
  runtime          = "provided.al2"
  architectures    = ["arm64"]
}

resource "aws_lambda_event_source_mapping" "kinesis_trigger" {
  event_source_arn               = aws_kinesis_stream.visits.arn
  function_name                  = aws_lambda_function.cloudfront_rt_logs_processor.arn
  starting_position              = "LATEST"
  batch_size                     = 100
  enabled                        = true
  parallelization_factor         = 2
  maximum_retry_attempts         = 3
  bisect_batch_on_function_error = true

  destination_config {
    on_failure {
      destination_arn = aws_sqs_queue.cloudfront_rt_logs_dlq.arn
    }
  }
}
