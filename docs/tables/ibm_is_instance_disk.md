---
title: "Steampipe Table: ibm_is_instance_disk - Query IBM Cloud Infrastructure Instance Disks using SQL"
description: "Allows users to query IBM Cloud Infrastructure Instance Disks, providing detailed information about each instance's disk usage, performance, and configuration."
---

# Table: ibm_is_instance_disk - Query IBM Cloud Infrastructure Instance Disks using SQL

IBM Cloud Infrastructure Instance Disks are a component of IBM's cloud computing platform that offer high-performance, reliable, and scalable block storage for virtual server instances. These disks can be used to store data, applications, and the operating system of the instance. They provide persistent, durable storage and can be customized to meet specific workload requirements.

## Table Usage Guide

The `ibm_is_instance_disk` table allows users to gain insights into the disk usage of each IBM Cloud Infrastructure instance. This table is particularly beneficial for system administrators and cloud engineers, providing detailed information on disk configuration, performance, and capacity. It can be used to monitor disk usage patterns, optimize storage allocation, and troubleshoot potential issues related to disk performance.

## Examples

### Basic info
Discover the segments that highlight the creation time and identifiers of specific instances within IBM's infrastructure service. This is useful for gaining insights into the lifecycle and management of your resources.

```sql
select
  name,
  id,
  instance_id,
  created_at
from
  ibm_is_instance_disk;
```

### List large disks (> 100 GB)
Explore which IBM instance disks are larger than 100GB. This is useful for managing storage resources and identifying potential areas for data optimization.

```sql
select
  name,
  id,
  instance_id,
  created_at
from
  ibm_is_instance_disk
where
  size > 100;
```

### List unused disks
Determine the areas in which unused disks are present by identifying the disks that are not associated with any running instances. This can help in optimizing resources and managing storage efficiently.

```sql
select
  d.name as disk_name,
  d.id as disk_id,
  i.name as instance_name,
  d.created_at
from
  ibm_is_instance_disk as d,
  ibm_is_instance as i
where
  d.instance_id = i.id
  and i.status <> 'running';
```