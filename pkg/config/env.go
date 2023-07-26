package config

import "fmt"

type EnvSchema struct {
	Usage EnvUsage
}

type EnvUsage struct {
	shell map[string]string
}

func (eu EnvUsage) Shell() string {
	if len(eu.shell) == 0 {
		return ""
	}

	var usage []byte
	for k, v := range eu.shell {
		usage = append(usage, []byte(fmt.Sprintf("%s=%s\n", k, v))...)
	}

	return string(usage)
}
