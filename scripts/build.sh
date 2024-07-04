#!/bin/bash

# gather input parameters
# -t tag

while getopts ":t:" opt; do
	case $opt in
	t)
		TAG="$OPTARG"
		;;
	\?)
		echo "Invalid option -$OPTARG" >&2
		;;
	esac
done

source .env

if [ -z "${TAG}" ]; then
	TAG="latest"
fi

echo "TAG = ${TAG}"

required_env_vars=("API_KEY" "PORT")

for var in "${required_env_vars[@]}"; do
	if [[ -z "${!var}" ]]; then
		echo "Required environment variable $var is missing"
		exit 1
	fi
done

go build -o email-alert-processor ./cmd/email-alert-processor
if [ $? -ne 0 ]; then
	echo "Failed to build email-alert-processor"
	exit 1
fi

docker build -t actlab.azurecr.io/email-alert-processor:${TAG} .
if [ $? -ne 0 ]; then
	echo "Failed to build docker image"
	exit 1
fi

rm email-alert-processor

az acr login --name actlab --subscription ACT-CSS-Readiness
docker push actlab.azurecr.io/email-alert-processor:${TAG}

docker tag actlab.azurecr.io/email-alert-processor:${TAG} ashishvermapu/email-alert-processor:${TAG}
docker push ashishvermapu/email-alert-processor:${TAG}
