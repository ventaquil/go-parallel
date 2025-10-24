# GitHub Actions Workflows

This directory contains GitHub Actions workflows for CI/CD automation.

## Workflows

### CI Workflow (`ci.yml`)

**Triggers:**
- On every push to any branch
- On every pull request to any branch

**What it does:**
- Tests the package on multiple Go versions (1.21, 1.22, 1.23)
- Downloads and verifies dependencies
- Builds the package
- Runs tests with race detection and coverage
- Runs `go vet` for static analysis
- Checks code formatting with `gofmt`
- Uploads coverage reports to Codecov (optional)

### Release Workflow (`release.yml`)

**Triggers:**
- When a tag matching `v*` pattern is pushed (e.g., `v0.1.0`, `v1.0.0`)

**What it does:**
1. **Release Job:**
   - Runs all tests to ensure quality
   - Creates a GitHub Release with release notes
   - Marks the release as published

2. **Publish Job:**
   - Triggers the Go proxy (proxy.golang.org) to index the new version
   - Verifies the package is available
   - The package becomes installable via `go get`

## How to Use

### Running Tests Automatically

Simply push commits or create pull requests. The CI workflow will automatically:
- Build your code
- Run all tests
- Check for formatting issues
- Report any problems

### Publishing a New Version

1. Ensure all tests pass on your branch
2. Create and push a version tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```
3. The release workflow will:
   - Create a GitHub release
   - Make the package available via Go proxy
   - Users can install with: `go get github.com/ventaquil/go-parallel@v0.1.0`

### Semantic Versioning

Follow semantic versioning for tags:
- `v0.x.x` - Initial development
- `vMAJOR.MINOR.PATCH` - Production releases
  - MAJOR: Breaking changes
  - MINOR: New features (backward compatible)
  - PATCH: Bug fixes (backward compatible)

## Go Proxy Publishing

When a version tag is pushed, the package is automatically published to the Go module proxy:
- **Proxy URL:** https://proxy.golang.org
- **Module Path:** github.com/ventaquil/go-parallel
- **Installation:** Users can install any tagged version with `go get`

The Go proxy caches modules and makes them available globally with checksums for security.
