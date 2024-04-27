#!/bin/bash

if [ -n "$1" ] && [ "$1" != "replace" ]; then
    echo "Error: ./blast/off.sh can be executed only with the \"replace\" option"
    exit 1
fi

if [ "$(realpath "$0" 2>/dev/null)" != "$(realpath ./blast/off.sh 2>/dev/null)" ]; then
    echo "Error: Please execute the script as ./blast/off.sh"
    exit 1
fi

echo -e "\e[32m✨ Starting blast/off.sh script...\e[0m"
PROJECT_ROOT_DIR=$(pwd)
TERRAFORM_DIR=$PROJECT_ROOT_DIR/terraform

if [ ! -f "config.json" ]; then
    echo "Error: config.json not found."
    exit 1
fi

echo -e "\e[33m🔍 Reading configuration from config.json...\e[0m"
GCP_PROJECT_ID=$(jq -r '.GCPProjectId' config.json)
GCP_REGION=$(jq -r '.GCPRegion' config.json)

echo -e "\e[36m🛠️ Creating a Google Artefact Repository...\e[0m"
cd $TERRAFORM_DIR
terraform apply \
    -auto-approve \
    -target=google_artifact_registry_repository.dailydashboardrepository

cd $PROJECT_ROOT_DIR
echo -e "\e[35m🐳 Building Docker image...\e[0m"
TAG=$GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/daily-dashboard-repository/daily-dashboard 
docker build -f Dockerfile-prod -t daily-dashboard . || { echo "Error: Docker build failed."; exit 1; }
docker tag daily-dashboard $TAG || { echo "Error: Docker tag failed."; exit 1; }

echo -e "\e[32m🔐 Configuring docker auth...\e[0m"
gcloud auth configure-docker $GCP_REGION-docker.pkg.dev

echo -e "\e[34m🚀 Pushing Docker image to the repository...\e[0m"
docker push $TAG || { echo "Error: Docker push failed."; exit 1; }

cd $TERRAFORM_DIR
echo -e "\e[36m🚀 Applying infrastructure...\e[0m"
terraform apply -auto-approve

if [ "$1" = "replace" ]; then
    echo -e "\e[36m🚀 Replacing with the new version...\e[0m"
    terraform apply -replace="google_cloud_run_service.dailydashboard"
else
    echo -e "\e[36m🚀 Applying infrastructure...\e[0m"
    terraform apply -auto-approve
fi

echo -e "\e[32m✅ Script execution completed successfully.\e[0m"
cd $PROJECT_ROOT_DIR
