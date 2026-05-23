# Third-Party Open Source Notices — agent-web

This notice enumerates the third-party open source components introduced by the `Agent/portal/deploy/helm/agent-web` Helm chart, both at build time of the distributed SPA container image and at runtime, together with their respective license information. The chart templates and the first-party web application source code (under `repos/Anybackup/Agent/portal`) are distributed under SSPL-1.0 together with the root [LICENSE](../../../../../LICENSE) of this repository, and are therefore out of scope for this file.

## 1. Scope

The `agent-web` chart and the `web-service` chart under [`deploy/helm/web`](../../../../../deploy/helm/web) deploy the same SPA image built from [`Agent/portal`](../../../../../Agent/portal). To avoid duplicate maintenance, the authoritative third-party component inventory and license texts for this image live in:

- [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](../../../../../deploy/helm/web/THIRD_PARTY_NOTICES.md)

This file incorporates that notice by reference and MUST be kept in sync with it.

In scope here (identical to the web chart):

- The base container images used to build and run the SPA image (`node:22-alpine`, `nginx:1.30-alpine`, Alpine Linux components).
- The Node.js runtime libraries bundled into the SPA by `vite build` and served as static assets inside the distributed image.

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- This chart talks to Auth Service and Conversation Service through the shared API Gateway over HTTP. This is a runtime protocol interop and does not constitute a dependency on those services' source code or binaries.
- Keycloak is accessed indirectly through the API Gateway over the OIDC protocol; its notice lives in [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](../../../../../deploy/helm/auth/THIRD_PARTY_NOTICES.md).

## 2. Third-Party Component Inventory

The authoritative inventory table is maintained in [`deploy/helm/web/THIRD_PARTY_NOTICES.md` §2](../../../../../deploy/helm/web/THIRD_PARTY_NOTICES.md#2-third-party-component-inventory). The entries cover:

- Build-time base image: `node:22-alpine`.
- Runtime base image: `nginx:1.30-alpine`.
- Alpine Linux base package set.
- SPA runtime npm dependencies from [`Agent/portal/package.json`](../../../package.json): `react`, `react-dom`, `react-router-dom`, `framer-motion`, `lucide-react`, `react-markdown`, `remark-gfm`, `class-variance-authority`, `clsx`, `tailwind-merge`, `tailwindcss-animate`, `zustand`, `github-markdown-css`.

When a new dependency is introduced or an existing one is upgraded in `Agent/portal`, the change MUST be reflected in the authoritative file first, and this file MUST then be re-checked to ensure the cross-reference is still accurate.

## 3. License Notice and References

The canonical license texts and per-component references are reproduced in [`deploy/helm/web/THIRD_PARTY_NOTICES.md` §3](../../../../../deploy/helm/web/THIRD_PARTY_NOTICES.md#3-license-notice-and-references). Any redistribution of images built from this chart MUST carry those license texts alongside this notice.

## 4. Compliance and Redistribution

- This chart does NOT modify Node.js, Nginx, the Alpine base image or any of the listed JavaScript libraries. The agent-web image simply runs a static SPA bundle on top of `nginx:1.30-alpine`.
- Any redistribution of this chart, or of images/artifacts built on top of this chart, MUST carry this notice together with the authoritative license texts listed in [`deploy/helm/web/THIRD_PARTY_NOTICES.md` §3](../../../../../deploy/helm/web/THIRD_PARTY_NOTICES.md#3-license-notice-and-references).
- This chart is NOT considered a derivative work of Node.js, Nginx, Alpine or any listed library. Copyright and liability for those components remain with their upstream owners and contributors.

## 5. Maintenance Rules

- Changes to the SPA dependency set MUST first be recorded in [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](../../../../../deploy/helm/web/THIRD_PARTY_NOTICES.md); this file is then validated for drift in the same change.
- When the base image tag or any npm dependency in the authoritative Inventory is upgraded, the web notice's "Version" column MUST be updated in the same change; this file requires no textual update as long as the cross-reference remains valid.
- If the two charts diverge in the future and deploy different images, the authoritative Inventory and License sections MUST be duplicated into this file and §1–§3 above updated accordingly.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
