FROM goland:1.20-alpine AS builder

WORKDIR .

RUN apk --no-cache add bash git make gcc gettext musl-dev

copy ["go.mod",