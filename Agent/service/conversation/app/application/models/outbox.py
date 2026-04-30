from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True, slots=True)
class OutboxPublishMessage:
    exchange: str
    routing_key: str
    body: dict[str, Any]
    headers: dict[str, str]
    message_id: str
    correlation_id: str


@dataclass(frozen=True, slots=True)
class OutboxPublishRunResult:
    lock_acquired: bool
    scanned: int = 0
    published: int = 0
    failed: int = 0
    dlq: int = 0
