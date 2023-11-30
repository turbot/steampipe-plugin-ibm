---
title: "Steampipe Table: ibm_kms_key_ring - Query IBM Key Protect Key Rings using SQL"
description: "Allows users to query Key Rings in IBM Key Protect, specifically the key ring ID, name, and creation date, providing insights into key management and security."
---

# Table: ibm_kms_key_ring - Query IBM Key Protect Key Rings using SQL

IBM Key Protect is a cloud-based service designed to manage and protect cryptographic keys used in IBM cloud services. It provides a secure way to generate, manage, and destroy encryption keys, which are used to protect data-at-rest in IBM cloud services. Key Rings in IBM Key Protect are used to organize keys and control access to them.

## Table Usage Guide

The `ibm_kms_key_ring` table provides insights into Key Rings within IBM Key Protect. As a security engineer, explore key ring-specific details through this table, including key ring ID, name, and creation date. Utilize it to manage and monitor your cryptographic keys, ensuring secure access to your IBM cloud services.

## Examples

### Basic info
Explore the creation dates and titles of key rings within a specific instance in IBM's Key Management Service. This can be beneficial in understanding the organization and timeline of your security infrastructure.

```sql
select
  title,
  instance_id,
  creation_date
from
  ibm_kms_key_ring
where
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f';
```