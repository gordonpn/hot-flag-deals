FROM golang:1.14-alpine AS builder
ARG TARGETARCH=amd64
ARG TARGETOS=linux
RUN test -n "$TARGETARCH"
RUN test -n "$TARGETOS"
RUN echo "Building for $TARGETOS $TARGETARCH"
RUN apk update
RUN apk add --no-cache \
  ca-certificates \
  git \
  tzdata
RUN update-ca-certificates
RUN mkdir /build
WORKDIR /build
RUN adduser \
  --disabled-password \
  --gecos "" \
  --no-create-home \
  --shell /bin/bash \
  --system \
  --uid 1000 \
  appuser
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify
COPY . .
ENV \
  CGO_ENABLED=0 \
  GO111MODULE=on \
  GOARCH=$TARGETARCH \
  GOOS=$TARGETOS
RUN go build -a -installsuffix cgo -o main .
FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /app/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
WORKDIR /app
USER appuser
CMD ["./main"]
