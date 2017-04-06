package docker

import (
	"log"

	"github.com/docker/docker/client"

	"github.com/weaveworks/flux/platform"
)

type Swarm struct {
	client
	logger log.Logger
}

func NewSwarm(logger log.Logger) (*Swarm, error) {
	cli, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	c := &Swarm{
		client: cli,
		logger: logger,
	}

	return c, nil
}

func (c *Swarm) Apply(defs []platform.ServiceDefinition) error {

}

func (c *Swarm) Ping() error {

}

func (c *Swarm) Version() (string, error) {
	return c.Version
}
