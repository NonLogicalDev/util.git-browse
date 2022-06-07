package config

// TODO: https://github.com/whilp/git-urls/blob/master/urls.go

type templateConfig struct {
	urlDir  string
	urlFile string

	base string
}

type templateURLContext struct {
	WithRange bool
	WithLine  bool

	Path  string
	LineS int
	LineE int

	Base string
	Repo string

	Fields map[string]string
}

var _defaultServiceTemplate = map[string]templateConfig{
	"sourcegraph": {
		urlDir:  `{{.Base}}/{{.Repo}}/-/tree/{{.Path}}`,
		urlFile: `{{.Base}}/{{.Repo}}/-/blob/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-{{.LineE}}{{end}}`,

		base: "https://sourcegraph.com",
	},
	"github": {
		urlDir:  `{{.Base}}/{{.Repo}}/tree/HEAD/{{.Path}}`,
		urlFile: `{{.Base}}/{{.Repo}}/blob/HEAD/{{.Path}}{{if .WithLine}}#L{{.LineS}}{{else if .WithRange}}#L{{.LineS}}-L{{.LineE}}{{end}}`,

		base: "https://github.com",
	},
}

func InitDefaultConfig(cfg *Config) {
	cfg.Default = "github"
	cfg.Services = map[string]ConfigService{}

	for svcName, svcConfig := range _defaultServiceTemplate {
		serviceConfig, ok := cfg.Services[svcName]
		if !ok {
			serviceConfig = ConfigService{
				Fields: map[string]string{},
			}
		}

		serviceConfig.URLTemplateFile = svcConfig.urlFile
		serviceConfig.Fields["url-file"] = svcConfig.urlFile

		serviceConfig.URLTemplateDir = svcConfig.urlDir
		serviceConfig.Fields["url-dir"] = svcConfig.urlDir

		serviceConfig.Base = svcConfig.base
		serviceConfig.Fields["domain"] = svcConfig.base

		cfg.Services[svcName] = serviceConfig
	}
}
