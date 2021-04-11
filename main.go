package main

import (
	"context"
	"flag"

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

	log.Infof("%v", deps)
}
