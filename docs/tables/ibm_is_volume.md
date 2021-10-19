# Table: ibm_is_volume

Block Storage for VPC provides hypervisor-mounted, high-performance data storage for your virtual server instances that you can provision within an IBM Cloud™ Virtual Private Cloud (VPC). The VPC infrastructure provides rapid scaling across zones and extra performance and security.


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

### Instance count in each availability zone

```sql
select
  zone ->> 'name' as zone_name,
  count(*)
from
  ibm_is_volume
group by
  zone_name;
```
