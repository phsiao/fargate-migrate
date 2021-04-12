package python

import (
	"bytes"
	"regexp"
	"text/template"

	corev1 "k8s.io/api/core/v1"
)

type TaskDefinitionStatementGenerator struct {
	Name       string
	Containers []corev1.Container
}

func NewTaskDefinitionStatementGenerator(containers []corev1.Container) *TaskDefinitionStatementGenerator {
	rval := TaskDefinitionStatementGenerator{
		Name:       "TaskDefinition",
		Containers: containers,
	}

	return &rval
}

func (g TaskDefinitionStatementGenerator) Generate() (string, error) {
	tmpl, err := template.New("taskdefinition").Funcs(
		template.FuncMap{
			"python": func(input string) string {
				re := regexp.MustCompile(`\n`)
				return re.ReplaceAllString(input, "\\n")
			},
		}).Parse(`
        task_def = ecs.FargateTaskDefinition(self, "{{.Name}}",
            cpu=256,
            memory_limit_mib=512)
{{range .Containers}}
        env = {}
{{range .Env}}
        env["{{.Name}}"] = "{{.Value|python}}"{{end}}

        task_def.add_container("{{.Name}}",
            image=ecs.ContainerImage.from_asset("docker/{{.Name}}"),
            port_mappings=[
{{range .Ports}}
              ecs.PortMapping(container_port={{.ContainerPort}}),{{end}}
            ],
            environment=env,
            logging=ecs.AwsLogDriver(stream_prefix="{{.Name}}")
        ){{end}}
`,
	)
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

var _ PythonCodeSnippetGenerator = TaskDefinitionStatementGenerator{}
