---
title: "Steampipe Table: ibm_account - Query IBM Cloud Accounts using SQL"
description: "Allows users to query IBM Cloud Accounts, specifically the details of each account, providing insights into account configurations and settings."
---

# Table: ibm_account - Query IBM Cloud Accounts using SQL

IBM Cloud Accounts are the foundational entities in IBM Cloud that provide the access to services, resources, and applications within IBM Cloud. They are used to organize resources, manage permissions, and control billing. IBM Cloud Accounts enable users to manage their cloud resources in a secure and efficient manner.

## Table Usage Guide

The `ibm_account` table provides insights into IBM Cloud Accounts. As a cloud administrator or DevOps engineer, explore account-specific details through this table, including account status, owner's identity, and associated metadata. Utilize it to uncover information about accounts, such as their creation time, resource group ID, and the state of the account.

## Examples

### Basic info
Explore the basic details of your IBM account such as name, status, and owner ID. This can be useful in understanding the current state of your account, and identifying the appropriate owner.

```sql+postgres
select
  name,
  guid as id,
  state,
  owner_user_id
from
  ibm_account;
```

```sql+sqlite
select
  name,
  guid as id,
  state,
  owner_user_id
from
  ibm_account;
```

### Get details about account owner
Explore which IBM accounts are linked with their respective owners, allowing you to identify instances where account ownership needs to be updated or verified. This provides useful insights into account management and ensures the correct assignment of resources.

```sql+postgres
select
  acc.name,
  acc.guid as id,
  acc.state,
  u.first_name || ' ' || u.last_name as owner_full_name
from
  ibm_account as acc,
  ibm_iam_user as u
where
  acc.owner_user_id = u.user_id;
```

```sql+sqlite
select
  acc.name,
  acc.guid as id,
  acc.state,
  u.first_name || ' ' || u.last_name as owner_full_name
from
  ibm_account as acc,
  ibm_iam_user as u
where
  acc.owner_user_id = u.user_id;
```