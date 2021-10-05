# Table: ibm_kms_key_ring

A key ring organizes keys in a specific IBM Cloud location and allows you to manage access control on groups of keys.

The `ibm_kms_key_ring` table can be used to query key rings in an specific service instance, and you must specify which service instance in the where or join clause (`where instance_id=`, `join ibm_kms_key_ring on instance_id=`).

## Examples

### Basic info

```sql
select
  name,
  instance_id,
  creation_date
from
  ibm_kms_key_ring
where
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f';
```
