---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/ibm.svg"
brand_color: "#00b050"
display_name: "IBM Cloud"
short_name: "ibm"
description: "Steampipe plugin to query resources, users and more from IBM Cloud."
og_description: "Query IBM Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/ibm-social-graphic.png"
---

# IBM Cloud + Steampipe

[IBM Cloud](https://www.ibm.com/cloud) is a set of cloud computing services for business including IaaS and PaaS.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

List VPCs in your IBM Cloud account:

```sql
select
  id,
  name,
  status,
  region
from
  ibm_is_vpc
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

TODO

### Configuration

Installing the latest ibm plugin will create a config file (`~/.steampipe/config/ibm.spc`) with a single connection named `ibm`:

```hcl
connection "ibm" {
  plugin  = "ibm"
  api_key = "0hrqaLNt-Nc831AW5k7z10CcwOGk_ttqTpOSWYJ2rnwi"
}
```

- `api_key` - API Key from IBM Cloud.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-ibm
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
