# Table: ibm_is_vpc

A VPC is a public cloud offering that lets an enterprise establish its own private cloud-like computing environment on shared public cloud infrastructure.

## Examples

### Basic info

```sql
select
  id,
  name,
  crn,
  classic_access,
  cse_source_ips
from
  ibm_is_vpc;
```

### List all vpcs having classic access

```sql
select
  id,
  name,
  crn,
  classic_access,
  cse_source_ips
from
  ibm_is_vpc
where
  classic_access;
```