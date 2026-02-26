resource "aws_iam_role" "cloudfront_rt_logs_processor_role" {
  name = "cloudfront_rt_logs_processor_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "log_processor_kinesis_policy" {
  role = aws_iam_role.cloudfront_rt_logs_processor_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "kinesis:GetRecords",
        "kinesis:GetShardIterator",
        "kinesis:DescribeStream",
        "kinesis:ListShards",
      ],
      Resource = aws_kinesis_stream.visits.arn
    }]
  })
}

resource "aws_iam_role_policy" "log_processor_sqs_policy" {
  role = aws_iam_role.cloudfront_rt_logs_processor_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "sqs:SendMessage",
      ],
      Resource = aws_sqs_queue.cloudfront_rt_logs_dlq.arn
    }]
  })
}

resource "aws_iam_role_policy" "visit_url_dynamodb_policy" {
  role = aws_iam_role.cloudfront_rt_logs_processor_role.name
  name = "lambda_dynamodb_visit_url"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "dynamodb:TransactWriteItems",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
      ],
      Resource = aws_dynamodb_table.url_shortener.arn
    }]
  })
}


resource "aws_lambda_function" "cloudfront_rt_logs_processor" {
  function_name    = "cloudfront_rt_logs_processor"
  role             = aws_iam_role.cloudfront_rt_logs_processor_role.arn
  filename         = "../dist/processor/function.zip"
  source_code_hash = filebase64sha256("../dist/processor/function.zip")
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

resource "aws_iam_role_policy_attachment" "log_processor_execution_role_attachment" {
  role       = aws_iam_role.cloudfront_rt_logs_processor_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
