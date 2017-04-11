package main

import (
	"github.com/ContainerSolutions/flux"
)

func parseServiceOption(s string) (flux.ServiceSpec, error) {
	if s == "" {
		return flux.ServiceSpecAll, nil
	}
	return flux.ParseServiceSpec(s)
}
