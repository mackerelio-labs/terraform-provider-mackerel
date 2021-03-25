# Terraform provider for mackerel.io

[![CI](https://github.com/mackerelio-labs/terraform-provider-mackerel/actions/workflows/ci.yml/badge.svg)](https://github.com/mackerelio-labs/terraform-provider-mackerel/actions/workflows/ci.yml)

A [Terraform](https://www.terraform.io/) provider for [mackerel.io](https://mackerel.io/).

- Terraform Website: https://terraform.io
- Slack Workspace: https://mackerel-ug-slackin.herokuapp.com/

## Usage example

Terraform 0.14 and later

```
terraform {
  required_providers {
    mackerel = {
      source  = "Mackerel/mackerel"
      version = "~> 0.0.1"
    }
  }
}

resource "mackerel_service" "app" {
  name = "app"
}

resource "mackerel_role" "compute" {
  service = mackerel_service.app.name
  name    = "ecs"
}
```

## Requirements

Terraform >= v0.14

## Acknowledgements

We thank @xcezx and @kjmkznr for contributing to terraform-provider-mackerel.

## License

Copyright 2021 Tsuyoshi Maekawa

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
