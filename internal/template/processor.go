package template

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/go-envparse"
	"github.com/tomdoesdev/knox/internal/secrets"
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

const DefaultTemplateFileName = ".env.template"

type Processor struct {
	secretStore secrets.SecretReader
}

func NewProcessor(secretStore secrets.SecretReader) *Processor {
	return &Processor{
		secretStore: secretStore,
	}
}

func (p *Processor) ProcessFile(filePath string) (map[string]string, error) {
	if filePath == "" {
		filePath = DefaultTemplateFileName
	}

	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errs.Wrap(err, TemplateFileNotFoundCode, "failed to get current working directory")
		}
		filePath = filepath.Join(cwd, filePath)
	}

	exists, err := fs.IsExist(filePath)
	if err != nil {
		return nil, errs.Wrap(err, TemplateFileNotFoundCode, "failed to check if template file exists")
	}

	if !exists {
		return nil, ErrTemplateFileNotFound
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err, TemplateReadFailureCode, "failed to read template file")
	}

	return p.ProcessTemplate(string(content), filePath)
}

func (p *Processor) ProcessFileOptional(filePath string) (map[string]string, error) {
	envVars, err := p.ProcessFile(filePath)
	if err != nil {
		if errs.Is(err, TemplateFileNotFoundCode) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	return envVars, nil
}

func (p *Processor) ProcessTemplate(templateContent, name string) (map[string]string, error) {
	ctx := NewContext(p.secretStore)

	tmpl, err := template.New(name).Funcs(ctx.FuncMap()).Parse(templateContent)
	if err != nil {
		return nil, errs.Wrap(err, TemplateParseFailureCode, "failed to parse template")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ctx)
	if err != nil {
		return nil, errs.Wrap(err, TemplateExecuteFailureCode, "failed to execute template")
	}

	return p.parseEnvironmentVariables(buf.String())
}

func (p *Processor) parseEnvironmentVariables(content string) (map[string]string, error) {
	env, err := envparse.Parse(strings.NewReader(content))
	if err != nil {
		return nil, errs.Wrap(err, TemplateParseFailureCode, "failed to parse environment variables")
	}
	return env, nil
}
