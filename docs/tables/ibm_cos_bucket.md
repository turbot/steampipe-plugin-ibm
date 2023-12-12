---
title: "Steampipe Table: ibm_cos_bucket - Query IBM Cloud Object Storage Buckets using SQL"
description: "Allows users to query IBM Cloud Object Storage Buckets, specifically retrieving details about the bucket's configuration, access policies, and geographical location."
---

# Table: ibm_cos_bucket - Query IBM Cloud Object Storage Buckets using SQL

IBM Cloud Object Storage is a highly scalable cloud storage service, designed for high durability, resiliency and security. It allows users to store, manage and access their data in a simple, cost-effective way. It is ideal for storing large amounts of data, such as images, video files, and backups.

## Table Usage Guide

The `ibm_cos_bucket` table provides insights into the configuration and access policies of IBM Cloud Object Storage Buckets. As a data engineer, explore bucket-specific details through this table, including bucket location, storage class, and access policies. Utilize it to monitor and manage your data storage, ensuring optimal performance and security.

## Examples

### Basic info
Explore which IBM COS buckets are in use, their respective regions and when they were created. This can be beneficial for understanding the distribution and timeline of your storage resources.

```sql+postgres
select
  name,
  region,
  creation_date
from
  ibm_cos_bucket;
```

```sql+sqlite
select
  name,
  region,
  creation_date
from
  ibm_cos_bucket;
```

### List unencrypted buckets
Determine the areas in which data stored in IBM Cloud Object Storage buckets are potentially at risk due to lack of encryption. This allows for a quick assessment of security vulnerabilities and aids in prioritizing necessary protective measures.

```sql+postgres
select
  name,
  region,
  creation_date,
  sse_kp_enabled
from
  ibm_cos_bucket
where
  not sse_kp_enabled;
```

```sql+sqlite
select
  name,
  region,
  creation_date,
  sse_kp_enabled
from
  ibm_cos_bucket
where
  not sse_kp_enabled;
```

### List buckets with versioning disabled
Explore which IBM COS buckets have versioning disabled to gain insights into potential data loss risks. This is useful in scenarios where maintaining different versions of objects in a bucket is crucial for data recovery and backup purposes.

```sql+postgres
select
  name,
  region,
  creation_date,
  versioning_enabled
from
  ibm_cos_bucket
where
  not versioning_enabled;
```

```sql+sqlite
select
  name,
  region,
  creation_date,
  versioning_enabled
from
  ibm_cos_bucket
where
  not versioning_enabled;
```