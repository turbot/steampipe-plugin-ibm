# Table: ibm_is_volume

Block Storage for VPC provides hypervisor-mounted, high-performance data storage for your virtual server instances that you can provision within an IBM Cloud Virtual Private Cloud (VPC). The VPC infrastructure provides rapid scaling across zones and extra performance and security.
Block Storage for VPC offers block-level volumes that are attached to an instance as a boot volume when the instance is created or attached as secondary data volumes.

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume;
```

### List volumes by name

```sql
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_volume
where
  name = 'steampipe01';
```

### List of volumes with size more than 100GB

```sql
select
  name,
  id,
  crn,
  capacity
from
  ibm_is_volume
where
  capacity > 100;
```

### List volumes not encrypted using user-managed key

```sql
select
  name,
  id,
  crn,
  encryption,
  encryption_key
from
  ibm_is_volume
where
  encryption <> 'user_managed';
```

### Volume count in each availability zone

```sql
select
  zone ->> 'name' as zone_name,
  count(*)
from
  ibm_is_volume
group by
  zone_name;
```
