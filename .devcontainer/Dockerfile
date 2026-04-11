FROM mcr.microsoft.com/devcontainers/go:dev-1.24-bookworm

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       git \
       curl \
       ca-certificates \
       awscli \
    && rm -rf /var/lib/apt/lists/*