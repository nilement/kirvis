FROM alpine:3.13

COPY ./k8s/experiments.yaml .
COPY bin/node /usr/local/bin/node

RUN mkdir -p ./backups

ENTRYPOINT [ "node" ]
