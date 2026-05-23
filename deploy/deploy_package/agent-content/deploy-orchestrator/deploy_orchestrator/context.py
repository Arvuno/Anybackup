from __future__ import annotations

from dataclasses import dataclass, field


@dataclass
class DeploymentContext:
    values: dict[str, dict[str, dict[str, object]]] = field(
        default_factory=lambda: {
            "skills": {},
            "bkns": {},
            "agents": {},
            "bindings": {},
        }
    )

    def set_resource(self, kind: str, name: str, payload: dict[str, object]) -> None:
        bucket_name = f"{kind}s" if not kind.endswith("s") else kind
        self.values.setdefault(bucket_name, {})
        self.values[bucket_name][name] = payload
