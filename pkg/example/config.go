package example

import (
	"fmt"
	"strconv"

	"github.com/drornir/factor3/pkg/config"
)

//go:generate easytags $GOFILE json:snake,yaml:snake

type Config struct {
	Version string `json:"version" yaml:"version"` // v0
	Log     Log    `flag:"log" json:"log" yaml:"log"`
	Github  Github `json:"github,omitempty" yaml:"github,omitempty"`
	String  string `json:"string" yaml:"string"`
}

type Log struct {
	Level  string `flag:"level" json:"level" yaml:"level"`
	Format string `flag:"format" json:"format" yaml:"format"`
}

type Github struct {
	Token config.SecretString `json:"token,omitempty" yaml:"token,omitempty"`
	App   GithubApp           `json:"app,omitempty" yaml:"app,omitempty"`
}

type GithubApp struct {
	ClientID       string `json:"client_id" yaml:"client_id"`
	PemFile        string `json:"pem_file" yaml:"pem_file"`
	InstallationID string `json:"installation_id" yaml:"installation_id"`
}

func (a GithubApp) InstallationIDMustInt64() int64 {
	i, err := strconv.ParseInt(a.InstallationID, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("github installation id %q is not parsable as integer: %s", a.InstallationID, err.Error()))
	}
	return i
}
