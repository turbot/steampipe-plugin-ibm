---
title: "Steampipe Table: ibm_is_instance - Query IBM Cloud Infrastructure Instances using SQL"
description: "Allows users to query IBM Cloud Infrastructure Instances, specifically providing details about the instances such as status, VPC, zone, profile, and resources. It aids in gaining insights into the instance configurations and their current state."
---

# Table: ibm_is_instance - Query IBM Cloud Infrastructure Instances using SQL

IBM Cloud Infrastructure Instances are virtual server instances deployed in IBM Cloud. They are a part of IBM's Infrastructure as a Service (IaaS) offering, providing scalable compute capacity for applications and workloads. These instances can be customized based on the compute power, memory, and storage requirements, and can be managed and accessed over the internet.

## Table Usage Guide

The `ibm_is_instance` table provides insights into instances within IBM Cloud Infrastructure. As a system administrator or a DevOps engineer, explore instance-specific details through this table, including status, VPC, zone, profile, and resources. Utilize it to uncover information about instances, such as their current state, associated resources, and the zones in which they are deployed.

## Examples

### Basic info
Discover the segments that are currently active, along with their unique identifiers and creation dates, to gain insights into your IBM cloud instances. This can help in managing and tracking the status of your instances.

```sql
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_instance;
```

### List instances by name
This query is used to identify specific instances by their name, in this case 'steampipe01'. It's useful for quickly locating specific instances, allowing for efficient management and monitoring of their status and other details.

```sql
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_instance
where
  name = 'steampipe01';
```

### Instance count in each availability zone
Explore which availability zones are hosting the most instances. This can help in understanding the distribution of resources and identifying any potential zones that may be underutilized or overloaded.

```sql
select
  zone ->> 'name' as zone_name,
  count(*)
from
  ibm_is_instance
group by
  zone_name;
```

### Get instance disks attached with instance
Analyze the settings to understand the association between instances and their attached disks, including the size of each disk. This is useful in managing storage resources, ensuring adequate disk space for each instance.

```sql
select
  name as instance_name,
  d ->> 'name' as instance_disk_name,
  d ->> 'size' as disk_size
from
  ibm_is_instance,
  jsonb_array_elements(disks) as d;
```

### Get floating ips associated to the instances
Explore which instances have floating IP addresses associated with them. This is useful for understanding the network configuration and resource allocation within your IBM cloud infrastructure.

```sql
select 
  name,
  fip -> 'target' ->> 'id' as network_interface_id,
  fip ->> 'address' as floating_ip,
  fip ->> 'created_at' as create_time 
from 
  ibm_is_instance,
  jsonb_array_elements(floating_ips) as fip;
```