# Table: ibm_is_subnet

A Subnet is a logical subdivision of an IP network. It enables dividing a network into two or more networks.

## Examples

### Basic info

```sql
select
  id,
  name,
  status,
  ipv4_cidr_block,
  total_ipv4_address_count,
  vpc
from
  ibm_is_subnet;
```

### List all subnets with fewer than 251 available IPv4 addresses

```sql
select
  id,
  name,
  status,
  ipv4_cidr_block,
  available_ipv4_address_count
from
  ibm_is_subnet
where
  available_ipv4_address_count < 251;
```
