package mackerel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

func Test_AWSIntegration_fromAPI(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		api     mackerel.AWSIntegration
		model   AWSIntegrationModel
		wantErr bool
	}{
		"basic": {
			api: mackerel.AWSIntegration{
				ID:           "aaaabbbb",
				Name:         "aws-integration",
				Memo:         "This resource is managed by Terraform.",
				Key:          "",
				RoleArn:      "arn:aws:iam::11111111:role/MackerelAWSIntegrationRole",
				ExternalID:   "ccccddddd",
				Region:       "ap-northeast-1",
				IncludedTags: "Name:production-server,Environment:production",
				ExcludedTags: "Name:staging-server,Environment:staging",
				Services: map[string]*mackerel.AWSIntegrationService{
					"EC2": {
						Enable:              true,
						Role:                nil,
						RetireAutomatically: true,
						ExcludedMetrics:     []string{},
					},
					"ELB": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"ALB": {
						Enable:          true,
						Role:            ptr("service: role"),
						ExcludedMetrics: []string{"alb.request.count", "alb.bytes.processed"},
					},
					"RDS": {
						Enable:              true,
						Role:                ptr("service: role"),
						ExcludedMetrics:     []string{"rds.cpu.used"},
						RetireAutomatically: false,
					},
					"Redshift": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"ElastiCache": {
						Enable:              false,
						Role:                nil,
						ExcludedMetrics:     []string{},
						RetireAutomatically: false,
					},
					"SQS": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Lambda": {
						Enable:          true,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"NLB": {
						Enable:          true,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"DynamoDB": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"CloudFront": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"APIGateway": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Kinesis": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"S3": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"ES": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"ECSCluster": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"SES": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"States": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"EFS": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Firehose": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Batch": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"WAF": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Billing": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Route53": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"Connect": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"DocDB": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					"CodeBuild": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					// AWS Integration supports Athena, but the terraform provider does not.
					"Athena": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
					// Unsupported services are ignored when they are empty.
					"SomeNewService": {
						Enable:          false,
						Role:            nil,
						ExcludedMetrics: []string{},
					},
				},
			},
			model: AWSIntegrationModel{
				ID:           types.StringValue("aaaabbbb"),
				Name:         types.StringValue("aws-integration"),
				Memo:         types.StringValue("This resource is managed by Terraform."),
				Key:          types.StringValue(""),
				SecretKey:    types.StringNull(),
				RoleARN:      types.StringValue("arn:aws:iam::11111111:role/MackerelAWSIntegrationRole"),
				ExternalID:   types.StringValue("ccccddddd"),
				Region:       types.StringValue("ap-northeast-1"),
				IncludedTags: types.StringValue("Name:production-server,Environment:production"),
				ExcludedTags: types.StringValue("Name:staging-server,Environment:staging"),

				EC2: []AWSIntegrationServiceWithRetireAutomatically{{
					Enable:              types.BoolValue(true),
					ExcludedMetrics:     []string{},
					RetireAutomatically: types.BoolValue(true),
				}},
				ELB: []AWSIntegrationService{},
				ALB: []AWSIntegrationService{{
					Enable:              types.BoolValue(true),
					Role:                types.StringValue("service: role"),
					ExcludedMetrics:     []string{"alb.request.count", "alb.bytes.processed"},
					RetireAutomatically: types.BoolValue(false),
				}},
				NLB: []AWSIntegrationService{{
					Enable:              types.BoolValue(true),
					Role:                types.StringNull(),
					ExcludedMetrics:     []string{},
					RetireAutomatically: types.BoolValue(false),
				}},
				RDS: []AWSIntegrationServiceWithRetireAutomatically{{
					Enable:              types.BoolValue(true),
					Role:                types.StringValue("service: role"),
					ExcludedMetrics:     []string{"rds.cpu.used"},
					RetireAutomatically: types.BoolValue(false),
				}},
				Redshift:    []AWSIntegrationService{},
				ElastiCache: []AWSIntegrationServiceWithRetireAutomatically{},
				SQS:         []AWSIntegrationService{},
				Lambda: []AWSIntegrationService{{
					Enable:              types.BoolValue(true),
					Role:                types.StringNull(),
					ExcludedMetrics:     []string{},
					RetireAutomatically: types.BoolValue(false),
				}},
				DynamoDB:   []AWSIntegrationService{},
				CloudFront: []AWSIntegrationService{},
				APIGateway: []AWSIntegrationService{},
				Kinesis:    []AWSIntegrationService{},
				S3:         []AWSIntegrationService{},
				ES:         []AWSIntegrationService{},
				ECSCluster: []AWSIntegrationService{},
				SES:        []AWSIntegrationService{},
				States:     []AWSIntegrationService{},
				EFS:        []AWSIntegrationService{},
				Firehose:   []AWSIntegrationService{},
				Batch:      []AWSIntegrationService{},
				WAF:        []AWSIntegrationService{},
				Billing:    []AWSIntegrationService{},
				Route53:    []AWSIntegrationService{},
				Connect:    []AWSIntegrationService{},
				DocDB:      []AWSIntegrationService{},
				CodeBuild:  []AWSIntegrationService{},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			model, err := newAWSIntegrationModel(tt.api)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %+v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(model, &tt.model); diff != "" {
				t.Error(diff)
			}
		})
	}
}
