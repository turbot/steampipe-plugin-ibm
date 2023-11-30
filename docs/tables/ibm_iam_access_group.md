---
title: "Steampipe Table: ibm_iam_access_group - Query IBM Cloud IAM Access Groups using SQL"
description: "Allows users to query IBM Cloud IAM Access Groups, providing insights into group details and associated policies."
---

# Table: ibm_iam_access_group - Query IBM Cloud IAM Access Groups using SQL

IBM Cloud Identity and Access Management (IAM) Access Groups are collections of users and service IDs, where access policies can be applied to grant or restrict access to resources in the IBM Cloud. Access Groups simplify the task of managing access to resources by allowing you to assign policies to groups, rather than individual users or service IDs. This grouping mechanism helps in ensuring proper access management across your IBM Cloud resources.

## Table Usage Guide

The `ibm_iam_access_group` table provides insights into IAM access groups within IBM Cloud Identity and Access Management (IAM). As a security administrator, you can explore group-specific details through this table, including group metadata, associated policies, and access details. Utilize it to uncover information about groups, such as those with specific access policies, the users and service IDs associated with each group, and the verification of access rights.

## Examples

### Basic info
Explore the fundamental details of your IBM IAM access groups to better understand their creation timeline and federation status. This can aid in managing access control and understanding the group's history.

```sql
select
  name,
  id,
  is_federated,
  href,
  created_by_id,
  created_at
from
  ibm_iam_access_group;
```