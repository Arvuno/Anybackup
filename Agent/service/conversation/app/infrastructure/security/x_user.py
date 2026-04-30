import json
from typing import Any

from app.application.models.conversation import AuthenticatedUser


class XUserParseError(ValueError):
    pass


def parse_x_user_header(header_value: str | None) -> AuthenticatedUser:
    if header_value is None or not header_value.strip():
        raise XUserParseError("missing X-User header")
    try:
        raw = json.loads(header_value)
    except json.JSONDecodeError as exc:
        raise XUserParseError("invalid X-User JSON") from exc
    if not isinstance(raw, dict):
        raise XUserParseError("X-User must be a JSON object")

    subject = raw.get("sub")
    if not isinstance(subject, str) or not subject.strip():
        raise XUserParseError("X-User sub is required")

    return AuthenticatedUser(
        user_id=subject,
        preferred_username=_optional_str(raw.get("preferred_username")),
        name=_optional_str(raw.get("name")),
        email=_optional_str(raw.get("email")),
        email_verified=_optional_bool(raw.get("email_verified")),
        roles=_roles(raw.get("roles")),
    )


def _optional_str(value: Any) -> str | None:
    if isinstance(value, str):
        return value
    return None


def _optional_bool(value: Any) -> bool | None:
    if isinstance(value, bool):
        return value
    return None


def _roles(value: Any) -> tuple[str, ...]:
    if not isinstance(value, list):
        return ()
    return tuple(role for role in value if isinstance(role, str))
