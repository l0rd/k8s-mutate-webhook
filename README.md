# k8s-validate-webhook

A fork of <https://github.com/alex-leonhardt/k8s-mutate-webhook> based on [this blog post](https://dev.to/ineedale/writing-a-very-basic-kubernetes-mutating-admission-webhook-5b1).

A playground to try build a crude k8s mutating webhook; the goal is to mutate a Pod CREATE request to _always_ use a debian image and by doing this, learning more about
the k8s api, objects, etc. - eventually figure out how scalable this is (could be made) if one had 1000 pods to schedule (concurrently)

## todo

1.

## concerns

1.

## build

```bash
make
```

## test

```bash
make test
```

## ssl/tls

the `ssl/` dir contains a script to create a self-signed certificate, not sure this will even work when running in k8s but that's part of figuring this out I guess

_NOTE: the app expects the cert/key to be in `ssl/` dir relative to where the app is running/started and currently is hardcoded to `mutateme.{key,pem}`_

```bash
cd ssl/
make
```

## docker

to create a docker image .. 

```bash
make docker
```

it'll be tagged with the current git commit (short `ref`) and `:latest`

don't forget to update `IMAGE_PREFIX` in the Makefile or set it when running `make`

### images

[`mariolet/k8s-mutate-webhook`](https://cloud.docker.com/repository/docker/mariolet/k8s-mutate-webhook:latest)

```bash
docker push mariolet/k8s-mutate-webhook:latest
```

## deploy

```bash
# Deploy the webhook
k apply -f deploy/webhook.yaml

# k rollout restart -n default deploy mutateme

# Setup a sample resource to validate
k create namespace admission-test
k patch namespace admission-test --patch '{"metadata": {"labels": {"mutateme": "disabled"}}}'
k patch namespace admission-test --patch '{"metadata": {"labels": {"validateme": "enabled"}}}'
kubens admission-test

curl -sSL https://k8s.io/examples/application/deployment.yaml | yq . | jq '.spec.replicas = 1' | jq '.spec.template.metadata.labels.valid = "false"' | k apply -f -
```

## cleanup

```bash
k delete deploy nginx-deployment
k delete namespace admission-test
k delete -f deploy/
```

## watcher

useful during devving ...

```bash
watcher -watch github.com/alex-leonhardt/k8s-mutate-webhook -run github.com/alex-leonhardt/k8s-mutate-webhook/cmd/
```

## kudos

- [https://github.com/morvencao/kube-mutating-webhook-tutorial](https://github.com/morvencao/kube-mutating-webhook-tutorial)
