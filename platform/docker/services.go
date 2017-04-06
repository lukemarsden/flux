package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/weaveworks/flux"
	"github.com/weaveworks/flux/platform"
)

func (c *Swarm) AllServices(namespace string, ignore flux.ServiceIDSet) ([]platform.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := c.client.ServiceList(ctx, types.ServiceListOptions{})
	pss := make([]platform.Service, len(s))
	if err != nil {
		return pss, err
	}
	for k, v := range s {
		ps := platform.Service{
			ID:         flux.MakeServiceID(namespace, v.Spec.Annotations.Name),
			IP:         "?",
			Metadata:   v.Spec.Annotations.Labels,
			Status:     string(v.UpdateStatus.State),
			Containers: platform.ContainersOrExcuse{},
		}
		if ignore.Contains(ps.ID) {
			continue
		}
		args := filters.NewArgs()
		args.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%v", v.Spec.Annotations.Name))
		cs, err := c.client.ContainerList(ctx, types.ContainerListOptions{Filters: args})
		if err != nil {
			return pss, err
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
	return pss, nil
}

func (c *Swarm) SomeServices(ids []flux.ServiceID) (res []platform.Service, err error) {
	namespace := "default_swarm"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := filters.NewArgs()
	for _, v := range ids {
		_, n := v.Components()
		args.Add("name", n)
	}
	s, err := c.client.ServiceList(ctx, types.ServiceListOptions{args})
	pss := make([]platform.Service, len(s))
	if err != nil {
		return pss, err
	}
	for k, v := range s {
		ps := platform.Service{
			ID:         flux.MakeServiceID(namespace, v.Spec.Annotations.Name),
			IP:         "?",
			Metadata:   v.Spec.Annotations.Labels,
			Status:     string(v.UpdateStatus.State),
			Containers: platform.ContainersOrExcuse{},
		}
		args := filters.NewArgs()
		args.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%v", v.Spec.Annotations.Name))
		cs, err := c.client.ContainerList(ctx, types.ContainerListOptions{Filters: args})
		if err != nil {
			return pss, err
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
	return pss, nil
}
