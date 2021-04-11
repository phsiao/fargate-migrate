package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostedZoneStatementGenerator(t *testing.T) {
	g := HostedZoneStatementGenerator{
		zoneName:   "TestZone",
		domainName: "example.com",
	}
	actual, err := g.Generate()
	assert.NoError(t, err)
	assert.Equal(t, []string{"route53.HostedZone.from_lookup(self, \"TestZone\", domain_name=\"example.com\")"}, actual)
}
