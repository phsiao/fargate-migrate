package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagedVPCStatementGenerator(t *testing.T) {
	g := ManagedVPCStatementGenerator{
		Name:   "foo",
		MaxAZs: 3,
		CIDR:   "10.0.0.0/28",
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, `ec2.Vpc(self, "foo", cidr="10.0.0.0/28", max_azs=3)`, actual)
}
