# Table: ibm_is_region

A region is a specific geographical location where you can deploy apps, services, and other IBM cloud resources. Regions consist of one or more zones, which are physical data centers that house the compute, network, and storage resources, with related cooling and power, for host services and applications.

## Examples

### Basic info

```sql
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region;
```

### List all european regions

```sql
select
  name,
  endpoint,
  status,
  href
from
  ibm_is_region
where
  name like 'eu-%'
```