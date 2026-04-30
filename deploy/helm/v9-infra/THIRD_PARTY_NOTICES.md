# Third-Party Open Source Notices — v9-infra

This notice enumerates the third-party open source components introduced by the `deploy/helm/v9-infra` Helm chart at runtime, together with their respective license information. This chart contains no proprietary application code; its templates and values are distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository.

## 1. Scope

In scope:

- All upstream container images referenced by this chart and deployed into the cluster.

Out of scope:

- Services in other charts that connect to the infrastructure deployed here (e.g., conversation-service connecting to PostgreSQL) do so over standard wire protocols only. Their notices live in their respective chart directories.

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location | Upstream Project |
|---|---|---|---|---|---|
| PostgreSQL | 17 | PostgreSQL License (permissive) | Runtime container image | [values.yaml](./values.yaml) (`postgres.image: library/postgres:17`) | https://www.postgresql.org/ |
| RabbitMQ (with management plugin) | 3-management | Mozilla Public License 2.0 (MPL-2.0) | Runtime container image | [values.yaml](./values.yaml) (`rabbitmq.image: library/rabbitmq:3-management`) | https://www.rabbitmq.com/ |
| Redis | 7 | BSD 3-Clause License (≤ 7.2.x) / RSALv2 + SSPLv1 (≥ 7.4.x) — see §3.3 | Runtime container image | [values.yaml](./values.yaml) (`redis.image: library/redis:7`) | https://redis.io/ |
| OpenSearch | 3.6.0 | Apache License 2.0 | Runtime container image | [values.yaml](./values.yaml) (`opensearch.image: opensearchproject/opensearch:3.6.0`) | https://opensearch.org/ |
| busybox | latest (init container) | GPL v2 | Runtime init container image for OpenSearch | [values.yaml](./values.yaml) (`opensearch.initImage: library/busybox:latest`) | https://busybox.net/ |

## 3. License Notices and References

### 3.1 PostgreSQL — PostgreSQL License

Copyright © 1996–2026 The PostgreSQL Global Development Group
Copyright © 1994 The Regents of the University of California

Permission to use, copy, modify, and distribute this software and its documentation for any purpose, without fee, and without a written agreement is hereby granted, provided that the above copyright notice and this paragraph and the following two paragraphs appear in all copies.

Canonical license text: https://www.postgresql.org/about/licence/

### 3.2 RabbitMQ — Mozilla Public License 2.0

Copyright © 2007-2026 Broadcom. All Rights Reserved. The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

RabbitMQ is distributed under the Mozilla Public License 2.0. The MPL-2.0 is a file-level copyleft license; any modifications to covered source files are subject to the MPL-2.0. Distributing unmodified binaries (as this chart does) does not require releasing source code. The full license text is available at:

    https://www.mozilla.org/en-US/MPL/2.0/

This chart does NOT modify RabbitMQ. Keycloak is deployed as its original container image.

Canonical license text: https://github.com/rabbitmq/rabbitmq-server/blob/main/LICENSE

### 3.3 Redis — License note (version-dependent)

**Redis ≤ 7.2.x** is distributed under the BSD 3-Clause License.  
**Redis ≥ 7.4.0** is distributed under a dual RSALv2 + SSPLv1 license.

The floating tag `redis:7` used in this chart may resolve to any 7.x patch release. To lock the applicable license, **pin the image to a specific patch version** in `values.yaml` (e.g., `library/redis:7.2.7`) and update this file accordingly.

BSD 3-Clause canonical text: https://redis.io/docs/latest/operate/oss_and_stack/license/  
RSALv2 canonical text: https://redis.io/legal/rsalv2-agreement/  
SSPLv1 canonical text: https://redis.io/legal/server-side-public-license-v1/

### 3.4 OpenSearch — Apache License 2.0

Copyright OpenSearch Contributors

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

Canonical license text: https://github.com/opensearch-project/OpenSearch/blob/main/LICENSE.txt

### 3.5 busybox — GPL v2

busybox is licensed under the GNU General Public License version 2.

    https://www.gnu.org/licenses/old-licenses/gpl-2.0.html

This chart uses busybox solely as a disposable init container (for sysctl/chown setup before the OpenSearch pod starts). The busybox binary is not embedded in any application image built by this repository; it is pulled directly from the upstream image registry at deployment time.

Canonical license text: https://git.busybox.net/busybox/tree/LICENSE

## 4. Compliance and Redistribution

- This chart does NOT modify any of the distributed images. Each component is deployed using its original upstream container image; only runtime configuration is passed via Kubernetes environment variables and `ConfigMap` objects.
- Any redistribution of this chart, or of deployment artifacts derived from it, MUST carry this notice together with the license texts of all components in scope.
- In accordance with Apache License 2.0 §4, the OpenSearch image is redistributed as an unmodified binary; this chart is NOT considered a derivative work of OpenSearch. Copyright and liability for OpenSearch remain with the upstream contributors.
- The `redis:7` floating tag introduces license uncertainty (see §3.3). Before redistributing in production, pin to a specific patch version and verify the applicable license.

## 5. Maintenance Rules

- When any image tag or version in the Inventory is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this chart introduces a new third-party dependency (for example, a new sidecar, init container, or Helm subchart), the component, version, license and upstream link MUST be registered in this file.
- The `redis:7` floating tag MUST be replaced with a pinned patch version (e.g., `redis:7.2.7`) to remove license ambiguity; update §3.3 accordingly at that time.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
