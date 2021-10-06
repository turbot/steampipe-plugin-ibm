# Table: ibm_is_security_group

An IBM cloud security group is a set of IP filter rules that define how to handle incoming (ingress) and outgoing (egress) traffic to both the public and private interfaces of a virtual server instance.

## Examples

### Basic info

```sql
select
  id,
  name,
  crn,
  rules,
  targets,
  vpc
from
  ibm_is_security_group;
```