package docker

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/log"
)

func TestGetStacks(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}
	swarm, err := NewSwarm(logger)
	if err != nil {
		t.Error(err)
	}
	s, err := swarm.GetStacks()
	if err != nil {
		t.Error(err)
	}
	spew.Dump(s)

}
