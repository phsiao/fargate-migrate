package python

import (
	"bytes"
	"text/template"
)

type FargateClusterStatementGenerator struct {
	Name string
}

func NewFargateClusterStatementGenerator() *FargateClusterStatementGenerator {
	rval := FargateClusterStatementGenerator{
		Name: "ManagedCluster",
	}

	return &rval
}

func (g FargateClusterStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("cluster").Parse(`ecs.Cluster(self, "{{.Name}}", vpc=vpc, enable_fargate_capacity_providers=True)`)
	if err != nil {
		return "", err
	}
	writer := &bytes.Buffer{}
	err = tmpl.Execute(writer, &g)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}

var _ PythonCodeSnippetGenerator = FargateClusterStatementGenerator{}
