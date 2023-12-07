---
title: "Steampipe Table: ibm_is_volume - Query IBM Cloud Infrastructure Volumes using SQL"
description: "Allows users to query IBM Cloud Infrastructure Volumes, providing insights into volume details including size, profile, and status."
---

# Table: ibm_is_volume - Query IBM Cloud Infrastructure Volumes using SQL

IBM Cloud Infrastructure Volumes is a block storage service that provides scalable and reliable storage for your virtual server instances. It offers a variety of volume types that are optimized for different types of workloads, including transactional and big data workloads. IBM Cloud Infrastructure Volumes ensures data availability and durability by automatically replicating data across multiple physical drives.

## Table Usage Guide

The `ibm_is_volume` table provides insights into volumes within IBM Cloud Infrastructure. As a system administrator, explore volume-specific details through this table, including size, profile, and status. Utilize it to uncover information about volumes, such as those with specific profiles, the capacity of the volumes, and the status of the volumes.

## Examples

### Basic info
Explore the status and creation date of various resources in your IBM Cloud infrastructure. This allows you to monitor the health and lifecycle of your resources effectively.

```sql+postgres
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume;
```

```sql+sqlite
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume;
```

### List volumes by name
Determine the status and creation date of a specific volume in an IBM cloud infrastructure. This can be useful for tracking the lifecycle of your resources and managing their use effectively.

```sql+postgres
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume
where
  name = 'steampipe01';
```

```sql+sqlite
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume
where
  name = 'steampipe01';
```

### List of volumes with size more than 100GB
Analyze the settings to understand which storage volumes exceed a capacity of 100GB. This can be useful for managing storage resources and identifying potential areas for data optimization.

```sql+postgres
select
  name,
  id,
  crn,
  capacity
from
  ibm_is_volume
where
  capacity > 100;
```

```sql+sqlite
select
  name,
  id,
  crn,
  capacity
from
  ibm_is_volume
where
  capacity > 100;
```

### List volumes not encrypted using user-managed key
Explore which storage volumes are not using user-managed encryption. This can help assess the security measures in place and identify potential areas for improvement.

```sql+postgres
select
  name,
  id,
  crn,
  encryption,
  encryption_key
from
  ibm_is_volume
where
  encryption <> 'user_managed';
```

```sql+sqlite
select
  name,
  id,
  crn,
  encryption,
  encryption_key
from
  ibm_is_volume
where
  encryption <> 'user_managed';
```

### Volume count in each availability zone
Explore which availability zones have the most volumes to better manage and distribute your resources. This could be particularly useful in balancing workloads and optimizing performance across different zones.

```sql+postgres
select
  zone ->> 'name' as zone_name,
  count(*)
from
  ibm_is_volume
group by
  zone_name;
```

```sql+sqlite
select
  json_extract(zone, '$.name') as zone_name,
  count(*)
from
  ibm_is_volume
group by
  zone_name;
```