# k9 command line interface
The k9-cli helps you analyze effective AWS IAM access from the command-line using k9.

The CLI synchronizes reports locally and then helps you answer common questions such as:
                                                
* who can administer IAM?
* what principal or resource access has changed?

This is a 'preview' release intended for early adopters.  Only the commands documented in the **Usage** section have been implemented.  The other commands are stubs that show where the CLI is going. 

## Get Started
Download one of the released binaries, rename the file to k9 or k9.exe and place it in your execution path. By default, the k9 CLI will expect the report database to be homed in your current working directory.

Everything is working if you can run the following command and it reports version such as `v0.0.2`.

```sh
k9 version
```

## Usage

Start by `list`ing the k9 customers, AWS accounts, and reports available in your [secure inbox](https://k9security.io/docs/how-k9-works/).  Then `sync` reports to your local directory.  Finally analyze your IAM configuration with the `query` and `diff` commands. 

> **Note**  
> The `list` and `sync` commands require valid AWS credentials, which are resolved using the [standard AWS credential provider chain](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html#credentials-default).  You will also need access to the secure s3 inbox.

### List Customers

Whether you need to look up your own k9 Security customer ID or you're managing an inbox for multiple k9 customers, you can use this command to list all the k9 customers you have available in the specified S3 bucket.

```sh
export K9_SECURE_S3_INBOX=CHANGE-ME-secure-inbox
k9 list \
    --bucket $K9_SECURE_S3_INBOX 
```

This will display a list of k9 customer ids whose reports are in the inbox, e.g.

```text
C123456
```

### List Monitored AWS Accounts

List the AWS accounts that k9 monitors for that customer:
```sh
export K9_CUSTOMER_ID=C123456
k9 list \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID
```

Sample output:

```text
123456789012
012345678901
```

### List Reports

List the reports available for a specific AWS account:

```sh
export K9_ACCOUNT_ID=123456789012
k9 list \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID
```

Sample output:

```text
2022-04-29-0715
2022-04-30-0714
2022-05-01-0713
```

### Download Reports

This version of the `k9` CLI operates on local copies of reports. You can use the following command to download reports for a specific account.

```sh
k9 sync \
    --bucket $K9_SECURE_S3_INBOX \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    -v # verbose output for monitoring progress
```

Sample output:

```text
... snip ...
customers/C10001/reports/aws/123456789012/2021/05/resource-access-summaries.2022-05-01-0714.csv
customers/C10001/reports/aws/123456789012/2022/05/principals.2022-05-01-0714.csv
customers/C10001/reports/aws/123456789012/2022/05/principal-access-summaries.2022-05-01-0714.csv
```

### Query the IAM Admins

Run the following command to query the set of IAM Admins in a customer account at a point in time.

```sh
k9 query risks iam-admins \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date 2022-04-29 \
    --format json # or csv
```

Sample output showing IAM admins, simplified by piping through `jq '.[] | .principal_name'`:

```text
"ci"
"training"
"AccountAdminAccessRole-Sandbox"
"AWS-CodePipeline-Service"
"AWSReservedSSO_AdministratorAccess_437be9d757c9ea2f"
"AWSServiceRoleForOrganizations"
"k9-dev-appeng"
```

### Query Principals at a Point in Time

You can use the `k9` CLI to query the set of principals for an account at a point in time (or from the latest report).

```sh
k9 query principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date 2022-04-29
```

You can use a combination of the `arns` and `names` flags to qualify a list of exact matching records.

```sh
k9 query principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --arns $SOME_ROLE_ARN,$SOME_USER_ARN --arns $ANOTHER_USER_ARN \
    --names $A_ROLE_NAME
```

### Query Resources at a Point in Time

You can use the `k9` CLI to query the set of resources for an account at a point in time (or from the latest report).

```sh
k9 query resources \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date 2022-04-29
```

You can use a combination of the `arns` and `names` flags to qualify a list of exact matching records.

```sh
k9 query resources \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --arns $SOME_BUCKET_ARN,$SOME_MACHINE_ARN --arns $ANOTHER_BUCKET_ARN \
    --names $A_KMS_KEY_NAME
```

### Changes to Principals or Resources Over Time

You can use the `k9` CLI to determine what has changed in an account! Run the following command to generate a diff report between a historical analysis date and the latest report.

```sh
k9 diff principals \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date 2022-04-29
```

Sample output:

```csv
type,principal_arn,before_principal_name,before_principal_type,before_principal_is_iam_admin,before_principal_last_used,before_principal_tag_business_unit,before_principal_tag_environment,before_principal_tag_used_by,before_principal_tags,before_password_last_used,before_password_last_rotated,before_password_state,before_access_key_1_last_used,before_access_key_1_last_rotated,before_access_key_1_state,before_access_key_2_last_used,before_access_key_2_last_rotated,before_access_key_2_state,after_principal_name,after_principal_type,after_principal_is_iam_admin,after_principal_last_used,after_principal_tag_business_unit,after_principal_tag_environment,after_principal_tag_used_by,after_principal_tags,after_password_last_used,after_password_last_rotated,after_password_state,after_access_key_1_last_used,after_access_key_1_last_rotated,after_access_key_1_state,after_access_key_2_last_used,after_access_key_2_last_rotated,after_access_key_2_state
changed,arn:aws:iam::123456789012:user/ci,,,false,2022-04-15 17:51:00+00:00,,,,,,,,2022-04-15 17:51:00+00:00,,,,,,,,false,2022-05-18 16:07:00+00:00,,,,,,,,2022-05-18 16:07:00+00:00,,,,,
changed,arn:aws:iam::123456789012:user/skuenzli,,,false,2022-04-26 22:12:00+00:00,,,,,,,,2022-04-26 22:12:00+00:00,,,,,,,,false,2022-05-23 07:02:00+00:00,,,,,,,,2022-05-23 07:02:00+00:00,,,,,
changed,arn:aws:iam::123456789012:role/k9-auditor,,,false,2022-04-28 21:49:01+00:00,,,,,,,,,,,,,,,,false,2022-05-22 21:35:31+00:00,,,,,,,,,,,,,
changed,arn:aws:iam::123456789012:role/k9-backend-dev,,,false,2022-04-28 23:20:46+00:00,,,,,,,,,,,,,,,,false,2022-05-22 23:20:45+00:00,,,,,,,,,,,,,
```

Or run this one to determine how resources have changed between the two reports.

```sh
k9 diff resources \
    --customer_id $K9_CUSTOMER_ID \
    --account $K9_ACCOUNT_ID \
    --analysis-date 2022-04-29
```

Sample output:

```csv
type,resource_arn,before_resource_name,before_resource_type,before_resource_tag_business_unit,before_resource_tag_environment,before_resource_tag_owner,before_resource_tag_confidentiality,before_resource_tag_integrity,before_resource_tag_availability,before_resource_tags,after_resource_name,after_resource_type,after_resource_tag_business_unit,after_resource_tag_environment,after_resource_tag_owner,after_resource_tag_confidentiality,after_resource_tag_integrity,after_resource_tag_availability,after_resource_tags
added,arn:aws:iam::123456789012:role/cdk-hnb659fds-deploy-role-123456789012-us-east-1,,,,,,,,,,cdk-hnb659fds-deploy-role-123456789012-us-east-1,IAMRole,,,,,,,{}
added,arn:aws:iam::123456789012:role/cdk-hnb659fds-file-publishing-role-123456789012-us-east-1,,,,,,,,,,cdk-hnb659fds-file-publishing-role-123456789012-us-east-1,IAMRole,,,,,,,{}
added,arn:aws:iam::123456789012:role/cdk-hnb659fds-image-publishing-role-123456789012-us-east-1,,,,,,,,,,cdk-hnb659fds-image-publishing-role-123456789012-us-east-1,IAMRole,,,,,,,{}
added,arn:aws:iam::123456789012:role/cdk-hnb659fds-lookup-role-123456789012-us-east-1,,,,,,,,,,cdk-hnb659fds-lookup-role-123456789012-us-east-1,IAMRole,,,,,,,{}
```
