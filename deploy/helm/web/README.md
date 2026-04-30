# web-service Helm chart

This repository contains the standalone Helm chart for deploying the frontend web application built from the `web_app_agent` repository.

## What this chart deploys

- One frontend `Deployment`
- One Kubernetes `Service`
- One Traefik `IngressRoute`

The container image is expected to be built from the frontend runtime Dockerfile in `web_app_agent`. This chart does not build the image. It only deploys it.

## Runtime behavior

The frontend container serves a static SPA with `nginx` on port `80`.
This chart does not configure backend API proxy routes. Business APIs are exposed by their own service charts through the shared API Gateway Traefik entry, and API route priority must keep those routes ahead of the frontend catch-all `/` route.

The frontend route is public. It does not attach authentication middleware.

## Default values to review before install

Update these values for your environment before installing:

- `image.repository`
- `image.tag`
- `service.type`
- `service.nodePort`
- `ingressRoute.enabled`
- `ingressRoute.pathPrefix`
- `ingressRoute.priority`

## Install with API Gateway route

```powershell
helm upgrade --install web-service . `
  --namespace anybackup-ai `
  --create-namespace `
  --set image.repository=registry.example.com/anybackup/web-service `
  --set image.tag=0.1.0
```

## Install with NodePort for quick testing

```powershell
helm upgrade --install web-service . `
  --namespace anybackup-ai `
  --create-namespace `
  --set image.repository=registry.example.com/anybackup/web-service `
  --set image.tag=0.1.0 `
  --set service.type=NodePort `
  --set service.nodePort=30080
```

Then access the app with:

```text
http://<node-ip>:30080
```

NodePort is intended for static UI smoke testing only. Platform API access should use the API Gateway route, where each business service owns its own routing and middleware policy.

## Notes

- Keep this chart aligned with the runtime contract defined in `web_app_agent`.
- This chart is intended to be linked into `disaster_recovery_agent` as a dedicated chart repository.

## License and Third-Party Notices

- This chart's source is distributed under the repository root [LICENSE](../../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../../NOTICE).
- The base images (`node:22-alpine`, `nginx:1.30-alpine`, Alpine Linux) and frontend npm runtime dependencies (React, framer-motion, zustand, lucide-react, etc.) bundled into the frontend image shipped by this chart are declared uniformly in [`THIRD_PARTY_NOTICES.md`](./THIRD_PARTY_NOTICES.md), which also serves as the authoritative notice for the same image and is cross-referenced by `Agent/portal/deploy/helm/agent-web`.
- When upgrading the base image or any frontend runtime dependency, the Inventory table of the above notice file MUST be updated in the same commit.
- `devDependencies` used inside the frontend build container (build-time only) are not redistributed with the runtime image. They are described only in §1 of the notice file as scope context, and are not listed in the Inventory table.
