#!/bin/bash

echo "Initializing 'url_shortener' dynamodb table..."

awslocal dynamodb create-table \
  --table-name url_shortener \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \

echo "dynamodb table 'url_shortener' created successfully"
