# Table: ibm_is_instance_disk

An instance storage disk provides fast, affordable, temporary storage that can improve the performance of some cloud native workloads (or apps or services).

## Examples

### Basic info

```sql
select
  name,
  id,
  instance_id,
  created_at
from
  ibm_is_instance_disk;
```

### List large disks (> 100 GB)

```sql
select
  name,
  id,
  instance_id,
  created_at
from
  ibm_is_instance_disk
where
  size > 100;
```

### List unused disks

```sql
select
  d.name as disk_name,
  d.id as disk_id,
  i.name as instance_name,
  d.created_at
from
  ibm_is_instance_disk as d,
  ibm_is_instance as i
where
  d.instance_id = i.id
  and i.status <> 'running';
```
