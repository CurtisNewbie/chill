#!/bin/bash
build="mini-fstore_build"
repo="mini-fstore"
service="mini-fstore"

cd ~/git/$repo \
    && git fetch \
    && git pull \
    && CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" /usr/local/go/bin/go build -o "$build" cmd/main.go \
    && cp "$build" ~/services/$service/build/ \
    && cp ./prod-conf.yml ~/services/$service/config/conf.yml \
    && cp ./Dockerfile_local ~/services/$service/build/Dockerfile \
    && cd ~/services \
    && /usr/bin/docker-compose up -d --build $service \
    && rm ~/git/$repo/$build
