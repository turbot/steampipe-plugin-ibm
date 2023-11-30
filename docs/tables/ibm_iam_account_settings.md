---
title: "Steampipe Table: ibm_iam_account_settings - Query IBM IAM Account Settings using SQL"
description: "Allows users to query IBM IAM Account Settings, providing information about the configuration of their IBM Cloud account."
---

# Table: ibm_iam_account_settings - Query IBM IAM Account Settings using SQL

IBM Identity and Access Management (IAM) is a service that helps users to securely manage their IBM Cloud resources. It offers features to control who can access the resources, which actions they can perform, and how they manage these permissions. IBM IAM Account Settings are a part of this service, containing the configuration details of the user's IBM Cloud account.

## Table Usage Guide

The `ibm_iam_account_settings` table provides insights into the IBM IAM Account Settings within IBM Identity and Access Management (IAM). As a cloud administrator, you can explore account-specific details through this table, including the account's API key, session duration, and associated metadata. It can be utilized to uncover information about the account's settings, such as the account's MFA status, the password settings, and the account's access groups.

## Examples

### Basic info
Explore which IBM account settings have restrictions on creating service IDs and platform API keys. This can be useful in maintaining security by understanding where creation of these elements is limited and which IP addresses are permitted.

```sql
select
  account_id,
  restrict_create_service_id,
  restrict_create_platform_api_key,
  allowed_ip_addresses
from
  ibm_iam_account_settings;
```