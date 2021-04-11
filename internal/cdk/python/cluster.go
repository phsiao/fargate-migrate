package python

import "fmt"

type FargateClusterStatementGenerator struct {
	clusterName string
}

func NewFargateClusterStatementGenerator() *FargateClusterStatementGenerator {
	rval := FargateClusterStatementGenerator{
		clusterName: "ManagedCluster",
	}

	return &rval
}

func (g FargateClusterStatementGenerator) Generate() ([]string, error) {
	rval := []string{
		fmt.Sprintf("ecs.Cluster(self, \"%s\", vpc=vpc, capacity_providers=[\"FARGATE\", \"FARGATE_SPOT\"])",
			g.clusterName,
		),
	}
	return rval, nil
}

var _ PythonCodeSnippetGenerator = FargateClusterStatementGenerator{}
