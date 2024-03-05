#!/bin/bash

# Load environment variables from .env file
export $(egrep -v '^#' .prod.containerapp.env | xargs)

env_status=$(az containerapp env show --name actlabs-hub-env-eastus \
  --resource-group actlabs-app \
  --subscription ACT-CSS-Readiness \
  --query properties.provisioningState \
  --output tsv)

if [ "$env_status" != "Succeeded" ]; then
  # Create the environment
  az containerapp env create --name actlabs-hub-env-eastus \
    --resource-group actlabs-app \
    --subscription ACT-CSS-Readiness \
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
  --resource-group actlabs-app \
  --subscription ACT-CSS-Readiness \
  --environment actlabs-hub-env-eastus \
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
