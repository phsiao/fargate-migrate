package python

import (
	"bytes"
	"text/template"
)

type HostedZoneStatementGenerator struct {
	Name       string
	DomainName string
}

func NewHostedZoneStatementGenerator(domainName string) *HostedZoneStatementGenerator {
	rval := HostedZoneStatementGenerator{
		Name:       "HostedZone",
		DomainName: domainName,
	}

	return &rval
}

func (g HostedZoneStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("zone").Parse(`route53.HostedZone.from_lookup(self, "{{.Name}}", domain_name="{{.DomainName}}")`)
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

var _ PythonCodeSnippetGenerator = HostedZoneStatementGenerator{}
