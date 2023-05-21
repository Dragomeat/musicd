FROM golang:1.20-alpine as tools

RUN apk add --no-cache build-base

RUN go install github.com/cespare/reflex@v0.3
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.8
RUN go install github.com/go-task/task/v3/cmd/task@latest

FROM golang:1.20-alpine

WORKDIR /go/src/app

RUN apk add --no-cache make build-base
COPY --from=tools /go/bin/reflex /usr/bin/reflex
COPY --from=tools /go/bin/sqlc /usr/bin/sqlc
COPY --from=tools /go/bin/task /usr/bin/task

CMD ["task", "start-api"]
