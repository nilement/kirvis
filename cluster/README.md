# Test environment

<p>This is the test environment for the security operator.</p>

Two variants: minikube and kind.

Kind is easier to run but minikube can be used with Virtualbox to simulate a whole virtualized node, which makes some configurations easier to test.

<p>
Test cluster to run experiments on 

1. `sh kind/cluster-with-registry.sh`
2. `kubectl apply -f nginx/deployment.yaml nginx/lb.yaml`
3. `kubectl port-forward service/example-service 13337:13337`
</p>

Current experiments: 

- Changing Kubelet service file permission on a Node. To be detected by kube-hunter.
- Uploading a different image to the image registry after deployment, before scaling up. Difference detected in logs.
