---
title: "Steampipe Table: ibm_resource_group - Query IBM Resource Group using SQL"
description: "Allows users to query IBM Resource Groups, providing details about the group's ID, name, state, and more."
---

# Table: ibm_resource_group - Query IBM Resource Group using SQL

An IBM Resource Group is a way to manage and organize resources in an IBM Cloud account. It allows users to manage access to resources, assign resources to different teams, and track costs. Resource Groups are similar to tags, but they provide a higher level of organization and control.

## Table Usage Guide

The `ibm_resource_group` table provides insights into Resource Groups within IBM Cloud. As a cloud administrator, you can explore group-specific details through this table, including the group's ID, name, state, and more. Utilize it to manage and organize your resources effectively, assign resources to different teams, and track costs.

## Examples

### Basic info
Explore the status and creation date of various resources within your IBM account. This can be useful for understanding the distribution and organization of resources, as well as identifying any potential issues or anomalies.

```sql
select
  name,
  id,
  crn,
  state,
  created_at,
  account_id
from
  ibm_resource_group;
```

### List default resource groups
Explore which resource groups have been set as default in your IBM cloud setup. This can help streamline resource management and optimize cloud operations.

```sql
select
  name,
  id,
  crn,
  state,
  created_at
from
  ibm_resource_group
where
  is_default;
```

### List resource groups by name
This query helps you pinpoint specific resource groups in your IBM account by their name. This can be particularly useful in managing resources and understanding their allocation within your infrastructure.

```sql
select
  name,
  id,
  crn,
  state,
  created_at,
  account_id
from
  ibm_resource_group
where
  name = 'Default';
```