# Table: ibm_iam_api_key

API keys are long-term credentials for an IAM user or the IBM Cloud account owner. You can use API keys to sign programmatic requests to the IBM CLoud CLI or IBM CLoud API.

To list all of **your** API keys use the `ibm_iam_my_api_key` table instead.

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  created_at,
  iam_id as user_iam_id
from
  ibm_iam_api_key;
```

### Access key count by user name

```sql
select
  u.user_id,
  count (key.id) as api_key_count
from
  ibm_iam_api_key as key,
  ibm_iam_user as u
where
  u.iam_id = key.iam_id
group by
  u.user_id;
```

### List keys older than 90 days

```sql
select
  key.id as api_key_id,
  key.name as api_key_name,
  u.user_id,
  extract(day from current_timestamp - key.created_at) as age,
  key.account_id
from
  ibm_iam_api_key as key,
  ibm_iam_user as u
where
  key.iam_id = u.iam_id
  and extract(day from current_timestamp - key.created_at) > 90;
```
