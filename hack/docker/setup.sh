#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
SRC=$GOPATH/src
BIN=$GOPATH/bin
REPO_ROOT=$GOPATH/src/github.com/appscode/stash

source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

APPSCODE_ENV=${APPSCODE_ENV:-dev}
IMG=stash
RESTIC_VER=${RESTIC_VER:-SOURCE}

DIST=$REPO_ROOT/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
	export $(cat $DIST/.tag | xargs)
fi

clean() {
    pushd $REPO_ROOT/hack/docker
    rm -rf restic stash Dockerfile
    popd
}

build_binary() {
    pushd $REPO_ROOT
    ./hack/builddeps.sh
    ./hack/make.py build stash
    detect_tag $DIST/.tag

    if [ $RESTIC_VER = 'SOURCE' ]; then
        rm -rf $DIST/restic
        cd $DIST
        clone https://github.com/appscode/restic.git
        cd restic
        checkout master
        gb build
        mv bin/restic restic
        rm -rf $DIST/restic/src $DIST/restic/vendor
    else
        # Download restic
        rm -rf $DIST/restic
        mkdir $DIST/restic
        cd $DIST/restic
        wget https://github.com/restic/restic/releases/download/v${RESTIC_VER}/restic_${RESTIC_VER}_linux_amd64.bz2
        bzip2 -d restic_${RESTIC_VER}_linux_amd64.bz2
        mv restic_${RESTIC_VER}_linux_amd64 restic
    fi

    popd
}

build_docker() {
    pushd $REPO_ROOT/hack/docker

    # Download restic
    cp $DIST/stash/stash-linux-amd64 stash
    chmod 755 stash

    cp $DIST/restic/restic restic
    chmod 755 restic

    cat >Dockerfile <<EOL
FROM alpine

RUN set -x \
  && apk update \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/*

COPY restic /restic
COPY stash /stash

ENTRYPOINT ["/stash"]
EXPOSE 56790
EOL
    local cmd="docker build -t appscode/$IMG:$TAG ."
    echo $cmd; $cmd

    rm stash Dockerfile restic
    popd
}

build() {
    build_binary
    build_docker
}

docker_push() {
    if [ "$APPSCODE_ENV" = "prod" ]; then
        echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"
        exit 1
    fi
    if [ "$TAG_STRATEGY" = "git_tag" ]; then
        echo "Are you trying to 'release' binaries to prod?"
        exit 1
    fi
    hub_canary
}

docker_release() {
    if [ "$APPSCODE_ENV" != "prod" ]; then
        echo "'release' only works in PROD env."
        exit 1
    fi
    if [ "$TAG_STRATEGY" != "git_tag" ]; then
        echo "'apply_tag' to release binaries and/or docker images."
        exit 1
    fi
    hub_up
}

source_repo $@
