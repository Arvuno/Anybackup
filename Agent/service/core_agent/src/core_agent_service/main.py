from __future__ import annotations

from core_agent_service.bootstrap import build_app


def main() -> int:
    app = build_app()
    app.get("consumer_runner", app["consumer"]).run()
    return 0
