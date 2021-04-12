package python

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type FargateServiceStack struct {
	Name        string
	ServiceName string
	AccountID   string
	Region      string
	TaskPort    int32
	TaskAsset   string
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
	options ...Option) *FargateServiceStack {
	s := &FargateServiceStack{
		Name:        stackName,
		ServiceName: serviceName,
		AccountID:   accountID,
		Region:      region,
		TaskPort:    80,
		importGen:   &CommonImportStatementGenerator{},
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func WithContainerPort(port int32) Option {
	return func(s *FargateServiceStack) {
		s.TaskPort = port
	}
}

func WithAsset(asset string) Option {
	return func(s *FargateServiceStack) {
		s.TaskAsset = asset
	}
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

func (g CommonImportStatementGenerator) Generate() (string, error) {
	rval :=
		`
from aws_cdk import (core as cdk, aws_ec2 as ec2, aws_ecs as ecs,
                      aws_ecs_patterns as ecs_patterns,
                      aws_route53 as route53, aws_elasticloadbalancingv2 as elbv2)
`
	return rval, nil
}

var _ PythonCodeSnippetGenerator = CommonImportStatementGenerator{}

func (s FargateServiceStack) Generate() (string, error) {
	var err error

	stackTmpl :=
		`
{{import}}

class {{.Name}}Stack(cdk.Stack):
    def __init__(self, scope: cdk.Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)
        zone = {{zone}}
        vpc = {{vpc}}
        cluster = {{cluster}}

        svc = ecs_patterns.ApplicationLoadBalancedFargateService(self, "{{.Name}}Service",
            cluster=cluster,
            redirect_http=True,
            desired_count=1,
            memory_limit_mib=512,
            public_load_balancer=True,
            protocol=elbv2.ApplicationProtocol.HTTPS,
            domain_zone=zone,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                container_port={{.TaskPort}},
                image=ecs.ContainerImage.from_asset("{{.TaskAsset}}")),
                domain_name="{{.ServiceName}}"
            )

app = cdk.App()
{{.Name}}Stack(app, "{{.Name}}Stack", env=cdk.Environment(account="{{.AccountID}}", region="{{.Region}}"))
app.synth()
`

	writer := &bytes.Buffer{}

	tmpl, err := template.New("stack").Funcs(
		template.FuncMap{
			"import": func() string {
				stmt, err := s.importGen.Generate()
				if err != nil {
					log.Fatal(err)
				}
				return stmt
			},
			"zone": func() string {
				stmt, err := s.dnsGen.Generate()
				if err != nil {
					log.Fatal(err)
				}
				return stmt
			},
			"vpc": func() string {
				stmt, err := s.vpcGen.Generate()
				if err != nil {
					log.Fatal(err)
				}
				return stmt
			},
			"cluster": func() string {
				stmt, err := s.clusterGen.Generate()
				if err != nil {
					log.Fatal(err)
				}
				return stmt
			},
		},
	).Parse(stackTmpl)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(writer, &s)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}
