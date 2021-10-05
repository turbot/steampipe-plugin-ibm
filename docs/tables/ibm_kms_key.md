# Table: ibm_kms_key

A KMS key is a named object containing one or more key versions, along with metadata for the key. A key exists on exactly one key ring tied to a specific location.

The `ibm_kms_key` table can be used to query keys in an specific service instance, and you must specify which service instance in the where or join clause (`where instance_id=`, `join ibm_kms_key on instance_id=`).

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
  ibm_kms_key
where
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f';
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
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f'
  and creation_date <= (current_date - interval '30' day)
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
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f'
  and key_ring_id = 'steampipe';
```
