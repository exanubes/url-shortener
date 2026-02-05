resource "aws_apigatewayv2_api" "url_shortener_api" {
  name          = "url_shortener_api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.url_shortener_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "create_short_url_integration" {
  api_id           = aws_apigatewayv2_api.url_shortener_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.create_short_url.invoke_arn

  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "create_short_url_route" {
  api_id    = aws_apigatewayv2_api.url_shortener_api.id
  route_key = "POST /short_url"
  target    = "integrations/${aws_apigatewayv2_integration.create_short_url_integration.id}"
}

resource "aws_lambda_permission" "allow_apigw_create_short_url" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.create_short_url.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.url_shortener_api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_integration" "resolve_url_integration" {
  api_id           = aws_apigatewayv2_api.url_shortener_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.resolve_url.invoke_arn

  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "resolve_url_route" {
  api_id    = aws_apigatewayv2_api.url_shortener_api.id
  route_key = "GET /{short_code}"
  target    = "integrations/${aws_apigatewayv2_integration.resolve_url_integration.id}"
}

resource "aws_apigatewayv2_route" "resolve_url_head_route" {
  api_id    = aws_apigatewayv2_api.url_shortener_api.id
  route_key = "HEAD /{short_code}"
  target    = "integrations/${aws_apigatewayv2_integration.resolve_url_integration.id}"
}

resource "aws_lambda_permission" "allow_apigw_resolve_url" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.resolve_url.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.url_shortener_api.execution_arn}/*/*"
}

