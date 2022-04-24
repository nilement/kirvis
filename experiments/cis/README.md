## CIS Benchmark experiments

These experiments test enabling/disabling of CIS benchmarks as in
https://www.cisecurity.org/benchmark/kubernetes/ , developed according to version 1.6.0

Most of the tests are simplistic, changing permissions/ownership of files to invalidated the benchmarks.
They can be then validated through usage of kube-bench (https://github.com/aquasecurity/kube-bench).

## Experiment workflow

1. Apply experiment pod to desired node in your Kubernetes cluster
2. Validate that you can observe changes being introduced 
3. Delete the pod which will restore the relevant configuration that was detected before the experiment