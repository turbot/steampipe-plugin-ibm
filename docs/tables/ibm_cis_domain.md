# Table: ibm_cis_domain

IBM Cloud Internet Services (CIS), powered by Cloudflare, provides a fast, highly performant, reliable, and secure internet service for customers running their business on IBM Cloud.
CIS gets you going quickly by establishing defaults for you, which you can change by using the API or UI.

## Examples

### Basic info

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