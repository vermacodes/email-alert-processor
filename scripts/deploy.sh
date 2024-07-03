#!/bin/bash

# Load environment variables from .env file
export $(egrep -v '^#' .prod.containerapp.env | xargs)

env_status=$(az containerapp env show --name email-alert-processor-env \
	--resource-group email-alert-processor \
	--subscription ACT-CSS-Readiness-NPRD \
	--query properties.provisioningState \
	--output tsv)

if [ "$env_status" != "Succeeded" ]; then
	# Create the environment
	az containerapp env create --name email-alert-processor-env \
		--resource-group email-alert-processor \
		--subscription ACT-CSS-Readiness-NPRD \
		--location eastus \
		--logs-destination none

	if [ $? -ne 0 ]; then
		echo "Failed to create environment"
		exit 1
	fi
else
	echo "Environment already exists"
fi

# Check if 'beta' argument is passed
if [ "$1" == "beta" ]; then
	APP_NAME="email-alert-processor-beta"
	IMAGE="ashishvermapu/email-alert-processor:beta"
else
	APP_NAME="email-alert-processor"
	IMAGE="ashishvermapu/email-alert-processor:latest"
fi

# Deploy the Container App
az containerapp create \
	--name $APP_NAME \
	--resource-group email-alert-processor \
	--subscription ACT-CSS-Readiness-NPRD \
	--environment email-alert-processor-env \
	--allow-insecure false \
	--image $IMAGE \
	--ingress 'external' \
	--min-replicas 1 \
	--max-replicas 1 \
	--target-port $PORT \
	--env-vars \
	"PORT=$PORT" \
	"API_KEY=$API_KEY"

if [ $? -ne 0 ]; then
	echo "Failed to create container app"
	exit 1
fi
