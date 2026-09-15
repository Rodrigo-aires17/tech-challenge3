variable "name_prefix" {
  type = string
}

variable "identifier" {
  description = "Nome curto do banco (ex: auth, flag, targeting)."
  type        = string
}

variable "engine_version" {
  type = string
}

variable "instance_class" {
  type = string
}

variable "allocated_storage" {
  type = number
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "allowed_security_group_ids" {
  description = "Security groups autorizados a acessar o banco (ex: SG do node group do EKS)."
  type        = list(string)
  default     = []
}

variable "vpc_cidr_block" {
  description = "Usado como fallback de acesso quando allowed_security_group_ids não é fornecido (ex: primeiro apply)."
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
