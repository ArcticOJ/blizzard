FROM golang:alpine AS builder
WORKDIR /usr/src/app

COPY . .
RUN go mod download

RUN go build -o ./out/blizzard -ldflags "-s -w" main.go

FROM alpine AS runner
WORKDIR /blizzard

COPY --from=builder /usr/src/app/out/blizzard ./

ENTRYPOINT ["/blizzard/blizzard"]

