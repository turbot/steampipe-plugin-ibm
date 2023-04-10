# Table: ibm_kms_key_ring

A key ring organizes keys in a specific IBM Cloud location and allows you to manage access control on groups of keys.

## Examples

### Basic info

```sql
select
  title,
  instance_id,
  creation_date
from
  ibm_kms_key_ring
where
  instance_id = '148be70a-ee65-4149-8222-1bf0ff45542f';
```
