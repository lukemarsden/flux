# For Luke

after cloning

```
gvt restore
make
```
This of course requires setup on your machine

clone this repo somewhere

```
git clone https://github.com/ContainerSolutions/flux-demo
```
All the services in the repo need to be launched individually

```
docker deploy -c finename.yml default_swarm
```
So K8s had namespaces in the definitions, swarm does not appear to have that, so for now we are hard coding. 
Will be putting it into the config file, once I find how to save in DB.
*So stack needs to be default_swarm*

you can use the docker-compose in this to launch the fluxsvc and fluxd

example flux.conf
```
git:
  URL: "https://github.com/ContainerSolutions/flux-demo.git"
  path: ""
  branch: "master"
  key: ""
slack:
  hookURL: ""
  username: ""
  releaseTemplate: ""
registry:
  auths: {}
```

then set the config in fluxctl

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
