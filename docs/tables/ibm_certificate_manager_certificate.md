# Table: ibm_certificate_manager_certificate

A Certificate Manager helps you to obtain, store and manage SSL/TLS certificates that you use for IBM Cloud deployments, or other cloud and on-premises deployments.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  issuer
from
  ibm_certificate_manager_certificate;
```

### List all imported certificates

```sql
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
