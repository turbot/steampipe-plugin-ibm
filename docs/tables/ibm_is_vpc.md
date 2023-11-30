---
title: "Steampipe Table: ibm_is_vpc - Query IBM Cloud Virtual Private Clouds using SQL"
description: "Allows users to query Virtual Private Clouds in IBM Cloud, particularly the VPC details, providing insights into network infrastructure and configurations."
---

# Table: ibm_is_vpc - Query IBM Cloud Virtual Private Clouds using SQL

A Virtual Private Cloud (VPC) in IBM Cloud is a secure, isolated virtual network where you can define and control a network space for your own applications and services that run on Virtual Server Instances. VPC provides advanced networking features, including custom subnetting, Network ACLs, and Security Groups. It offers a high degree of control and flexibility over your cloud environment.

## Table Usage Guide

The `ibm_is_vpc` table provides insights into Virtual Private Clouds within IBM Cloud. As a network administrator or cloud engineer, explore VPC-specific details through this table, including network configurations, security settings, and associated metadata. Utilize it to uncover information about VPCs, such as those with specific security settings, the network configurations of each VPC, and the overall structure of your cloud network.

## Examples

### Basic info
Discover the segments that have classic access in your IBM cloud virtual private cloud (VPC) settings to understand potential security implications and to enhance the overall network configuration.

```sql
select
  id,
  name,
  crn,
  classic_access,
  cse_source_ips
from
  ibm_is_vpc;
```

### List all VPCs with classic access
Discover the segments that have classic access within your virtual private cloud (VPC) settings. This allows you to identify potential security risks and manage access controls more effectively.

```sql
select
  id,
  name,
  crn,
  classic_access,
  cse_source_ips
from
  ibm_is_vpc
where
  classic_access;
```

### List address prefix details for VPCs
Analyze the settings to understand the details of address prefixes for Virtual Private Clouds (VPCs). This is useful to manage network configurations and identify default settings.

```sql
select
  name,
  addressp ->> 'cidr' as "cidr",
  addressp -> 'zone' ->> 'name' as "zone",
  addressp ->> 'created_at' as "created_at",
  addressp ->> 'is_default' as "is_default"
from
  ibm_is_vpc,
  jsonb_array_elements(address_prefixes) addressp;
```