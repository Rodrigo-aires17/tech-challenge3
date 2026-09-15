variable "chart_version" {
  type = string
}

variable "gitops_repo_url" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
