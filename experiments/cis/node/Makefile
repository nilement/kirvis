build:
	GOOS=linux go build -ldflags "-s -w" -o bin/node .
	docker build -t localhost:5000/cispod:local .

push:
	docker push localhost:5000/cispod:local

reapply:
	- kubectl delete -f ./k8s/worker-node.yaml
	kubectl apply -f ./k8s/worker-node.yaml

bench:
	- kubectl delete -f kube-bench/kube-bench.yaml
	kubectl apply -f kube-bench/kube-bench.yaml

retry: build push reapply
