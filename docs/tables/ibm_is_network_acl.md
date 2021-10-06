# Table: ibm_is_network_acl

A Network ACL is an optional layer of security for your VPC that acts as a firewall for controlling traffic in and out of one or more subnets.

## Examples

### Basic info

```sql
select
  name,
  crn,
  vpc ->> 'name' as vpc_name,
  region,
  account_id
from
  ibm_is_network_acl;
```

### List the default NACL associated with the VPCs

```sql
select
  acl.name,
  acl.crn,
  vpc.name as vpc_name,
  acl.region,
  acl.account_id
from
  ibm_is_network_acl as acl,
  ibm_is_vpc as vpc
where
  acl.id = vpc.default_network_acl ->> 'id';
```

### Subnet associated with each network ACL

```sql
select
  name,
  crn,
  vpc ->> 'name' as vpc_name,
  subnet ->> 'id' as subnet_id,
  subnet ->> 'name' as subnet_name,
  region,
  account_id
from
  ibm_is_network_acl,
  jsonb_array_elements(subnets) as subnet;
```
