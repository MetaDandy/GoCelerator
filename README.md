# GoCelerator

**GoCelerator** is a CLI tool to bootstrap and manage Go projects with a consistent, modular architecture. Whether you prefer the standard `net/http` server or the Fiber v2 framework, GoCelerator generates your project structure, basic handlers, configuration and migration setup in seconds.

## Features

- `init <name>`: scaffold a new project (choose `net/http` or Fiber v2)  
- `serve [--watch]`: start the dev server (`go run ./cmd/main` or with `air` hot-reload)  
- `generate module <name>` / `generate core <name>`: create new domain modules or core packages with model, DTO, repository, service and handler stubs  
- `add route <module> <method> <path>`: inject new endpoint into handler and container  
- `migrate create <name>` / `migrate up` / `migrate down`: manage Gorm or SQL migrations  
- `seed run [name]`: execute data seeders  
- `env init`: generate `.env` from template and guide you through API keys  
- `docker build` / `docker compose up`: containerize your app and bring up dependencies  
- `test`: run `go test ./... -cover` and report coverage  
- `doctor`: check Go ≥1.21, PostgreSQL, env vars, and give tips  
- `version`: show GoCelerator and Go versions  

## Installation

1. Install Go 1.21 (latest): https://go.dev/doc/go1.21  
2. Install Cobra CLI v1.9.0:  
   ```bash
   go install github.com/spf13/cobra-cli@v1.9.0
