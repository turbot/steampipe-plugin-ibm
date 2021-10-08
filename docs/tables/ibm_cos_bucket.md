# Table: ibm_cos_bucket

An IBM Cloud storage bucket. Cloud Object Storage stores encrypted and dispersed data across multiple geographic locations.

## Examples

### Basic info

```sql
select
  name,
  region,
  creation_date
from
  ibm_cos_bucket;
```

### List unencrypted buckets

```sql
select
  name,
  region,
  creation_date,
  sse_kp_enabled
from
  ibm_cos_bucket
where
  not sse_kp_enabled;
```

### List buckets with versioning disabled

```sql
select
  name,
  region,
  creation_date,
  versioning_enabled
from
  ibm_cos_bucket
where
  not versioning_enabled;
```
