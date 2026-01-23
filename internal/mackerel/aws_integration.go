package mackerel

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mackerelio/mackerel-client-go"
)

type AWSIntegrationModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Memo         types.String `tfsdk:"memo"`
	Key          types.String `tfsdk:"key"`
	SecretKey    types.String `tfsdk:"secret_key"`
	RoleARN      types.String `tfsdk:"role_arn"`
	ExternalID   types.String `tfsdk:"external_id"`
	Region       types.String `tfsdk:"region"`
	IncludedTags types.String `tfsdk:"included_tags"`
	ExcludedTags types.String `tfsdk:"excluded_tags"`

	AWSIntegrationSerfvices
}

type AWSIntegrationDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Memo         types.String `tfsdk:"memo"`
	Key          types.String `tfsdk:"key"`
	SecretKey    types.String `tfsdk:"-"`
	RoleARN      types.String `tfsdk:"role_arn"`
	ExternalID   types.String `tfsdk:"external_id"`
	Region       types.String `tfsdk:"region"`
	IncludedTags types.String `tfsdk:"included_tags"`
	ExcludedTags types.String `tfsdk:"excluded_tags"`

	AWSIntegrationSerfvices
}

type AWSIntegrationSerfvices struct {
	EC2         AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"ec2"`
	ELB         AWSIntegrationServiceOpt                        `tfsdk:"elb"`
	ALB         AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"alb"`
	NLB         AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"nlb"`
	RDS         AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"rds"`
	Redshift    AWSIntegrationServiceOpt                        `tfsdk:"redshift"`
	ElastiCache AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"elasticache"`
	SQS         AWSIntegrationServiceOpt                        `tfsdk:"sqs"`
	Lambda      AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"lambda"`
	DynamoDB    AWSIntegrationServiceOpt                        `tfsdk:"dynamodb"`
	CloudFront  AWSIntegrationServiceOpt                        `tfsdk:"cloudfront"`
	APIGateway  AWSIntegrationServiceOpt                        `tfsdk:"api_gateway"`
	Kinesis     AWSIntegrationServiceOpt                        `tfsdk:"kinesis"`
	S3          AWSIntegrationServiceOpt                        `tfsdk:"s3"`
	ES          AWSIntegrationServiceOpt                        `tfsdk:"es"`
	ECSCluster  AWSIntegrationServiceWithRetireAutomaticallyOpt `tfsdk:"ecs_cluster"`
	SES         AWSIntegrationServiceOpt                        `tfsdk:"ses"`
	States      AWSIntegrationServiceOpt                        `tfsdk:"states"`
	EFS         AWSIntegrationServiceOpt                        `tfsdk:"efs"`
	Firehose    AWSIntegrationServiceOpt                        `tfsdk:"firehose"`
	Batch       AWSIntegrationServiceOpt                        `tfsdk:"batch"`
	WAF         AWSIntegrationServiceOpt                        `tfsdk:"waf"`
	Billing     AWSIntegrationServiceOpt                        `tfsdk:"billing"`
	Route53     AWSIntegrationServiceOpt                        `tfsdk:"route53"`
	Connect     AWSIntegrationServiceOpt                        `tfsdk:"connect"`
	DocDB       AWSIntegrationServiceOpt                        `tfsdk:"docdb"`
	CodeBuild   AWSIntegrationServiceOpt                        `tfsdk:"codebuild"`
}

type AWSIntegrationService struct {
	Enable              types.Bool   `tfsdk:"enable"`
	Role                types.String `tfsdk:"role"`
	ExcludedMetrics     []string     `tfsdk:"excluded_metrics"`
	RetireAutomatically types.Bool   `tfsdk:"-"`
}

func (s *AWSIntegrationService) isEmpty() bool {
	return !s.Enable.ValueBool() &&
		s.Role.IsNull() &&
		len(s.ExcludedMetrics) == 0 &&
		!s.RetireAutomatically.ValueBool()
}

type AWSIntegrationServiceOpt []AWSIntegrationService // length <= 1

type AWSIntegrationServiceWithRetireAutomatically struct {
	Enable              types.Bool   `tfsdk:"enable"`
	Role                types.String `tfsdk:"role"`
	ExcludedMetrics     []string     `tfsdk:"excluded_metrics"`
	RetireAutomatically types.Bool   `tfsdk:"retire_automatically"`
}

type AWSIntegrationServiceWithRetireAutomaticallyOpt []AWSIntegrationServiceWithRetireAutomatically // length <= 1

func ReadAWSIntegration(_ context.Context, client *Client, id string) (*AWSIntegrationModel, error) {
	return readAWSIntegration(client, id)
}

func readAWSIntegration(client *Client, id string) (*AWSIntegrationModel, error) {
	mackerelAWSIntegration, err := client.FindAWSIntegration(id)
	if err != nil {
		return nil, err
	}
	return newAWSIntegrationModel(*mackerelAWSIntegration)
}

func (m *AWSIntegrationModel) Create(_ context.Context, client *Client) error {
	newIntegration, err := client.CreateAWSIntegration(m.createParam())
	if err != nil {
		return err
	}

	m.ID = types.StringValue(newIntegration.ID)
	return nil
}

func (m *AWSIntegrationModel) Read(_ context.Context, client *Client) error {
	integration, err := readAWSIntegration(client, m.ID.ValueString())
	if err != nil {
		return err
	}

	// Inherit secret key from the existing integration
	integration.SecretKey = m.SecretKey

	// Copy existing services if they are empty
	oldSvcs := map[string]*AWSIntegrationService{}
	m.each(func(name string, service *AWSIntegrationService) *AWSIntegrationService {
		if service != nil {
			oldSvcs[name] = service
		}
		return service
	})
	integration.each(func(name string, service *AWSIntegrationService) *AWSIntegrationService {
		if service == nil {
			if oldSvc, ok := oldSvcs[name]; ok && oldSvc.isEmpty() {
				return oldSvc
			}
		}
		return service
	})

	*m = *integration
	return nil
}

func (m *AWSIntegrationModel) Update(_ context.Context, client *Client) error {
	if _, err := client.UpdateAWSIntegration(m.ID.ValueString(), m.updateParam()); err != nil {
		return err
	}
	return nil
}

func (m *AWSIntegrationModel) Delete(_ context.Context, client *Client) error {
	if _, err := client.DeleteAWSIntegration(m.ID.ValueString()); err != nil {
		return err
	}
	return nil
}

func newAWSIntegrationModel(aws mackerel.AWSIntegration) (*AWSIntegrationModel, error) {
	model := &AWSIntegrationModel{
		ID:           types.StringValue(aws.ID),
		Name:         types.StringValue(aws.Name),
		Memo:         types.StringValue(aws.Memo),
		Key:          types.StringValue(aws.Key),
		RoleARN:      types.StringValue(aws.RoleArn),
		ExternalID:   types.StringValue(aws.ExternalID),
		Region:       types.StringValue(aws.Region),
		IncludedTags: types.StringValue(aws.IncludedTags),
		ExcludedTags: types.StringValue(aws.ExcludedTags),
	}

	svcs := make(map[string]AWSIntegrationService, len(aws.Services))
	for name, awsService := range aws.Services {
		if /* nil */ !awsService.Enable &&
			awsService.Role == nil &&
			len(awsService.ExcludedMetrics) == 0 &&
			len(awsService.IncludedMetrics) == 0 &&
			!awsService.RetireAutomatically {
			continue
		}
		if len(awsService.IncludedMetrics) != 0 {
			return nil, fmt.Errorf("%s: IncludedMetrics is not supported", name)
		}

		svcs[name] = AWSIntegrationService{
			Enable:              types.BoolValue(awsService.Enable),
			Role:                types.StringPointerValue(awsService.Role),
			ExcludedMetrics:     awsService.ExcludedMetrics,
			RetireAutomatically: types.BoolValue(awsService.RetireAutomatically),
		}
	}
	model.each(func(name string, _ *AWSIntegrationService) *AWSIntegrationService {
		svc, ok := svcs[name]
		if ok {
			delete(svcs, name)
			return &svc
		} else {
			return nil
		}
	})
	if len(svcs) != 0 {
		unsupportedServiceNames := make([]string, 0, len(svcs))
		for name := range svcs {
			unsupportedServiceNames = append(unsupportedServiceNames, name)
		}
		slices.SortStableFunc(unsupportedServiceNames, strings.Compare)
		return nil, fmt.Errorf("unsupported AWS integration service(s): %s",
			strings.Join(unsupportedServiceNames, ","))
	}

	return model, nil
}

func (m *AWSIntegrationModel) createParam() *mackerel.CreateAWSIntegrationParam {
	mackerelServices := make(map[string]*mackerel.AWSIntegrationService)
	m.each(func(name string, service *AWSIntegrationService) *AWSIntegrationService {
		var mackerelService mackerel.AWSIntegrationService
		if service != nil {
			mackerelService = mackerel.AWSIntegrationService{
				Enable:              service.Enable.ValueBool(),
				Role:                nil,
				ExcludedMetrics:     nil,
				IncludedMetrics:     nil,
				RetireAutomatically: service.RetireAutomatically.ValueBool(),
			}
			if role := service.Role.ValueString(); role != "" {
				mackerelService.Role = &role
			}
			if service.ExcludedMetrics != nil {
				mackerelService.ExcludedMetrics = service.ExcludedMetrics
			} else {
				mackerelService.ExcludedMetrics = []string{}
			}
		} else {
			mackerelService = mackerel.AWSIntegrationService{
				Enable:              false,
				Role:                nil,
				ExcludedMetrics:     []string{},
				IncludedMetrics:     nil,
				RetireAutomatically: false,
			}
		}
		mackerelServices[name] = &mackerelService
		return service
	})

	return &mackerel.CreateAWSIntegrationParam{
		Name:         m.Name.ValueString(),
		Memo:         m.Memo.ValueString(),
		Key:          m.Key.ValueString(),
		SecretKey:    m.SecretKey.ValueString(),
		RoleArn:      m.RoleARN.ValueString(),
		ExternalID:   m.ExternalID.ValueString(),
		Region:       m.Region.ValueString(),
		IncludedTags: m.IncludedTags.ValueString(),
		ExcludedTags: m.ExcludedTags.ValueString(),
		Services:     mackerelServices,
	}
}

func (m *AWSIntegrationModel) updateParam() *mackerel.UpdateAWSIntegrationParam {
	return (*mackerel.UpdateAWSIntegrationParam)(m.createParam())
}

type awsServiceEachFunc func(name string, service *AWSIntegrationService) *AWSIntegrationService

// Iterates and updates over services by name
func (m *AWSIntegrationSerfvices) each(fn awsServiceEachFunc) {
	m.EC2.each("EC2", fn)
	m.ELB.each("ELB", fn)
	m.ALB.each("ALB", fn)
	m.NLB.each("NLB", fn)
	m.RDS.each("RDS", fn)
	m.Redshift.each("Redshift", fn)
	m.ElastiCache.each("ElastiCache", fn)
	m.SQS.each("SQS", fn)
	m.Lambda.each("Lambda", fn)
	m.DynamoDB.each("DynamoDB", fn)
	m.CloudFront.each("CloudFront", fn)
	m.APIGateway.each("APIGateway", fn)
	m.Kinesis.each("Kinesis", fn)
	m.S3.each("S3", fn)
	m.ES.each("ES", fn)
	m.ECSCluster.each("ECSCluster", fn)
	m.SES.each("SES", fn)
	m.States.each("States", fn)
	m.EFS.each("EFS", fn)
	m.Firehose.each("Firehose", fn)
	m.Batch.each("Batch", fn)
	m.WAF.each("WAF", fn)
	m.Billing.each("Billing", fn)
	m.Route53.each("Route53", fn)
	m.Connect.each("Connect", fn)
	m.DocDB.each("DocDB", fn)
	m.CodeBuild.each("CodeBuild", fn)
}

func (s *AWSIntegrationServiceOpt) each(name string, fn awsServiceEachFunc) {
	var svc *AWSIntegrationService
	if len(*s) != 0 {
		svc = &(*s)[0]
	}

	newSvc := fn(name, svc)
	if newSvc != nil {
		*s = []AWSIntegrationService{*newSvc}
	} else {
		*s = AWSIntegrationServiceOpt{}
	}
}

func (s *AWSIntegrationServiceWithRetireAutomaticallyOpt) each(name string, fn awsServiceEachFunc) {
	var svc *AWSIntegrationService
	if len(*s) != 0 {
		baseSvc := AWSIntegrationService((*s)[0])
		svc = &baseSvc
	}

	newSvc := fn(name, svc)
	if newSvc != nil {
		*s = []AWSIntegrationServiceWithRetireAutomatically{
			AWSIntegrationServiceWithRetireAutomatically(*newSvc),
		}
	} else {
		*s = AWSIntegrationServiceWithRetireAutomaticallyOpt{}
	}
}
