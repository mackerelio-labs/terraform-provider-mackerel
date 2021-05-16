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

### Environment Variables

You can provide your Mackerel API key using environment variables, `MACKEREL_API_KEY`.

## Argument Reference

* `api_key` - (Optional) Mackerel API Key. It must be provided, but it can also be sourced from the `MACKEREL_API_KEY` environment variable.
