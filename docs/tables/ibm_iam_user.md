# Table: ibm_iam_user

An IBM Identity and Access Management (IAM) user is an entity that you create in IBM Cloud to represent the person that uses it to interact with.

## Examples

### Basic info

```sql
select
  first_name,
  last_name,
  user_id,
  email,
  account_id
from
  ibm_iam_user;
```

### List inactive users

```sql
select
  first_name,
  last_name,
  user_id,
  email,
  state
from
  ibm_iam_user
where
  state <> 'ACTIVE';
```

### List users with no primary contact phone number

```sql
select
  first_name,
  last_name,
  user_id,
  phonenumber
from
  ibm_iam_user
where
  phonenumber is null;
```
