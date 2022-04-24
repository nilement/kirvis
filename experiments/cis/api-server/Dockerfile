FROM alpine:3.13

COPY ./k8s/experiments.yaml .
COPY bin/apiserver /usr/local/bin/apiserver

RUN mkdir -p ./backups

ENTRYPOINT [ "apiserver" ]
