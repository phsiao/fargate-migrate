package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagedVPCStatementGenerator(t *testing.T) {
	g := ManagedVPCStatementGenerator{
		vpcName: "foo",
		maxAZs:  3,
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ec2.Vpc(self, \"foo\", max_azs=3)"}, actual)
}
