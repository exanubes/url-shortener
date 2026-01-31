resource "aws_iam_role" "create_short_url" {
  name = "create_short_url_exec_role"


  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "create_short_url_dynamodb_policy" {
  name = "lambda_dynamodb_create_short_url"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "dynamodb:PutItem"
      ],
      Resource = aws_dynamodb_table.url_shortener.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "create_short_url_dynamodb_policy_attachment" {
  role       = aws_iam_role.create_short_url.name
  policy_arn = aws_iam_policy.create_short_url_dynamodb_policy.arn
}

resource "aws_iam_role_policy_attachment" "create_short_url_execution_role_attachment" {
  role       = aws_iam_role.create_short_url.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "create_short_url" {
  function_name    = "create_short_url"
  role             = aws_iam_role.create_short_url.arn
  filename         = "../dist/create/function.zip"
  source_code_hash = filebase64sha256("../dist/create/function.zip")
  handler          = "bootstrap"
  runtime          = "provided.al2"
  architectures    = ["arm64"]
}
