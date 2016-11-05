// Package registry provides domain abstractions over container registries.
package registry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	dockerregistry "github.com/heroku/docker-registry-client/registry"
	"golang.org/x/net/publicsuffix"

	"github.com/weaveworks/flux"
)

const (
	dockerHubHost    = "index.docker.io"
	dockerHubLibrary = "library"
)

type creds struct {
	username, password string
}

// Credentials to a (Docker) registry.
type Credentials struct {
	m map[string]creds
}

// Client is a handle to a registry.
type Client struct {
	Credentials Credentials
	Logger      log.Logger
}

type roundtripperFunc func(*http.Request) (*http.Response, error)

func (f roundtripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// GetRepository yields a repository matching the given name, if any exists.
// Repository may be of various forms, in which case omitted elements take
// assumed defaults.
//
//   helloworld             -> index.docker.io/library/helloworld
//   foo/helloworld         -> index.docker.io/foo/helloworld
//   quay.io/foo/helloworld -> quay.io/foo/helloworld
//
func (c *Client) GetRepository(repository string) ([]flux.ImageDescription, error) {
	var host, org, image string
	parts := strings.Split(repository, "/")
	switch len(parts) {
	case 1:
		host = dockerHubHost
		org = dockerHubLibrary
		image = parts[0]
	case 2:
		host = dockerHubHost
		org = parts[0]
		image = parts[1]
	case 3:
		host = parts[0]
		org = parts[1]
		image = parts[2]
	default:
		return nil, fmt.Errorf(`expected image name as either "<host>/<org>/<image>", "<org>/<image>", or "<image>"`)
	}

	hostlessImageName := fmt.Sprintf("%s/%s", org, image)
	httphost := "http://" + host

	// quay.io wants us to use cookies for authorisation, so we have
	// to construct one (the default client has none). This means a
	// bit more constructing things to be able to make a registry
	// client literal, rather than calling .New()
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	auth := c.Credentials.credsFor(host)

	// A context we'll use to cancel requests on error
	ctx, cancel := context.WithCancel(context.Background())

	// Use the wrapper to fix headers for quay.io, and remember bearer tokens
	var transport http.RoundTripper = &wwwAuthenticateFixer{transport: http.DefaultTransport}
	// Now the auth-handling wrappers that come with the library
	transport = dockerregistry.WrapTransport(transport, httphost, auth.username, auth.password)

	client := &dockerregistry.Registry{
		URL: httphost,
		Client: &http.Client{
			Transport: roundtripperFunc(func(r *http.Request) (*http.Response, error) {
				return transport.RoundTrip(r.WithContext(ctx))
			}),
			Jar: jar,
		},
		Logf: dockerregistry.Quiet,
	}

	tags, err := client.Tags(hostlessImageName)
	if err != nil {
		cancel()
		return nil, err
	}

	return c.tagsToRepository(cancel, client, repository, tags)
}

func (c *Client) lookupImage(client *dockerregistry.Registry, repoName, tag string) (flux.ImageDescription, error) {
	// Minor cheat: this will work whether the registry (e.g.,
	// quay.io) is part of the repo name given or not.
	id := flux.MakeImageID("", repoName, tag)
	img := flux.ImageDescription{ID: id}

	_, name, _ := id.Components()
	meta, err := client.Manifest(name, tag)
	if err != nil {
		return img, err
	}
	// the manifest includes some v1-backwards-compatibility data,
	// oddly called "History", which are layer metadata as JSON
	// strings; these appear most-recent (i.e., topmost layer) first,
	// so happily we can just decode the first entry to get a created
	// time.
	type v1image struct {
		Created time.Time `json:"created"`
	}
	var topmost v1image
	if err = json.Unmarshal([]byte(meta.History[0].V1Compatibility), &topmost); err == nil {
		img.CreatedAt = topmost.Created
	}

	return img, err
}

func (c *Client) tagsToRepository(cancel func(), client *dockerregistry.Registry, repoName string, tags []string) ([]flux.ImageDescription, error) {
	// one way or another, we'll be finishing all requests
	defer cancel()

	type result struct {
		image flux.ImageDescription
		err   error
	}

	fetched := make(chan result, len(tags))

	for _, tag := range tags {
		go func(t string) {
			img, err := c.lookupImage(client, repoName, t)
			if err != nil {
				c.Logger.Log("registry-metadata-err", err)
			}
			fetched <- result{img, err}
		}(tag)
	}

	images := make([]flux.ImageDescription, cap(fetched))
	for i := 0; i < cap(fetched); i++ {
		res := <-fetched
		if res.err != nil {
			return nil, res.err
		}
		images[i] = res.image
	}

	sort.Sort(byCreatedDesc(images))
	return images, nil
}

// --- Credentials

// NoCredentials returns a usable but empty credentials object.
func NoCredentials() Credentials {
	return Credentials{
		m: map[string]creds{},
	}
}

// CredentialsFromFile returns a credentials object parsed from the given
// filepath.
func CredentialsFromFile(path string) (Credentials, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Credentials{}, err
	}

	type dockerConfig struct {
		Auths map[string]struct {
			Auth  string `json:"auth"`
			Email string `json:"email"`
		} `json:"auths"`
	}

	var config dockerConfig
	if err = json.Unmarshal(bytes, &config); err != nil {
		return Credentials{}, err
	}

	m := map[string]creds{}
	for host, entry := range config.Auths {
		decodedAuth, err := base64.StdEncoding.DecodeString(entry.Auth)
		if err != nil {
			return Credentials{}, err
		}
		authParts := strings.SplitN(string(decodedAuth), ":", 2)
		m[host] = creds{
			username: authParts[0],
			password: authParts[1],
		}
	}
	return Credentials{m: m}, nil
}

func CredentialsFromConfig(config flux.InstanceConfig) (Credentials, error) {
	m := map[string]creds{}
	for host, entry := range config.Registry.Auths {
		decodedAuth, err := base64.StdEncoding.DecodeString(entry.Auth)
		if err != nil {
			return Credentials{}, err
		}
		authParts := strings.SplitN(string(decodedAuth), ":", 2)
		m[host] = creds{
			username: authParts[0],
			password: authParts[1],
		}
	}
	return Credentials{m: m}, nil
}

// For yields an authenticator for a specific host.
func (cs Credentials) credsFor(host string) creds {
	if cred, found := cs.m[host]; found {
		return cred
	}
	if cred, found := cs.m[fmt.Sprintf("https://%s/v1/", host)]; found {
		return cred
	}
	return creds{}
}

// Hosts returns all of the hosts available in these credentials.
func (cs Credentials) Hosts() []string {
	hosts := []string{}
	for host := range cs.m {
		hosts = append(hosts, host)
	}
	return hosts
}

// -----

type byCreatedDesc []flux.ImageDescription

func (is byCreatedDesc) Len() int      { return len(is) }
func (is byCreatedDesc) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is byCreatedDesc) Less(i, j int) bool {
	if is[i].CreatedAt.Equal(is[j].CreatedAt) {
		return is[i].ID < is[j].ID
	}
	return is[i].CreatedAt.After(is[j].CreatedAt)
}

// Log requests as they go through, and responses as they come back.
// transport = logTransport{
// 	transport: transport,
// 	log: func(format string, args ...interface{}) {
// 		c.Logger.Log("registry-client-log", fmt.Sprintf(format, args...))
// 	},
// }
type logTransport struct {
	log       func(string, ...interface{})
	transport http.RoundTripper
}

func (t logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.log("Request %s %#v", req.URL, req)
	res, err := t.transport.RoundTrip(req)
	t.log("Response %#v", res)
	if err != nil {
		t.log("Error %s", err)
	}
	return res, err
}
