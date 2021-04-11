package python

import "fmt"

type HostedZoneStatementGenerator struct {
	zoneName   string
	domainName string
}

func NewHostedZoneStatementGenerator(domainName string) *HostedZoneStatementGenerator {
	rval := HostedZoneStatementGenerator{
		zoneName:   "HostedZone",
		domainName: domainName,
	}

	return &rval
}

func (g HostedZoneStatementGenerator) Generate() ([]string, error) {
	rval := []string{
		fmt.Sprintf("route53.HostedZone.from_lookup(self, \"%s\", domain_name=\"%s\")",
			g.zoneName, g.domainName,
		),
	}
	return rval, nil
}

var _ PythonCodeSnippetGenerator = HostedZoneStatementGenerator{}
