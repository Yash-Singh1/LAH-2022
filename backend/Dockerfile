# k's funny go dockerfile
FROM golang:1.19-alpine as build

WORKDIR /go/src/app
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -o out/backend \
    -v \
    -ldflags="-X main.CommitHash=$(cat .git/$(cat .git/HEAD | cut -d' ' -f2) | cut -c1-8)" \
    # delve
    -gcflags "all=-N -l" \
    lah-2022/backend

FROM alpine

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build --chown=appuser:appgroup /go/bin/dlv /dlv
COPY --from=build --chown=appuser:appgroup /go/src/app/out/ /opt/app/

WORKDIR /opt/app

USER appuser
ENTRYPOINT [ "/opt/app/backend" ]