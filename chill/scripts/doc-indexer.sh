build="docindexer_build"

cd ~/git/doc-indexer \
    && git fetch \
    && git pull \
    && CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" /usr/local/go/bin/go build -o "$build" \
    && cp "$build" ~/services/docindexer/build/ \
    && cp ./prod-conf.yml ~/services/docindexer/config/conf.yml \
    && cp ./Dockerfile_local ~/services/docindexer/build/Dockerfile \
    && cd ~/services \
    && /usr/bin/docker-compose up -d --build docindexer \
    && rm ~/git/doc-indexer/$build

