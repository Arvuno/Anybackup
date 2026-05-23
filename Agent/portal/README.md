English | [中文](README_zh.md)

# AnyBackup Agent Portal Frontend

`portal` is the web frontend workspace of AnyBackup Agent.  
It focuses on the conversation workspace and related operational UI for agent-based backup workflows.

## Positioning

- **What it is**: frontend implementation for Agent portal UX.
- **Primary users**: operators and developers validating conversation-driven workflows.
- **What it is not**: not a standalone product repository and not the full AnyBackup platform.

## Open Source Scope

Currently open in this repo:

- Conversation workspace UI (chat panel, conversation list, interaction flow)
- Structured rich content rendering (including AG-UI layout-tree based rendering path)
- Frontend service layer, state management, and local development toolchain
- Unit tests for critical chat/store/service paths

## Roadmap

Planned incremental open improvements:

- Better extensibility for multimodal cards and action handlers
- Clearer API contract examples and mock data recipes for contributors
- More stable test baselines (component + integration-level)
- Continued cleanup of legacy placeholder routes

## Quick Start

Recommended:

```bash
Node.js 20.x
npm 10.x
```

Install and run:

```bash
npm install
npm run dev
```

Build and test:

```bash
npm run build
npm run test
```

## Development Configuration

For local API proxy, add to `.env.local`:

```bash
VITE_AUTH_SERVICE_PROXY_TARGET=http://<auth-service-host>
VITE_CONVERSATION_SERVICE_PROXY_TARGET=http://<conversation-service-host>
```

Proxy paths:

- `/api/auth_service`
- `/api/conversation_service`

## Project Structure

```text
src/
|-- app/          # app composition and routing
|-- pages/        # route pages
|-- components/   # reusable UI/business components
|-- services/     # API adapters and protocol parsing
|-- store/        # Zustand state and conversation flow
|-- config/       # app configuration
|-- test/         # test setup and helpers
`-- ...
```

## Documentation

- Start with `docs/README.md`
- Engineering baseline and architecture docs are under `docs/`
- Detailed implementation notes are intentionally kept in docs, not this top-level README

## Related Project

- [Anybackup](https://github.com/anybackup-ai/Anybackup): full AI-native data resilience platform built on this ecosystem

## License

This project is distributed under the terms of the repository root [LICENSE](../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../NOTICE).

## Third-Party Notices

When this SPA is packaged into a container image and shipped via a Helm chart, the third-party base images (`node:22-alpine` for build, `nginx:1.30-alpine` for runtime, Alpine Linux) and the npm runtime dependencies bundled by `vite build` are declared in:

- [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](../../deploy/helm/web/THIRD_PARTY_NOTICES.md) — authoritative notice for the shared SPA image.
- [`deploy/helm/agent-web/THIRD_PARTY_NOTICES.md`](./deploy/helm/agent-web/THIRD_PARTY_NOTICES.md) — cross-reference notice for the `agent-web` chart.

Any change to the runtime dependency set in `package.json` MUST be reflected in the authoritative chart notice in the same commit.
