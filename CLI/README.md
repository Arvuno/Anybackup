# foundation-cli

`foundation-cli` is a unified command-line tool for the Foundation control plane, providing standardized access to capabilities for protected objects, policies, jobs, clients, storage, networks, and application-specific features. It is designed for manual operations, script integration, and Agent orchestration calls.

## Features
- Unified entry: Consolidates control capabilities scattered across multiple backend interfaces into a single CLI.
- AK/SK authentication: Uniformly uses `endpoint + tenant-id + ak + sk` to initiate signed requests.
- Agent-oriented: Returns stable structures, suitable for scripts, automated tasks, and upper-layer Agent calls.
- Clear business domains: Organizes commands by domains such as `protect`, `mysql`, `vmware`, `host`, `job`, `client`, etc.
- Raw result friendly: Retains the backend `status/error/responseData` structure by default for further orchestration.

## Feature List

| Category | Capabilities |
|---|---|
| Client Management | Deploy client agents, view host list, deploy jobs, view deployment job execution output, view client and agent lists |
| Job Management | Stop jobs, delete jobs, view job list, backup job details, job execution output, sub-job list |
| MYSQL | MySQL object and data source query, configure backup settings, initiate recovery and backup, view backup configuration/recovery point-in-time/recovery and backup job configuration details, authorization |
| Policy Management | Create, delete backup policies, view policy list |
| Storage Pool Management | Create, delete storage pools, get storage pool list |
| Point-in-Time Management | Clean up point-in-time copies, view point-in-time list |

## Installation

### Environment Requirements
- Go 1.22 or higher
- Access to the target Foundation control plane
- Valid `ak`, `sk` and target `endpoint` (`tenant-id` is optional)

### Build from Source
Windows:

```powershell
go build -o foundation-cli.exe .\cmd\foundation-cli
```

macOS / Linux:

```bash
go build -o foundation-cli ./cmd/foundation-cli
```

### Check Version

```bash
foundation-cli version
```

The default target version supported by the current source code is `9.0.9.0`.

## Quick Start

### 1. View Root Command Help
```bash
foundation-cli --help
```

### 2. Execute a Read-Only Command
```bash
foundation-cli job logs \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```

### 3. Debug Unstandardized Capabilities Using Passthrough API
```bash
foundation-cli api \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --path "/job_center/v1/business_types"
```

## Global Parameters

Most business commands inherit this set of remote access parameters:

| Parameter | Description |
|---|---|
| `--tenant-id` | Tenant ID (optional) |
| `--endpoint` | Foundation control plane address |
| `--ak` | Access Key |
| `--sk` | Secret Key |
| `--target-version` | Target version, default `9.0.9.0` |

## Command Overview

The top-level command domains in the current source code are as follows:

| Command Domain | Description |
|---|---|
| `protect` | General protection operation capabilities |
| `policy` | Protection policy query and management |
| `timepoint` | General object point-in-time capabilities |
| `job` | Runtime observation capabilities such as job list, details, logs, sub-job list, etc. |
| `client` | Client, deployment, Runner, data source capabilities |
| `mysql` | MySQL-specific objects, recovery, point-in-time, data source capabilities |
| `network` | Network subnet and node query capabilities |
| `storage` | Storage service, storage pool, node, device capabilities |
| `version` | Output CLI version information |

You can continue to use `foundation-cli <domain> --help` to view the specific subcommands of a command domain.

## Common Usage Scenarios

### Query Job List
```bash
foundation-cli job list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

### Query Sub-Job List
```bash
foundation-cli job child list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```

### Query MySQL Point-in-Time
```bash
foundation-cli mysql timepoint list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

### Query Storage Service
```bash
foundation-cli storage service list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## Output and Error Handling
- Usually returns `status=success` when successful.
- Retains the backend `errorCode / errorArgs` when failed.
- The command exit code is uniformly managed by the CLI, suitable for script execution result judgment.

Example:

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [],
    "totalNum": 0
  }
}
```

## Project Structure

```text
cmd/        Cobra entry and root commands
internal/   Business commands, request construction, signature and output implementation
skill/      Agent-oriented command routing and single command documentation
scripts/    Validation and auxiliary scripts
tests/      Unit tests
docs/       Architecture, design and reference materials
```

## Agent / Skill Documentation

If you want Agents to call this tool more stably, it is recommended to directly view the entry documentation in the `skill/` directory, for example:

- [skill/job/SKILL.md](./skill/job/SKILL.md)
- [skill/mysql/SKILL.md](./skill/mysql/SKILL.md)
- [skill/protect/SKILL.md](./skill/protect/SKILL.md)

These documents are more suitable for Agents to perform command selection, parameter completion and single command navigation.

## Development

### Run Tests

```bash
go test ./...
```

If you only want to verify a certain command domain, it is recommended to run focused tests, for example:

```bash
go test ./tests/unit/business/domains -run '^TestJobReadMappings$' -count=1
```

### Documentation Validation

For single command documentation, you can use the `doc-checker` script for structure checking:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\doc-checker\references\check-command-doc.ps1 -Path .\skill\job\references\commands\child-list.md
```

## References

- [Foundation CLI Management Briefing](./docs/references/2026-04-08-foundation-cli-management-brief.md)

## License and Third-Party Notices

- This repository is distributed under the root repository [LICENSE](../LICENSE) (SSPL-1.0) together with the root [NOTICE](../NOTICE).
- Third-party Go components statically linked into the built CLI binary (`cobra`, `pflag`, `mousetrap`) are declared uniformly in [THIRD_PARTY_NOTICES.md](./THIRD_PARTY_NOTICES.md).
- When adding or upgrading any direct/indirect dependency in `go.mod`, `THIRD_PARTY_NOTICES.md` MUST be updated in the same commit.
