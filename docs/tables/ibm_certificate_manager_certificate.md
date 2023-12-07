---
title: "Steampipe Table: ibm_certificate_manager_certificate - Query IBM Certificate Manager Certificates using SQL"
description: "Allows users to query IBM Certificate Manager Certificates, specifically the details of each certificate, helping in managing and maintaining the certificates."
---

# Table: ibm_certificate_manager_certificate - Query IBM Certificate Manager Certificates using SQL

IBM Certificate Manager is a cloud service that lets you plan, manage, and deploy digital certificates to enable secure transactions and privacy over the internet. It helps you automate, enforce, and audit certificate lifecycle processes to prevent outages and meet audit and compliance requirements. IBM Certificate Manager provides a secure repository for storing and managing the keys and certificates that are used in cryptographic processes.

## Table Usage Guide

The `ibm_certificate_manager_certificate` table provides insights into the certificates within IBM Certificate Manager. As a Security Engineer, explore certificate-specific details through this table, including certificate status, expiration date, and associated metadata. Utilize it to manage and maintain the certificates, such as those nearing expiration, and to ensure the security of your applications.

## Examples

### Basic info
Analyze the status of various certificates to understand their validity and issuing authority. This could be useful in managing and maintaining the security and integrity of your system.

```sql+postgres
select
  name,
  id,
  status,
  issuer
from
  ibm_certificate_manager_certificate;
```

```sql+sqlite
select
  name,
  id,
  status,
  issuer
from
  ibm_certificate_manager_certificate;
```

### List all imported certificates
Discover the segments that have imported certificates to understand their status and issuer. This can help in maintaining the security and integrity of your system.

```sql+postgres
select
  name,
  id,
  status,
  issuer,
  imported
from
  ibm_certificate_manager_certificate
where
  imported;
```

```sql+sqlite
select
  name,
  id,
  status,
  issuer,
  imported
from
  ibm_certificate_manager_certificate
where
  imported = 1;
```