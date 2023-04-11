
#!/bin/bash

set -x

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/obaydullahmhs/crd-controller/pkg/client \
  github.com/obaydullahmhs/crd-controller/pkg/apis \
  aadee.apps:v1alpha1 \
  --go-header-file /home/appscodepc/go/src/github.com/obaydullahmhs/crd-controller/hack/boilerplate.go.txt

#generate yaml for crd

controller-gen rbac:roleName=controller-perms crd paths=github.com/obaydullahmhs/crd-controller/pkg/apis/aadee.apps/v1alpha1 crd:crdVersions=v1 output:crd:dir=$GOPATH/src/github.com/obaydullahmhs/crd-controller/manifests output:stdout