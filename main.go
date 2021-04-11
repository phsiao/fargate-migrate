package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cdkpython "github.com/phsiao/fargate-migrate/internal/cdk/python"
	"github.com/phsiao/fargate-migrate/internal/config"
	"github.com/phsiao/fargate-migrate/internal/kubernetes"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	configFile  string
	CDKPath     string = "cdk"
	APPPath     string = "app.py"
	CDKJsonPath string = "cdk.json"
	DockerPath  string = "docker"
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "fargate-migrate config file to use")
}

func writeCDKArtifacts(config *config.Config, service *corev1.Service, deployment *appsv1.Deployment) error {
	var err error

	// create CDK directory
	if _, err := os.Stat(CDKPath); os.IsNotExist(err) {
		os.Mkdir(CDKPath, 0744)
	}

	// write cdk.json file
	cdkJson := []string{
		`{`,
		`    "app": "python3 app.py"`,
		`}`,
	}
	err = ioutil.WriteFile(filepath.Join(CDKPath, CDKJsonPath), []byte(strings.Join(cdkJson, "\n")), 0644)
	if err != nil {
		return err
	}

	// create docker directory
	if _, err := os.Stat(filepath.Join(CDKPath, DockerPath)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(CDKPath, DockerPath), 0744)
	}

	var firstTaskAsset string
	for _, container := range deployment.Spec.Template.Spec.Containers {
		// create container docker directory
		if _, err := os.Stat(filepath.Join(CDKPath, DockerPath, container.Name)); os.IsNotExist(err) {
			os.Mkdir(filepath.Join(CDKPath, DockerPath, container.Name), 0744)
		}
		dockerfile := fmt.Sprintf("FROM %s", container.Image)
		err = ioutil.WriteFile(
			filepath.Join(CDKPath, DockerPath, container.Name, "Dockerfile"),
			[]byte(dockerfile), 0644,
		)
		if err != nil {
			return err
		}

		if firstTaskAsset == "" {
			// taskAssetPath is relative to CDKPath
			firstTaskAsset = filepath.Join(DockerPath, container.Name)
		}
	}

	stack := cdkpython.NewFargateServiceStack(
		config.Spec.FargateConfig.StackName,
		config.Spec.FargateConfig.ServiceName,
		config.Spec.FargateConfig.AccountID,
		config.Spec.FargateConfig.Region,
		int(service.Spec.Ports[0].Port),
		firstTaskAsset,
		cdkpython.WithVPC(cdkpython.NewManagedVPCStatementGenerator()),
		cdkpython.WithDomain(cdkpython.NewHostedZoneStatementGenerator(config.Spec.FargateConfig.DomainName)),
		cdkpython.WithCluster(cdkpython.NewFargateClusterStatementGenerator()),
	)
	rval, err := stack.Generate()
	if err != nil {
		return err
	}
	data := strings.Join(rval, "\n")
	err = ioutil.WriteFile(filepath.Join(CDKPath, APPPath), []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	config, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	svc, deps, err := kubernetes.LookupService(
		context.Background(),
		config.Spec.KubernetesConfig.Context,
		config.Spec.KubernetesConfig.Namespace,
		config.Spec.KubernetesConfig.Service,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(svc.Spec.Ports) != 1 {
		log.Fatal("only support service with exactly one port")
	}

	if len(deps) > 1 {
		log.Fatal("only support services backed by exactly one deployment")
	}

	dep := deps[0]
	if len(dep.Spec.Template.Spec.Containers) != 1 {
		log.Fatal("only support deployment with exactly one container")
	}

	err = writeCDKArtifacts(config, svc, deps[0])
	if err != nil {
		log.Fatal(err)
	}
}
