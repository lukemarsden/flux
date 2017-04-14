## Caveats

So K8s had namespaces in the definitions, swarm does not appear to have that, so for now we are hard coding.
Will be putting it into the config file, once I find how to save in DB.
*So stack needs to be default_swarm*

## Requirements

* docker 1.12+

## Installation

```
docker run --rm -it --net=host -v /var/run/docker.sock:/var/run/docker.sock golang:1.7.5 bash
```

### Within the docker container

```
curl -sSL https://get.docker.com/ | sh

curl -L "https://github.com/docker/compose/releases/download/1.11.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

curl -o /usr/local/bin/fluxctl -sSL https://github.com/weaveworks/flux/releases/download/master-0d109dd/fluxctl_linux_amd64
chmod +x /usr/local/bin/fluxctl

go get -u github.com/FiloSottile/gvt

go get github.com/ContainerSolutions/flux
cd $GOPATH/src/github.com/ContainerSolutions/flux
gvt restore
make
docker-compose up

cd ~/
git clone https://github.com/ContainerSolutions/flux-demo
cd ~/flux-demo

docker swarm init
for svc in *.yml; do docker deploy -c $svc default_swarm; done
```

Edit the flux.conf in flux-demo repo including the key given then set the config in fluxctl

```
fluxctl set-config -f flux.conf
```
Set env variable before executing any commands
```
export FLUX_URL=http://localhost:3030/api/flux
```
you can then list services
```
fluxctl list-services
```

and update a service
```
fluxctl release --service=default_swarm/orders --update-image=weaveworks/orders:master-ff176275
```

list images is problematic because io.docker seems to timeout.


# Flux

Flux is a tool for deploying container images to Kubernetes clusters.

![Flux Example](https://cloud.githubusercontent.com/assets/8793723/22978790/0d58861a-f38c-11e6-92d4-ce3f869e1ace.gif)

Please start by browsing through the documentation below.

[Introduction to Flux](/site/introduction.md)

[Installing Flux](/site/installing.md)

[Using Flux](/site/using.md)

[FAQ](/site/faq.md)

[Troubleshooting](/site/troubleshooting.md)

## Developer information

[Build documentation](/site/building.md)

[Release documentation](/internal_docs/releasing.md)

### Contribution

Flux follows a typical PR workflow.
All contributions should be made as PRs that satisfy the guidelines below.

### Guidelines

- All code must abide [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Names should abide [What's in a name](https://talks.golang.org/2014/names.slide#1)
- Code must build on both Linux and Darwin, via plain `go build`
- Code should have appropriate test coverage, invoked via plain `go test`

In addition, several mechanical checks are enforced.
See [the lint script](/lint) for details.

## <a name="help"></a>Getting Help

If you have any questions about Flux and continuous delivery:

- Invite yourself to the <a href="https://weaveworks.github.io/community-slack/" target="_blank"> #weave-community </a> slack channel.
- Ask a question on the <a href="https://weave-community.slack.com/messages/general/"> #weave-community</a> slack channel.
- Join the <a href="https://www.meetup.com/pro/Weave/"> Weave User Group </a> and get invited to online talks, hands-on training and meetups in your area.
- Send an email to <a href="mailto:weave-users@weave.works">weave-users@weave.works</a>
- <a href="https://github.com/ContainerSolutions/flux/issues/new">File an issue.</a>

Your feedback is always welcome!

curl -L "https://github.com/docker/compose/releases/download/1.11.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
