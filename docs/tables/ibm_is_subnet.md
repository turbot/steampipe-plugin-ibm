---
title: "Steampipe Table: ibm_is_subnet - Query IBM Cloud Subnets using SQL"
description: "Allows users to query IBM Cloud Subnets, providing detailed information about each subnet in an IBM Cloud VPC."
---

# Table: ibm_is_subnet - Query IBM Cloud Subnets using SQL

IBM Cloud Subnets are a range of IP addresses in a VPC that are divided into smaller blocks, each with its own set of security and network connection settings. They are used to segment the network, improve security and efficiency, and help organizations manage resources in a more granular way. Subnets are a critical component of the IBM Cloud architecture, providing the foundation for deploying resources in the VPC.

## Table Usage Guide

The `ibm_is_subnet` table provides insights into subnets within IBM Cloud VPC. As a network administrator, you can explore subnet-specific details through this table, including IP ranges, associated network ACLs, and public gateway information. Utilize it to uncover information about subnets, such as those with specific security settings, the connectivity options of each subnet, and the management of network resources.

## Examples

### Basic info
Explore the status and details of your network subnets, such as the range of IP addresses and associated virtual private cloud. This can be useful for managing and optimizing your network resources.

```sql
select
  id,
  name,
  status,
  ipv4_cidr_block,
  total_ipv4_address_count,
  vpc
from
  ibm_is_subnet;
```

### List all subnets with fewer than 251 available IPv4 addresses
Gain insights into subnets that are nearing their capacity by identifying those with fewer than 251 available IPv4 addresses. This can aid in planning and managing your network resources effectively.

```sql
select
  id,
  name,
  status,
  ipv4_cidr_block,
  available_ipv4_address_count
from
  ibm_is_subnet
where
  available_ipv4_address_count < 251;
```