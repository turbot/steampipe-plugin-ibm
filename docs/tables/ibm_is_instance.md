# Table: ibm_is_instance

A virtual server instances for VPC helps to provision instances with high network performance quickly. When you provision an instance, you select a profile that matches the amount of memory and compute power that you need for the application that you plan to run on the instance. Instances are available on the x86 architecture. After you provision an instance, you control and manage those infrastructure resources.

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
  ibm_is_instance;
```

### List instances by name

```sql
select
  name,
  id,
  crn,
  status,
  created_at
from
  ibm_is_instance
where
  name = 'steampipe01';
```

### Instance count in each availability zone

```sql
select
  zone ->> 'name' as zone_name,
  count(*)
from
  ibm_is_instance
group by
  zone_name;
```

### Get instance disks attached with instance

```sql
select
  name as instance_name,
  d ->> 'name' as instance_disk_name,
  d ->> 'size' as disk_size
from
  ibm_is_instance,
  jsonb_array_elements(disks) as d;
```

### Get floating ips associated to the instances

```sql
select 
  name,
  fip -> 'target' ->> 'id' as network_interface_id,
  fip ->> 'address' as floating_ip,
  fip ->> 'created_at' as create_time 
from 
  ibm_is_instance,
  jsonb_array_elements(floating_ips) as fip;
```