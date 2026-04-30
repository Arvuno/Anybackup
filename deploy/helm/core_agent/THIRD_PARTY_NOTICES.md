# Third-Party Open Source Notices — core-agent-service

This notice enumerates the third-party open source components introduced by the `deploy/helm/core_agent` Helm chart, both at build time of the distributed container image and at runtime, together with their respective license information. The chart templates and the core-agent-service application code are distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository, and are therefore out of scope for this file.

## 1. Scope

In scope:

- The base container image pulled and distributed with the core-agent-service image.
- Python runtime libraries installed into the distributed image according to [`pyproject.toml`](../../../Agent/service/core_agent/pyproject.toml).

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- External PostgreSQL and RabbitMQ instances are accessed only over their respective wire protocols. No source code or binaries of those systems are embedded in this chart or the core-agent-service image. They are therefore **not** declared here.
- KWeaver is accessed via `kweaver-sdk`, which is listed in the Inventory below.
- Anybackup Foundation is accessed over HTTPS via the Anybackup CLI; its dependencies are declared in [`CLI/THIRD_PARTY_NOTICES.md`](../../../CLI/THIRD_PARTY_NOTICES.md).

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location | Upstream Project |
|---|---|---|---|---|---|
| Python (CPython) | 3.12 (slim) | Python Software Foundation License 2.0 | Runtime base container image | [Dockerfile](../../../Agent/service/core_agent/Dockerfile) (`FROM python:3.12-slim`) | https://www.python.org/ |
| Debian slim base layer | Bookworm (tracks python:3.12-slim) | Various (GPL/LGPL/BSD/MIT per package) | Runtime base container image | [Dockerfile](../../../Agent/service/core_agent/Dockerfile) | https://www.debian.org/ |
| SQLAlchemy | >=2.0 | MIT License | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | https://github.com/sqlalchemy/sqlalchemy |
| psycopg (with `binary` extras) | >=3.2 | LGPL-3.0-or-later (see §4 for redistribution implications) | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | https://github.com/psycopg/psycopg |
| Pydantic | >=2.10 | MIT License | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | https://github.com/pydantic/pydantic |
| pydantic-settings | >=2.8 | MIT License | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | https://github.com/pydantic/pydantic-settings |
| kweaver-sdk | >=0.6.10 | See upstream distribution (to be confirmed with the upstream project) | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | Upstream project page (internal) |
| pika | >=1.3 | BSD 3-Clause License | Python runtime dependency | [pyproject.toml](../../../Agent/service/core_agent/pyproject.toml) | https://github.com/pika/pika |

## 3. License Notice and References

### 3.1 Python (CPython) — PSF License 2.0

Copyright © 2001-2026 Python Software Foundation; All Rights Reserved.

Canonical license text: https://docs.python.org/3/license.html

### 3.2 Debian base layer

The `python:3.12-slim` image is built on top of Debian and includes a number of GNU/Debian components under GPL, LGPL, BSD and MIT style licenses. When redistributing images built from this chart, the complete set of upstream notices shipped inside `/usr/share/doc/*/copyright` of the image MUST be preserved as the canonical source.

Upstream references:

- https://www.debian.org/legal/licenses/
- https://hub.docker.com/_/python

### 3.3 psycopg — LGPL-3.0-or-later

The `psycopg` PostgreSQL driver and its `binary` extras are distributed under the GNU Lesser General Public License, version 3 or later. Canonical license text:

    https://www.gnu.org/licenses/lgpl-3.0.txt

See §4 for the obligations this triggers when redistributing the core-agent-service image.

### 3.4 MIT / BSD components

Components listed with `MIT License` or `BSD 3-Clause License` retain the copyright notices of their respective upstream projects as linked in the Inventory table. The full license texts are available in the source distributions of each upstream project and MUST be carried with any redistribution of images built from this chart.

### 3.5 kweaver-sdk

`kweaver-sdk` is an internal SDK distributed through a private Python package index. Its license and redistribution terms MUST be confirmed with the upstream project owner before redistributing any image that embeds it, and the "License" column in the Inventory MUST be updated accordingly once confirmed.

## 4. Compliance and Redistribution

- This chart does NOT modify Python, the Debian base image or any of the listed runtime libraries. The core-agent-service image simply installs these components via `pip` on top of `python:3.12-slim`.
- Because `psycopg` is distributed under LGPL-3.0-or-later, any redistribution of images built from this chart MUST:
  - carry the full text of the LGPL-3.0 license and the upstream copyright notice;
  - preserve the ability for downstream users to replace the `psycopg` library with a modified version, in line with §4 of the LGPL.
- Any redistribution of this chart, or of images/artifacts built on top of this chart, MUST carry this notice together with the license texts of all components in scope.
- This chart is NOT considered a derivative work of Python, Debian or any listed library. Copyright and liability for those components remain with their upstream owners and contributors.

## 5. Maintenance Rules

- When the base image tag or any component in the Inventory is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this chart or the underlying image introduces a new third-party dependency (for example, a new Python package, a sidecar, init container, or Helm subchart), the component, version, license and upstream link MUST be registered in this file.
- If `psycopg[binary]` is ever replaced with a non-LGPL alternative (for example, a pure Python driver), §3.3 and §4 MUST be updated in the same change.
- Keep the split consistent with [`deploy/helm/conversation/THIRD_PARTY_NOTICES.md`](../conversation/THIRD_PARTY_NOTICES.md) and [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](../web/THIRD_PARTY_NOTICES.md): dependencies of the core-agent-service image live in this file; dependencies of other images live next to their own charts.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
