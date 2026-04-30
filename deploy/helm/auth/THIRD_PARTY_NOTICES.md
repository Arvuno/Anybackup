# Third-Party Open Source Notices — auth-service

This notice enumerates the third-party open source components introduced by the `deploy/helm/auth` Helm chart, both at build time and at runtime, together with their respective license information. The chart's own templates and configuration are distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository, and are therefore out of scope for this file.

## 1. Scope

This chart adopts the "Keycloak native Helm deployment" approach: it directly distributes and runs the official Keycloak container image in the cluster as the authentication backend. Keycloak is therefore treated as a distributed open source binary and **is** in scope for this notice.

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location in Chart | Upstream Project |
|---|---|---|---|---|---|
| Keycloak | 26.5.1 | Apache License 2.0 | Runtime container image | [Chart.yaml](./Chart.yaml) (`appVersion: 26.5.1`) / [values.yaml](./values.yaml) (`keycloak.image`) | https://github.com/keycloak/keycloak |

## 3. License Notice and References

### 3.1 Keycloak — Apache License 2.0

Copyright © Red Hat, Inc. and contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Canonical license text: https://github.com/keycloak/keycloak/blob/main/LICENSE.txt

## 4. Compliance and Redistribution

- This chart does NOT modify Keycloak source code or binaries. Keycloak is deployed as its original container image; only startup arguments and external PostgreSQL connection information are passed in via Helm values.
- In accordance with Section 4 of the Apache License 2.0, any redistribution of this chart, or of images/artifacts built on top of this chart, MUST carry this notice and the full text of the Apache License 2.0.
- This chart acts only as a packaging-and-deployment shell for Keycloak and is NOT considered a derivative work of Keycloak. Copyright and liability for Keycloak remain with the upstream owners and contributors.

## 5. Maintenance Rules

- When the Keycloak image `tag` or the chart `appVersion` is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this chart introduces a new third-party dependency (for example, a sidecar, init container, or external Helm subchart), the component, version, license and upstream link MUST be registered in this file.
- Keep the split consistent with [`deploy/helm/api_gateway/THIRD_PARTY_NOTICES.md`](../api_gateway/THIRD_PARTY_NOTICES.md): the Keycloak notice lives in this file; the Traefik notice lives in the api_gateway chart.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
