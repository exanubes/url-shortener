output "api_url" {
  value = aws_cloudfront_distribution.url_shortener.domain_name
}

output "gateway_url" {
  value = aws_apigatewayv2_api.url_shortener_api.api_endpoint
}


