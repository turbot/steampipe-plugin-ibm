---
title: "Steampipe Table: ibm_is_network_acl - Query IBM Cloud Network ACLs using SQL"
description: "Allows users to query Network ACLs in IBM Cloud, specifically providing details about the access control rules for virtual server instances and subnets."
---

# Table: ibm_is_network_acl - Query IBM Cloud Network ACLs using SQL

A Network ACL in IBM Cloud is a set of rules that control the inbound and outbound traffic for your virtual server instances and subnets. These rules act as a firewall at the subnet level, providing a security layer for your resources in the VPC. Network ACLs are stateless, meaning they evaluate each packet in isolation, without considering any related packets or the connection state.

## Table Usage Guide

The `ibm_is_network_acl` table provides insights into Network ACLs within IBM Cloud. As a network administrator or security analyst, explore ACL-specific details through this table, including rule actions, directions, and associated metadata. Utilize it to uncover information about ACLs, such as those with specific access rules, the direction of the rules (inbound or outbound), and the verification of rule priorities.

## Examples

### Basic info
Explore which network access control lists (ACLs) are associated with your IBM Cloud resources. This can help you manage access permissions, ensuring secure and efficient operations across your network.

```sql+postgres
select
  name,
  crn,
  vpc ->> 'name' as vpc_name,
  region,
  account_id
from
  ibm_is_network_acl;
```

```sql+sqlite
select
  name,
  crn,
  json_extract(vpc, '$.name') as vpc_name,
  region,
  account_id
from
  ibm_is_network_acl;
```

### List the default NACL associated with the VPCs
Determine the areas in which the default Network Access Control List (NACL) is associated with Virtual Private Clouds (VPCs). This query is beneficial to understand the security and networking configuration within your cloud environment.

```sql+postgres
select
  acl.name,
  acl.crn,
  vpc.name as vpc_name,
  acl.region,
  acl.account_id
from
  ibm_is_network_acl as acl,
  ibm_is_vpc as vpc
where
  acl.id = vpc.default_network_acl ->> 'id';
```

```sql+sqlite
select
  acl.name,
  acl.crn,
  vpc.name as vpc_name,
  acl.region,
  acl.account_id
from
  ibm_is_network_acl as acl,
  ibm_is_vpc as vpc
where
  acl.id = json_extract(vpc.default_network_acl, '$.id');
```

### Subnet associated with each network ACL
Explore which subnets are associated with each network ACL. This can help in network management by providing insights into the configuration and relationship between subnets and network ACLs.

```sql+postgres
select
  name,
  crn,
  vpc ->> 'name' as vpc_name,
  subnet ->> 'id' as subnet_id,
  subnet ->> 'name' as subnet_name,
  region,
  account_id
from
  ibm_is_network_acl,
  jsonb_array_elements(subnets) as subnet;
```

```sql+sqlite
select
  name,
  crn,
  json_extract(vpc, '$.name') as vpc_name,
  json_extract(subnet.value, '$.id') as subnet_id,
  json_extract(subnet.value, '$.name') as subnet_name,
  region,
  account_id
from
  ibm_is_network_acl,
  json_each(subnets) as subnet;
```