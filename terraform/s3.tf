resource "aws_s3_bucket" "logs" {
  bucket        = "url-shortener-cloudfront-real-time-logs"
  force_destroy = true
}
