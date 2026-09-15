# Backend remoto obrigatório (não pode ficar local).
# Os valores reais (bucket, key, region) são passados via `-backend-config=backend.hcl`
# para não fixar nomes de bucket específicos de cada aluno/grupo no código versionado.
terraform {
  backend "s3" {
    # bucket       = "<preenchido via backend.hcl>"
    # key          = "togglemaster/dev/terraform.tfstate"
    # region       = "<preenchido via backend.hcl>"
    # encrypt      = true
    # use_lockfile = true # lock nativo do backend S3 (Terraform >= 1.10), sem tabela DynamoDB extra
  }
}
