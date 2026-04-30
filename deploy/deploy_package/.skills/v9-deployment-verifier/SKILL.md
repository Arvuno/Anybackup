---
name: v9-deployment-verifier
description: Verify V9 deployment readiness and deployment regressions after changes to install scripts, Ansible roles, Helm charts, smoke checks, or middleware topology. Use when Codex modifies files under `install.sh`, `ansible/`, `helm-chart/`, `images/load-images.sh`, `foundation/`, `init-scripts/`, or deployment validation scripts and needs to run the standard local validation flow.
---

# V9 Deployment Verifier

Run the repository validator after every deployment-related change:

```bash
python scripts/validate_deploy.py --mode regression
```

Use release mode before claiming the package is customer-installable:

```bash
python scripts/validate_deploy.py --mode release
```

## What the validator covers

- Required deployment files exist.
- `install.sh` still points to existing companion scripts.
- Helm lint and template succeed for `v9-infra`, `core-agent-service`, `v9-services`, and the umbrella chart.
- `v9-infra` renders all four middleware components.
- `core-agent-service` still renders `POSTGRES_URL`, `RABBITMQ_URL`, `REDIS_URL`, and `OPENSEARCH_URL`.
- Ansible syntax is checked when `ansible-playbook` is installed.
- Missing local tools such as `ansible-playbook`, `kubectl`, `docker`, or `bash` are reported explicitly as skipped runtime checks.

## How to interpret results

- Treat any `FAIL` as a real regression or delivery blocker.
- Treat `WARN` in regression mode as a follow-up item, not as proof of deployability.
- Treat `SKIP` as "not verified on this machine"; mention the missing tool in the final response.
- Do not say the deployment is fully verified if release mode fails or if runtime checks are skipped.

## Guardrails

- Prefer the validator script over ad hoc one-off checks.
- If you change deployment logic, run the validator before replying.
- If the validator cannot cover a runtime path because the machine lacks tooling or cluster access, say exactly which checks were skipped.
- Keep the validator script updated when deployment entrypoints or expected middleware components change.
