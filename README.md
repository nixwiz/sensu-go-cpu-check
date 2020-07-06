[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-go-cpu-check)
![Go Test](https://github.com/nixwiz/sensu-go-cpu-check/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-go-cpu-check/workflows/goreleaser/badge.svg)

# sensu-go-cpu-check

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The sensu-go-cpu-check is a [Sensu Check][1] that provides a breakdown of CPU usage on
a Linux instance.

## Usage examples

```
Sensu Go CPU Check

Usage:
  sensu-go-cpu-check [flags]
  sensu-go-cpu-check [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -i, --interval int   How long to sleep between CPU usage samples, in seconds (default 1)
  -c, --critical int   Critical threshold (>=) for CPU usage (default 90)
  -w, --warning int    Warning threshold (>=) for CPU usage (default 80)
  -h, --help           help for sensu-go-cpu-check

Use "sensu-go-cpu-check [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][2] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add nixwiz/sensu-go-cpu-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][3].


### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-go-cpu-check
  namespace: default
spec:
  command: sensu-go-cpu-check -w 80 -c 95
  subscriptions:
  - system
  runtime_assets:
  - nixwiz/sensu-go-cpu-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-go-cpu-check repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[2]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[3]: https://bonsai.sensu.io/assets/nixwiz/sensu-go-cpu-check
[4]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
