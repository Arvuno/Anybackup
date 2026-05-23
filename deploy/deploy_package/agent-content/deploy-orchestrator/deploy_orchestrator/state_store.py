from __future__ import annotations

import json
from pathlib import Path

from deploy_orchestrator.logger import get_logger


class StateStore:
    def __init__(self, path: str | Path):
        self.path = Path(path)
        self.logger = get_logger(__name__)

    def load(self) -> dict[str, object]:
        if not self.path.exists():
            self.logger.info("State file does not exist at %s", self.path)
            return {}
        self.logger.info("Loading state file from %s", self.path)
        return json.loads(self.path.read_text(encoding="utf-8"))

    def save(self, data: dict[str, object]) -> None:
        self.path.parent.mkdir(parents=True, exist_ok=True)
        self.logger.info("Saving state file to %s", self.path)
        self.path.write_text(
            json.dumps(data, ensure_ascii=False, indent=2, sort_keys=True),
            encoding="utf-8",
        )
