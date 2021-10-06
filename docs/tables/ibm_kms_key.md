# Table: ibm_kms_key

A KMS key is a named object containing one or more key versions, along with metadata for the key. A key exists on exactly one key ring tied to a specific location.

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date
from
  ibm_kms_key;
```

### List keys older than 30 days

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date
from
  ibm_kms_key
where
  creation_date <= (current_date - interval '30' day)
order by
  creation_date;
```

### List keys by key ring

```sql
select
  name,
  id,
  crn,
  instance_id,
  state,
  creation_date,
  key_ring_id
from
  ibm_kms_key
where
  key_ring_id = 'steampipe';
```
