resource "aws_iam_role" "lambda_exec" {
  name = "lambda-exec-role"


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

resource "aws_iam_role_policy_attachment" "execution_role" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "testing" {
  function_name    = "testing"
  role             = aws_iam_role.lambda_exec.arn
  filename         = "../dist/function.zip"
  source_code_hash = filebase64sha256("../dist/function.zip")
  handler          = "bootstrap"
  runtime          = "provided.al2"
  architectures    = ["arm64"]
}
