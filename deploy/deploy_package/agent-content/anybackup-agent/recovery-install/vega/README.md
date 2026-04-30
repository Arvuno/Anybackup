# Vega Data Installer

This directory contains an Ubuntu-oriented installer for the recovery experience seed data used by Vega data virtualization.

## Contents

- `install.sh` creates the target database, creates five tables, and imports seed CSV data.
- `install.sh` also creates or updates the Vega PostgreSQL data connection, then scans it.
- `uninstall.sh` deletes the Vega data connection, removes the seed tables, and drops the target database.
- `data/` contains a copy of `examples/bkn/recovery_experience/database/data`.

## Supported Engines

- PostgreSQL, default
- MySQL
- MariaDB

## Usage

PostgreSQL:

```bash
cd install/vega
chmod +x install.sh uninstall.sh

# Uses the default database name: ExperienceBKNDB
./install.sh 127.0.0.1 5432 postgres 'your-password'

# Override the database name
./install.sh 127.0.0.1 5432 custom_experience_db postgres 'your-password'
```

MySQL or MariaDB:

```bash
cd install/vega
chmod +x install.sh uninstall.sh

# Uses the default database name: ExperienceBKNDB
./install.sh 127.0.0.1 3306 root 'your-password' --engine mysql

# Override the database name
./install.sh 127.0.0.1 3306 custom_experience_db root 'your-password' --engine mysql
```

Flag mode:

```bash
./install.sh \
  --host 127.0.0.1 \
  --port 5432 \
  --username postgres \
  --password 'your-password' \
  --engine postgresql
```

The host argument may also be `host:port`. If the port is omitted, it defaults to `5432` for PostgreSQL and `3306` for MySQL/MariaDB.

The default database name is `ExperienceBKNDB`. Override it with `--database <name>` or the four-position-argument form shown above.

## Vega Data Connection

For PostgreSQL installs, the installer creates or updates these Vega data connections after importing the CSV data:

| Connection name | Type | Database | Schema | Protocol |
|-----------------|------|----------|--------|----------|
| `恢复经验知识网络数据连接` | `postgresql` | `ExperienceBKNDB` by default | `public` by default | `jdbc` |
| `CommonServiceDB-client` | `postgresql` | `CommonServiceDB` | `public` by default | `jdbc` |
| `MultiStorageSvcMgmServiceDB-storageservice` | `postgresql` | `MultiStorageSvcMgmServiceDB` | `public` by default | `jdbc` |
| `StorageResMgmServiceDB-poolv8` | `postgresql` | `StorageResMgmServiceDB` | `public` by default | `jdbc` |
| `HyperBackupMgmServiceDB-protectobject` | `postgresql` | `HyperBackupMgmServiceDB` | `public` by default | `jdbc` |
| `HyperJobWorkerServiceDB-job` | `postgresql` | `HyperJobWorkerServiceDB` | `public` by default | `jdbc` |

Both connections use the host, port, username, and password from the installer inputs.

After each connection is created or updated successfully, the installer runs:

```bash
kweaver vega catalog discover <catalog-id> --wait
```

If a connection creation or update fails, scanning is not attempted for that connection and the installer stops.

The Vega connection and scan steps require the `kweaver` CLI to be installed and authenticated before running the installer. Use `--biz-domain <name>` when the target business domain is not `bd_public`.

To skip this step:

```bash
./install.sh 127.0.0.1 5432 postgres 'your-password' --skip-vega-catalog
```

## Uninstall

PostgreSQL:

```bash
cd install/vega

# Uses the default database name and engine:
# database: ExperienceBKNDB
# engine: postgresql
./uninstall.sh 127.0.0.1 5432 postgres 'your-password'

# Explicit database name
./uninstall.sh 127.0.0.1 5432 ExperienceBKNDB postgres 'your-password'
```

MySQL or MariaDB:

```bash
cd install/vega
./uninstall.sh 127.0.0.1 3306 root 'your-password' --engine mysql
```

To remove a non-default database, pass it as the second positional argument or use `--database <name>`:

```bash
./uninstall.sh 127.0.0.1 5432 custom_experience_db postgres 'your-password'
```

The uninstaller first deletes the Vega data connections named `恢复经验知识网络数据连接`, `CommonServiceDB-client`, `MultiStorageSvcMgmServiceDB-storageservice`, `StorageResMgmServiceDB-poolv8`, `HyperBackupMgmServiceDB-protectobject`, and `HyperJobWorkerServiceDB-job`, then drops the five recovery experience tables and the target database. Use `--skip-vega-catalog` to skip the KWeaver/Vega cleanup step.

## Tables

The installer creates and initializes these tables:

- `availability_checkpoint_template`
- `fault_pattern`
- `recovery_capability`
- `recovery_strategy_template`
- `risk_rule`

Re-running the installer truncates only these five target tables before importing CSV data again.

## Requirements

Install the client for the selected database engine before running the script:

```bash
sudo apt-get update
sudo apt-get install -y postgresql-client
```

For MySQL or MariaDB:

```bash
sudo apt-get update
sudo apt-get install -y mysql-client
```
