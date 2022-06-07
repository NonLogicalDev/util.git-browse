package config

import (
	"EXP/pkg/gitutils"
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
// browse.service.<name>.domain
// browse.service.<name>.repo
//
// browse.service.<name>.x-<field>

type Config struct {
	Default string // browse.default
	//Service string // browse.service

	Services map[string]ConfigService // browse.service.<service-name>
}

type ConfigService struct {
	Base string // browse.service.<service-name>.base
	Repo string // browse.service.<service-name>.repo

	URLTemplateFile string // browse.service.<service-name>.url-file
	URLTemplateDir  string // browse.service.<service-name>.url-dir

	// Fields fields.
	Fields map[string]string // browse.service.<service-name>.<field>
}

// TODO: ADD SPRIG Templating Lib
func (cfg Config) GetURL(service string, path string, isDir bool, lineS, lineE int) (string, error) {
	serviceConfig := cfg.Services[cfg.Default]
	//if len(cfg.Service) > 0 {
	//	serviceConfig = cfg.Services[cfg.Service]
	//}
	if len(service) > 0 {
		serviceConfig = cfg.Services[service]
	}

	var templateRaw string
	if isDir {
		templateRaw = serviceConfig.URLTemplateDir
	} else {
		templateRaw = serviceConfig.URLTemplateFile
	}

	tpl, err := template.New("url").Parse(templateRaw)
	if err != nil {
		return "", err
	}

	b := new(strings.Builder)
	err = tpl.Execute(b, templateURLContext{
		WithRange: lineS != 0 && lineE != 0,
		WithLine:  lineS != 0 && lineE == 0,

		Path:  path,
		LineS: lineS,
		LineE: lineE,

		Base:   serviceConfig.Base,
		Repo:   serviceConfig.Repo,
		Fields: serviceConfig.Fields,
	})
	return b.String(), err
}

func LoadConfig() (Config, error) {
	output, err := gitutils.GitExec("config", "--get-regexp", "browse\\..*")
	if err != nil {
		return Config{}, err
	}

	return ParseRawConfig(output)
}

func ParseRawConfig(rawConfig string) (Config, error) {
	cfg := Config{}
	InitDefaultConfig(&cfg)

	for _, line := range strings.Split(rawConfig, "\n") {
		if !strings.HasPrefix(line, "browse.") {
			continue
		}

		line = strings.TrimPrefix(line, "browse.")
		var (
			keyFull, value = parseConfigLine(line)
			keyParts, _    = parseConfigKey(keyFull, 2)
			keyTop         = keyParts[0]
			keySub         = keyParts[1]
		)

		switch keyTop {
		case "default":
			cfg.Default = value
		}

		if len(keySub) > 0 {
			serviceConfig, ok := cfg.Services[keyTop]
			if !ok {
				serviceConfig = ConfigService{
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

			cfg.Services[keyTop] = serviceConfig
		}
	}

	return cfg, nil
}
