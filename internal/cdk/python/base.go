package python

type PythonCodeSnippetGenerator interface {
	Generate() ([]string, error)
}
