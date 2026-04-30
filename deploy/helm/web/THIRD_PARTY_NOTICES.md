# Third-Party Open Source Notices — web-service

This notice enumerates the third-party open source components introduced by the `deploy/helm/web` Helm chart, both at build time of the distributed SPA container image and at runtime, together with their respective license information. The chart templates and the first-party web application source code (under `repos/Anybackup/Agent/portal`) are distributed under SSPL-1.0 together with the root [LICENSE](../../../LICENSE) of this repository, and are therefore out of scope for this file.

## 1. Scope

In scope:

- The base container images used to build and run the web-service SPA image.
- The Node.js runtime libraries bundled into the SPA during `vite build` and served as static assets inside the distributed image.

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- This chart talks to Auth Service and Conversation Service through the shared API Gateway over HTTP. This is a runtime protocol interop and does not constitute a dependency on those services' source code or binaries.
- Keycloak is accessed indirectly through the API Gateway over the OIDC protocol; its notice lives in [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](../auth/THIRD_PARTY_NOTICES.md).

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location | Upstream Project |
|---|---|---|---|---|---|
| Node.js | 22 (alpine) | MIT License (Node.js) + additional component licenses | Build-time base container image | [Dockerfile](../../../Agent/portal/Dockerfile) (`FROM node:22-alpine`) | https://github.com/nodejs/node |
| Alpine Linux | Tracks `node:22-alpine` and `nginx:1.30-alpine` | Various (MIT / BSD / GPL per package) | Build-time and runtime base container image | [Dockerfile](../../../Agent/portal/Dockerfile) / [Dockerfile.runtime](../../../Agent/portal/Dockerfile.runtime) | https://alpinelinux.org/ |
| Nginx | 1.30-alpine | BSD 2-Clause License | Runtime base container image | [Dockerfile.runtime](../../../Agent/portal/Dockerfile.runtime) (`FROM nginx:1.30-alpine`) | https://nginx.org/ |
| React | ^18.3.1 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/facebook/react |
| react-dom | ^18.3.1 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/facebook/react |
| react-router-dom | ^7.1.1 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/remix-run/react-router |
| framer-motion | ^11.18.2 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/framer/motion |
| lucide-react | ^0.468.0 | ISC License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/lucide-icons/lucide |
| react-markdown | ^10.1.0 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/remarkjs/react-markdown |
| remark-gfm | ^4.0.1 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/remarkjs/remark-gfm |
| class-variance-authority | ^0.7.1 | Apache License 2.0 | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/joe-bell/cva |
| clsx | ^2.1.1 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/lukeed/clsx |
| tailwind-merge | ^2.6.0 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/dcastil/tailwind-merge |
| tailwindcss-animate | ^1.0.7 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/jamiebuilds/tailwindcss-animate |
| zustand | ^5.0.3 | MIT License | Runtime dependency bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/pmndrs/zustand |
| github-markdown-css | ^5.9.0 | MIT License | Runtime stylesheet bundled into SPA | [package.json](../../../Agent/portal/package.json) | https://github.com/sindresorhus/github-markdown-css |

## 3. License Notice and References

### 3.1 Nginx — BSD 2-Clause License

Copyright (C) 2002-2026 Igor Sysoev, Nginx, Inc.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.

Canonical license text: https://nginx.org/LICENSE

### 3.2 Node.js and Alpine Linux base layers

- Node.js is distributed under the MIT License with additional component licenses for bundled dependencies (OpenSSL, libuv, V8, ICU, etc.). The canonical `LICENSE` file shipped with the Node.js source MUST be preserved when redistributing images built from this chart.
- Alpine Linux and its package set consist of multiple upstream projects under MIT, BSD, ISC, Apache-2.0 and GPL style licenses. When redistributing images built from this chart, the complete set of upstream notices shipped inside `/usr/share/licenses/*` and `/usr/share/doc/*` of the image MUST be preserved as the canonical source.

Upstream references:

- https://github.com/nodejs/node/blob/main/LICENSE
- https://alpinelinux.org/about/

### 3.3 Apache License 2.0 components

Components listed with `Apache License 2.0` above are governed by the standard Apache License 2.0 text, available at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

### 3.4 MIT / BSD / ISC components

Components listed with `MIT License`, `BSD 2-Clause License` or `ISC License` retain the copyright notices of their respective upstream projects as linked in the Inventory table. The full license texts are available in the source distributions of each upstream project and MUST be carried with any redistribution of images built from this chart; in practice, the aggregated third-party notice file produced by the frontend build (for example, the output of a license-checker step) SHOULD be bundled into the image and referenced from this notice.

## 4. Compliance and Redistribution

- This chart does NOT modify Node.js, Nginx, the Alpine base image or any of the listed JavaScript libraries. The web-service image simply runs a static SPA bundle on top of `nginx:1.30-alpine`.
- Any redistribution of this chart, or of images/artifacts built on top of this chart, MUST carry this notice together with the license texts of all components in scope.
- This chart is NOT considered a derivative work of Node.js, Nginx, Alpine or any listed library. Copyright and liability for those components remain with their upstream owners and contributors.

## 5. Maintenance Rules

- When the base image tag or any npm dependency in the Inventory is upgraded, the "Version" column in this file MUST be updated in the same change.
- When the frontend build introduces a new runtime npm dependency that ships in the SPA bundle, the component, version, license and upstream link MUST be registered in this file.
- Keep the split consistent with [`Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md`](../../../Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md): the two charts deploy the same SPA image and SHOULD converge on a single source of truth (this file) by cross-reference.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
