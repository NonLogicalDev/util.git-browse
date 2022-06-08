package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Config{}
	InitDefaultConfig(&cfg)

	cfg.Services["github"].Repo = "nonlogicaldev/cli.photo.philter"

	url, err := cfg.URLFor(URLOpt{
		Path:  ".",
		IsDir: true,
	})
	assert.NoError(t, err)

	t.Log(url)
}
