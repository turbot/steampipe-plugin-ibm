## v0.9.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#116](https://github.com/turbot/steampipe-plugin-ibm/pull/116))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#116](https://github.com/turbot/steampipe-plugin-ibm/pull/116))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-ibm/blob/main/docs/LICENSE). ([#116](https://github.com/turbot/steampipe-plugin-ibm/pull/116))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to  column, and fixing connection and potential divide-by-zero bugs. ([#115](https://github.com/turbot/steampipe-plugin-ibm/pull/115))

## v0.8.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#95](https://github.com/turbot/steampipe-plugin-ibm/pull/95))

## v0.8.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#90](https://github.com/turbot/steampipe-plugin-ibm/pull/90))
- Recompiled plugin with Go version `1.21`. ([#90](https://github.com/turbot/steampipe-plugin-ibm/pull/90))

## v0.7.0 [2023-07-17]

_Enhancements_

- Updated the `docs/index.md` file to include multi-account configuration examples. ([#78](https://github.com/turbot/steampipe-plugin-ibm/pull/78))

## v0.6.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#76](https://github.com/turbot/steampipe-plugin-ibm/pull/76))

## v0.5.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#74](https://github.com/turbot/steampipe-plugin-ibm/pull/74))

## v0.4.0 [2023-04-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#71](https://github.com/turbot/steampipe-plugin-ibm/pull/71))

## v0.3.1 [2023-02-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#68](https://github.com/turbot/steampipe-plugin-ibm/pull/68))

## v0.3.0 [2022-09-28]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#63](https://github.com/turbot/steampipe-plugin-ibm/pull/63))
- Recompiled plugin with Go version `1.19`. ([#62](https://github.com/turbot/steampipe-plugin-ibm/pull/62))

## v0.2.0 [2022-07-13]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#60](https://github.com/turbot/steampipe-plugin-ibm/pull/60))

## v0.1.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#10](https://github.com/turbot/steampipe-plugin-ibm/pull/10))

## v0.1.0 [2022-04-28]

_What's new?_

- New tables added
  - [ibm_is_flow_log](https://hub.steampipe.io/plugins/turbot/ibm/tables/ibm_is_flow_log) ([#50](https://github.com/turbot/steampipe-plugin-ibm/pull/50))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#46](https://github.com/turbot/steampipe-plugin-ibm/pull/46))
- Added support for native Linux ARM and Mac M1 builds. ([#52](https://github.com/turbot/steampipe-plugin-ibm/pull/52))
- Added column `floating_ips` to `ibm_is_instance` table. ([#47](https://github.com/turbot/steampipe-plugin-ibm/pull/47))
- Added column `address_prefixes` to `ibm_is_vpc` table. ([#48](https://github.com/turbot/steampipe-plugin-ibm/pull/48))

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
