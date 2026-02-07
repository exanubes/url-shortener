resource "aws_sqs_queue" "link_visited_dlq" {
  name = "link_visited_dead_letter"
}

resource "aws_sqs_queue" "link_visited_queue" {
  name                      = "link_visited_messages"
  delay_seconds             = 0
  max_message_size          = 1024
  message_retention_seconds = 60 * 60 # 1 hour
  receive_wait_time_seconds = 10
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.link_visited_dlq.arn
    maxReceiveCount     = 4
  })
}

resource "aws_sqs_queue_redrive_allow_policy" "sqs_redrive_policy" {
  queue_url = aws_sqs_queue.link_visited_dlq.id
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.link_visited_queue.arn]
  })
}


resource "aws_sqs_queue" "cloudfront_rt_logs_dlq" {
  name = "cloudfront_rt_logs"
}
