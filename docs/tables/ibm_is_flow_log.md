---
title: "Steampipe Table: ibm_is_flow_log - Query IBM Cloud VPC Flow Logs using SQL"
description: "Allows users to query IBM Cloud VPC Flow Logs, specifically the network traffic data, providing insights into network traffic patterns and potential anomalies."
---

# Table: ibm_is_flow_log - Query IBM Cloud VPC Flow Logs using SQL

IBM Cloud VPC Flow Logs is a feature that enables you to capture information about the IP traffic going to and from network interfaces in your VPC. This service helps in monitoring and troubleshooting connectivity issues, and it also provides valuable information for security and compliance purposes. Flow Logs data can be used to simplify the diagnostic tasks such as tracking down why certain traffic is reaching an instance.

## Table Usage Guide

The `ibm_is_flow_log` table provides insights into VPC Flow Logs within IBM Cloud. As a network administrator or security analyst, explore log-specific details through this table, including source and destination IP addresses, ports, protocols, and packet and byte counts. Utilize it to uncover information about network traffic, such as identifying patterns, potential security risks, and analyzing the overall network behavior.

## Examples

### Basic info
Gain insights into the basic information of your IBM flow logs, including their names, IDs, lifecycle states, and creation dates. This can be useful for assessing the status and history of your flow logs.

```sql+postgres
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log;
```

```sql+sqlite
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log;
```

### List flow log collectors by name
Explore the specific flow log collectors by name to assess their lifecycle state and creation time, which can help in managing and monitoring your IBM cloud resources. This can be particularly useful when you need to track the changes or status of a specific log collector named 'steampipe01'.

```sql+postgres
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  name = 'steampipe01';
```

```sql+sqlite
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  name = 'steampipe01';
```

### List all inactive flow log collectors
Discover the segments that contain inactive flow log collectors. This can be beneficial in optimizing resources by identifying unused or unnecessary elements within your system.

```sql+postgres
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  not active;
```

```sql+sqlite
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  active = 0;
```

### List all flow log collectors with auto delete disabled
Discover the segments that have auto-delete disabled in flow log collectors, which is essential for maintaining data security and ensuring no essential log data is lost unintentionally.

```sql+postgres
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  not auto_delete;
```

```sql+sqlite
select
  name,
  id,
  crn,
  lifecycle_state,
  created_at
from
  ibm_is_flow_log
where
  auto_delete = 0;
```

### List flow logs with their corresponding VPC details
Explore which flow logs are associated with specific VPCs to better manage network traffic and security in your IBM cloud infrastructure. This can help identify potential issues or bottlenecks in your network configuration.

```sql+postgres
select 
  id, 
  name, 
  vpc ->> 'id' as vpc_id, 
  vpc ->> 'name' as vpc_name 
from 
  ibm_is_flow_log;
```

```sql+sqlite
select 
  id, 
  name, 
  json_extract(vpc, '$.id') as vpc_id, 
  json_extract(vpc, '$.name') as vpc_name 
from 
  ibm_is_flow_log;
```