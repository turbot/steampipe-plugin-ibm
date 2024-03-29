---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/ibm.svg"
brand_color: "#466BB0"
display_name: "IBM Cloud"
short_name: "ibm"
description: "Steampipe plugin to query resources, users and more from IBM Cloud."
og_description: "Query IBM Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/ibm-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# IBM Cloud + Steampipe

[IBM Cloud](https://www.ibm.com/cloud) is a set of cloud computing services for business including IaaS and PaaS.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

List VPCs in your IBM Cloud account:

```sql
select
  id,
  name,
  status,
  region
from
  ibm_is_vpc;
```

```
+-------------------------------------------+------+-----------+----------+
| id                                        | name | status    | region   |
+-------------------------------------------+------+-----------+----------+
| r006-db18c7a8-0ccd-43eb-b9c1-4216c9206201 | prod | available | us-east  |
| r006-c89bee10-d788-4cbe-bf45-dcb940d663a5 | dev  | available | us-south |
+-------------------------------------------+------+-----------+----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/ibm/tables)**

## Get started

### Install

Download and install the latest IBM Cloud plugin:

```bash
steampipe plugin install ibm
```

### Credentials

| Item | Description |
| - | - |
| Credentials | [Create an API key](https://cloud.ibm.com/docs/account?topic=account-userapikey&interface=ui#manage-user-keys) from [IBM Cloud console](https://cloud.ibm.com/iam/apikeys). |
| Permissions | <ol><li>`Viewer` access for `All Account Management Services`</li><li>`Viewer` access for `All Identity and Access enabled services`</li></ol>|
| Radius | Each connection represents a single IBM cloud account. |
| Resolution | 1. `api_key` in steampipe config.<br />2. `IC_API_KEY` environment variable.<br />3. `IBMCLOUD_API_KEY` environment variable. |
| Region Resolution | 1. Regions set for the connection via the regions argument in the config file (~/.steampipe/config/ibm.spc).<br />2. The region specified in the `IC_REGION` or `IBMCLOUD_REGION` environment variable.|

### Configuration

Installing the latest ibm plugin will create a config file (`~/.steampipe/config/ibm.spc`) with a single connection named `ibm`:

```hcl
connection "ibm" {
  plugin  = "ibm"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using:
  # The `IC_REGION` or `IBMCLOUD_REGION` environment variable
  # regions     = ["us-south", "eu-de"]

  # API Key from IBM Cloud
  # api_key = "0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpOSWYJ2rnwi"
}
```

## Multi-Account Connections

You may create multiple IBM connections:

```hcl
connection "ibm_dev" {
  plugin    = "ibm"
  regions   = ["us-south", "eu-de"]
  api_key   = "0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpPOYYJ2rnwi"
}

connection "ibm_qa" {
  plugin    = "ibm"
  regions   = ["us-south", "eu-de"]
  api_key   = "0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpBBTYJ2rnwi"
}

connection "ibm_prod" {
  plugin    = "ibm"
  regions   = ["us-south", "eu-de"]
  api_key   = "0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpWRDYJ2rnwi"
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from ibm_qa.ibm_iam_user
```

You can create multi-account connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection.

```hcl
connection "ibm_all" {
  plugin      = "ibm"
  type        = "aggregator"
  connections = ["ibm_dev", "ibm_qa", "ibm_prod"]
}
```

Querying tables from this connection will return results from the `ibm_dev`, `ibm_qa`, and `ibm_prod` connections:

```sql
select * from ibm_all.ibm_iam_user
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/guides/search-path). It's a good idea to name your aggregator first alphabetically so that it is the first connection in the search path (i.e. `ibm_all` comes before `ibm_dev`):

```sql
select * from ibm_iam_user
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the IBM plugin connections whose names begin with `ibm_`:

```hcl
connection "ibm_all" {
  type        = "aggregator"
  plugin      = "ibm"
  connections = ["ibm_*"]
}
```

## Multi-Region Connections

You may also specify one or more regions with the `regions` argument:

```hcl
connection "ibm" {
  plugin  = "ibm"
  regions = ["au-syd", "eu-de", "eu-gb", "jp-osa", "jp-tok", "us-east", "us-south"]
}
```

The `region` argument supports wildcards:

- All regions

  ```hcl
  connection "ibm" {
    plugin  = "ibm"
    regions = ["*"]
  }
  ```

- All US and EU regions

  ```hcl
  connection "ibm" {
    plugin    = "ibm"
    regions   = ["us-*", "eu-*"]
  }
  ```

IBM multi-region connections are common, but be aware that performance may be impacted by the number of regions and the latency to them.

## Configuring IBM Credentials

### Credentials from Environment Variables

The IBM plugin will use the standard IBM environment variables to obtain credentials **only if other arguments (`api_key`, `regions`) are not specified** in the connection:

```sh
export IC_API_KEY=0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpOSWYJ2rnwi
export IBMCLOUD_API_KEY=0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpOSWYJ2rnwi
```

```hcl
connection "ibm" {
  plugin = "ibm"
}
```
