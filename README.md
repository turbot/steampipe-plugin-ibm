![image](https://hub.steampipe.io/images/plugins/turbot/ibm-social-graphic.png)

# IBM Cloud Plugin for Steampipe

Use SQL to query instances, domains and more from IBM Cloud.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/ibm)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/ibm/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-ibm/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install ibm
```

Run a query:

```sql
select
  id,
  name,
  status,
  region
from
  ibm_is_vpc;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-ibm.git
cd steampipe-plugin-ibm
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/ibm.spc
```

Try it!

```shell
steampipe query
> .inspect ibm
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-ibm/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [IBM Cloud Plugin](https://github.com/turbot/steampipe-plugin-ibm/labels/help%20wanted)
