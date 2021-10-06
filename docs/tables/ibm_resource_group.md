# Table: ibm_resource_group

A resource group is a way for you to organize your account resources in customizable groupings so that you can quickly assign users access to more than one resource at a time. Any account resource that is managed by using IBM Cloud Identity and Access Management (IAM) access control belongs to a resource group within your account. Cloud Foundry services are assigned to orgs and spaces and can't be added to a resource group.

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  state,
  created_at,
  account_id
from
  ibm_resource_group;
```

### List default resource groups

```sql
select
  name,
  id,
  crn,
  state,
  created_at
from
  ibm_resource_group
where
  is_default;
```

### List resource groups by name

```sql
select
  name,
  id,
  crn,
  state,
  created_at,
  account_id
from
  ibm_resource_group
where
  name = 'Default';
```
