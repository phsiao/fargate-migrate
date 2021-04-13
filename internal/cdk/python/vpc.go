package python

import (
	"bytes"
	"text/template"
)

type LookupVPCStatementGenerator struct {
	Name  string
	VPCID string
}

func NewLookupVPCStatementGenerator(vpcID string) *LookupVPCStatementGenerator {
	rval := LookupVPCStatementGenerator{
		Name:  "ExistingVPC",
		VPCID: vpcID,
	}

	return &rval
}

func (g LookupVPCStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("vpc").Parse(`ec2.Vpc.from_lookup(self, "{{.Name}}", vpc_id="{{.VPCID}}")`)
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

var _ PythonCodeSnippetGenerator = LookupVPCStatementGenerator{}

type ManagedVPCStatementGenerator struct {
	Name   string
	MaxAZs int
	CIDR   string
}

func NewManagedVPCStatementGenerator() *ManagedVPCStatementGenerator {
	rval := ManagedVPCStatementGenerator{
		Name:   "ManagedVPC",
		MaxAZs: 2,
		CIDR:   "10.0.0.0/26",
	}

	return &rval
}

func (g ManagedVPCStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("vpc").Parse(`ec2.Vpc(self, "{{.Name}}", cidr="{{.CIDR}}", max_azs={{.MaxAZs}})`)
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
