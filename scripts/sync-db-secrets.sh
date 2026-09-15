#!/usr/bin/env bash
# Materializa os Secrets do Kubernetes a partir do AWS Secrets Manager,
# evitando qualquer credencial em texto plano no repositório.
set -euo pipefail

NAMESPACE="${NAMESPACE:-togglemaster}"
PROJECT="${PROJECT:-togglemaster}"
ENVIRONMENT="${ENVIRONMENT:-dev}"
SERVICES=(auth flag targeting)

kubectl get namespace "$NAMESPACE" >/dev/null 2>&1 || kubectl create namespace "$NAMESPACE"

for svc in "${SERVICES[@]}"; do
  secret_id="${PROJECT}-${ENVIRONMENT}-${svc}-db"
  echo "Lendo segredo ${secret_id} do Secrets Manager..."

  secret_json=$(aws secretsmanager get-secret-value --secret-id "$secret_id" --query SecretString --output text)

  username=$(echo "$secret_json" | jq -r .username)
  password=$(echo "$secret_json" | jq -r .password)
  host=$(echo "$secret_json" | jq -r .host)
  port=$(echo "$secret_json" | jq -r .port)
  dbname=$(echo "$secret_json" | jq -r .dbname)

  kubectl -n "$NAMESPACE" create secret generic "${svc}-db-credentials" \
    --from-literal=DB_HOST="$host" \
    --from-literal=DB_PORT="$port" \
    --from-literal=DB_NAME="$dbname" \
    --from-literal=DB_USERNAME="$username" \
    --from-literal=DB_PASSWORD="$password" \
    --dry-run=client -o yaml | kubectl apply -f -

  echo "Secret ${svc}-db-credentials sincronizado no namespace ${NAMESPACE}."
done

# Segredo de assinatura JWT do serviço `auth` (não gerenciado pelo Secrets Manager
# por não estar ligado a nenhum recurso de infraestrutura específico).
if ! kubectl -n "$NAMESPACE" get secret auth-jwt-secret >/dev/null 2>&1; then
  echo "Gerando auth-jwt-secret..."
  kubectl -n "$NAMESPACE" create secret generic auth-jwt-secret \
    --from-literal=secret="$(openssl rand -base64 32)"
else
  echo "auth-jwt-secret já existe, mantendo o valor atual."
fi
