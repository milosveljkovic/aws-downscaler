FROM golang:1.24-alpine AS build
WORKDIR /src

ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags="-s -w -X main.version=${VERSION}-${BUILD_DATE}" \
    -o /out/aws-downscaler ./cmd/aws-downscaler

FROM alpine:3.19

RUN apk add --no-cache tzdata


RUN adduser -D appuser
USER appuser

WORKDIR /app
COPY --from=build /out/aws-downscaler /app/aws-downscaler

ENTRYPOINT ["/app/aws-downscaler"]
CMD ["-config", "/config/aws-downscaler.yaml"]