import logging
import re
from collections.abc import Mapping, MutableMapping
from typing import Any

import structlog

_CREDENTIAL_KEYS = {
    "authorization",
    "cookie",
    "set_cookie",
    "password",
    "token",
    "access_token",
    "refresh_token",
    "api_key",
    "secret",
    "private_key",
}
_EMAIL_PATTERN = re.compile(r"^([^@\s])[^@\s]*(@[^@\s]+)$")


def configure_structured_logging(level: str = "INFO") -> None:
    logging.basicConfig(level=getattr(logging, level.upper(), logging.INFO), format="%(message)s")
    structlog.configure(
        processors=[
            structlog.contextvars.merge_contextvars,
            structlog.stdlib.add_log_level,
            structlog.stdlib.add_logger_name,
            structlog.processors.TimeStamper(fmt="iso", utc=True),
            _structlog_redaction_processor,
            structlog.processors.JSONRenderer(),
        ],
        wrapper_class=structlog.stdlib.BoundLogger,
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )


def _structlog_redaction_processor(
    logger: logging.Logger,
    method_name: str,
    event_dict: MutableMapping[str, Any],
) -> Mapping[str, Any]:
    del logger, method_name
    return redact_event_dict(event_dict)


def redact_event_dict(event_dict: Mapping[str, Any]) -> dict[str, Any]:
    return {
        str(key): _redact_value(key=str(key), value=value)
        for key, value in event_dict.items()
    }


def _redact_value(*, key: str, value: Any) -> Any:
    normalized_key = key.lower()
    if any(token in normalized_key for token in _CREDENTIAL_KEYS):
        return "[REDACTED:credential]"
    if isinstance(value, dict):
        return {
            str(item_key): _redact_value(key=str(item_key), value=item_value)
            for item_key, item_value in value.items()
        }
    if isinstance(value, list):
        return [_redact_value(key=key, value=item) for item in value]
    if isinstance(value, str):
        cleaned = " ".join(value.replace("\r", " ").replace("\n", " ").split())
        email_match = _EMAIL_PATTERN.match(cleaned)
        if email_match is not None:
            return f"{email_match.group(1)}***{email_match.group(2)}"
        if _looks_like_sensitive_value(cleaned):
            return "[REDACTED:credential]"
        return cleaned
    return value


def _looks_like_sensitive_value(value: str) -> bool:
    normalized = value.lower()
    return any(
        token in normalized
        for token in (
            "access_token=",
            "refresh_token=",
            "password=",
            "api_key=",
            "private key",
            "system prompt",
            "internal reasoning",
        )
    )
