# TODO: Add Windoze support

FROM golang:alpine AS builder
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./out/blizzard -ldflags "-s -w" main.go

FROM alpine AS runner
WORKDIR /igloo

COPY --from=builder /usr/src/app/out/blizzard ./

ENTRYPOINT ["/blizzard/blizzard"]

