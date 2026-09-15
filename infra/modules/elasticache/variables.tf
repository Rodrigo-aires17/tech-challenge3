variable "name_prefix" {
  type = string
}

variable "node_type" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "allowed_security_group_ids" {
  type    = list(string)
  default = []
}

variable "vpc_cidr_block" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
