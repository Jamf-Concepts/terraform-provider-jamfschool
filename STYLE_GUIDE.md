# Style Guide

Code style conventions for the Terraform Provider for Jamf School.

## Go Conventions

- Follow standard Go conventions and idiomatic patterns.
- Run `make fmt` and `make lint` before committing.
- Use clear, descriptive names for variables, functions, and types.
- Every exported constant, function, variable set, and type must have a short comment describing its purpose.
- Do not add comments inside type definitions or function bodies.

## Dependencies

Only use native Go, `golang.org/x` packages, the [Jamf School Go SDK](https://github.com/Jamf-Concepts/jamfschool-go-sdk), and Terraform Plugin Framework packages. Do not introduce third-party dependencies without discussion.

## Resource Package File Conventions

Resource packages live under `internal/resources/<resource_name>/` and use resource-agnostic filenames:

| File                 | Purpose                                                   |
| -------------------- | --------------------------------------------------------- |
| `resource.go`        | Schema definition and boilerplate                         |
| `crud.go`            | Create, Read, Update, Delete, and ImportState             |
| `model_types.go`     | Terraform model structs                                   |
| `schema_types.go`    | Attribute type maps for `ObjectValue`/`ListValue` state   |
| `mappings.go`        | Lookup tables and name mappings                           |
| `input_builders.go`  | Build API request inputs from Terraform model data        |
| `state_builders.go`  | Map API responses to Terraform state                      |
| `helpers.go`         | Resource-specific helper functions                        |
| `plan_modifiers.go`  | Schema plan modifiers (if needed)                         |
| `validators.go`      | Schema validators (if needed)                             |
| `list_resource.go`   | List resource implementation                              |
| `data_source.go`     | Data source implementation                                |

### Optional split-outs for complex resources

- `nested_builders.go` / `nested_state.go` — for large nested payloads.

### Data-source-only packages

Packages that only contain a data source use `model_types.go` for their model structs and `data_source.go` for the implementation.

## Test File Conventions

| File                      | Purpose                                    |
| ------------------------- | ------------------------------------------ |
| `resource_test.go`        | Acceptance tests for the resource          |
| `data_source_test.go`     | Acceptance tests for the data source       |
| `list_resource_test.go`   | Unit tests for list resource metadata/schema |
| `helpers_test.go`         | Helper function tests                      |
| `input_builders_test.go`  | Input builder tests                        |
| `state_builders_test.go`  | State builder tests                        |
| `mappings_test.go`        | Mapping table tests                        |

Schema and metadata tests live in `internal/provider/provider_test.go`.

## Service Layer

The provider uses the [Jamf School Go SDK](https://github.com/Jamf-Concepts/jamfschool-go-sdk) (`jamfschool.Client`) for all API operations. The SDK handles authentication, HTTP transport, and provides typed CRUD methods per resource. Sentinel errors (`ErrAuthentication`, `ErrHTTP`, `ErrNotFound`) are defined by the SDK.

## Schema Guidelines

- Keep schemas inline and as flat as possible.
- Favor nested attributes (`SingleNestedAttribute`, `SetNestedAttribute`, `ListNestedAttribute`) over blocks.
- **Terraform attribute names should reflect the Jamf School UI labels**, not the API JSON keys.
- Use validators for all fields with constrained values (enums, ranges, formats).
- For write-only fields (API doesn't return them), use `Default` values and add them to `ImportStateVerifyIgnore` in tests.
- For fields where the API ignores updates, use `RequiresReplace()` plan modifier.

### Sets vs Lists

- **Sets** for user-supplied unordered collections where deduplication and order-independent comparison matter.
- **Lists** for computed API results that are read-only. Sets require element hashing which adds overhead with no benefit when the user doesn't control the values.

Data source attributes returning API data should always use lists. Sort API responses in data source state builders.

## Error Handling

- Use the SDK sentinel errors: `jamfschool.ErrAuthentication`, `jamfschool.ErrHTTP`, `jamfschool.ErrNotFound`.
- In `Read` methods, check for `ErrNotFound` and call `resp.State.RemoveResource()` to handle deleted-outside-Terraform scenarios.
- Wrap errors with `fmt.Errorf("context: %w", err)` to preserve the error chain.
- Use `resp.Diagnostics.AddError()` for all error reporting in CRUD methods.

## Naming Patterns

### Resources

Terraform resource type names follow `jamfschool_<resource>`:

- `jamfschool_user`
- `jamfschool_user_group`
- `jamfschool_device_group`
- `jamfschool_class`
- `jamfschool_ibeacon`

### Test names

Test functions use the pattern `TestAcc<Resource>Resource_<scenario>` for acceptance tests and `Test<Function>_<case>` for unit tests:

```go
func TestAccUserResource_basic(t *testing.T) { ... }
func TestAccClassResource_basic(t *testing.T) { ... }
func TestStringValueOrNull(t *testing.T) { ... }
```

### Acceptance test resource names

Use the `tf-acc-` prefix for all resources created during acceptance tests:

```go
rName := acctest.RandomWithPrefix("tf-acc-user")
rName := acctest.RandomWithPrefix("tf-acc-class")
```

## Helpers

Shared helper functions live in `internal/common/helpers/`. Always check for existing helpers before adding new code.

- `Ptr[T](v T)` — returns a pointer to a value.
- `StringValueOrNull(s string)` — converts a Go string to `types.String`, null if empty.
- `Int64ValueOrNull(v int64)` — converts a Go int64 to `types.Int64`, null if zero.
- `StringPtrValueOrNull(s *string)` — converts a `*string` to `types.String`, null if nil.
- `Int64PtrIfKnown(v types.Int64)` — converts a `types.Int64` to `*int64`, nil if null/unknown.
