resource "aws_sqs_queue" "dlq" {
  name                      = "${var.name_prefix}-analytics-events-dlq"
  message_retention_seconds = 1209600 # 14 dias
  sqs_managed_sse_enabled   = true
  tags                      = var.tags
}

resource "aws_sqs_queue" "analytics_events" {
  name                       = "${var.name_prefix}-analytics-events"
  visibility_timeout_seconds = 60
  message_retention_seconds  = 345600 # 4 dias
  sqs_managed_sse_enabled    = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 5
  })

  tags = var.tags
}
