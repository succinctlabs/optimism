#!/bin/bash

set -u
set -o pipefail

###
### gen_overrides.sh - generate or update the vault_overrides.yaml file
###
### Usage:
###   gen_overrides.sh [options]
###
### Options:
###   -d | --domain-name <name>   The DNS Domain Name of the nodes in the Vault Cluster
###   -r | --region-name <name>   The GCP Region that the resources are operating in
###   -p | --project-name <name>  The GCP Project that the resources are operating in
###   -c | --cluster-name <name>  The GKE Cluster Name
###   -v | --server-version <ver> The Vault Server version to install
###   -h | --help                 Show help / usage
###
###   --ui                       The Vault UI will be enabled (disabled is default)
###   --log-level                The Vault Server log level (info is default)
###   --num-replicas             How many Vault Server replicas should be created (default is 5)
###   --data-vol-size            How big should the data volume be (default is 200Gi)
###   --audit-vol-size           How big should the audit volume be (default is 100Gi)
### 
### Notes:
###   The default for -d/--domain-name is vault-internal.default.svc.cluster.local. You
###   probably only need to change this (to vault-internal) if running on minikube
###
###   If you have GCP_REGION set, that value will be the default for -r/--region-name
###   If you have GCP_PROJECT set, that value will be the default for -p/--project-name
###   If you have GKE_CLUSTER_NAME set, that value will be the default for -c/--cluster-name
###
###   The number of Vault Server replicas should be either 3 or 5
###

DOMAIN="vault-internal.default.svc.cluster.local"
REGION=${GCP_REGION:-}
PROJECT=${GCP_PROJECT:-}
CLUSTER=${GKE_CLUSTER_NAME:-}

VAULT_SERVER_VERSION="1.5.3"
VAULT_UI_ENABLED="false"
VAULT_LOG_LEVEL="info"
VAULT_REPLICAS="5"
VAULT_DATA_SIZE="200Gi"
VAULT_AUDIT_SIZE="100Gi"

# usage displays some helpful information about the script and any errors that need
# to be emitted
usage() {
	MESSAGE=${1:-}

	awk -F'### ' '/^###/ { print $2 }' $0 >&2

	if [[ "${MESSAGE}" != "" ]]; then
		echo "" >&2
		echo "${MESSAGE}" >&2
		echo "" >&2
	fi

	exit -1
}

# validate_config ensures that required variables are set
validate_config() {
	if [[ $(basename ${PWD}) != "infrastructure" ]]; then
		usage "Please execute this script from the \"infrastructure\" directory"
	fi

    if [[ "${DOMAIN}" == "" ]]; then
		usage "DNS Domain (-d) is required"
    fi

    if [[ "${REGION}" == "" ]]; then
		usage "GCP Region (-r) is required"
    fi

    if [[ "${PROJECT}" == "" ]]; then
		usage "GCP Project (-p) is required"
    fi

    if [[ "${CLUSTER}" == "" ]]; then
		usage "GKE Cluster (-c) is required"
    fi

    if [[ "${VAULT_REPLICAS}" != "3" && "${VAULT_REPLICAS}" != "5" ]]; then
		usage "Number of Vault Server replicas (--num-replicas) is invalid"
    fi
}

# gen_overrides updates the vault-overrides.json file
gen_overrides() {
	echo "> Generate Overrides" >&2

  read -r -d '' CONFIG<<EOF
    ui = ${VAULT_UI_ENABLED}
log_level = "${VAULT_LOG_LEVEL}"
cluster_name = "${CLUSTER}"

listener "tcp" {
    tls_disable = {{ .Values.global.tlsDisable }}
    tls_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    tls_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"

    address = "[::]:8200"
    cluster_address = "[::]:8201"
}

seal "gcpckms" {
    region      = "${REGION}"
    project     = "${PROJECT}"
    key_ring    = "omgnetwork-vault-keyring"
    crypto_key  = "omgnetwork-vault-unseal-key"
}

storage "raft" {
    path = "/vault/data"

    retry_join {
    leader_api_addr = "https://vault-0.${DOMAIN}:8200"
    leader_client_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    leader_client_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"
    leader_ca_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/ca.crt"
    }

    retry_join {
    leader_api_addr = "https://vault-1.${DOMAIN}:8200"
    leader_client_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    leader_client_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"
    leader_ca_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/ca.crt"
    }

    retry_join {
    leader_api_addr = "https://vault-2.${DOMAIN}:8200"
    leader_client_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    leader_client_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"
    leader_ca_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/ca.crt"
    }

    retry_join {
    leader_api_addr = "https://vault-3.${DOMAIN}:8200"
    leader_client_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    leader_client_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"
    leader_ca_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/ca.crt"
    }

    retry_join {
    leader_api_addr = "https://vault-4.${DOMAIN}:8200"
    leader_client_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.crt"
    leader_client_key_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/vault.key"
    leader_ca_cert_file = "/vault/userconfig/{{ .Values.global.certSecretName }}/ca.crt"
    }
}
EOF

  cd k8s

	yq w -i vault-overrides.yaml server.image.tag ${VAULT_SERVER_VERSION}
	yq w -i vault-overrides.yaml server.auditStorage.size ${VAULT_AUDIT_SIZE}
	yq w -i vault-overrides.yaml server.dataStorage.size ${VAULT_DATA_SIZE}
	yq w -i vault-overrides.yaml server.extraEnvironmentVars.GOOGLE_REGION ${REGION}
	yq w -i vault-overrides.yaml server.extraEnvironmentVars.GOOGLE_PROJECT ${PROJECT}
	yq w -i vault-overrides.yaml server.ha.raft.config "${CONFIG}"
	yq w -i vault-overrides.yaml server.ha.replicas ${VAULT_REPLICAS}
	yq w -i vault-overrides.yaml server.resources.requests.memory 256Mi
	yq w -i vault-overrides.yaml server.resources.requests.cpu 250m
	yq w -i vault-overrides.yaml server.resources.limits.memory 256Mi
	yq w -i vault-overrides.yaml server.resources.limits.cpu 250m

  cd ..
}

##
## main
##

while [[ $# -gt 0 ]]; do
	case $1 in 
	-d | --domain-name) 
		DOMAIN=$2
		shift
	;;
	-r | --region-name) 
		REGION=$2
		shift
	;;
	-p | --project-name) 
		PROJECT=$2
		shift
	;;
	-v | --server-version) 
		VAULT_SERVER_VERSION=$2
		shift
	;;
	--ui) 
		VAULT_UI_ENABLED=true
	;;
	--log-level) 
		VAULT_LOG_LEVEL=$2
		shift
	;;
	--num-replicas) 
		VAULT_REPLICAS=$2
		shift
	;;
	--audit-vol-size) 
		VAULT_AUDIT_SIZE=$2
		shift
	;;
	-h | --help) 
		usage
	;;
	--)
		shift 
		break
		;;
	-*) usage "Invalid argument: $1" 1>&2 ;;
	*) usage "Invalid argument: $1" 1>&2 ;;
	esac
	shift
done

validate_config
gen_overrides