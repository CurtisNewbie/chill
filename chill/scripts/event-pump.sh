build="event-pump_build"
repo="event-pump"
service="event-pump"

cd ~/git/$repo \
    && git fetch \
    && git pull \
    && CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" /usr/local/go/bin/go build -o "$build" main.go \
    && cp "$build" ~/services/$service/build/ \
    && cp ./conf-prod.yml ~/services/$service/config/conf.yml \
    && cp ./Dockerfile_local ~/services/$service/build/Dockerfile \
    && cd ~/services \
    && /usr/bin/docker-compose up -d --build $service \
    && rm ~/git/$repo/$build

