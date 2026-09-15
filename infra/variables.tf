variable "project_name" {
  description = "Nome curto do projeto, usado como prefixo dos recursos."
  type        = string
  default     = "togglemaster"
}

variable "environment" {
  description = "Nome do ambiente (dev, staging, prod)."
  type        = string
  default     = "dev"
}

variable "aws_region" {
  description = "Região AWS onde os recursos serão provisionados."
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block da VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability Zones utilizadas."
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

variable "public_subnet_cidrs" {
  description = "CIDRs das subnets públicas (uma por AZ)."
  type        = list(string)
  default     = ["10.0.0.0/24", "10.0.1.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDRs das subnets privadas (uma por AZ)."
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
}

variable "single_nat_gateway" {
  description = "Se true, cria um único NAT Gateway (economia); se false, um por AZ (alta disponibilidade)."
  type        = bool
  default     = true
}

variable "services" {
  description = "Lista dos 5 microsserviços do ToggleMaster."
  type        = list(string)
  default     = ["auth", "flag", "targeting", "evaluation", "analytics"]
}

variable "rds_services" {
  description = "Microsserviços que possuem um banco RDS PostgreSQL dedicado."
  type        = list(string)
  default     = ["auth", "flag", "targeting"]
}

# --- EKS ---
variable "cluster_version" {
  description = "Versão do Kubernetes no EKS."
  type        = string
  default     = "1.30"
}

variable "node_instance_types" {
  description = "Tipos de instância EC2 do Node Group."
  type        = list(string)
  default     = ["t3.medium"]
}

variable "node_desired_size" {
  type    = number
  default = 2
}

variable "node_min_size" {
  type    = number
  default = 2
}

variable "node_max_size" {
  type    = number
  default = 4
}

# --- RDS ---
variable "db_engine_version" {
  type    = string
  default = "16.4"
}

variable "db_instance_class" {
  type    = string
  default = "db.t3.micro"
}

variable "db_allocated_storage" {
  type    = number
  default = 20
}

# --- ElastiCache ---
variable "redis_node_type" {
  type    = string
  default = "cache.t3.micro"
}

# --- ArgoCD ---
variable "argocd_chart_version" {
  type    = string
  default = "7.7.11"
}

variable "gitops_repo_url" {
  description = "URL do repositório Git (HTTPS) usado pelo ArgoCD para sincronizar os manifestos em gitops/apps."
  type        = string
}

variable "gitops_target_revision" {
  description = "Branch/tag monitorada pelo ArgoCD."
  type        = string
  default     = "main"
}

variable "github_org" {
  description = "Organização/usuário do GitHub dono do repositório (para o OIDC do GitHub Actions)."
  type        = string
}

variable "github_repo" {
  description = "Nome do repositório GitHub (para o OIDC do GitHub Actions)."
  type        = string
  default     = "togglemaster"
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
