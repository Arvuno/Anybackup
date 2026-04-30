# Anybackup

<p align="center">
  <a href="https://github.com/anybackup-ai/Anybackup/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-SSPL--1.0-blue.svg" alt="License"></a>
  <a href="https://github.com/anybackup-ai/Anybackup/blob/main/VERSION.txt"><img src="https://img.shields.io/badge/version-9.0.0--alpha-orange.svg" alt="Version"></a>
</p>

<p align="center">
  <em>An AI-Native Data Resilience Platform — autonomous backup, recovery, and optimization. Built on an open-source model, delivering 35% lower TCO than traditional approaches.</em>
</p>

<p align="center">
  <a href="./README_zh.md">中文</a>
</p>

---

## 📖 Table of Contents

- [Overview](#overview)
- [Key Capabilities](#key-capabilities)
- [How It Works](#how-it-works)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Repository Structure](#repository-structure)
- [Community](#community)
- [Contributing](#contributing)
- [Related Projects](#related-projects)
- [License & Third-Party Notices](#license--third-party-notices)

---

## Overview

Anybackup V9 is a more cost-effective, intelligent data resilience platform that leverages an open-source business model to deliver the data protection businesses actually need. Its AI backup administrator — powered by the Anybackup Agent — enables autonomous backup, autonomous recovery, and autonomous optimization, cutting total cost of ownership by 35% and replacing reactive firefighting with proactive assurance.

The platform is re-architected with AI-native thinking, centered on an intelligent agent. It follows a **three-in-one architecture: Collect · Execute · Decide** — where each layer plays a distinct role:

- **Anybackup Agent** (Decide) — the thinking and planning agent: intent recognition, plan generation, risk assessment, approval enforcement, and full audit trail. Runs as SaaS.
- **Anybackup Foundation** (Execute) — the execution engine: full/incremental/log backup, retention management, and recovery orchestration. Deployed in your data center.
- **Anybackup Client** (Collect) — Per-workload data collectors on protected assets.

---

## Key Capabilities

- **Natural-language backup recommendations** — Describe your workload, receive 3 ranked backup plans sorted by RPO/RTO
- **Flexible backup granularity** — Full/incremental/log backup at minute, daily, and weekly cadences, adapting to different workloads
- **Application-aware routing** — When multiple instances exist under one resource, the system auto-identifies and guides you to pick the right target
- **Candidate-to-decision loop** — Pick one of 3 plans; the system records your choice and closes the recommendation loop
- **Failure-type recognition** — Recognizes instance crash, database drop, log damage, ransomware and other common fault types; recovery scope confirmed by the operator in conversation
- **Granular recovery paths** — Instance-level, database-level, point-in-time, and log-file recovery
- **Kubernetes-native deployment** — All Agent-side services and infrastructure deployed via Helm charts
- **MySQL first release** — XtraBackup, real-time log, and multi-channel backup; more workloads on the way

---

## How It Works

```
  You speak                Agent reasons              Foundation executes
  ─────────              ─────────────────          ──────────────────────
  "Recommend a       →   Intent recognition      →  Queries current config
   backup plan for       Plan generation              & Foundation capabilities
   our database"         Risk assessment
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
              3 candidate plans     You pick one
              ranked by RPO/RTO     → Agent records,
              & retention              dispatches to
                                       Foundation
```

**Example — Backup recommendation:**

```
You:   "Recommend a backup plan for our finance database."

Agent: 3 candidates based on your environment:

       A — Standard Protection
       Weekly full + daily incremental, 90-day retention.

       B — Low RPO (★ Recommended)
       Weekly full + daily incremental + minute-level log backup,
       5-year retention. Real-time log enabled.

       C — Maximum Simplicity
       Daily full only, 90-day retention.
```

**Example — Recovery from instance crash:**

```
You:   "Production database crashed — won't start. All databases unreachable."

Agent: Instance-level failure.

       Recovery: restore full instance to latest valid backup on original host.
       ⚠️  This overwrites current data. [ Confirm ]  [ Cancel ]
```

---

## Architecture

```
  ┌───────────────────────────────────────────────────┐
  │         Anybackup Agent  ·  SaaS                 │
  │                                                   │
  │  Intent Recognition · Plan Generation · Risk     │
  │  Approval · Human Review · Audit Trail            │
  └──────────────────────▬───────────────────────────┘
                         │
             skill + foundation-cli
                         │
                         ▼
  ┌───────────────────────────────────────────────────┐
  │      Anybackup Foundation  ·  Your Data Center   │
  │                                                   │
  │  Backup: Full · Incremental · Log · Retention    │
  │  Recovery: Instance · Database · Point-in-Time    │
  └──────────────────────▬───────────────────────────┘
                         │
                      Ingest
                         │
                         ▼
  ┌───────────────────────────────────────────────────┐
  │        Anybackup Client  ·  Protected Assets     │
  │                                                   │
  │  MySQL: XtraBackup · Real-time Log · Multi-      │
  │  Channel Backup                                   │
  │  More workloads on the way                        │
  └───────────────────────────────────────────────────┘
```

---

## Getting Started

### Prerequisites

- **Foundation**: deployed and accessible in your data center (contact us for deployment details)
- **Kubernetes**: cluster running for Agent-side services
- **Helm**: 3.x or later
- **Target assets**: MySQL instances running with connectivity to Foundation

### Installation

**Step 1 — Deploy infrastructure services**

```bash
helm install v9-infra ./deploy/helm/v9-infra
```

This provisions Postgres, Redis, RabbitMQ, and OpenSearch for the Agent stack.

**Step 2 — Deploy Agent services**

```bash
# Core agent engine
helm install core-agent ./deploy/helm/core_agent

# Conversation service
helm install conversation ./deploy/helm/conversation

# Auth service (Keycloak)
helm install auth ./deploy/helm/auth

# API Gateway (Traefik)
helm install api-gateway ./deploy/helm/api_gateway

# Web portal
helm install web ./deploy/helm/web
```

**Step 3 — Install the CLI**

```bash
cd CLI
go build -o foundation-cli ./cmd/foundation-cli
```

**Step 4 — Verify connectivity**

```bash
foundation-cli version
```

**Step 5 — Start a conversation**

Open the Agent Web Portal and describe your backup need in plain language, or use the CLI directly:

```bash
anybackup chat "Recommend a backup plan for my production database"
```

For detailed configuration options, see each chart's `values.yaml` in [`deploy/helm/`](./deploy/helm/).

---

## Repository Structure

```
Anybackup/
├── Agent/                     # AI decision layer
│   ├── portal/                #   Web frontend (React)
│   ├── service/               #   Backend services
│   │   ├── conversation/      #     Conversation orchestration
│   │   └── core_agent/        #     AI agent engine
│   ├── skills/                #   Agent skills & AG-UI design patterns
│   └── knowledge/             #   Knowledge networks (backup & recovery BKN)
├── CLI/                       # foundation-cli — command-line tooling
├── deploy/                    # Helm charts & deployment configuration
│   └── helm/
│       ├── v9-infra/          #   Infrastructure (Postgres, Redis, RabbitMQ, OpenSearch)
│       ├── api_gateway/       #   API Gateway (Traefik)
│       ├── auth/              #   Auth service (Keycloak)
│       ├── conversation/      #   Conversation service
│       ├── core_agent/        #   Core agent service
│       └── web/               #   Web SPA
├── LICENSE                    # SSPL-1.0
├── NOTICE                     # Copyright notice
└── VERSION.txt                # Current version
```

---

## Community (coming soon)

- **Issues**: [GitHub Issues](https://github.com/anybackup-ai/Anybackup/issues)
- **Discussions**: [GitHub Discussions](https://github.com/anybackup-ai/Anybackup/discussions)

---

## Contributing

We welcome contributions. Before submitting a pull request:

1. Open an issue to discuss the change or feature.
2. Follow the repository coding conventions.
3. Update relevant documentation and third-party notices if dependencies change.
4. Ensure all tests pass.

For contribution guidelines, see [CONTRIBUTING.md](./CONTRIBUTING.md) (coming soon).

---

## Related Projects

| Project | Relationship |
|---|---|
| **[Agent Portal](./Agent/portal/)** | Web frontend for conversation-based backup workflows |
| **[foundation-cli](./CLI/)** | CLI for Foundation control-plane operations |

---

## License & Third-Party Notices

- Source code: [LICENSE](./LICENSE) (SSPL-1.0) with [NOTICE](./NOTICE)
- Third-party components are declared per distribution unit:
  - [`CLI/THIRD_PARTY_NOTICES.md`](./CLI/THIRD_PARTY_NOTICES.md) — Go dependencies (CLI binary)
  - [`deploy/helm/v9-infra/THIRD_PARTY_NOTICES.md`](./deploy/helm/v9-infra/THIRD_PARTY_NOTICES.md) — Infrastructure (PostgreSQL, RabbitMQ, Redis, OpenSearch)
  - [`deploy/helm/api_gateway/THIRD_PARTY_NOTICES.md`](./deploy/helm/api_gateway/THIRD_PARTY_NOTICES.md) — API Gateway (Traefik)
  - [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](./deploy/helm/auth/THIRD_PARTY_NOTICES.md) — Auth (Keycloak)
  - [`deploy/helm/conversation/THIRD_PARTY_NOTICES.md`](./deploy/helm/conversation/THIRD_PARTY_NOTICES.md) — Conversation Service
  - [`deploy/helm/core_agent/THIRD_PARTY_NOTICES.md`](./deploy/helm/core_agent/THIRD_PARTY_NOTICES.md) — Core Agent Service
  - [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](./deploy/helm/web/THIRD_PARTY_NOTICES.md) — Web SPA
  - [`Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md`](./Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md) — Agent Web portal

When adding new images or upgrading dependencies, update the corresponding notice in the same commit.
