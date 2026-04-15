# Contributing to Docker Deps

Thanks for your interest in contributing to **Docker Deps**! We welcome contributions of all kinds, including bug fixes, new features, documentation improvements, and bug reports.

## 📝 Code of Conduct

Please follow our [Code of Conduct](CODE_OF_CONDUCT.md) (standard professional behavior) when interacting with the project and other contributors.

## 🛠 How to Contribute

1. **Fork the Repository**: Start by forking the `docker-deps` repository on GitHub.
2. **Clone the Repository**: Clone your forked repository to your local machine.
3. **Create a New Branch**: Create a descriptive branch name for your changes (e.g., `feature/add-labels`).
4. **Make Changes**: Implement your changes and ensure they adhere to our coding standards.
5. **Write Tests**: Include tests for any new functionality or bug fixes.
6. **Verify Changes**: Run `make lint` and `make test` to ensure everything is correct.
7. **Submit a Pull Request**: Submit a PR to the `main` branch with a clear description of your changes.

## 🔨 Development Workflow

### Build and Run Local Tests

```bash
make test
```

### Linter Check

Ensure your code is clean:

```bash
make lint
```

## 🐛 Bug Reports

If you find a bug, please create an issue on GitHub with:
- A clear description of the problem.
- Steps to reproduce the issue.
- Details about your environment (Docker version, Go version, OS).

## 💡 Feature Requests

We're always looking for new ideas! Please open an issue to discuss any new features you'd like to see implemented.
