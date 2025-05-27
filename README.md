# GoCelerator

**GoCelerator** is a CLI tool to bootstrap and manage Go projects with a consistent, modular architecture. Whether you prefer the standard `net/http` server or the Fiber v2 framework, GoCelerator generates your project structure, basic handlers, configuration and migration setup in seconds.

## Features

- `init <name>`: scaffold a new project (choose `net/http` or Fiber v2)  
- `version`: show GoCelerator and Go versions  
- `serve [--watch]`: start the dev server (`go run ./cmd/main` or with `air` hot-reload)  
- `docker` / `docker compose up`: containerize your app and bring up dependencies  
- `generate module <name>` / `generate core <name>`: create new domain modules or core packages with model, DTO, repository, service and handler stubs  
- `add route <module> <method> <path>`: inject new endpoint into handler and container  
- `migrate create <name>` / `migrate up` / `migrate down`: manage Gorm or SQL migrations  
- `seed run [name]`: execute data seeders  
- `env init`: generate `.env` from template and guide you through API keys  
- `test`: run `go test ./... -cover` and report coverage  
- `doctor`: check Go ≥1.21, PostgreSQL, env vars, and give tips  

## Installation

### Prerequisites

- **Go 1.21 or later**  
  Download and install from https://go.dev/doc/go1.21  
- **Zsh shell** (or adjust below to your shell’s rc file)

### Install GoCelerator

1. Ensure your Go environment variables are set so that `$(go env GOPATH)/bin` is on your `PATH`. In `~/.zshrc`, add:

   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. Reload your shell:

   ```bash
   source ~/.zshrc
   ```

3. Install the Cobra CLI (for generating your own commands later):

   ```bash
   go install github.com/spf13/cobra-cli@v1.9.0
   ```

4. Install GoCelerator itself:

   ```bash
   go install github.com/Metadandy/GoCelerator@latest
   ```

5. Verify installation:

   ```bash
   goce version
   ```

   You should see something like this:

   ```bash
   goce version: v0.1.0
   Go version:   go1.24.1
   ```
   Note: if you prefer a specific version, replace @latest with @v0.1.0 (or any other tag).

## Quick Start

### Scaffold a new project
   ```bash
   goce init myapp --fiber
   cd myapp
   go mod tidy
   ```

### Start the server with hot-reload
   ```bash
   goce serve --watch
   ```

### Spin up production Docker stack
   ```bash
   goce docker
   ```

### Spin up development Docker stack
   ```bash
   goce docker --dev
   ```

## Contributing
PRs welcome! Please open issues or pull requests for bugs and feature requests.