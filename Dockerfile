FROM golang:1.10.4 AS build
WORKDIR /go/src/github.com/owensengoku/pixie

RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/pixie -ldflags="-w -s" -v github.com/owensengoku/pixie/cmd/pixie

FROM alpine:3.8 AS final
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/pixie /bin/pixie
