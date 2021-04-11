package python

import (
	"fmt"
	"strings"
)

type FargateServiceStack struct {
	stackName   string
	serviceName string
	accountID   string
	region      string
	taskPort    int
	taskAsset   string
	importGen   *CommonImportStatementGenerator
	vpcGen      *ManagedVPCStatementGenerator
	dnsGen      *HostedZoneStatementGenerator
	clusterGen  *FargateClusterStatementGenerator
}

type Option func(*FargateServiceStack)

func NewFargateServiceStack(
	stackName string,
	serviceName string,
	accountID string,
	region string,
	taskPort int,
	taskAsset string,
	options ...Option) *FargateServiceStack {
	s := &FargateServiceStack{
		stackName:   stackName,
		serviceName: serviceName,
		accountID:   accountID,
		region:      region,
		taskPort:    taskPort,
		taskAsset:   taskAsset,
		importGen:   &CommonImportStatementGenerator{},
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func WithVPC(vpcGen *ManagedVPCStatementGenerator) Option {
	return func(s *FargateServiceStack) {
		s.vpcGen = vpcGen
	}
}

func WithDomain(zoneGen *HostedZoneStatementGenerator) Option {
	return func(s *FargateServiceStack) {
		s.dnsGen = zoneGen
	}
}

func WithCluster(clusterGen *FargateClusterStatementGenerator) Option {
	return func(s *FargateServiceStack) {
		s.clusterGen = clusterGen
	}
}

var _ PythonCodeSnippetGenerator = FargateServiceStack{}

type CommonImportStatementGenerator struct {
}

func (g CommonImportStatementGenerator) Generate() ([]string, error) {
	return []string{
		"from aws_cdk import (core as cdk, aws_ec2 as ec2, aws_ecs as ecs,",
		"                     aws_ecs_patterns as ecs_patterns,",
		"                     aws_route53 as route53, aws_elasticloadbalancingv2 as elbv2)",
	}, nil
}

var _ PythonCodeSnippetGenerator = CommonImportStatementGenerator{}

func (s FargateServiceStack) Generate() ([]string, error) {
	output := []string{}
	var snippets []string
	var err error

	snippets, err = s.importGen.Generate()
	if err != nil {
		return nil, err
	}
	output = append(output, snippets...)
	output = append(output, "\n")

	output = append(output, fmt.Sprintf("class %sStack(cdk.Stack):", s.stackName))
	output = append(output, "\n")
	output = append(output, "    def __init__(self, scope: cdk.Construct, construct_id: str, **kwargs) -> None:")
	output = append(output, "\n")
	output = append(output, "        super().__init__(scope, construct_id, **kwargs)")

	snippets, err = s.dnsGen.Generate()
	if err != nil {
		return nil, err
	}
	output = append(output, fmt.Sprintf("        zone = %s", strings.Join(snippets, "\n")))

	snippets, err = s.vpcGen.Generate()
	if err != nil {
		return nil, err
	}
	output = append(output, fmt.Sprintf("        vpc = %s", strings.Join(snippets, "\n")))

	snippets, err = s.clusterGen.Generate()
	if err != nil {
		return nil, err
	}
	output = append(output, fmt.Sprintf("        cluster = %s", strings.Join(snippets, "\n")))

	output = append(output, fmt.Sprintf(`        svc = ecs_patterns.ApplicationLoadBalancedFargateService(self, "%sService",`, s.stackName))
	output = append(output, "            cluster=cluster,")
	output = append(output, "            redirect_http=True,")
	output = append(output, "            desired_count=1,")
	output = append(output, "            memory_limit_mib=512,")
	output = append(output, "            public_load_balancer=True,")
	output = append(output, "            protocol=elbv2.ApplicationProtocol.HTTPS,")
	output = append(output, "            domain_zone=zone,")
	output = append(output, "            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(")
	output = append(output, fmt.Sprintf(`                container_port=%d,`, s.taskPort))
	output = append(output, fmt.Sprintf(`                image=ecs.ContainerImage.from_asset("%s")),`, s.taskAsset))
	output = append(output, fmt.Sprintf(`            domain_name="%s"`, s.serviceName))
	output = append(output, "            )")

	output = append(output, "\n")

	output = append(output, "app = cdk.App()")
	output = append(output, fmt.Sprintf(
		`%sStack(app, "%sStack", env=cdk.Environment(account="%s", region="%s"))`,
		s.stackName,
		s.stackName,
		s.accountID,
		s.region))
	output = append(output, "app.synth()")

	return output, nil
}
