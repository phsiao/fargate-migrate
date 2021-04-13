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

func TestLookupVPCStatementGenerator(t *testing.T) {
	g := LookupVPCStatementGenerator{
		Name:  "foo",
		VPCID: "foobar",
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, `ec2.Vpc.from_lookup(self, "foo", vpc_id="foobar")`, actual)
}
