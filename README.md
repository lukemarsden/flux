# For Luke

Docker needs to be setup with experimental flag enabled

and swarm mode turned on

```
docker swarm init
```

Go and gvt needs to be setup on your machine


get flux-swarm

```
go get github.com/ContainerSolutions/flux
```

after cloning

```
cd $GOPATH/src/github.com/ContainerSolutions/flux
gvt restore
make


```

you can use the docker-compose in this folder to launch fluxsvc and fluxd

```
docker-compose up
```

and set the flux url config
```
export FLUX_URL=http://localhost:3030/api/flux
```

clone the deploy scripts repo somewhere

```
cd ~/
git clone https://github.com/ContainerSolutions/flux-demo
```
All the services in the repo need to be launched individually

```
cd ~/flux-demo
for svc in *.yml; do docker deploy -c $svc default_swarm; done
```

So K8s had namespaces in the definitions, swarm does not appear to have that, so for now we are hard coding. 
Will be putting it into the config file, once I find how to save in DB.
*So stack needs to be default_swarm*


create a flux.conf with these contents and the deploy key Jason can provide
```
git:
  URL: git@github.com:ContainerSolutions/flux-demo.git
  path: 
  branch: master
  key: |
      #put key I gave you here
slack:
  hookURL: ""
  username: ""
  releaseTemplate: ""
registry:
  auths: {}
```

then set the config in fluxctl
```
fluxctl set-config -f flux.conf
```
you can list services

and update a service

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
