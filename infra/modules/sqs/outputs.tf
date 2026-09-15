output "queue_url" {
  value = aws_sqs_queue.analytics_events.url
}

output "queue_arn" {
  value = aws_sqs_queue.analytics_events.arn
}

output "dlq_arn" {
  value = aws_sqs_queue.dlq.arn
}
