# BKN Installer MVP

This directory contains a minimal Ubuntu-oriented installer for importing and deleting two fixed recovery knowledge networks.

## Scope

The BKN installer does the following:

1. Login to KWeaver with the CLI
2. Validate the local BKN directories
3. Import two fixed knowledge networks into KWeaver
4. Configure indexes for every data property in every imported object type

The BKN uninstaller does the following:

1. Login to KWeaver with the CLI
2. Resolve the two fixed knowledge networks by exact name
3. Verify the expected ID
4. Delete the matching knowledge networks

## Fixed BKN Sources

The installer imports these local directories:

- `examples/bkn/recovery_experience/database`
- `examples/bkn/recovery_run/database`

The imported knowledge networks are:

- `MySQL数据库恢复经验知识网络` (`mysql_recovery_experience`)
- `MySQL数据库恢复运行知识网络` (`mysql_recovery_run`)

## Required Inputs

Both install and uninstall scripts require these inputs for normal use:

1. KWeaver IP or base URL
2. KWeaver username
3. KWeaver password

The scripts expose these arguments:

- `--kweaver-ip`
- `--username`
- `--password`

You can also pass:

- `--biz-domain` (optional, default: `bd_public`)
- `--index-small-model` (optional, vector index small model ID/name; omitted means KWeaver default)
- `--insecure` (optional)

The installer applies this index configuration to every object type data property:

- Keyword index: enabled, length 512 bytes
- Full-text index: enabled, standard tokenizer/analyzer
- Vector index: enabled, small model; `--index-small-model` selects a specific small model, otherwise KWeaver uses its default small model

## Quick Start

### Install the BKNs

```bash
cd install/recovery-bkn

./install.sh \
  --kweaver-ip "115.190.186.186" \
  --username "admin" \
  --password "your-password" \
  --biz-domain "bd_public" \
  --index-small-model "your-small-model-id" \
  --insecure
```

### Uninstall the BKNs

```bash
cd install/recovery-bkn

./uninstall.sh \
  --kweaver-ip "115.190.186.186" \
  --username "admin" \
  --password "your-password" \
  --biz-domain "bd_public" \
  --insecure
```

## Notes

- The scripts use fixed local paths and fixed knowledge network IDs.
- The installer verifies the imported ID after each `bkn push`.
- Index configuration runs only after a BKN import succeeds.
- The uninstaller skips a network if it is already missing.
