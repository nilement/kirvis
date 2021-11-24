FROM golang:1.16-buster 

RUN mkdir -p /src/backups
WORKDIR /src

COPY . .
RUN go build -o main

FROM alpine:3.11

COPY --from=0 /src/main .
COPY ./experiments.yaml .

ENTRYPOINT [ "./main" ]