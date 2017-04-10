package docker

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestFindDefinedServices(t *testing.T) {
	ss, err := FindDefinedServices("../../")
	if err != nil {
		t.Error(err)
	}
	spew.Dump(ss)
}
