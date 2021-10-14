#!/usr/bin/env bash

ORANGE='\033[0;33m'
NC='\033[0m' # No Color

# Add Bitnami repo for Postgres etc.
helm repo add bitnami https://charts.bitnami.com/bitnami

# Add hashicorp repo for Vault
helm repo add hashicorp https://helm.releases.hashicorp.com

if [[ $(helm list | grep postgres) ]]; then
    echo -e "${ORANGE}Postgres already installed, skipping. To re-install run 'helm uninstall postgres' and then re-run this script${NC}"
else
    # Install postgres using overrides (https://github.com/bitnami/charts/tree/master/bitnami/postgresql/#installing-the-chart)
    helm install postgres bitnami/postgresql --values infra/postgres/override-values.yaml
fi

if [[ $(helm list | grep redis) ]]; then
    echo -e "${ORANGE}Redis already installed, skipping. To re-install run 'helm uninstall redis' and then re-run this script${NC}"
else
    # Install Redis using overrides (https://github.com/bitnami/charts/tree/master/bitnami/redis)
    helm install redis bitnami/redis --values infra/redis/override-values.yaml
fi

if [[ $(helm list | grep kafka) ]]; then
    echo -e "${ORANGE}Kafka already installed, skipping. To re-install run 'helm uninstall kafka' and then re-run this script${NC}"
else
    # Install Kafka using overrides (https://github.com/bitnami/charts/tree/master/bitnami/kafka)
    helm install kafka bitnami/kafka --values infra/kafka/override-values.yaml
fi

if [[ $(helm list | grep vault) ]]; then
    echo -e "${ORANGE}Vault already installed, skipping. To re-install run 'helm uninstall vault' and then re-run this script${NC}"
else
    # Install Vault using overrides (https://github.com/hashicorp/vault-helm)
    helm install vault hashicorp/vault --values infra/vault/override-values.yaml
fi
