---
title: "Steampipe Table: ibm_iam_user - Query IBM Cloud IAM Users using SQL"
description: "Allows users to query IBM Cloud IAM Users, providing details such as user ID, email, account ID, and created at timestamp."
---

# Table: ibm_iam_user - Query IBM Cloud IAM Users using SQL

IBM Cloud Identity and Access Management (IAM) is a service that helps secure access to IBM Cloud resources. It enables the management of identities and access, allowing users to control who has access to their IBM Cloud resources and what actions they can perform. With IAM, you can manage access to your resources by creating policies and assigning them to IAM identities (users, groups, and service IDs).

## Table Usage Guide

The `ibm_iam_user` table provides insights into users within IBM Cloud Identity and Access Management (IAM). As a security officer or administrator, you can explore user-specific details through this table, including user ID, email, account ID, and created at timestamp. Utilize it to uncover information about users, such as their access levels, assigned roles, and other related metadata.

## Examples

### Basic info
Explore the basic user information from an IBM IAM user list to gain insights into user details and their associated account IDs. This can be useful for user management and account audits.

```sql
select
  first_name,
  last_name,
  user_id,
  email,
  account_id
from
  ibm_iam_user;
```

### List inactive users
This query helps identify users who are not currently active within the IBM IAM system. It is useful in auditing user activity and assessing the need for potential clean-up of inactive accounts.

```sql
select
  first_name,
  last_name,
  user_id,
  email,
  state
from
  ibm_iam_user
where
  state <> 'ACTIVE';
```

### List users with no primary contact phone number
Discover the users who lack a primary contact phone number, allowing you to identify gaps in your contact information and reach out for updates. This can be particularly useful in maintaining effective communication channels with all users.

```sql
select
  first_name,
  last_name,
  user_id,
  phonenumber
from
  ibm_iam_user
where
  phonenumber is null;
```