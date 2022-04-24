# Test environment

Kirvis can be tested with kind by using [kind](https://kind.sigs.k8s.io/), but feel free to use Minikube or any other Kubernetes distribution, as the project depends on Kubernetes version.

For some experiments, interaction with the control plane is required.

<p>
Test cluster to run experiments on 

1. `./kind/cluster-with-registry.sh` - Kind cluster with two nodes and an image registry.
2. (optional) `./observ/loki-prom-grafana.sh` - Loki-Promtail-Grafana observability stack.