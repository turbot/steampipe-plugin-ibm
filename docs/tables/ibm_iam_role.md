# Table: ibm_iam_role

An IAM role is an Identity and Access Management (IAM) entity with permissions to make IBM cloud service requests.

## Examples

### Basic info

```sql
select
  name,
  id,
  crn,
  actions
from
  ibm_iam_role;
```