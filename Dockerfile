FROM golang:1.19.6 as build

RUN mkdir /gosrc
WORKDIR /gosrc
COPY go.mod go.sum ./

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go mod download

COPY . .

RUN go build main.go

FROM scratch

ARG COMMIT_SHA=""
ENV COMMIT_SHA="$COMMIT_SHA"

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /gosrc/main /

ENTRYPOINT ["/main"]
