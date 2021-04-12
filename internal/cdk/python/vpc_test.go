package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagedVPCStatementGenerator(t *testing.T) {
	g := ManagedVPCStatementGenerator{
		Name:   "foo",
		MaxAZs: 3,
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, `ec2.Vpc(self, "foo", max_azs=3)`, actual)
}
