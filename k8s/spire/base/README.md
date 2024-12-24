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

## Joining an agent to a spire server

IF you need to run a local agent or an agent outside of kubernetes, you'll need to register it.
You can do this by generating a join token from the `spire-server` instance in k8s.

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server token generate -spiffeID spiffe://tunnel.threadr.ai/ns/eventrunner/sa/api \
```

This will create some output, a join token.

```shell
export JOIN_TOKEN=...
```

**Note**: that is creating a token for the `api` pod in the `eventrunner` namespace (ns)

## Registering a workload based on a UNIX id

in this scenario, we have an spire-agent running on a server or laptop,
our API connects to it over a unix:/// socket and gets a SVID. First we need
to register 

```shell
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
    -spiffeID "spiffe://tunnel.threadr.ai/my-local-workload" \
    -parentID "spiffe://tunnel.threadr.ai/spire/agent/join_token/$JOIN_TOKEN" \
    -selector "unix:uid:502"
```

