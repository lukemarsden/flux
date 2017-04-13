package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

type stack struct {
	// Name is the name of the stack
	Name string
	// Services is the number of the services
	Services int
}

const (
	// LabelNamespace is the label used to track stack resources
	LabelNamespace = "com.docker.stack.namespace"
)

func (c *Swarm) GetStacks() ([]*stack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services, err := c.client.ServiceList(
		ctx,
		types.ServiceListOptions{},
	)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*stack, 0)
	for _, service := range services {
		labels := service.Spec.Labels
		name, ok := labels[LabelNamespace]
		if !ok {
			return nil, errors.Errorf("cannot get label %s for service %s",
				LabelNamespace, service.ID)
		}
		ztack, ok := m[name]
		if !ok {
			m[name] = &stack{
				Name:     name,
				Services: 1,
			}
		} else {
			ztack.Services++
		}
	}
	var stacks []*stack
	for _, stack := range m {
		stacks = append(stacks, stack)
	}
	return stacks, nil
}
