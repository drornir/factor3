package config

import "fmt"

type EnvSchema struct {
	Usage map[string]string
}

func (eu EnvSchema) ShellUsage() string {
	if len(eu.Usage) == 0 {
		return ""
	}

	var usage []byte
	for k, v := range eu.Usage {
		usage = append(usage, []byte(fmt.Sprintf("%s=%s\n", k, v))...)
	}

	return string(usage)
}
