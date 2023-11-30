---
title: "Steampipe Table: ibm_kms_key - Query IBM Key Protect Keys using SQL"
description: "Allows users to query IBM Key Protect Keys, specifically the key ID, key name, and key creation date, providing insights into key management and security."
---

# Table: ibm_kms_key - Query IBM Key Protect Keys using SQL

IBM Key Protect is a cloud-based security service that provides lifecycle management for encryption keys that are used in IBM Cloud services or customer-built applications. The service provides a simple and scalable way to manage keys, including creating, importing, storing, and disposing of them. IBM Key Protect helps to facilitate secure cloud data protection and key management at scale.

## Table Usage Guide

The `ibm_kms_key` table provides insights into keys within IBM Key Protect. As a security or DevOps engineer, explore key-specific details through this table, including key ID, key name, and key creation date. Utilize it to uncover information about keys, such as their lifecycle status, the associated instances, and the verification of key policies.

## Examples

### Basic info
Analyze the settings to understand the status and creation date of IBM Key Management Service keys, which can be useful in managing and auditing key usage across your IBM Cloud services.

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date
from
  ibm_kms_key;
```

### List keys older than 30 days
Explore which encryption keys have been in existence for over a month. This can be useful for managing and auditing key lifecycles, ensuring old and potentially vulnerable keys are replaced or retired.

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date
from
  ibm_kms_key
where
  creation_date <= (current_date - interval '30' day)
order by
  creation_date;
```

### List keys by key ring
Determine the areas in which specific keys are associated with a given key ring. This can help in managing and organizing your encryption keys, enhancing your security strategy.

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date,
  key_ring_id
from
  ibm_kms_key
where
  key_ring_id = 'steampipe';
```