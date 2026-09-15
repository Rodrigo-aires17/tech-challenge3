# Relatório de Entrega — Tech Challenge Fase 3

## Participantes

- Nome completo — RM
- Nome completo — RM

## Links

- Repositório: `<url-do-repo-github>`
- Vídeo de demonstração: `<url-do-video>`
- Documentação adicional: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## Resumo dos desafios encontrados e decisões tomadas

- **IAM**: como não utilizamos o AWS Academy, o Terraform cria diretamente as
  roles necessárias (cluster EKS, node group, IRSA para os pods `evaluation` e
  `analytics`), aplicando o princípio de menor privilégio (cada role só acessa
  os recursos que o serviço correspondente precisa).
- **Credenciais de banco**: eliminamos o problema de credenciais em texto plano
  gerando senhas aleatórias com `random_password` e armazenando-as no AWS
  Secrets Manager; um script (`scripts/sync-db-secrets.sh`) materializa os
  Secrets do Kubernetes a partir do Secrets Manager, nunca do Git.
- **Estado do Terraform**: configuramos backend remoto em S3 com
  `use_lockfile` para lock nativo, evitando estado local e conflitos entre
  desenvolvedores.
- **GitOps no monorepo**: optamos por uma pasta `gitops/` dentro do mesmo
  repositório (em vez de um repositório separado) para simplificar o fluxo de
  entrega da fase, com path-filters no CI para não disparar loops de pipeline
  quando o próprio CI atualiza a tag da imagem.
- **DevSecOps**: o pipeline falha automaticamente caso o Trivy (SCA/Container
  Scan) ou o gosec (SAST) encontrem vulnerabilidades classificadas como
  CRITICAL, conforme demonstrado no vídeo (inserimos uma dependência
  vulnerável propositalmente, mostramos a falha, corrigimos e mostramos o
  pipeline passando).
- **Outros desafios**: `<descreva aqui dificuldades específicas do seu grupo,
  ex: custos de NAT Gateway, tempo de provisionamento do EKS, troubleshooting
  do ArgoCD, etc.>`

## Estimativa de custos AWS

`<insira aqui o print da AWS Pricing Calculator / Cost Explorer>`
