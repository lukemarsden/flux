package docker

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ContainerSolutions/flux"
	yaml "gopkg.in/yaml.v2"
)

func TestTryUpdate(t *testing.T) {
	fc, err := ioutil.ReadFile("../../demo/docker-compose-carts-db.yml")
	if err != nil {
		t.Error(err)
	}
	var def minimalCompose
	err = yaml.Unmarshal(fc, &def)
	if err != nil {
		t.Error(err)
	}
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	err = tryUpdate(&def, flux.ImageID{"test", "test", "test", "test"}, foo)
	if err != nil {
		t.Error(err)
	}
	var f []byte
	f, err = yaml.Marshal(def)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(string(f), "image: test/test/test:test") {
		t.Error("Expected image 'image: test/test/test:test'")
	}
}
