#!/bin/bash

IMAGE_NAME=$1
PROJECT_ID=$2
REGION=$3
COMMIT=$4

echo "Deploying ESP Gateway .."
output=$(gcloud run deploy $IMAGE_NAME --image="gcr.io/endpoints-release/endpoints-runtime-serverless:2" --allow-unauthenticated --platform managed --project=$PROJECT_ID --region $REGION)
CLOUD_RUN_HOSTNAME=$(echo $output | tr ' ' '\n' | grep "run.app")
CLOUD_RUN_HOSTNAME=${CLOUD_RUN_HOSTNAME#"https://"}

echo "Preparing OpenAPI Specification .."
scripts/substitute.py openapi-functions.yaml --output openapi-render.yaml --values host=$CLOUD_RUN_HOSTNAME region=$REGION projectid=$PROJECT_ID commit=$COMMIT

echo "Enabling endpoint serving .."
gcloud endpoints services deploy openapi-render.yaml \
    --project $PROJECT_ID
