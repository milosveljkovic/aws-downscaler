# aws-downscaler

A Go-based tool for scaling AWS resources up or down based on configuration (tags, schedules, etc.).

---

## Usage

### Config example

Define your config.yaml (example in [aws-downscaler.yaml](./aws-downscaler.yaml))

Note: Only ec2 scale-up/down supported!

```yaml
# cloud provider, aws for now, maybe will be extended to support azure/google cloud as well
aws:
  eu-west-1:
    ec2:
      - name: qa-environment
        tags:
          - name: downscale
            value: 'true'
        downtime: 'Sat-Sun 00:00-24:00 UTC'
      - name: night-instances
        tags:
          - name: downscale
            value: 'true'
        downtime: 'Mon-Sun 20:00-08:00 Europe/Belgrade'
interval: 60s
log: info
```

Based on the config your aws-downscaler will perform scale-up/down of your aws services (ec2 supported only for now).

- Use tags to filter ec2 instances you want to scale-up/down:

```yaml
tags:
  - name: downscale
    value: 'true'
  - name: team
    value: 'qa'
```

- Set one or many downtimes separated by `,` character.

For example: During working days scale down instances between 20:00 and 08:00 and during weekend keep them off.

```yaml
downtime: 'Mon-Fri 20:00-08:00 Europe/Belgrade, Sat-Sun 00:00-24:00 Europe/Belgrade'
```

Try it, play with it, adapt to your need.

---

### Try it as docker

Available docker images available here: [DOCKER-IMAGES](https://github.com/milosveljkovic/aws-downscaler/pkgs/container/aws-downscaler)

```sh
docker run \
   -v $(pwd)/aws-downscaler.yaml:/config/aws-downscaler.yaml \
   -v "$HOME/.aws:/home/appuser/.aws:ro" \
   ghcr.io/milosveljkovic/aws-downscaler:latest

# OR

docker run \
   -v $(pwd)/aws-downscaler.yaml:/config/other-config.yaml \
   -v "$HOME/.aws:/home/appuser/.aws:ro" \
   ghcr.io/milosveljkovic/aws-downscaler:latest --config /config/other-config.yaml
```

### Or download from release page

Download binaries from [release-page](https://github.com/milosveljkovic/aws-downscaler/releases).

## 📦 Project Structure

```txt
cmd/aws-downscaler/main.go          # Application entrypoint
aws-downscaler.yaml                 # Configuration file
```

---

## ⚙️ Prerequisites

- Go (>= 1.24 recommended)
- AWS CLI installed
- VS Code (for debugging)
- Dev Container (recommended setup)

---

## 🔐 AWS Authentication

Before running or debugging the app, you must authenticate with AWS.

Login using your profile:

```bash
aws sso login --profile dev
```

Export the profile:

```bash
export AWS_PROFILE=dev
```

Verify credentials:

```bash
aws sts get-caller-identity
```

---

## 🐳 Run in Dev Container

1. Open the project in VS Code
2. Reopen in Dev Container
3. Inside the container, run:

```bash
aws sso login --profile dev
export AWS_PROFILE=dev
```

---

## ▶️ Run the Application

```bash
make run
```

This executes:

```bash
go run ./cmd/aws-downscaler
```

---

## 🐞 Debugging (VS Code)

1. Ensure:

   - Dev container is running
   - AWS is authenticated
   - `AWS_PROFILE` is set

2. Press: F5

VS Code will use `launch.json` to start debugging.

---

## 🧪 Run Tests

```bash
make test
```

This runs:

```bash
go test -v -timeout 30s ./...
```

---

## 🏗️ Build Binary

```bash
make build
```

This runs:

```bash
go build -o aws-downscaler ./cmd/aws-downscaler
```

Run the binary:

```bash
./aws-downscaler --help
```

---

## 🛠️ Makefile Commands

```bash
make run     # Run application
make test    # Run all tests with timeout
make build   # Build binary
```

---

## ⚠️ Notes

- Always authenticate with AWS before running/debugging
- `AWS_PROFILE` must be set (e.g. `dev`)
- Use `AWS_EC2_METADATA_DISABLED=true` inside containers to avoid hangs
- Ensure your `config.yaml` is in the project root

---

## 🚀 Future Improvements

- Cli flags to have more control over the app
- Dry-run mode for safe testing
- Metrics/logging integration
- Support for other AWS services

---
