# Docker Deps

[![CI](https://github.com/filoucrackeur/docker-deps/actions/workflows/ci.yml/badge.svg)](https://github.com/filoucrackeur/docker-deps/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/filoucrackeur/docker-deps)](https://github.com/filoucrackeur/docker-deps)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Contributors](https://img.shields.io/github/contributors/filoucrackeur/docker-deps)](https://github.com/filoucrackeur/docker-deps/graphs/contributors)

A Docker CLI plugin for managing project dependencies. Think `npm`, `composer`, or `pip` — but for Docker images.

---

## Features

- **Manifest-based**: Declare image dependencies in a `docker-deps.json` file
- **CLI commands**: `add`, `remove`, `update`, `list`, `info`, `why`, `install`
- **Bulk install**: Pull all dependencies at once with progress display
- **Docker native**: Integrated as a first-class Docker CLI plugin

---

## Installation

### Prerequisites

- Go 1.26+
- Docker CLI

### From source

```bash
git clone https://github.com/filoucrackeur/docker-deps.git
cd docker-deps
make build
make install
```

Verify:
```bash
docker help | grep deps
```

---

## Usage

```bash
# Initialize a project
docker deps init my-project 1.0.0

# Add dependencies
docker deps add redis 7.2.0
docker deps add nginx 1.25.1

# Install all dependencies (pull images)
docker deps install

# List dependencies
docker deps list

# Update a dependency
docker deps update redis 7.2.4

# Get dependency info
docker deps info redis
docker deps why redis

# Remove a dependency
docker deps remove nginx
```

---

## Manifest

`docker-deps.json`:
```json
{
  "name": "my-project",
  "version": "1.0.0",
  "dependencies": {
    "nginx": "1.25.1",
    "redis": "7.2.0"
  }
}
```

---

## Development

```bash
make test      # Run tests
make lint      # Run linter
make coverage  # Check coverage
```

---

## License

MIT © Philippe Court

---

## Contributors

[![Contributors](https://contrib.rocks/image?repo=filoucrackeur/docker-deps)](https://github.com/filoucrackeur/docker-deps/graphs/contributors)
