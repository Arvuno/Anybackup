# Third-Party Open Source Notices — api-gateway-service

This notice enumerates the third-party open source components introduced by the `deploy/helm/api_gateway` Helm chart, both at build time and at runtime, together with their respective license information. The chart itself, as well as the in-house Go plugin source code under `plugins/token-to-x-user/`, is distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository, and is therefore out of scope for this file.

## 1. Scope

In scope:

- Upstream Helm charts, container images and embedded open source components that this chart depends on.
- Third-party software binaries deployed to the cluster by this chart.

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- This chart talks to Auth Service through the standard OIDC protocol (an HTTP call to the `userinfo` endpoint). This is a runtime protocol interop and does not constitute a dependency on Keycloak source code or binaries. Keycloak is therefore **not** declared here; its open source notice lives in [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](../auth/THIRD_PARTY_NOTICES.md).

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location in Chart | Upstream Project |
|---|---|---|---|---|---|
| Traefik (Helm Chart) | 39.0.8 | MIT License | Helm chart dependency | [Chart.yaml](./Chart.yaml) / [Chart.lock](./Chart.lock) / `charts/traefik-39.0.8.tgz` | https://github.com/traefik/traefik-helm-chart |
| Traefik (Container Image) | Determined by the `appVersion` of the upstream chart | MIT License | Runtime container image | [values.yaml](./values.yaml) (`traefik.image.repository: library/traefik`) | https://github.com/traefik/traefik |

## 3. License Texts and References

### 3.1 Traefik — MIT License

Copyright (c) 2016-2026 Containous SAS; 2020-2026 Traefik Labs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Canonical license text: https://github.com/traefik/traefik/blob/master/LICENSE.md

## 4. Maintenance Rules

- When the Traefik chart or container image version is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this chart introduces a new third-party dependency (for example, a new Helm subchart or sidecar image), the component, version, license and upstream link MUST be registered in this file.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
