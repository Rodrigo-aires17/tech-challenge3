variable "name_prefix" {
  type = string
}

variable "github_org" {
  description = "Organização/usuário do GitHub dono do repositório (ex: meu-usuario)."
  type        = string
}

variable "github_repo" {
  description = "Nome do repositório GitHub (ex: togglemaster)."
  type        = string
}

variable "ecr_repository_arns" {
  description = "ARNs dos repositórios ECR que o pipeline de CI pode publicar imagens."
  type        = list(string)
}

variable "tags" {
  type    = map(string)
  default = {}
}
