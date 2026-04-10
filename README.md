# aws-downscaler

A Go-based tool for scaling AWS EC2 resources up or down based on configuration (tags, schedules, etc.).

---

## 📦 Project Structure

```
cmd/aws-downscaler/main.go          # Application entrypoint
aws-downscaler.yaml                 # Configuration file
```

---

## ⚙️ Prerequisites

* Go (>= 1.24 recommended)
* AWS CLI installed
* VS Code (for debugging)
* Dev Container (recommended setup)

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

   * Dev container is running
   * AWS is authenticated
   * `AWS_PROFILE` is set

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
./aws-downscaler
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

* Always authenticate with AWS before running/debugging
* `AWS_PROFILE` must be set (e.g. `dev`)
* Use `AWS_EC2_METADATA_DISABLED=true` inside containers to avoid hangs
* Ensure your `config.yaml` is in the project root

---

## 🚀 Future Improvements

* Cli flags to have more control over the app
* Dry-run mode for safe testing
* Metrics/logging integration
* Support for other AWS services

---
