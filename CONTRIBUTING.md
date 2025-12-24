# Contributing Guide

Thank you for your interest in contributing to **go-autoupdater** 🎉  
All kinds of contributions are welcome: bug reports, feature requests, documentation improvements, and code contributions.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Project Overview](#project-overview)
- [Development Environment](#development-environment)
- [Repository Structure](#repository-structure)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Requesting Features](#requesting-features)
  - [Submitting Pull Requests](#submitting-pull-requests)
- [Coding Guidelines](#coding-guidelines)
- [Commit Message Convention](#commit-message-convention)
- [Testing](#testing)
- [Security Considerations](#security-considerations)
- [License](#license)

---

## Code of Conduct

By participating in this project, you agree to uphold a respectful and inclusive environment.

- Be respectful and constructive
- Assume good intentions
- No harassment, discrimination, or abusive behavior

Serious or repeated violations may result in removal from the project.

---

## Project Overview

**go-autoupdater** is a cross-platform self-update framework written in Go that supports:

- Manifest-based updates (`manifest.json`)
- SHA256 integrity verification
- Windows, Linux, macOS, and embedded Linux/STB
- Standalone CLI and embeddable library usage
- Service-aware updates (NSSM, SC, systemd, launchd)

The project is designed with **modularity and portability** in mind.

---

## Development Environment

### Requirements

- Go **1.22+** recommended
- Git
- (Optional) Docker or CI for cross-platform builds

### Setup

```bash
git clone https://github.com/<your-username>/go-autoupdater.git
cd go-autoupdater
go mod tidy
go test ./...
```

---

## Repository Structure

```
cmd/
  updaterctl/        # CLI entrypoint
  updater-helper/    # Windows-only helper (binary swap)
pkg/
  updater/           # Core update engine (platform-agnostic)
  source/            # Update sources (HTTP manifest)
  verify/            # SHA256 verification
  apply/             # Apply/swap logic (OS-specific)
  service/           # Service controllers (OS-specific)
  util/              # Utilities (logging, fs helpers)
```

### Important Rule
- **Core logic must remain platform-agnostic**
- OS-specific code must use Go build tags:
  - `//go:build windows`
  - `//go:build linux`
  - `//go:build darwin`

---

## How to Contribute

### Reporting Bugs

Please open an issue with:

- OS and architecture
- Go version
- Steps to reproduce
- Expected vs actual behavior
- Logs (if available, redact sensitive data)

### Requesting Features

Feature requests are welcome.  
When opening an issue, please include:

- Problem you are trying to solve
- Why the current implementation is insufficient
- Proposed solution or alternatives (if any)

---

### Submitting Pull Requests

1. Fork the repository
2. Create a new branch:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Make your changes
4. Run tests:
   ```bash
   go test ./...
   ```
5. Commit and push:
   ```bash
   git commit -m "feat: add new updater source"
   git push origin feature/my-feature
   ```
6. Open a Pull Request (PR)

#### PR Requirements

- Clear description of changes
- One logical change per PR (avoid mixing unrelated changes)
- Tests or explanation if tests are not applicable
- No secrets, credentials, or private URLs

---

## Coding Guidelines

- Follow standard Go formatting (`gofmt`)
- Prefer explicit error handling
- Avoid unnecessary dependencies
- Keep functions small and focused
- Public APIs should be documented with comments

---

## Commit Message Convention

Use a simple conventional format:

```
type: short description

(optional body)
```

Examples:
- `fix: handle service restart failure`
- `feat: add HTTP manifest source`
- `docs: update README`
- `refactor: simplify swap logic`

---

## Testing

Before submitting a PR:

```bash
go test ./...
```

Recommended:
- Test at least on your local OS
- For Windows-related changes, ensure code is guarded by build tags

---

## Security Considerations

⚠️ **Do not submit security vulnerabilities as public issues.**

If you discover a security issue:
- Do **not** open a public issue
- Contact the maintainer privately or open a GitHub Security Advisory

Current updater verifies **SHA256 only**.  
Signature-based verification (e.g., Ed25519) is planned but not yet implemented.

---

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License, unless stated otherwise).

---

Thank you for contributing 🙏  
Your help makes this project better for everyone.
