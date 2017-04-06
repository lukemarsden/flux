package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ContainerSolutions/flux/platform"
	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/weaveworks/flux"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cl, _ := client.NewEnvClient()
	s, err := cl.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	pss := make([]platform.Service, len(s))
	for k, v := range s {
		args := filters.NewArgs()
		args.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%v", v.Spec.Annotations.Name))
		cs, err := cl.ContainerList(ctx, types.ContainerListOptions{Filters: args})
		spew.Dump(cs)
		if err != nil {
			log.Fatal(err)
		}
		ps := platform.Service{
			ID:         flux.MakeServiceID("default-swarm", v.Spec.Annotations.Name),
			IP:         "?",
			Metadata:   v.Spec.Annotations.Labels,
			Status:     string(v.UpdateStatus.State),
			Containers: platform.ContainersOrExcuse{},
		}
		pcs := make([]platform.Container, len(cs))
		for k, v := range cs {
			pcs[k] = platform.Container{
				Name:  v.Names[0],
				Image: v.Image,
			}
		}
		ps.Containers.Containers = pcs
		pss[k] = ps
	}
	spew.Dump(pss)
}
