---
title: "Steampipe Table: ibm_is_security_group - Query IBM Cloud Infrastructure Security Groups using SQL"
description: "Allows users to query Security Groups in IBM Cloud Infrastructure, providing insights into the configuration and status of security groups within a virtual private cloud (VPC)."
---

# Table: ibm_is_security_group - Query IBM Cloud Infrastructure Security Groups using SQL

IBM Cloud Infrastructure Security Groups are a set of IP filter rules that define how to handle inbound and outbound IP traffic to both the virtual server instance and the network interfaces on the virtual server instance. These groups act as a virtual firewall for your virtual server instances to control inbound and outbound traffic. Security groups in a VPC specify which traffic is allowed to or from resources connected to the VPC.

## Table Usage Guide

The `ibm_is_security_group` table provides insights into Security Groups within IBM Cloud Infrastructure. As a network administrator or security analyst, you can explore security group-specific details through this table, including the attached network interfaces, associated rules, and other metadata. Utilize it to uncover information about security groups, such as their configuration, the resources they are protecting, and the rules that govern their behavior.

## Examples

### Basic info
Explore the basic details of security groups within your IBM cloud infrastructure. This can help you understand the security configurations and rules applied, and identify any potential vulnerabilities or misconfigurations.

```sql
select
  id,
  name,
  crn,
  rules,
  targets,
  vpc
from
  ibm_is_security_group;
```