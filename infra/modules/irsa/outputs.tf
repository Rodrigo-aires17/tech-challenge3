output "evaluation_role_arn" {
  value = aws_iam_role.evaluation.arn
}

output "analytics_role_arn" {
  value = aws_iam_role.analytics.arn
}
