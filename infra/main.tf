module "networking" {
  source = "./modules/networking"

  name_prefix          = local.name_prefix
  vpc_cidr             = var.vpc_cidr
  availability_zones   = var.availability_zones
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  single_nat_gateway   = var.single_nat_gateway
  cluster_name         = "${local.name_prefix}-eks"
  tags                 = local.common_tags
}

module "iam" {
  source = "./modules/iam"

  name_prefix = local.name_prefix
  tags        = local.common_tags
}

module "eks" {
  source = "./modules/eks"

  name_prefix         = local.name_prefix
  cluster_version     = var.cluster_version
  cluster_role_arn    = module.iam.cluster_role_arn
  node_role_arn       = module.iam.node_role_arn
  public_subnet_ids   = module.networking.public_subnet_ids
  private_subnet_ids  = module.networking.private_subnet_ids
  node_instance_types = var.node_instance_types
  node_desired_size   = var.node_desired_size
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
  tags                = local.common_tags
}

module "sqs" {
  source = "./modules/sqs"

  name_prefix = local.name_prefix
  tags        = local.common_tags
}

module "dynamodb" {
  source = "./modules/dynamodb"

  table_name = "ToggleMasterAnalytics"
  tags       = local.common_tags
}

module "irsa" {
  source = "./modules/irsa"

  name_prefix        = local.name_prefix
  namespace          = "togglemaster"
  oidc_provider_arn  = module.eks.oidc_provider_arn
  oidc_provider_url  = module.eks.oidc_provider_url
  sqs_queue_arn      = module.sqs.queue_arn
  dynamodb_table_arn = module.dynamodb.table_arn
  tags               = local.common_tags
}

module "rds" {
  source   = "./modules/rds"
  for_each = toset(var.rds_services)

  name_prefix       = local.name_prefix
  identifier        = each.value
  engine_version    = var.db_engine_version
  instance_class    = var.db_instance_class
  allocated_storage = var.db_allocated_storage
  vpc_id            = module.networking.vpc_id
  subnet_ids        = module.networking.private_subnet_ids
  vpc_cidr_block    = module.networking.vpc_cidr_block
  tags              = local.common_tags
}

module "elasticache" {
  source = "./modules/elasticache"

  name_prefix    = local.name_prefix
  node_type      = var.redis_node_type
  vpc_id         = module.networking.vpc_id
  subnet_ids     = module.networking.private_subnet_ids
  vpc_cidr_block = module.networking.vpc_cidr_block
  tags           = local.common_tags
}

module "ecr" {
  source = "./modules/ecr"

  services = var.services
  tags     = local.common_tags
}

module "argocd" {
  source = "./modules/argocd"

  chart_version   = var.argocd_chart_version
  gitops_repo_url = var.gitops_repo_url
  tags            = local.common_tags

  depends_on = [module.eks]
}

module "github_oidc" {
  source = "./modules/github_oidc"

  name_prefix         = local.name_prefix
  github_org          = var.github_org
  github_repo         = var.github_repo
  ecr_repository_arns = [for name in var.services : "arn:aws:ecr:${var.aws_region}:${data.aws_caller_identity.current.account_id}:repository/togglemaster/${name}"]
  tags                = local.common_tags
}

data "aws_caller_identity" "current" {}
