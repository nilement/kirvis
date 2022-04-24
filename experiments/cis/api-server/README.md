## API Server experiments

These experiments test changing the run-time parameters of K8s API Server.

## Experiment workflow

In `k8s/experiments.yaml` file, values for CIS Benchmark relevant parameters are provided.
Each experiment consists of a:
- key (string) - behaves as an ID of an experiment. It is mapped from the CIS Benchmark 1.6.0 Control-plane configuration items.
- parameter (string) - run-time parameter of K8s API Server that will be affected. (ex. `--anonymous-auth`, `--authorization-mode`)
- action (string) - decides which way the configuration should be applied (`set` - add the provided parameter-value pair if they do not exist, `setValue` - sets a value for an existing parameter, `pushValue` - adds a new value to an existing parameter, `removeValue` - removes a value from an existing parameters).
- value (string) - value to be used in the action.

To use the experiment

1. Apply experiment pod on a control-plane in your cluster. Example configuration in `k8s/pod.yaml` can be used.
2. Upon startup, it will backup the `kube-apiserver.yaml` manifest to it's local storage. (as of this version)
3. If the kube-apiserver manifest edit succeeds, it will check the `kube-apiserver` process for the changed values, to validate it is successful. (therefore `hostPID: true` is neccessary for the pod).
4. Upon deletion, the pod will restore the backed up `kube-apiserver.yaml` manifest.