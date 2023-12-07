---
title: "Steampipe Table: ibm_iam_user_policy - Query IBM IAM User Policies using SQL"
description: "Allows users to query IBM IAM User Policies, providing insights into the permissions and roles assigned to each user."
---

# Table: ibm_iam_user_policy - Query IBM IAM User Policies using SQL

IBM Identity and Access Management (IAM) is a service within IBM Cloud that allows you to manage access to resources and applications. It provides a centralized way to manage users, roles, and permissions across your IBM Cloud resources. IAM helps you ensure that only authorized users have access to specific resources and can perform specific actions.

## Table Usage Guide

The `ibm_iam_user_policy` table provides insights into user policies within IBM Identity and Access Management (IAM). As a security engineer, explore policy-specific details through this table, including policy roles, resources, and associated metadata. Utilize it to uncover information about policies, such as those with specific permissions, the relationships between users and policies, and the verification of policy roles.

## Examples

### Basic info
Explore which user policies are in effect within your IBM IAM setup. This allows you to identify instances where permissions may be overly broad or insufficient, enhancing overall security and compliance.

```sql+postgres
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_user_policy;
```

```sql+sqlite
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_user_policy;
```

### List all system created policies
Explore which policies have been created automatically by the system. This is useful for understanding system-level permissions and roles.

```sql+postgres
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_user_policy
where
  created_by_id = 'system';
```

```sql+sqlite
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_user_policy
where
  created_by_id = 'system';
```