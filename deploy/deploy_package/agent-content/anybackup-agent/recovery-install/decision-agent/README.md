# Decision Agent Installer MVP

This directory contains a minimal Ubuntu-oriented installer for a BKN-backed decision agent.

## Scope

The MVP installer does the following:

1. Login to KWeaver with the CLI
2. Check two fixed business knowledge networks by exact name and expected ID
3. Resolve and check the fixed `contextloader工具集` skill
4. Load the long `system_prompt` and `dolphin` texts from files
5. Generate a temporary agent config JSON
6. Create the decision agent
7. Publish the decision agent

The agent identity is fixed in this MVP:

- Name: `MySQL恢复Agent`
- Profile: `MySQL数据库恢复智能化 Agent`

The knowledge source is also fixed in this MVP:

- `MySQL数据库恢复经验知识网络` (`mysql_recovery_experience`)
- `MySQL数据库恢复运行知识网络` (`mysql_recovery_run`)

The attached agent skill is also fixed in this MVP:

- `contextloader工具集` (resolved from the published skill list at install time)

The task planning switch is also fixed in this MVP:

- `任务规划` = enabled

The MVP uninstaller removes the agent only. It does not delete the linked knowledge network.

## Directory Layout

```text
install/decision-agent/
- install.sh
- uninstall.sh
- README.md
- prompts/
  - system_prompt.md
  - dolphin.txt
```

## Files You Must Customize

Before running the installer, replace the placeholder content in:

- `prompts/system_prompt.md`
- `prompts/dolphin.txt`

The installer stops if either file still contains the token `TODO_REPLACE_ME`.

## Prerequisites

- Ubuntu or another Linux distribution with Bash
- `kweaver` CLI available in `PATH`
- `python3` available in `PATH`
- KWeaver account credentials
- A published skill named `contextloader工具集` must already exist in the target business domain
- The following knowledge networks must already exist in the target business domain:
  - `MySQL数据库恢复经验知识网络` (`mysql_recovery_experience`)
  - `MySQL数据库恢复运行知识网络` (`mysql_recovery_run`)

## Quick Start

### Install the agent

```bash
cd install/decision-agent

./install.sh \
  --kweaver-ip "115.190.186.186" \
  --username "admin" \
  --password "your-password" \
  --biz-domain "bd_public" \
  --insecure
```

## Environment Variable Mode

You can also provide values through environment variables:

```bash
export KWEAVER_IP="115.190.186.186"
export KWEAVER_USERNAME="admin"
export KWEAVER_PASSWORD="your-password"
export KWEAVER_BUSINESS_DOMAIN="bd_public"

./install.sh --insecure
```

## Uninstall

```bash
cd install/decision-agent

./uninstall.sh \
  --kweaver-ip "115.190.186.186" \
  --username "admin" \
  --password "your-password" \
  --biz-domain "bd_public" \
  --insecure
```

The uninstaller resolves the created decision agent by the fixed name `MySQL恢复Agent`.
If needed, you can still override that resolution with `--agent-id <id>`.

## Current MVP Assumptions

- The installer binds two fixed knowledge networks to the agent.
- The installer resolves one fixed published skill by exact name: `contextloader工具集`.
- The installer checks both knowledge network name and expected ID before creating the agent.
- The installer enables `plan_mode.is_enabled = true`.
- The installer assumes the two fixed knowledge networks already exist in the target KWeaver environment.
- The installer publishes the agent unless `--no-publish` is used.
- The current CLI publish flow is safest with the default business domain `bd_public`.
- Troubleshooting is based on direct CLI stdout/stderr instead of extra runtime logs.
- The uninstaller resolves the installed decision agent by exact name inside the target business domain.
