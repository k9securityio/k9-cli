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

### List Customers

Whether you need to look up your own k9 Security customer ID or you're managing an inbox for multiple k9 customers, you can use this command to list all the k9 customers you have available in the specified S3 bucket.

```sh
k9 list \
    --bucket <YOUR_SECURE_S3_INBOX> 
```

### List Monitored AWS Accounts

```sh
k9 list \
    --bucket <YOUR_SECURE_S3_INBOX> \
    --customer_id <YOUR_K9_CUSTOMER_ID>
```

### List Reports

```sh
k9 list \
    --bucket <YOUR_SECURE_S3_INBOX> \
    --customer_id <YOUR_K9_CUSTOMER_ID> \
    --account <YOUR_AWS_ACCOUNT_ID>
```

### Download Reports

This version of the `k9` CLI operates on local copies of reports. You can use the following command to download reports for a specific account.

```sh
k9 sync \
    --bucket <YOUR_SECURE_S3_INBOX> \
    --customer_id <YOUR_K9_CUSTOMER_ID> \
    --account <YOUR_AWS_ACCOUNT_ID> \
    -v # verbose output for monitoring progress
```

### Query the IAM Admins

Run the following command to query the set of IAM Admins in a customer account at a point in time.

```sh
k9 query risks iam-admins \
    --customer_id C10001 \
    --account 139710491120 \
    --analysis-date 2022-04-29 \
    --format json
```

### Changes to Principals or Resources Over Time

You can use the `k9` CLI to determine what has changed in an account! Run the following command to generate a diff report between a historical analysis date and the latest report.

```sh
k9 diff principals \
    --customer_id C10001 \
    --account 139710491120 \
    --analysis-date 2022-04-29
```

Or run this one to determine how resources have changed between the two reports.

```sh
k9 diff resources \
    --customer_id C10001 \
    --account 139710491120 \
    --analysis-date 2022-04-29
```

