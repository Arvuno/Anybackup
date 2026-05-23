from dataclasses import dataclass
from typing import Protocol

from app.domain.shared.errors import ErrorReason


@dataclass(frozen=True, slots=True)
class ErrorDescriptor:
    code: str
    http_status: int
    retryable: bool
    message_key: str
    scope: str
    source: str = "local_fallback"


class ErrorCatalog(Protocol):
    def resolve_code(self, code: str) -> ErrorDescriptor:
        raise NotImplementedError

    def resolve_reason(self, reason: ErrorReason) -> ErrorDescriptor:
        raise NotImplementedError
