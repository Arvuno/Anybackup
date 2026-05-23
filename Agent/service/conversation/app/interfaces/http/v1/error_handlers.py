from uuid import uuid4

from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from starlette.exceptions import HTTPException as StarletteHTTPException

from app.domain.shared.errors import DomainError
from app.interfaces.http.v1.i18n import negotiate_locale
from app.interfaces.http.v1.schemas import ErrorCode, ErrorResponse, VisibleError


def register_error_handlers(app: FastAPI) -> None:
    @app.exception_handler(DomainError)
    async def domain_error_handler(
        request: Request,
        exc: DomainError,
    ) -> JSONResponse:
        descriptor = request.app.state.container.error_catalog().resolve_reason(exc.reason)
        return _error_response(
            request=request,
            status_code=descriptor.http_status,
            code=_error_code(descriptor.code),
            message_key=descriptor.message_key,
            retryable=descriptor.retryable,
            scope=descriptor.scope,
        )

    @app.exception_handler(RequestValidationError)
    async def validation_error_handler(
        request: Request,
        exc: RequestValidationError,
    ) -> JSONResponse:
        del exc
        return _error_response_for_code(request=request, code="VALIDATION_FAILED")

    @app.exception_handler(StarletteHTTPException)
    async def http_error_handler(
        request: Request,
        exc: StarletteHTTPException,
    ) -> JSONResponse:
        code = _http_status_to_error_code(exc.status_code)
        return _error_response_for_code(request=request, code=code.value)


def _error_response(
    *,
    request: Request,
    status_code: int,
    code: ErrorCode,
    message_key: str,
    retryable: bool,
    scope: str,
) -> JSONResponse:
    locale = negotiate_locale(request.headers.get("Accept-Language"))
    message = request.app.state.container.message_catalog().render(message_key, locale)
    request_id = getattr(request.state, "request_id", None) or request.headers.get(
        "X-Request-Id"
    ) or str(uuid4())
    trace_id = getattr(request.state, "trace_id", None)
    body = ErrorResponse(
        error=VisibleError(
            code=code,
            message=message,
            retryable=retryable,
            details={"scope": scope},
        ),
        request_id=request_id,
        trace_id=trace_id,
    )
    return JSONResponse(
        status_code=status_code,
        content=body.model_dump(mode="json"),
        headers={"Content-Language": locale, "X-Request-Id": request_id},
    )


def _error_response_for_code(*, request: Request, code: str) -> JSONResponse:
    descriptor = request.app.state.container.error_catalog().resolve_code(code)
    return _error_response(
        request=request,
        status_code=descriptor.http_status,
        code=_error_code(descriptor.code),
        message_key=descriptor.message_key,
        retryable=descriptor.retryable,
        scope=descriptor.scope,
    )


def _error_code(value: str) -> ErrorCode:
    try:
        return ErrorCode(value)
    except ValueError:
        return ErrorCode.INTERNAL_ERROR


def _http_status_to_error_code(status_code: int) -> ErrorCode:
    if status_code == 400:
        return ErrorCode.BAD_REQUEST
    if status_code == 401:
        return ErrorCode.UNAUTHORIZED
    if status_code == 403:
        return ErrorCode.FORBIDDEN
    if status_code == 404:
        return ErrorCode.NOT_FOUND
    if status_code == 409:
        return ErrorCode.IDEMPOTENCY_CONFLICT
    if status_code == 422:
        return ErrorCode.VALIDATION_FAILED
    return ErrorCode.INTERNAL_ERROR
