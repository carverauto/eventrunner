# spire server setup

## Setup secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: spire-postgres
  namespace: spire
type: Opaque
stringData:
  DB_PASSWORD: ""
```

## Create server

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://tunnel.threadr.ai/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:threadr-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent \
    -node
```

## Setup workload

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://tunnel.threadr.ai/ns/default/sa/default \
    -parentID spiffe://tunnel.threadr.ai/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default
```


