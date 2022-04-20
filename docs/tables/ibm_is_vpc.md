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

### List all VPCs with classic access

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

### List address prefix details for VPCs

```sql
select
  name,
  addressp ->> 'cidr' as "cidr",
  addressp -> 'zone' ->> 'name' as "zone",
  addressp ->> 'created_at' as "created_at",
  addressp ->> 'is_default' as "is_default"
from
  ibm_is_vpc,
  jsonb_array_elements(address_prefixes) addressp;
```