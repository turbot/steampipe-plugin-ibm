# Table: ibm_iam_access_group_policy

Managing group access by using IAM policies. To assign group's access to resources you must be an administrator on all services in the account, or the assigned administrator for the particular service or service instance.

## Examples

### Basic info

```sql
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_access_group_policy;
```

### List all system created policies

```sql
select
  id,
  type,
  created_by_id,
  href,
  roles
from
  ibm_iam_access_group_policy
where
  created_by_id = 'system';
```