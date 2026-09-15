output "vpc_id" {
  value = module.networking.vpc_id
}

output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "ecr_repository_urls" {
  value = module.ecr.repository_urls
}

output "sqs_queue_url" {
  value = module.sqs.queue_url
}

output "dynamodb_table_name" {
  value = module.dynamodb.table_name
}

output "redis_primary_endpoint" {
  value = module.elasticache.primary_endpoint_address
}

output "rds_endpoints" {
  value = { for k, v in module.rds : k => v.endpoint }
}

output "evaluation_irsa_role_arn" {
  value = module.irsa.evaluation_role_arn
}

output "analytics_irsa_role_arn" {
  value = module.irsa.analytics_role_arn
}

output "github_actions_role_arn" {
  description = "Role ARN a ser configurada no secret AWS_ROLE_TO_ASSUME do GitHub Actions."
  value       = module.github_oidc.role_arn
}
