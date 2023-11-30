---
title: "Steampipe Table: ibm_cis_domain - Query IBM Cloud Internet Services Domains using SQL"
description: "Allows users to query IBM Cloud Internet Services Domains, providing insights into domain-specific details and configurations."
---

# Table: ibm_cis_domain - Query IBM Cloud Internet Services Domains using SQL

IBM Cloud Internet Services is a set of edge network services for applications hosted on IBM Cloud. It combines the power of IBM's global network with Cloudflare's edge computing solutions, providing a suite of robust network services including domains. These domains are part of the Cloud Internet Services offering and are used to manage and secure your applications.

## Table Usage Guide

The `ibm_cis_domain` table provides insights into domains within IBM Cloud Internet Services. As a Network Administrator, explore domain-specific details through this table, including domain name, status, and associated metadata. Utilize it to uncover information about domains, such as their current status, the time they were created, and the time they were last modified.

## Examples

### Basic info
Explore which domains are active and their respective security levels by assessing their status and minimum TLS version. This is particularly useful for maintaining security standards and ensuring all domains are operating as expected.

```sql
select
  name,
  id,
  status,
  minimum_tls_version
from
  ibm_cis_domain;
```

### List pending domains
Identify domains that are currently in a pending status to monitor and manage their progress effectively.

```sql
select
  name,
  id,
  status,
  minimum_tls_version
from
  ibm_cis_domain
where
  status = 'pending';
```

### List domains where web_application_firewall not enabled
Identify domains where the web application firewall is not enabled. This is useful for enhancing security by pinpointing potential vulnerabilities in your network infrastructure.

```sql
select
  name,
  id,
  status,
  minimum_tls_version
from
  ibm_cis_domain
where
  web_application_firewall = 'off';
```