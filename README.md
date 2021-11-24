# SCE Kubernetes experimental project

This is a project containing several experiments for trying out Security Chaos Engineering in Kubernetes. It contains several experiments in Kubernetes that can be deployed as pods in a Kubernetes cluster. Upon deletion, using a preStop hook, they will attempt to revert any caused changes to the cluster, without any guarantees. Use at your own caution.
