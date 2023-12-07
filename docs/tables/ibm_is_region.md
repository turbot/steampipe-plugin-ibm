---
title: "Steampipe Table: ibm_is_region - Query IBM Cloud Regions using SQL"
description: "Allows users to query IBM Cloud Regions, specifically providing details about the name, status, and geographical location of each region."
---

# Table: ibm_is_region - Query IBM Cloud Regions using SQL

IBM Cloud Regions are geographically dispersed data centers where cloud resources are hosted. Each region is designed to be isolated from the others, ensuring a high level of fault tolerance and stability. Regions are grouped into geographic areas, providing users with the flexibility to host their resources close to their business or customers.

## Table Usage Guide

The `ibm_is_region` table offers insights into the regions within IBM Cloud infrastructure. As a cloud architect or system administrator, you can explore region-specific details through this table, including the name, status, and geographical location of each region. Use it to understand the geographic distribution of your IBM Cloud resources, to plan for disaster recovery, or to make informed decisions when deploying new resources.

## Examples

### Basic info
Explore the status and endpoints of different regions within your IBM cloud infrastructure. This can help you understand the overall health and accessibility of your resources across different geographical locations.

```sql+postgres
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region;
```

```sql+sqlite
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region;
```

### List all European regions
Uncover the details of all IBM cloud regions located in Europe, including their status and endpoint information. This can be useful for managing resources or troubleshooting issues specific to these regions.

```sql+postgres
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region
where
  name like 'eu-%'
```

```sql+sqlite
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region
where
  name like 'eu-%'
```