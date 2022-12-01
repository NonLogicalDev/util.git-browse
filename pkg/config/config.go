package config

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	"golang.org/x/xerrors"
	"strings"
	"text/template"
)

// Configuration Format:
//
// browse.default
//
// browse.service.<name>.url-file
// browse.service.<name>.url-dir
//
// browse.service.<name>.base
// browse.service.<name>.repo
//
// browse.service.<name>.<field>

type Config struct {
	Default string // browse.default

	Services map[string]*ConfigService // browse.service.<service-name>
}

const cfgDumpTemplateFile = `
[browse]
	default = {{toJson .Default}}
{{range $svc, $svcCfg := .Services }}
[browse "service.{{$svc}}"]
	base = {{toJson $svcCfg.Base}}
	repo = {{toJson $svcCfg.Repo}}
	url-file = {{toJson $svcCfg.URLTemplateFile}}
	url-dir = {{toJson $svcCfg.URLTemplateDir}}
	{{range $key, $val := $svcCfg.Fields }}{{$key}} = {{toJson $val}}
{{end}}{{end}}
`

const cfgDumpTemplateCmds = `
git config browse.default {{quote .Default}}
{{range $svc, $svcCfg := .Services }}
git config browse.service.{{$svc}}.base {{quote $svcCfg.Base}}
git config browse.service.{{$svc}}.repo {{quote $svcCfg.Repo}}
git config browse.service.{{$svc}}.url-file {{quote $svcCfg.URLTemplateFile}}
git config browse.service.{{$svc}}.url-file {{quote $svcCfg.URLTemplateDir}}
{{range $key, $val := $svcCfg.Fields }}{{$key}} = {{quote $val}}
git config browse.service.{{$svc}}.{{$key}} {{quote $val}}
{{end}}{{end}}
`

type DumpType int

const (
	DumpGitconfig DumpType = iota
	DumpCli
	DumpJson
)

func (cfg Config) Dump(typ DumpType) string {
	switch typ {
	case DumpGitconfig:
		b := new(strings.Builder)
		tpl := template.Must(
			template.New("cfg").
				Funcs(sprig.TxtFuncMap()).
				Parse(strings.Trim(cfgDumpTemplateFile, "\n")),
		)
		if err := tpl.Execute(b, cfg); err != nil {
			panic(err)
		}
		return b.String()
	case DumpCli:
		b := new(strings.Builder)
		tpl := template.Must(
			template.New("cfg").
				Funcs(sprig.TxtFuncMap()).
				Parse(strings.Trim(cfgDumpTemplateCmds, "\n")),
		)
		if err := tpl.Execute(b, cfg); err != nil {
			panic(err)
		}
		return b.String()
	case DumpJson:
		out, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			panic(err)
		}
		return string(out)
	default:
		return ""
	}
}

type ConfigService struct {
	Base string // browse.service.<name>.base
	Repo string // browse.service.<name>.repo
	Ref  string // browse.service.<name>.ref

	URLTemplateFile string // browse.service.<name>.url-file
	URLTemplateDir  string // browse.service.<name>.url-dir

	// Fields fields.
	Fields map[string]string // browse.service.<service-name>.<field>
}

type URLOpt struct {
	Service string
	Ref     string
	Path    string

	IsDir bool

	LineS int
	LineE int
}
type URLOption func()

func (cfg Config) URLFor(o URLOpt) (string, error) {
	targetService := cfg.Default
	if len(o.Service) > 0 {
		targetService = o.Service
	}

	serviceConfig, ok := cfg.Services[targetService]
	if !ok {
		return "", xerrors.Errorf("service(%v): not configured", targetService)
	}

	targetRef := serviceConfig.Ref
	if len(o.Ref) > 0 {
		targetRef = o.Ref
	}

	var templateRaw string
	if o.IsDir {
		templateRaw = serviceConfig.URLTemplateDir
	} else {
		templateRaw = serviceConfig.URLTemplateFile
	}

	tpl, err := template.New("url").
		Funcs(sprig.TxtFuncMap()).
		Parse(templateRaw)
	if err != nil {
		return "", xerrors.Errorf("service(%v): failed parsing template '%v': %w", targetService, templateRaw, err)
	}

	c := templateURLContext{
		WithRange: o.LineS != 0 && o.LineE != 0,
		WithLine:  o.LineS != 0 && o.LineE == 0,

		Path:  o.Path,
		LineS: o.LineS,
		LineE: o.LineE,

		Base: serviceConfig.Base,
		Repo: serviceConfig.Repo,
		Ref:  targetRef,

		Fields: serviceConfig.Fields,
	}
	return c.Render(tpl)
}

func ParseGitConfig(gitConfig string) (Config, error) {
	cfg := Config{}
	InitDefaultConfig(&cfg)

	for _, line := range strings.Split(gitConfig, "\n") {
		if !strings.HasPrefix(line, "browse.") {
			continue
		}

		line = strings.TrimPrefix(line, "browse.")
		var (
			keyFull, value = parseConfigLine(line)
			keyParts, _    = parseConfigKey(keyFull, 3)
			keyTop         = keyParts[0]
			keyMid         = keyParts[1]
			keySub         = keyParts[2]
		)

		switch keyTop {
		case "default":
			cfg.Default = value
		case "service":
			if len(keySub) > 0 {
				serviceConfig, ok := cfg.Services[keyMid]
				if !ok {
					serviceConfig = &ConfigService{
						Fields: map[string]string{},
					}
				}

				switch keySub {
				case "url-file":
					serviceConfig.URLTemplateFile = value
				case "url-dir":
					serviceConfig.URLTemplateDir = value
				case "base", "domain":
					serviceConfig.Base = value
				case "repo":
					serviceConfig.Repo = value
				default:
					serviceConfig.Fields[keySub] = value
				}

				cfg.Services[keyMid] = serviceConfig
			}
		}
	}

	return cfg, nil
}
