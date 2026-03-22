# Contributing

Thank you for your interest in contributing to the Terraform Provider for Jamf School.

## Prerequisites

- **Go** >= 1.26 (see `go.mod` for the exact version)
- **Terraform** >= 1.0
- **golangci-lint** for linting
- A Jamf School tenant with API credentials (for acceptance tests only)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/Jamf-Concepts/terraform-provider-jamfschool.git
cd terraform-provider-jamfschool

# Build, lint, and generate docs
make
```

## Development Workflow

1. Create a feature branch from `main`.
2. Make your changes following the conventions in the [Style Guide](STYLEGUIDE.md).
3. Run formatting, linting, and tests before committing:

   ```bash
   make fmt
   make lint
   make test
   ```

4. Regenerate documentation if schema descriptions changed:

   ```bash
   make generate
   ```

5. Open a pull request against `main`.

## Make Targets

| Target     | Description                                           |
| ---------- | ----------------------------------------------------- |
| `build`    | Build the provider                                    |
| `install`  | Build and install the provider locally                |
| `fmt`      | Format Go source files                                |
| `lint`     | Run golangci-lint                                     |
| `generate` | Generate provider documentation with tfplugindocs     |
| `test`     | Run unit tests                                        |
| `testacc`  | Run acceptance tests (requires environment variables) |

The default target runs: `fmt lint install generate`.

## Adding a New Resource

1. Add types and service methods in the [jamfschool-go-sdk](https://github.com/Jamf-Concepts/jamfschool-go-sdk) (`jamfschool/<resource>.go`).
2. Create the resource package under `internal/resources/<resource_name>/` following the [file conventions](STYLEGUIDE.md#resource-package-file-conventions).
3. Register the resource in `internal/provider/provider.go` in the `Resources()` method.
4. Add tests:
   - Acceptance tests in `internal/resources/<resource_name>/resource_test.go`.
   - Data source tests in `internal/resources/<resource_name>/data_source_test.go`.
5. Add example `.tf` files under `examples/resources/jamfschool_<resource_name>/`.
6. Run `make generate` to regenerate documentation.
7. Run `make test` to verify all tests pass.

## Adding a New Data Source

Follow the same pattern as resources, but implement `datasource.DataSource` instead of `resource.Resource`. Place the data source in the same resource package (e.g., `internal/resources/<resource_name>/data_source.go`).

## Project Structure

| Directory              | Purpose                                                  |
| ---------------------- | -------------------------------------------------------- |
| `internal/provider/`   | Provider configuration and resource registration         |
| `internal/resources/`  | Resource, data source implementations                    |
| `internal/actions/`    | Action implementations (device commands)                 |
| `internal/common/`     | Shared helpers and validators                            |
| `examples/`            | Example `.tf` configurations                             |
| `templates/`           | tfplugindocs templates for doc generation                 |
| `docs/`                | Auto-generated provider documentation                    |

The REST client and service layer live in the standalone [jamfschool-go-sdk](https://github.com/Jamf-Concepts/jamfschool-go-sdk).

## Testing

### Unit Tests

```bash
make test
```

Unit tests cover the client, helpers, provider schema, and metadata. No live API access is needed.

### Acceptance Tests

```bash
export JAMFSCHOOL_URL="https://myschool.jamfcloud.com"
export JAMFSCHOOL_NETWORK_ID="your-network-id"
export JAMFSCHOOL_API_KEY="your-api-key"
make testacc
```

Acceptance tests create and destroy real resources in Jamf School. For read-only data source tests (device, app, profile, location, DEP device), set the additional environment variables referenced in the test files (e.g., `JAMFSCHOOL_TEST_DEVICE_UDID`).

### Writing Tests

- Resource tests should cover create, import, and update steps.
- Verify both required and computed attributes.
- Data source tests should verify all returned attributes.
- Use `acctest.RandomWithPrefix("tf-acc-")` for unique test resource names.

## Dependencies

This project uses native Go and Terraform Plugin Framework packages. Do not introduce third-party dependencies without discussion.

## Commit Messages

Use [conventional commit](https://www.conventionalcommits.org/) style messages:

- `feat: add device_group import support`
- `fix: handle nil response in state builder`
- `test: add schema validation for user resource`
- `refactor: extract common helpers`
- `docs: update README with new resource examples`

## Pull Requests

- Keep PRs focused — one feature or fix per PR.
- Include unit tests for new code.
- Include acceptance tests for new resources and data sources.
- Update `examples/` for new Terraform constructs.
- Run `make generate` if schema descriptions changed (to update docs).
- Linting must pass before merge.

## Reporting Issues

Open an issue on GitHub with:

- Provider version and Terraform version.
- Relevant Terraform configuration (redact credentials).
- Expected vs actual behaviour.
- Any error messages or logs.
