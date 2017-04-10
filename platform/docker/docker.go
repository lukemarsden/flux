package docker

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"os"
	"os/exec"

	"github.com/weaveworks/flux/platform"

	"golang.org/x/net/context"
)

type Swarm struct {
	client *client.Client
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
	bin, err := findBinary("docker")
	cmd := exec.Command(bin, "deploy -c docker-compose.yml dockerswarm")
	err = cmd.Run()
	return err
}

func (c *Swarm) Sync(spec platform.SyncDef) error {
	return nil
}

func (c *Swarm) Ping() error {
	return nil
}

func (c *Swarm) Export() ([]byte, error) {
	return nil, nil
}

func (c *Swarm) Version() (string, error) {
	ctx := context.Background()
	version, err := c.client.ServerVersion(ctx)
	return version.Version, err
}

func findBinary(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	localBin := filepath.Join(cwd, name)
	if _, err := os.Stat(localBin); err == nil {
		return localBin, nil
	}
	if pathBin, err := exec.LookPath(name); err == nil {
		return pathBin, nil
	}
	return "", fmt.Errorf("%s not found", name)
}
