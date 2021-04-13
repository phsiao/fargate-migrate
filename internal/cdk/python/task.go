package python

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/phsiao/fargate-migrate/internal/fargate"
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
	accumCPU := int(0)
	accumMemory := int(0)

	for _, container := range g.Containers {
		accumCPU += int(container.Resources.Requests.Cpu().MilliValue())
		accumMemory += int(container.Resources.Requests.Memory().MilliValue() / (1000 * 1024 * 1024))
	}

	configCPU, configMemory := fargate.MinCPUMemroyConfiguration(accumCPU, accumMemory)

	tmpl, err := template.New("taskdefinition").Funcs(
		template.FuncMap{
			"python": func(input string) string {
				re := regexp.MustCompile(`\n`)
				return re.ReplaceAllString(input, "\\n")
			},
			"configCPU": func() string {
				return fmt.Sprintf("%d", configCPU)
			},
			"configMemory": func() string {
				return fmt.Sprintf("%d", configMemory)
			},
		}).Parse(`
        task_def = ecs.FargateTaskDefinition(self, "{{.Name}}",
            cpu={{configCPU}},
            memory_limit_mib={{configMemory}})
{{range .Containers}}
        env = {}{{range .Env}}
        env["{{.Name}}"] = "{{.Value|python}}"{{end}}

        task_def.add_container("{{.Name}}",
            image=ecs.ContainerImage.from_asset("docker/{{.Name}}"),
            port_mappings=[{{range .Ports}}
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
