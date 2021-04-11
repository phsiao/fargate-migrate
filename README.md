# fargate-migrate

A tool to migrate a Kubernetes service to run on AWS Fargate.

It inspects the configuration of a Kubernetes Service object,
in particular the pods that is backing the service, and generates
corresponding Fargate configuration using [CDK](https://aws.amazon.com/cdk/).
A user can then use CDK commands to manage the Fargate deployment.