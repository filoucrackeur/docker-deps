# Docker Deps

[![Go Report Card](https://goreportcard.com/badge/github.com/philippe/docker-deps)](https://goreportcard.com/report/github.com/philippe/docker-deps)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Docker Deps** is a Docker CLI plugin that simplifies dependency management for your containerized projects. It works like `npm`, `composer`, or `pip`, but for your Docker images. It manages a `docker-deps.json` manifest to track, update, and install the specific image versions your project requires.

---

## 🚀 Key Features

- **Project Manifest**: Declare all your image dependencies in a single `docker-deps.json` file.
- **Dependency Management**: Easily `add`, `remove`, and `update` image requirements via the CLI.
- **Bulk Installation**: Run `docker deps install` to pull all required images at once with a clean progress display.
- **Visibility**: Get detailed `info` and `list` all your project requirements instantly.
- **Docker Native**: Integrated directly into the Docker CLI as a first-class plugin.

---

## 📦 Installation

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21 or higher
- [Docker](https://docs.docker.com/get-docker/) installed and running

### Step 1: Clone and Build

```bash
git clone https://github.com/your-username/docker-deps.git
cd docker-deps
make build
```

### Step 2: Install as a Docker Plugin

```bash
make install
```
*Note: This moves the binary to `~/.docker/cli-plugins/docker-deps`.*

### Step 3: Verify Installation

```bash
docker help | grep deps
```

---

## 🛠 Usage

### Initialize your project
```bash
docker deps init my-awesome-project 1.0.0
```

### Add a dependency
```bash
docker deps add redis 7.2.0
docker deps add nginx 1.25.1
```

### Install all dependencies (Pull images)
```bash
docker deps install
```

### List all dependencies
```bash
docker deps list
```

### Update a dependency
```bash
docker deps update redis 7.2.4
```

### Remove a dependency
```bash
docker deps remove nginx
```

---

## 📄 Manifest Structure (`docker-deps.json`)

```json
{
  "name": "my-awesome-project",
  "version": "1.0.0",
  "dependencies": {
    "nginx": "1.25.1",
    "redis": "7.2.0"
  }
}
```

---

## 🛠 Development & Testing

### Run Tests with Coverage
```bash
make test
```

### Linting
```bash
make lint
```

---

## 🤝 Contributing

Contributions are welcome! Please check our [Contributing Guide](CONTRIBUTING.md) for more details.

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.

---

## 👨‍💻 Author

**Philippe Court** - [philippe.court@gmail.com](mailto:philippe.court@gmail.com)
