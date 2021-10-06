# Table: ibm_iam_account_settings

Describes the account settings information of IBM Cloud.

## Examples

### Basic info

```sql
select
  account_id,
  restrict_create_service_id,
  restrict_create_platform_api_key,
  allowed_ip_addresses
from
  ibm_iam_account_settings;
```
