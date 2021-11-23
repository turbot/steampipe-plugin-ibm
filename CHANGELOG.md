## v0.0.3 [2021-11-23]

_Enhancements_

- Recompiled plugin with Go version 1.17 ([#39](https://github.com/turbot/steampipe-plugin-ibm/pull/39))
- Recompile plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#38](https://github.com/turbot/steampipe-plugin-ibm/pull/38))

_Bug fixes_

- Fixed the typo in the `ibm_is_volume` table docs ([#36](https://github.com/turbot/steampipe-plugin-ibm/pull/36))

## v0.0.2 [2021-10-21]

_What's new?_

- New tables added
  - [ibm_cis_domain](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_cis_domain) ([#34](https://github.com/turbot/steampipe-plugin-ibm/pull/34))
  - [ibm_cos_bucket](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_cos_bucket) ([#13](https://github.com/turbot/steampipe-plugin-ibm/pull/13))
  - [ibm_iam_access_group_policy](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_access_group_policy) ([#19](https://github.com/turbot/steampipe-plugin-ibm/pull/19))
  - [ibm_is_volume](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_volume`) ([#24](https://github.com/turbot/steampipe-plugin-ibm/pull/24))

_Enhancements_

- Updated: Add column `rotation_policy` to the `ibm_kms_key` table ([#31](https://github.com/turbot/steampipe-plugin-ibm/pull/31))
- Updated: Add columns `order_policy_name` and `auto_renew_enabled` to the `ibm_certificate_manager_certificate` table ([#27](https://github.com/turbot/steampipe-plugin-ibm/pull/27))
- Updated: Add column `iam_id` to the `ibm_iam_user_policy` table ([#17](https://github.com/turbot/steampipe-plugin-ibm/pull/17))

## v0.0.1 [2021-10-06]

_What's new?_

- New tables added
  - [ibm_account](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_account)
  - [ibm_certificate_manager_certificate](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_certificate_manager_certificate)
  - [ibm_iam_access_group](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_access_group)
  - [ibm_iam_account_settings](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_account_settings)
  - [ibm_iam_api_key](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_api_key)
  - [ibm_iam_my_api_key](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_my_api_key)
  - [ibm_iam_role](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_role)
  - [ibm_iam_user](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_user)
  - [ibm_iam_user_policy](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_iam_user_policy)
  - [ibm_is_instance](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_instance)
  - [ibm_is_instance_disk](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_instance_disk)
  - [ibm_is_network_acl](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_network_acl)
  - [ibm_is_region](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_region)
  - [ibm_is_security_group](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_security_group)
  - [ibm_is_subnet](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_subnet)
  - [ibm_is_vpc](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_vpc)
  - [ibm_is_vpc](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_vpc)
  - [ibm_kms_key](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_kms_key)
  - [ibm_kms_key_ring](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_kms_key_ring)
  - [ibm_resource_group](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_resource_group)
