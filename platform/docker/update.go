package docker

import (
	"errors"
	"io"

	yaml "gopkg.in/yaml.v2"

	"github.com/ContainerSolutions/flux"
)

// UpdatePodController takes the body of a ReplicationController or Deployment
// resource definition (specified in YAML) and the name of the new image that
// should be put in the definition (in the format "repo.org/group/name:tag"). It
// returns a new resource definition body where all references to the old image
// have been replaced with the new one.
//
// This function has many additional requirements that are likely in flux. Read
// the source to learn about them.
func UpdatePodController(def []byte, newImageID flux.ImageID, trace io.Writer) (ret []byte, err error) {
	var mc minimalCompose
	err = yaml.Unmarshal(def, &mc)
	if err != nil {
		return
	}
	err = tryUpdate(&mc, newImageID, trace)
	if err != nil {
		return
	}
	ret, err = yaml.Marshal(mc)
	return
}

func tryUpdate(mc *minimalCompose, newImage flux.ImageID, trace io.Writer) (err error) {
	for _, mcv := range mc.Services {
		var m map[interface{}]interface{}
		var ok bool
		if m, ok = mcv.(map[interface{}]interface{}); !ok {
			return errors.New("Possible malformed yaml")
		}
		for mk, _ := range m {
			if i, ok := mk.(string); ok {
				if i == "image" {
					m[mk] = newImage.FullID()
					mcv = m
				}
			}
		}
	}
	return
}
