# k8s-validate-webhook

A fork of <https://github.com/alex-leonhardt/k8s-mutate-webhook> based on [this blog post](https://dev.to/ineedale/writing-a-very-basic-kubernetes-mutating-admission-webhook-5b1).

A playground to try build a the k8s mutating and validating webhook needed to deny PodExec CONNECT from users different from the original CR user. The following webhook need to be defined:

- A mutating webhook that adds users details to the DevWorkspace meta at creation time
- A validating webhook that denies changes to the user details on the DevWorkspace CR
- A validating webhook that denies changes to the link to a DevWorkspace CR on a container meta
- A validating webhook that denies PodExec CONNECT if the user should does not match the DevWorkspace CR creator

## todo

(0. Deploy a DevWorkspace CRD https://github.com/che-incubator/che-workspace-crd-operator/pull/3 and decide if that's worth using it or if we should use something simpler as https://github.com/l0rd/operator-hello-world)
1. Create a mutating webhook that adds OpenShift user details when a CR of type DevWorkspace is created
2. Create a validating webhook that denies changes to the user details on the DevWorkspace CR
3. Create a validating webhook that denies changes to the link to a DevWorkspace CR on a container meta (a label? an annotation?)
4. Create a validating webhook that denies PodExec CONNECT if the user should does not match the DevWorkspace CR creator

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
