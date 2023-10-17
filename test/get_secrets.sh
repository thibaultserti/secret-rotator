#!/usr/bin/bash

# Script used by CI to populate files from github secret

echo "$GOTEST_GCP_SA_KEY_JSON" > test/sa.json


echo "GOOGLE_APPLICATION_CREDENTIALS=sa.json" >> test/.env
echo "GITLAB_TOKEN=$GOTEST_GITLAB_TOKEN" >> test/.env
