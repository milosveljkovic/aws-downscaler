FROM golang:1.24 AS build
WORKDIR /src
ARG VERSION=dev
ARG BUILD_DATE=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X main.version=${VERSION}-${BUILD_DATE}" \
  -o /out/aws-downscaler ./cmd/aws-downscaler

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/aws-downscaler /app/aws-downscaler
ENTRYPOINT ["/app/aws-downscaler"]
CMD ["-config", "/config/aws-downscaler.yaml"]