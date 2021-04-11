package python

import "fmt"

type ManagedVPCStatementGenerator struct {
	vpcName string
	maxAZs  int
}

func NewManagedVPCStatementGenerator() *ManagedVPCStatementGenerator {
	rval := ManagedVPCStatementGenerator{
		vpcName: "ManagedVPC",
		maxAZs:  2,
	}

	return &rval
}

func (g ManagedVPCStatementGenerator) Generate() ([]string, error) {
	rval := []string{
		fmt.Sprintf("ec2.Vpc(self, \"%s\", max_azs=%d)", g.vpcName, g.maxAZs),
	}
	return rval, nil
}

var _ PythonCodeSnippetGenerator = ManagedVPCStatementGenerator{}
