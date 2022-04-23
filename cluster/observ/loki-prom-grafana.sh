#!/usr/bin/env bash

kubectl create ns loki
kubectl create ns grafana

helm install loki --namespace=loki --version '2.10.1' grafana/loki
helm install promtail --namespace=loki -f promtail-values.yaml --version '3.11.0' grafana/promtail
helm install grafana --namespace=grafana -f grafana-values.yaml --version '6.24.1' grafana/grafana
