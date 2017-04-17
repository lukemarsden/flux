package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/ContainerSolutions/flux"
	"github.com/ContainerSolutions/flux/platform"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
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
			ID:       flux.MakeServiceID("default_swarm", v.Spec.Networks[0].Aliases[0]),
			IP:       "?",
			Metadata: v.Spec.Annotations.Labels,
			//			Status:     string(v.UpdateStatus.State),
			Containers: platform.ContainersOrExcuse{},
		}
		if ignore.Contains(ps.ID) {
			continue
		}
		args := filters.NewArgs()
		args.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%v", v.Spec.Annotations.Name))
		/// OLD
		cs, err := c.client.ContainerList(ctx, types.ContainerListOptions{Filters: args})
		if err != nil {
			return pss, err
		}
		oldPcs := make([]platform.Container, len(cs))
		for k, v := range cs {
			oldPcs[k] = platform.Container{
				Name:  v.Labels["com.docker.swarm.task.name"],
				Image: v.Image,
			}
		}
		/// NEW
		ts, err := c.client.TaskList(ctx, types.TaskListOptions{Filters: args})
		if err != nil {
			return pss, err
		}
		pcs := make([]platform.Container, len(ts))
		for k, t := range ts {
			pcs[k] = platform.Container{
				Name:  t.Name,
				Image: t.Spec.ContainerSpec.Image,
			}
		}

		/// LOG
		fmt.Printf(
			"[AllServices] %s: OLD: %v, NEW: %v\n",
			v.ServiceSpec.Name, oldPcs, pcs,
		)

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
		args.Add("name", fmt.Sprintf("%s_%s", namespace, n))
	}
	s, err := c.client.ServiceList(ctx, types.ServiceListOptions{args})

	pss := make([]platform.Service, 0)
	if err != nil {
		return pss, err
	}

	// Filter out excessive services since ServiceList doesn't match explicitly
	d := make([]swarm.Service, 0)
	for _, v := range s {
		for _, k := range ids {
			_, n := k.Components()
			if n == v.Spec.Networks[0].Aliases[0] {
				d = append(d, v)
			}
		}
	}

	for _, v := range d {
		ps := platform.Service{
			ID:       flux.MakeServiceID(namespace, v.Spec.Networks[0].Aliases[0]),
			IP:       "?",
			Metadata: v.Spec.Annotations.Labels,
			//Status:     string(v.UpdateStatus.State),
			Containers: platform.ContainersOrExcuse{},
		}
		args := filters.NewArgs()
		args.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%v", v.Spec.Annotations.Name))
		/// OLD
		cs, err := c.client.ContainerList(ctx, types.ContainerListOptions{Filters: args})
		if err != nil {
			return pss, err
		}
		oldPcs := make([]platform.Container, len(cs))
		for k, v := range cs {
			oldPcs[k] = platform.Container{
				Name:  v.Names[0],
				Image: v.Image,
			}
		}
		/// NEW
		ts, err := c.client.TaskList(ctx, types.TaskListOptions{Filters: args})
		if err != nil {
			return pss, err
		}
		pcs := make([]platform.Container, len(ts))
		for k, t := range ts {
			pcs[k] = platform.Container{
				Name:  t.Name,
				Image: t.Spec.ContainerSpec.Image,
			}
		}
		/// LOG
		fmt.Printf(
			"[SomeServices] %s: OLD: %v, NEW: %v\n",
			v.ServiceSpec.Name, oldPcs, pcs,
		)

		ps.Containers.Containers = pcs
		pss = append(pss, ps)
	}
	return pss, nil
}
