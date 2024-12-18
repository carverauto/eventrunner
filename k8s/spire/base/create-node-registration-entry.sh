#/bin/bash

set -e

bb=$(tput bold)
nn=$(tput sgr0)


echo "${bb}Creating registration entry for the node...${nn}"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -node  \
    -spiffeID spiffe://tunnel.threadr.ai/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:threadr-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent
