package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFargateClusterStatementGenerator(t *testing.T) {
	g := FargateClusterStatementGenerator{
		Name: "TestCluster",
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, `ecs.Cluster(self, "TestCluster", vpc=vpc, capacity_providers=["FARGATE", "FARGATE_SPOT"])`, actual)
}
