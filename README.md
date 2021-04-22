# fargate-migrate

A tool to migrate a Kubernetes service to run on AWS Fargate.

It inspects the configuration of a Kubernetes Service object,
in particular the pods that is backing the service, and generates
corresponding Fargate configuration using [CDK](https://aws.amazon.com/cdk/).
A user can then use CDK commands to manage the Fargate deployment.

## Usage

You can install it with

```
$ go get github.com/phsiao/fargate-migrate
```

The command expects a config file (default is `config.yaml` in the
current directory) as input, and the content of the config file
is defined in [`internal/config/config.go`](internal/config/config.go).

Upon successful execution, it outpus a CDK stack in the `cdk/`
directory of current working directory.  The content of the `cdk`
directory is ready for you to run `cdk deploy` against your account.

A very basic example of config file is:

```yaml
spec:
  kubernetesConfig:
    context: mycluster
    namespace: mynamespace
    service: myservice
  fargateConfig:
    name: myservice
    accountID: 123456789012
    region: us-east-1
    serviceName: myservice.mydomain.com
    domainName: mydomain.com
```