# Third-Party Open Source Notices — conversation-service

This notice enumerates the third-party open source components introduced by the `deploy/helm/conversation` Helm chart, both at build time of the distributed container image and at runtime, together with their respective license information. The chart templates and the conversation-service application code are distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository, and are therefore out of scope for this file.

## 1. Scope

In scope:

- The base container image pulled and distributed with the conversation-service image.
- Python runtime libraries installed into the distributed image according to [`requirements.txt`](../../../Agent/service/conversation/requirements.txt) / [`pyproject.toml`](../../../Agent/service/conversation/pyproject.toml).

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- External PostgreSQL, RabbitMQ and Redis instances are accessed only over their respective wire protocols. No source code or binaries of those systems are embedded in this chart or the conversation-service image. They are therefore **not** declared here.
- Keycloak is accessed indirectly through the API Gateway over the OIDC protocol; its notice lives in [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](../auth/THIRD_PARTY_NOTICES.md).

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location | Upstream Project |
|---|---|---|---|---|---|
| Python (CPython) | 3.12 (slim) | Python Software Foundation License 2.0 | Runtime base container image | [Dockerfile](../../../Agent/service/conversation/Dockerfile) (`FROM python:3.12-slim`) | https://www.python.org/ |
| Debian slim base layer | Bookworm (tracks python:3.12-slim) | Various (GPL/LGPL/BSD/MIT per package) | Runtime base container image | [Dockerfile](../../../Agent/service/conversation/Dockerfile) | https://www.debian.org/ |
| FastAPI | >=0.128.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/fastapi/fastapi |
| Uvicorn (with `standard` extras) | >=0.38.0 | BSD 3-Clause License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/encode/uvicorn |
| Pydantic | >=2.12.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/pydantic/pydantic |
| pydantic-settings | >=2.12.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/pydantic/pydantic-settings |
| PyYAML | >=6.0.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/yaml/pyyaml |
| dependency-injector | >=4.48.0 | BSD 3-Clause License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/ets-labs/python-dependency-injector |
| structlog | >=25.0.0 | MIT / Apache License 2.0 (dual) | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/hynek/structlog |
| opentelemetry-api | >=1.39.0 | Apache License 2.0 | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/open-telemetry/opentelemetry-python |
| opentelemetry-sdk | >=1.39.0 | Apache License 2.0 | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/open-telemetry/opentelemetry-python |
| redis-py | >=7.0.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/redis/redis-py |
| aio-pika | >=9.5.0 | Apache License 2.0 | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/mosquito/aio-pika |
| SQLAlchemy (with `asyncio` extras) | >=2.0.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/sqlalchemy/sqlalchemy |
| Alembic | >=1.17.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/sqlalchemy/alembic |
| aiosqlite | >=0.21.0 | MIT License | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/omnilib/aiosqlite |
| asyncpg | >=0.31.0 | Apache License 2.0 | Python runtime dependency | [requirements.txt](../../../Agent/service/conversation/requirements.txt) | https://github.com/MagicStack/asyncpg |

## 3. License Notice and References

### 3.1 Python (CPython) — PSF License 2.0

Copyright © 2001-2026 Python Software Foundation; All Rights Reserved.

Canonical license text: https://docs.python.org/3/license.html

### 3.2 Debian base layer

The `python:3.12-slim` image is built on top of Debian and includes a number of GNU/Debian components under GPL, LGPL, BSD and MIT style licenses. When redistributing images built from this chart, the complete set of upstream notices shipped inside `/usr/share/doc/*/copyright` of the image MUST be preserved as the canonical source.

Upstream references:

- https://www.debian.org/legal/licenses/
- https://hub.docker.com/_/python

### 3.3 Apache License 2.0 components

Components listed with `Apache License 2.0` above are governed by the standard Apache License 2.0 text, available at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

### 3.4 MIT / BSD components

Components listed with `MIT License` or `BSD 3-Clause License` retain the copyright notices of their respective upstream projects as linked in the Inventory table. The full license texts are available in the source distributions of each upstream project and MUST be carried with any redistribution of images built from this chart.

## 4. Compliance and Redistribution

- This chart does NOT modify Python, the Debian base image or any of the listed runtime libraries. The conversation-service image simply installs these components via `pip` on top of `python:3.12-slim`.
- Any redistribution of this chart, or of images/artifacts built on top of this chart, MUST carry this notice together with the license texts of all components in scope.
- This chart is NOT considered a derivative work of Python, Debian or any listed library. Copyright and liability for those components remain with their upstream owners and contributors.

## 5. Maintenance Rules

- When the base image tag or any component in the Inventory is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this chart or the underlying image introduces a new third-party dependency (for example, a new Python package in `requirements.txt`, a sidecar, init container, or Helm subchart), the component, version, license and upstream link MUST be registered in this file.
- Keep the split consistent with [`deploy/helm/core_agent/THIRD_PARTY_NOTICES.md`](../core_agent/THIRD_PARTY_NOTICES.md) and [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](../web/THIRD_PARTY_NOTICES.md): dependencies of the conversation-service image live in this file; dependencies of other images live next to their own charts.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
