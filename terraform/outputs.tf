output "api_url" {
  value = aws_apigatewayv2_api.http_api.api_endpoint
}


output "hello_endpoint_curl" {
  value = "curl ${aws_apigatewayv2_api.http_api.api_endpoint}/hello"
}

