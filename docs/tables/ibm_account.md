# Table: aws_iam_user

An IBM account is a container for your IBM Cloud resources. You create and manage your IBM Cloud resources in an IBM account, and the IBM account provides administrative capabilities for access and billing.

## Examples

### Basic info

```sql
select
  name,
  guid as id,
  state,
  owner_user_id
from
  ibm_account;
```

### Get details about account owner

```sql
select
  acc.name,
  acc.guid as id,
  acc.state,
  u.first_name || ' ' || u.last_name as owner_full_name
from
  ibm_account as acc,
  ibm_iam_user as u
where
  acc.owner_user_id = u.user_id;
```
