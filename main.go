package main

import (
	"context"
	"flag"
	"strings"

	cdkpython "github.com/phsiao/fargate-migrate/internal/cdk/python"
	"github.com/phsiao/fargate-migrate/internal/config"
	"github.com/phsiao/fargate-migrate/internal/kubernetes"
	log "github.com/sirupsen/logrus"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "fargate-migrate config file to use")
}

func main() {
	flag.Parse()

	config, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	deps, err := kubernetes.LookupService(
		context.Background(),
		config.Spec.KubernetesConfig.Context,
		config.Spec.KubernetesConfig.Namespace,
		config.Spec.KubernetesConfig.Service,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("%v", deps)

	stack := cdkpython.NewFargateServiceStack(
		config.Spec.FargateConfig.StackName,
		config.Spec.FargateConfig.ServiceName,
		config.Spec.FargateConfig.AccountID,
		config.Spec.FargateConfig.Region,
		cdkpython.WithVPC(cdkpython.NewManagedVPCStatementGenerator()),
		cdkpython.WithDomain(cdkpython.NewHostedZoneStatementGenerator(config.Spec.FargateConfig.DomainName)),
		cdkpython.WithCluster(cdkpython.NewFargateClusterStatementGenerator()),
	)
	rval, err := stack.Generate()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("\n%s", strings.Join(rval, "\n"))
}
