package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Config{}
	InitDefaultConfig(&cfg)

	cfg.Services["github"].Repo = "github/github"

	url, err := cfg.URLFor(URLOpt{
		Path:  ".",
		IsDir: true,
	})
	assert.NoError(t, err)

	t.Log(url)
}

func Test1(t *testing.T) {
	fmt.Errorf("hi")
}
