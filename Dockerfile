ARG GOLANG_VERSION=1.24
FROM golang:${GOLANG_VERSION}-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION


ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV GOARM=${TARGETVARIANT}
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
ENV GOSUMDB=sum.golang.org
ENV LDFLAGS="-s -w"


WORKDIR /go/src/github.com/technicaldomain/aws-secret-to-file
COPY . .

RUN export GOARM=$( echo ${TARGETVARIANT} | cut -c2- )
RUN go build -mod=vendor -ldflags "${LDFLAGS}" -o /go/bin/aws-secret-to-file

FROM scratch
COPY --from=builder /go/bin/aws-secret-to-file /bin/aws-secret-to-file
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
