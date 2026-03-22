# Repository Guidelines

## Overview

This is a Terraform provider for [Jamf School](https://www.jamf.com/products/jamf-school/), built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) v1.18.0 with Protocol v6. The Go module path is `github.com/Jamf-Concepts/terraform-provider-jamfschool`.

## Tooling

- Use `make` for build, lint, test, and doc generation. See `GNUmakefile` for available targets.
- Go >= 1.26, Terraform >= 1.0.

### Available make targets

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

## Jamf School API

- The classic Jamf School API spec is in `api-spec.yml` (copy-pasted from the rendered docs at `school.jamfcloud.com/api/docs/`).
- Base URL: `https://{yourDomain}.jamfcloud.com` (e.g. `https://myschool.jamfcloud.com`).
- Authentication: **HTTP Basic Auth** — the username is the **Network ID** (found at Devices > Enroll Device(s) in Jamf School) and the password is the **API Key** (generated at Organization > Settings > API in Jamf School).
- Request/response format: JSON. Content-Type `application/json`.
- API versioning: defaults to v1. Use `X-Server-Protocol-Version` header for newer versions (provider uses version 3).
- Standard error responses return `{ "code": <int>, "message": "<string>" }`.
- The API is read-only for some resource types (Apps, Profiles, Locations, Devices) and supports full CRUD for others (Users, Groups, Device Groups, Classes, iBeacons).

### Known API Discrepancies

The API spec and live API behaviour do not always match. Always verify response shapes by probing the live API:

- **Groups ACL**: Spec shows `acl` as an array; live API returns a single object. Only `teacher` and `parent` fields are returned; self-service fields are accepted on write but never returned.
- **iBeacon major/minor**: Spec says Number; live API sometimes returns strings. The provider has a custom `UnmarshalJSON` to handle both.
- **User domain**: The API ignores user-supplied values and always returns a server-set default. This field is Computed-only in the provider.
- **User storePassword**: Write-only. The API does not return this field in responses.
- **User memberOf**: The API spec only documents `memberOf` on POST, but PUT also supports it with full replacement semantics. Accepts mixed array of group IDs (int) and group names (string). Empty array `[]` clears all memberships. Response includes `groupIds` (int array) and `groups` (string array).
- **DeviceGroup information**: Persists on create, returned on read, but updates are ignored. Uses `RequiresReplace`.
- **DeviceGroup collectionType**: Create-only, write-only. The response `type` field is the group type (normal/class/smart), not the collection type.
- **Device model/os**: Spec shows arrays; live API returns single objects.
- **App trash**: The trash endpoint (`POST /apps/:id/trash`) returns 200 but deletion propagation across backend pods is unreliable. Apps must be manually removed from trash in the Jamf School UI before Terraform can recreate them with the same Adam ID. The `inTrash` list filter does not work for apps. Create is idempotent — calling create with an existing Adam ID returns the existing app.
- **App create/trash protocol version**: These endpoints require `X-Server-Protocol-Version: 4` (not 3). The provider uses `WithProtocolVersion("4")` for these calls.
- **Eventual consistency**: The Jamf School API runs across multiple Kubernetes pods (`x-server-instance` header) with no sticky session mechanism. Reads after writes may hit a different pod and return stale data. This affects CheckDestroy reliability in acceptance tests.

### API Resource Summary

| API Section       | Endpoint Prefix      | CRUD Support | Provider Mapping                     |
| ----------------- | -------------------- | ------------ | ------------------------------------ |
| Users             | `/users`             | Full CRUD    | Resource + Data Source               |
| Groups            | `/users/groups`      | Full CRUD    | Resource + Data Source               |
| Device Groups     | `/devices/groups`    | Full CRUD    | Resource + Data Source               |
| Classes           | `/classes`           | Full CRUD    | Resource + Data Source               |
| iBeacons          | `/ibeacons`          | Full CRUD    | Resource + Data Source               |
| Devices           | `/devices`           | Read + Cmds  | Data Source only                     |
| Apps              | `/apps`              | Read only    | Data Source only                     |
| Profiles          | `/profiles`          | Read only    | Data Source only                     |
| Locations         | `/locations`         | Read only    | Data Source only                     |
| DEP Devices       | `/dep`               | Read + Update| Data Source only                     |

## Project Structure

```text
main.go                          # Provider entry point (registry.terraform.io/Jamf-Concepts/jamfschool)
internal/
  common/
    helpers/                     # Shared helper utilities (Ptr, Int64PtrIfKnown, etc.)
  provider/                      # Provider wiring + schema validation tests
  resources/                     # Per-resource packages (resource + data source)
    user/                        # Users (CRUD + data source)
    user_group/                  # User Groups (CRUD + data source)
    device_group/                # Device Groups (CRUD + data source)
    class/                       # Classes (CRUD + data source)
    ibeacon/                     # iBeacons (CRUD + data source)
    device/                      # Data source only (enrolled devices)
    app/                         # Data source only (apps)
    profile/                     # Data source only (configuration profiles)
    location/                    # Data source only (locations)
    dep_device/                  # Data source only (Automated Device Enrollment)
  actions/
    device/                      # Device command actions (erase, restart, refresh, etc.)
tools/                           # Go generate tooling (tfplugindocs, copywrite)
docs/                            # Generated provider documentation (resources + data sources)
examples/
  resources/                     # Example .tf files for resources
  data-sources/                  # Example .tf files for data sources
  provider/                      # Example provider configuration
templates/                       # tfplugindocs templates for doc generation
local_testing/                   # Manual testing project (gitignored)
api-spec.yml                     # Classic Jamf School API spec (copy-pasted from rendered docs)
```

## Provider Development

- Terraform Plugin Framework code lives in `internal/`.
- The REST client and service layer have been extracted to the standalone SDK at `github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool`. The SDK provides `*jamfschool.Client` with typed methods per API resource (e.g. `GetUser`, `CreateUser`, `UpdateUser`, `DeleteUser`) and sentinel errors (`ErrAuthentication`, `ErrHTTP`, `ErrNotFound`).
- Resource implementations are grouped by package in `internal/resources/<resource>` with files split by concern.
- Run formatting and linting before committing: `make fmt` and `make lint`.
- Generate docs with `make generate`.
- Run tests with `make test`; acceptance tests with `make testacc` (requires real tenant).

## Code Organization Guidelines

- Look for opportunities to create reusable packages (helper/utility functions) instead of duplicating logic in resource packages.
- Keep packages split by concern with focused files.
- Always look for existing helper functions that can be reused before adding new code.

## Code Style Guidelines

- Follow Go conventions and idiomatic patterns.
- Favor clear and descriptive naming for variables, functions, and types.
- Always ensure constants, functions, variable sets and types have a short comment describing their purpose.
- Do not add comments inside type definitions or function bodies.

### Resource Package File Conventions

Use resource-agnostic filenames and helper names so the same structure can apply to all resources:

- `resource.go`: schema and boilerplate.
- `crud.go`: Create/Read/Update/Delete and import.
- `model_types.go`: Terraform model structs only.
- `input_builders.go`: build API inputs from Terraform model data.
- `state_builders.go`: map API responses to Terraform state.
- `validators.go`: schema validators (if needed).
- `data_source.go`: for data sources implementing `datasource.DataSource`.
- `resource_test.go`: acceptance tests for the resource.
- `data_source_test.go`: acceptance tests for the data source.

## Schema Guidelines

- Schemas should be inline and as flat as possible.
- **Terraform attribute names should reflect the Jamf School UI labels**, not the API JSON keys. For example, the API's `teacher` ACL field is named `jamf_school_teacher` in the provider to match the UI label "Jamf School Teacher".
- Use validators for all fields with constrained values (enums, ranges, formats).
- For write-only fields (API doesn't return them), use `Default` values and add them to `ImportStateVerifyIgnore` in tests.
- For fields where the API ignores updates, use `RequiresReplace()` plan modifier.
- For fields where the API may return a different value than sent, either make the field Computed-only or use the `stringValueOrKeep` pattern to preserve plan values when the API returns empty.

## Environment Variables

- `JAMFSCHOOL_URL` — Base URL of the Jamf School instance (e.g. `https://myschool.jamfcloud.com`).
- `JAMFSCHOOL_NETWORK_ID` — Network ID used as the HTTP Basic Auth username (found at Devices > Enroll Device(s)).
- `JAMFSCHOOL_API_KEY` — API key used as the HTTP Basic Auth password (generated at Organization > Settings > API).
- These can also be set in the provider block in Terraform configuration.

## Testing

### Test types

- **Unit tests**: `make test` — runs schema validation, metadata, client tests, service layer tests (mock HTTP), and helper tests (no real API needed).
- **Service acceptance tests**: `go test -run TestAccService ./internal/jamfschool/` — tests service layer CRUD against a live API. Requires env vars but not `TF_ACC`.
- **Acceptance tests**: `make testacc` — creates real Terraform resources against a Jamf School instance. Requires `JAMFSCHOOL_URL`, `JAMFSCHOOL_NETWORK_ID`, and `JAMFSCHOOL_API_KEY`.
- **Local testing**: `local_testing/` contains a comprehensive Terraform project exercising all resources, data sources, actions, and list resources. See its README.md for setup.

### Required tests for new code

When adding **any** new functionality, the following tests are required:

**Service layer methods** (`internal/jamfschool/`):

- Unit test per method in `service_test.go` using `httptest.NewServer` to mock the API. Verify HTTP method, path, request body, and response parsing.
- Service acceptance test in `service_acc_test.go` for CRUD operations against the live API. Use `t.Cleanup()` to ensure resources are deleted even if the test fails.

**Resources** (`internal/resources/<name>/`):

- `resource_test.go`: Acceptance test with steps for create (verify all attributes), import (`ImportStateVerifyIgnore` for write-only fields), and update (verify changed + unchanged attributes).

**Data sources** (`internal/resources/<name>/`):

- `data_source_test.go`: Acceptance test that creates a fixture resource then reads it via the data source. Verify all computed attributes.

**List resources** (`internal/resources/<name>/`):

- `list_resource_test.go` (internal package): Metadata and Schema unit tests verifying type name and `name_prefix` attribute.

**Actions** (`internal/actions/<domain>/`):

- `schema_test.go` (internal package): Metadata test (verify type name) and Schema test (verify all expected attributes) for each action.

**Provider** (`internal/provider/`):

- Update `provider_test.go` counts: `TestProviderResources`, `TestProviderDataSources`, `TestProviderActions`, `TestProviderListResources`.

### Test conventions

- Test files follow the `*_test.go` convention next to the code they test.
- Write-only fields (API doesn't return them) must be added to `ImportStateVerifyIgnore`.
- Use `acctest.RandomWithPrefix("tf-acc-")` for unique resource names.
- Service acceptance tests use `t.Cleanup()` for resource teardown.
- CheckDestroy is unreliable due to API eventual consistency across pods — only use for resources where the API reliably returns 404 after delete (e.g. iBeacons).

## Adding a New Resource

1. Create a new package under `internal/resources/<resource_name>/` following the file conventions above.
2. Add service layer methods to `internal/jamfschool/<resource>.go` with unit tests.
3. Implement `resource.Resource` with CRUD + `ImportState` + `IdentitySchema` in `resource.go` and `crud.go`.
4. Set identity in Create, Read, and Update via `resp.Identity.SetAttribute()`.
5. Register the resource in `internal/provider/provider.go` -> `Resources()`.
6. Add a list resource in `list_resource.go` and register in `ListResources()`.
7. Add all required tests (see Testing section above).
8. Add example `.tf` files under `examples/resources/<resource_name>/` and `examples/list-resources/<resource_name>/`.
9. Update provider test counts.
10. Run `make test` to ensure tests pass.
11. Run `make generate` to generate documentation from schema descriptions.

## Adding a New Action

Actions are side-effect operations (device commands, etc.) that don't manage state. They require Terraform 1.14+.

### Action Package Structure

Actions live in `internal/actions/<domain>/` (e.g. `internal/actions/device/`):

- `helpers.go`: shared base struct with `configure()`, `ensureService()`, and identifier resolution helpers.
- `<action_name>.go`: one file per action implementing `action.Action` and `action.ActionWithConfigure`.
- `schema_test.go`: unit tests for Metadata and Schema of all actions in the package.

### Implementation Pattern

Each action follows this structure:

```go
var _ action.Action = (*MyAction)(nil)
var _ action.ActionWithConfigure = (*MyAction)(nil)

type MyAction struct {
    deviceAction  // embed shared base
}

type MyActionModel struct {
    UDID         types.String `tfsdk:"udid"`
    SerialNumber types.String `tfsdk:"serial_number"`
    // action-specific fields
}

func NewMyAction() action.Action { return &MyAction{} }

func (a *MyAction) Metadata(...)  { resp.TypeName = req.ProviderTypeName + "_my_action" }
func (a *MyAction) Schema(...)    { /* use actionschema.Schema with actionschema.Attribute */ }
func (a *MyAction) Configure(...) { a.configure(ctx, req, resp) }
func (a *MyAction) Invoke(...)    {
    // 1. ensureService
    // 2. read config
    // 3. resolve device identifier
    // 4. send progress
    // 5. call service method
    // 6. handle errors
    // 7. send completion progress
}
```

### Key Differences from Resources

- Use `actionschema` package (not `resource/schema`) for schema definitions.
- Actions have `Invoke` (not Create/Read/Update/Delete).
- Actions cannot modify state — they execute and report diagnostics + progress.
- Use `resp.SendProgress()` for status updates during execution.
- Register via `provider.ProviderWithActions` interface and `Actions()` method.
- Set `resp.ActionData = svc` in the provider's `Configure` method.

### Registration

1. Add constructor to `internal/provider/provider.go` -> `Actions()`.
2. Import the action package with an alias (e.g. `deviceactions`).

## Adding a New Data Source

1. Create `data_source.go` in the relevant package under `internal/resources/<resource_name>/`.
2. Implement `datasource.DataSource` with a `Read` method.
3. Register in `internal/provider/provider.go` -> `DataSources()`.
4. Add acceptance tests.
5. Add example `.tf` files under `examples/data-sources/<data_source_name>/`.
6. Run `make generate` to generate documentation.

## Jamf School UI ↔ Terraform Field Mapping

Terraform attribute names should match the Jamf School admin UI wherever possible.

### Resource Navigation Paths

| Terraform Resource         | UI Navigation Path                              |
| -------------------------- | ----------------------------------------------- |
| `jamfschool_user`          | Users > Users > + Add User                      |
| `jamfschool_user_group`    | Users > Groups > + Add Group                    |
| `jamfschool_device_group`  | Devices > Device Groups > + Add Group            |
| `jamfschool_class`         | Classes > Add Class                              |
| `jamfschool_ibeacon`       | Organization > Settings > iBeacons > Add iBeacon |

### Key Field Mappings

| Terraform Attribute              | API JSON Key        | UI Label                                    |
| -------------------------------- | ------------------- | ------------------------------------------- |
| **User**                         |                     |                                             |
| `exclude`                        | `exclude`           | "Don't apply Teacher restrictions"          |
| `store_mail_contacts_calendars`  | `storePassword`     | "Store Mail, Contacts, Calendars credentials locally" |
| `mail_contacts_calendars_domain` | `domain`            | "Domain" (under Mail, Contacts, Calendars)  |
| `member_of`                      | `memberOf`/`groupIds` | Group membership (write: `memberOf`, read: `groupIds`) |
| **User Group**                   |                     |                                             |
| `jamf_school_teacher`            | `acl.teacher`       | "Jamf School Teacher" (Features dropdown)   |
| `jamf_parent`                    | `acl.parent`        | "Jamf Parent" (Features dropdown)           |
| **Device Group**                 |                     |                                             |
| `information`                    | `information`       | "Information"                               |
| `show_in_ios_app`                | `collectionType`    | "Show in iOS app" (radio buttons)           |
| `shared`                         | `shared`            | Location sharing toggle                     |
| **iBeacon**                      |                     |                                             |
| `uuid`                           | `UUID`              | "UUID"                                      |
| `major`                          | `major`             | "Major" (0-65535)                           |
| `minor`                          | `minor`             | "Minor" (0-65535)                           |

### show_in_ios_app Value Mapping

| Terraform Value    | API Value      | UI Label                                      |
| ------------------ | -------------- | --------------------------------------------- |
| `none`             | `none`         | "Do not show"                                 |
| `article`          | `article`      | "Show as article"                             |
| `list`             | `list`         | "Show as list of apps, documents and profiles"|
| `animated_icons`   | `runningTiles` | "Show as animated icons"                      |
