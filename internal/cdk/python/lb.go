package python

import (
	"bytes"
	"text/template"
)

type HealthCheckStatementGenerator struct {
	Name string
	Path string
}

func NewHealthCheckStatementGenerator(path string) *HealthCheckStatementGenerator {
	rval := HealthCheckStatementGenerator{
		Name: "HealthCheck",
		Path: path,
	}

	return &rval
}

func (g HealthCheckStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("healthcheck").Parse(`svc.target_group.configure_health_check(path="{{.Path}}")`)
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

var _ PythonCodeSnippetGenerator = HealthCheckStatementGenerator{}
