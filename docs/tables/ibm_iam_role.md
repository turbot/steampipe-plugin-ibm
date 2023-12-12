---
title: "Steampipe Table: ibm_iam_role - Query IBM IAM Roles using SQL"
description: "Allows users to query IBM IAM Roles, providing insights into the roles, permissions, trust policies and associated metadata in IBM Cloud."
---

# Table: ibm_iam_role - Query IBM IAM Roles using SQL

IBM Identity and Access Management (IAM) is a service that helps in managing access to IBM Cloud services. It allows you to create and manage identities, and set policies so that you can control who has access to what. IAM Roles are an IAM identity that you can create and use to delegate permissions to AWS service that needs to interact with your resources.

## Table Usage Guide

The `ibm_iam_role` table provides insights into IAM roles within IBM Identity and Access Management (IAM). As a DevOps engineer, explore role-specific details through this table, including permissions, trust policies, and associated metadata. Utilize it to uncover information about roles, such as those with wildcard permissions, the trust relationships between roles, and the verification of trust policies.

## Examples

### Basic info
Explore which IBM IAM roles are currently in use to understand their actions and assess their elements within your system. This can help pinpoint specific areas for security improvement or optimization.

```sql+postgres
select
  name,
  id,
  crn,
  actions
from
  ibm_iam_role;
```

```sql+sqlite
select
  name,
  id,
  crn,
  actions
from
  ibm_iam_role;
```