Hydra gets installed first, it needs its secret created though:

```shell
kubectl create secret generic kratos-secrets \
  --namespace auth \
  --from-literal='dsn=postgresql://kratos:'"${KRATOS_POSTGRES_PASSWORD}"'@cluster-pg-rw.cnpg-system.svc.cluster.local:5432/kratos?sslmode=disable' \
  --from-literal="smtpConnectionURI=${KRATOS_SMTP_URI}" \
  --from-literal="github_client_id=${KRATOS_GITHUB_CLIENT_ID}" \
  --from-literal="github_client_secret=${KRATOS_GITHUB_CLIENT_SECRET}"
```

```shell
kubectl create secret generic hydra-secrets \
  --namespace auth \
  --from-literal=dsn='postgresql://hydra:'"${HYDRA_POSTGRES_PASSWORD}"'@cluster-pg-rw.cnpg-system.svc.cluster.local:5432/hydra?sslmode=disable' \
  --from-literal=secretsCookie=$(openssl rand -base64 32) \
  --from-literal=secretsSystem=$(openssl rand -base64 32)
secret/hydra-secrets created
```
