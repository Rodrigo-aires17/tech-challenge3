variable "name_prefix" {
  type = string
}

variable "namespace" {
  description = "Namespace Kubernetes onde os Service Accounts residem."
  type        = string
  default     = "togglemaster"
}

variable "oidc_provider_arn" {
  type = string
}

variable "oidc_provider_url" {
  type = string
}

variable "sqs_queue_arn" {
  type = string
}

variable "dynamodb_table_arn" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
