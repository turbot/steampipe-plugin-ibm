---
title: "Steampipe Table: ibm_iam_my_api_key - Query IBM IAM API Keys using SQL"
description: "Allows users to query IBM IAM API Keys, specifically the details of the API key of the user, providing insights into the user's account access, permissions, and potential security risks."
---

# Table: ibm_iam_my_api_key - Query IBM IAM API Keys using SQL

IBM Identity and Access Management (IAM) is a service within IBM Cloud that manages access to resources and applications. It provides a centralized way to manage API keys, service IDs, access groups, and policies. IBM IAM helps you control who has access to your IBM cloud resources and services, and what actions they can perform.

## Table Usage Guide

The `ibm_iam_my_api_key` table provides insights into API keys within IBM Identity and Access Management (IAM). As a security or DevOps engineer, explore API key-specific details through this table, including permissions, creation time, and associated metadata. Utilize it to uncover information about API keys, such as those with unrestricted permissions and the verification of access policies.

**Important Notes**
- To query all API keys in an account, use the `ibm_iam_api_key` table.

## Examples

### Basic info
Discover the segments that help you understand the creation and user details of your IBM IAM API keys. This can be useful to track key creation and usage patterns for security and auditing purposes.

```sql
select
  name,
  id,
  crn,
  created_at,
  iam_id as user_iam_id
from
  ibm_iam_my_api_key;
```

### Access key count by user name
Assess the elements within your IBM IAM system to understand the distribution of API keys among users. This can be useful for identifying users with an unusually high number of keys, which could suggest a potential security risk.

```sql
select
  u.user_id,
  count (key.id) as api_key_count
from
  ibm_iam_my_api_key as key,
  ibm_iam_user as u
where
  u.iam_id = key.iam_id
group by
  u.user_id;
```

### List keys older than 90 days
Discover the segments that have API keys older than 90 days to maintain security and ensure timely key rotation. This helps in managing outdated keys which may pose potential security risks.

```sql
select
  key.id as api_key_id,
  key.name as api_key_name,
  u.user_id,
  extract(day from current_timestamp - key.created_at) as age,
  key.account_id
from
  ibm_iam_my_api_key as key,
  ibm_iam_user as u
where
  key.iam_id = u.iam_id
  and extract(day from current_timestamp - key.created_at) > 90;
```