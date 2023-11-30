---
title: "Steampipe Table: ibm_iam_api_key - Query IBM IAM API Keys using SQL"
description: "Allows users to query IBM IAM API Keys, providing access to details about API keys associated with IBM Cloud Identity and Access Management (IAM)."
---

# Table: ibm_iam_api_key - Query IBM IAM API Keys using SQL

IBM IAM API Keys are a type of credentials in IBM Cloud that clients can use to authenticate with IBM Cloud services. These API keys are associated with IBM Cloud Identity and Access Management (IAM) and can be used to make programmatic calls to the IBM Cloud. They provide a secure way to manage authentication and authorization for IBM Cloud services.

## Table Usage Guide

The `ibm_iam_api_key` table provides insights into API keys within IBM Cloud Identity and Access Management (IAM). As a security analyst, explore API key-specific details through this table, including account IDs, creation timestamps, descriptions, and associated metadata. Utilize it to uncover information about API keys, such as their status, the services they have access to, and the duration for which they are valid.

## Examples

### Basic info
Explore which API keys were created at what time and by which IAM user within IBM's IAM service. This can be particularly useful for auditing purposes or to track key creation in your environment.

```sql
select
  name,
  id,
  crn,
  created_at,
  iam_id as user_iam_id
from
  ibm_iam_api_key;
```

### Access key count by user name
Assess the distribution of access keys across different users to understand their individual API usage. This is useful for auditing purposes and to ensure appropriate access control.

```sql
select
  u.user_id,
  count (key.id) as api_key_count
from
  ibm_iam_api_key as key,
  ibm_iam_user as u
where
  u.iam_id = key.iam_id
group by
  u.user_id;
```

### List keys older than 90 days
Determine the API keys that have been in use for more than 90 days. This query can help identify potentially outdated keys for review, enhancing security and access management.

```sql
select
  key.id as api_key_id,
  key.name as api_key_name,
  u.user_id,
  extract(day from current_timestamp - key.created_at) as age,
  key.account_id
from
  ibm_iam_api_key as key,
  ibm_iam_user as u
where
  key.iam_id = u.iam_id
  and extract(day from current_timestamp - key.created_at) > 90;
```