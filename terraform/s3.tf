resource "aws_s3_bucket" "logs" {
  bucket        = "url-shortener-cloudfront-real-time-logs"
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "archive" {
  bucket = aws_s3_bucket.logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "archive" {
  bucket = aws_s3_bucket.logs.id

  rule {
    id     = "archive-tiering"
    status = "Enabled"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 90
      storage_class = "DEEP_ARCHIVE"
    }

  }
}
