package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

var awsIntegrationServiceDataResourceWithRetireAutomatically = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"enable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"excluded_metrics": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"retire_automatically": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	},
}

var awsIntegrationServiceDataSchemaWithRetireAutomatically = &schema.Schema{
	Type:     schema.TypeSet,
	Computed: true,
	Elem:     awsIntegrationServiceDataResourceWithRetireAutomatically,
}

var awsIntegrationServiceDataResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"enable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"excluded_metrics": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
}

var awsIntegrationServiceDataSchema = &schema.Schema{
	Type:     schema.TypeSet,
	Computed: true,
	Elem:     awsIntegrationServiceDataResource,
}

func dataSourceMackerelAWSIntegration() *schema.Resource {
	resource := &schema.Resource{
		ReadContext: dataSourceMackerelAWSIntegrationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"included_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"excluded_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	var supportedRetireAutomatically = map[string]bool{"ec2": true}
	for schemaKey := range awsIntegrationServicesKey {
		if supportedRetireAutomatically[schemaKey] {
			resource.Schema[schemaKey] = awsIntegrationServiceDataSchemaWithRetireAutomatically
		} else if !supportedRetireAutomatically[schemaKey] {
			resource.Schema[schemaKey] = awsIntegrationServiceDataSchema
		}
	}
	return resource
}

func dataSourceMackerelAWSIntegrationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	client := m.(*mackerel.Client)

	awsIntegrations, err := client.FindAWSIntegrations()
	if err != nil {
		return diag.FromErr(err)
	}
	var awsIntegration *mackerel.AWSIntegration
	for _, a := range awsIntegrations {
		if a.ID == id {
			awsIntegration = a
			break
		}
	}
	if awsIntegration == nil {
		return diag.Errorf(`the id '%s' does not match any aws integration in mackerel.io`, id)
	}
	d.SetId(awsIntegration.ID)
	return flattenAWSIntegration(awsIntegration, d)
}
