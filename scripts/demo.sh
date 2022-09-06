#!/usr/bin/env bash
set -eo pipefail

export ANALYSIS_DATE='2022-09-01'

export K9_SECURE_S3_INBOX=qm-dev-k9-reports
echo "list customers in ${K9_SECURE_S3_INBOX}"
k9 list \
    --bucket $K9_SECURE_S3_INBOX 

export K9_CUSTOMER_ID=C10001
echo "list accounts for customer: ${K9_CUSTOMER_ID}"
k9 list \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID

export K9_ACCOUNT_ID=139710491120
echo "list reports for ${K9_CUSTOMER_ID} account ${K9_ACCOUNT_ID}"
k9 list \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID

echo "synchronizing reports for ${K9_CUSTOMER_ID} account ${K9_ACCOUNT_ID}"
k9 sync \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    -v # verbose output for monitoring progress

echo "Querying Risks: IAM admins"
k9 query risks iam-admins \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --format json | jq '.[] | .principal_name'

# See changes to principals over time
echo "Seeing how principals change over time"
k9 diff principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE
    
# Query specific principals
echo "Query specific principals"
k9 query principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --arns arn:aws:iam::139710491120:role/k9-dev-appeng \
    --names k9-auditor \
    --format json | jq '.'

# Query principal access summaries for specific principals
echo "Query specific principal access"
k9 query principal-access \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --arns arn:aws:iam::139710491120:role/k9-dev-appeng \
    --names k9-auditor \
    --format json | jq '.'
    
# Query specific resources
echo "Query specific resources"
k9 query resources \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --arns arn:aws:rds:us-east-1:139710491120:cluster:int-test-pg-01 \
    --names qm-dev-k9-reports \
    --format json | jq '.'

# Query resource access summaries for specific resources
echo "Query specific resource access"
k9 query resource-access \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --arns arn:aws:rds:us-east-1:139710491120:cluster:int-test-pg-01 \
    --names qm-dev-k9-reports \
    --format json \
      | jq '.'

echo "Query for over-permissioned principals"
k9 query risks over-permissioned-principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --format json \
    --service S3 \
    --max-admin 5 \
      | jq '.[].principal_arn'

echo "Query for over-accessible resources"
k9 query risks over-accessible-resources \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date $ANALYSIS_DATE \
    --format json \
    --service S3 \
    --max-read 3 \
      | jq '.[].resource_arn'