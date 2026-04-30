from typing import Annotated

from fastapi import Header, HTTPException, Request

from app.application.models.conversation import AuthenticatedUser
from app.infrastructure.security.x_user import XUserParseError, parse_x_user_header


async def require_user_context(
    x_user: Annotated[str | None, Header(alias="X-User")] = None,
) -> AuthenticatedUser:
    try:
        return parse_x_user_header(x_user)
    except XUserParseError as exc:
        raise HTTPException(status_code=401, detail="Unauthorized") from exc


async def require_core_agent_service_identity(
    request: Request,
    service_name: Annotated[str | None, Header(alias="X-Service-Name")] = None,
    service_token: Annotated[str | None, Header(alias="X-Service-Token")] = None,
) -> str:
    settings = request.app.state.container.settings()
    if service_name is None or service_token is None:
        raise HTTPException(status_code=401, detail="Service identity required")
    if service_name != settings.core_agent_service_name:
        raise HTTPException(status_code=403, detail="Service identity is not allowed")
    if service_token != settings.core_agent_service_token:
        raise HTTPException(status_code=403, detail="Service token is invalid")
    return service_name
