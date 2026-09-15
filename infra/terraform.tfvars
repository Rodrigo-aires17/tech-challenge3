# Copie para terraform.tfvars e ajuste conforme sua conta AWS pessoal.

aws_region   = "us-east-1"
project_name = "togglemaster"
environment  = "dev"

availability_zones   = ["us-east-1a", "us-east-1b"]
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.0.0/24", "10.0.1.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24"]
single_nat_gateway   = true

cluster_version      = "1.33"
node_instance_types  = ["t3.medium"]
node_desired_size    = 2
node_min_size        = 2
node_max_size         = 4

# Aponte para o seu fork/repositório GitHub (usado pelo ArgoCD para GitOps)
gitops_repo_url        = "https://github.com/Rodrigo-aires17/tech-challenge3.git"
gitops_target_revision = "main"

# Usado para restringir o OIDC do GitHub Actions (AWS_ROLE_TO_ASSUME) a este repositório
github_org  = "Rodrigo-aires17"
github_repo = "tech-challenge3"
