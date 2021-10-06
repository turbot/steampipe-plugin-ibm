# Table: ibm_iam_access_group

Access groups can be used to define a set of permissions that you want to grant to a group of users.

## Examples

### Basic info

```sql
select
  name,
  id,
  is_federated,
  href,
  created_by_id,
  created_at
from
  ibm_iam_access_group;
```
