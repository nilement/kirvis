build:
	GOOS=linux go build -ldflags "-s -w" -o bin/apiserver .
	docker build -t localhost:5000/apiserver:local .

push:
	docker push localhost:5000/apiserver:local

reapply:
	- kubectl delete -f ./k8s/pod.yaml
	kubectl apply -f ./k8s/pod.yaml

retry: build push reapply
