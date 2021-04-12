package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostedZoneStatementGenerator(t *testing.T) {
	g := HostedZoneStatementGenerator{
		Name:       "TestZone",
		DomainName: "example.com",
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, `route53.HostedZone.from_lookup(self, "TestZone", domain_name="example.com")`, actual)
}
