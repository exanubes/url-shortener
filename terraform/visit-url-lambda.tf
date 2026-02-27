resource "aws_iam_role" "visit_url" {
  name = "visit_url_exec_role"


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

resource "aws_iam_role_policy_attachment" "visit_url_execution_role_attachment" {
  role       = aws_iam_role.visit_url.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "visit_url_dynamodb_policy" {
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

resource "aws_iam_role_policy_attachment" "visit_url_dynamodb_policy_attachment" {
  role       = aws_iam_role.visit_url.name
  policy_arn = aws_iam_policy.visit_url_dynamodb_policy.arn
}

resource "aws_lambda_function" "visit_url" {
  function_name    = "visit_url"
  role             = aws_iam_role.visit_url.arn
  filename         = "../dist/visit/function.zip"
  source_code_hash = filebase64sha256("../dist/visit/function.zip")
  handler          = "bootstrap"
  runtime          = "provided.al2"
  architectures    = ["arm64"]
}
