# Copie para backend.hcl (arquivo ignorado pelo git) e ajuste os valores.
# Crie o bucket previamente, ex:
#   aws s3api create-bucket --bucket tfstate-fiap-tc3 --region us-east-1
#   aws s3api put-bucket-versioning --bucket <seu-bucket-tfstate> --versioning-configuration Status=Enabled
#
# terraform init -backend-config=backend.hcl

bucket       = "tfstate-fiap-tc3"
key          = "togglemaster/dev/terraform.tfstate"
region       = "us-east-1"
encrypt      = true
