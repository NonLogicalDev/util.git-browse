package config

import (
	"strings"
	"text/template"
)

// TODO: https://github.com/whilp/git-urls/blob/master/urls.go

type templateConfig struct {
	urlDir  string
	urlFile string

	base string
	ref  string
}

type templateURLContext struct {
	WithRange bool
	WithLine  bool

	Path  string
	LineS int
	LineE int

	Base string
	Repo string
	Ref  string

	Fields map[string]string
}

var _defaultServiceTemplate = map[string]templateConfig{
	"sourcegraph": {
		urlDir:  `{{.Base}}/{{.Repo}}/{{.Ref}}/tree/{{.Path}}`,
		urlFile: `{{.Base}}/{{.Repo}}/{{.Ref}}/blob/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-{{.LineE}}{{end}}`,

		base: "https://sourcegraph.com",
		ref:  "-",
	},
	"github": {
		urlDir:  `{{.Base}}/{{.Repo}}/tree/{{.Ref}}/{{.Path}}`,
		urlFile: `{{.Base}}/{{.Repo}}/blob/{{.Ref}}/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-L{{.LineE}}{{end}}`,

		base: "https://github.com",
		ref:  "HEAD",
	},
}

func (c templateURLContext) Render(tpl *template.Template) (string, error) {
	b := new(strings.Builder)
	err := tpl.Execute(b, c)
	return b.String(), err
}

func InitDefaultConfig(cfg *Config) {
	cfg.Default = "github"
	cfg.Services = map[string]*ConfigService{}

	for svcName, svcConfig := range _defaultServiceTemplate {
		serviceConfig, ok := cfg.Services[svcName]
		if !ok {
			serviceConfig = &ConfigService{
				Fields: map[string]string{},
			}
		}

		serviceConfig.URLTemplateFile = svcConfig.urlFile
		serviceConfig.URLTemplateDir = svcConfig.urlDir
		serviceConfig.Base = svcConfig.base
		serviceConfig.Ref = svcConfig.ref

		cfg.Services[svcName] = serviceConfig
	}
}
