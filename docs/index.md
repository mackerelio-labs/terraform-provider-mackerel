---
page_title: "Provider: Mackerel"
subcategory: ""
description: |-
  The terraform provider for mackerel.io
---

# Mackerel Provider

The Mackerel provider provides resources to interact with a Mackerel API.

## Example Usage

Terraform 0.14 and later:

```terraform
terraform {
  required_providers {
    mackerel = {
      source  = "mackerelio-labs/mackerel"
      version = "~> 0.7"
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

## Authentication

Mackerel terraform provider offers two ways of setting credential.

### Environment Variables

You can provide your Mackerel API key using environment variables, `MACKEREL_APIKEY` or `MACKEREL_API_KEY`.

### Static credentials

Static credentials can be provided by adding `mackerel_api_key`.

Usage:

```terraform
variable "mackerel_api_key" {
}

provider "mackerel" {
  api_key = var.mackerel_api_key
}
```

## Argument Reference

* `api_key` - (Optional) Mackerel API Key. It must be provided, but it can also be sourced either from the `MACKEREL_APIKEY` or from the `MACKEREL_API_KEY` environment variable.
* `api_base` - (Optional) Mackerel API Endpoint. It can also be sourced from the `API_BASE` environment variable.
