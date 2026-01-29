output "api_url" {
  value = aws_apigatewayv2_api.url_shortener_api.api_endpoint
}


output "hello_endpoint_curl" {
  value = "curl -X POST ${aws_apigatewayv2_api.url_shortener_api.api_endpoint}/short_url -H 'Content-Type: application/json' -d '{\"url\": \"https://exanubes.com\", \"one_time_link\": false}'"
}

