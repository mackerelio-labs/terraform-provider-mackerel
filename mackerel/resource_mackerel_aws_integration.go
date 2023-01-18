package mackerel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

var awsIntegrationServiceEC2Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"role": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"excluded_metrics": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"retire_automatically": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	},
}

var awsIntegrationServiceResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"role": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"excluded_metrics": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
}

var awsIntegrationServiceSchema = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	Elem:     awsIntegrationServiceResource,
}

var awsIntegrationServicesKey = map[string]string{
	"ec2":         "EC2",
	"elb":         "ELB",
	"alb":         "ALB",
	"nlb":         "NLB",
	"rds":         "RDS",
	"redshift":    "Redshift",
	"elasticache": "ElastiCache",
	"sqs":         "SQS",
	"lambda":      "Lambda",
	"dynamodb":    "DynamoDB",
	"cloudfront":  "CloudFront",
	"api_gateway": "APIGateway",
	"kinesis":     "Kinesis",
	"s3":          "S3",
	"es":          "ES",
	"ecs_cluster": "ECSCluster",
	"ses":         "SES",
	"states":      "States",
	"efs":         "EFS",
	"firehose":    "Firehose",
	"batch":       "Batch",
	"waf":         "WAF",
	"billing":     "Billing",
	"route53":     "Route53",
	"connect":     "Connect",
	"docdb":       "DocDB",
	"codebuild":   "CodeBuild",
}

func resourceMackerelAWSIntegration() *schema.Resource {
	resource := &schema.Resource{
		CreateContext: resourceMackerelAWSIntegrationCreate,
		ReadContext:   resourceMackerelAWSIntegrationRead,
		UpdateContext: resourceMackerelAWSIntegrationUpdate,
		DeleteContext: resourceMackerelAWSIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"role_arn": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"external_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"included_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"excluded_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ec2": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     awsIntegrationServiceEC2Resource,
			},
		},
	}
	for schemaKey := range awsIntegrationServicesKey {
		if schemaKey != "ec2" {
			resource.Schema[schemaKey] = awsIntegrationServiceSchema
		}
	}
	return resource
}

func resourceMackerelAWSIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	awsIntegration, err := client.CreateAWSIntegration(expandCrateAWSIntegrationParam(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(awsIntegration.ID)
	return resourceMackerelAWSIntegrationRead(ctx, d, m)
}

func resourceMackerelAWSIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	awsIntegration, err := client.FindAWSIntegration(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return flattenAWSIntegration(awsIntegration, d)
}

func resourceMackerelAWSIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mackerel.Client)
	_, err := client.UpdateAWSIntegration(d.Id(), expandUpdateAWSIntegrationParam(d))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMackerelAWSIntegrationRead(ctx, d, m)
}

func resourceMackerelAWSIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mackerel.Client)
	_, err := client.DeleteAWSIntegration(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandCrateAWSIntegrationParam(d *schema.ResourceData) *mackerel.CreateAWSIntegrationParam {
	awsIntegration := &mackerel.CreateAWSIntegrationParam{
		Name:         d.Get("name").(string),
		Memo:         d.Get("memo").(string),
		Key:          d.Get("key").(string),
		SecretKey:    d.Get("secret_key").(string),
		RoleArn:      d.Get("role_arn").(string),
		ExternalID:   d.Get("external_id").(string),
		Region:       d.Get("region").(string),
		IncludedTags: d.Get("included_tags").(string),
		ExcludedTags: d.Get("excluded_tags").(string),
		Services:     expandAWSIntegrationServicesSet(d),
	}
	return awsIntegration
}

func expandUpdateAWSIntegrationParam(d *schema.ResourceData) *mackerel.UpdateAWSIntegrationParam {
	awsIntegration := &mackerel.UpdateAWSIntegrationParam{
		Name:         d.Get("name").(string),
		Memo:         d.Get("memo").(string),
		Key:          d.Get("key").(string),
		SecretKey:    d.Get("secret_key").(string),
		RoleArn:      d.Get("role_arn").(string),
		ExternalID:   d.Get("external_id").(string),
		Region:       d.Get("region").(string),
		IncludedTags: d.Get("included_tags").(string),
		ExcludedTags: d.Get("excluded_tags").(string),
		Services:     expandAWSIntegrationServicesSet(d),
	}
	return awsIntegration
}

func expandAWSIntegrationServicesSet(d *schema.ResourceData) map[string]*mackerel.AWSIntegrationService {
	services := make(map[string]*mackerel.AWSIntegrationService)
	for schemaKey, mapKey := range awsIntegrationServicesKey {
		if _, ok := d.GetOk(schemaKey); ok {
			l := d.Get(schemaKey).(*schema.Set).List()
			service := l[0].(map[string]interface{})
			services[mapKey] = &mackerel.AWSIntegrationService{
				Enable:          service["enable"].(bool),
				Role:            toPointer(service["role"].(string)),
				ExcludedMetrics: toSliceString(service["excluded_metrics"].([]interface{})),
			}
			if schemaKey == "ec2" {
				services[mapKey].RetireAutomatically = service["retire_automatically"].(bool)
			}
		}
	}
	return deleteAWSIntegrationDisableService(services)
}

func toPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toSliceString(target []interface{}) []string {
	var s []string
	for _, v := range target {
		if v == nil {
			continue
		} else {
			s = append(s, v.(string))
		}
	}
	if s == nil {
		s = make([]string, 0)
	}
	return s
}

func toSliceInterface(target []string) []interface{} {
	var s []interface{}
	for _, v := range target {
		s = append(s, v)
	}
	return s
}

func toAWSIntegrationServicesSchemaKey(target string) string {
	for hashKey, mapKey := range awsIntegrationServicesKey {
		if mapKey == target {
			return hashKey
		}
	}
	return ""
}

func deleteAWSIntegrationDisableService(s map[string]*mackerel.AWSIntegrationService) map[string]*mackerel.AWSIntegrationService {
	services := make(map[string]*mackerel.AWSIntegrationService)
	for key, service := range s {
		if service.Enable {
			services[key] = service
		}
	}
	return services
}
