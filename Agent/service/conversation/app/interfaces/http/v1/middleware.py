from collections.abc import Awaitable, Callable
from uuid import uuid4

from fastapi import FastAPI, Request, Response

from app.infrastructure.observability.tracing import request_trace_id


def register_request_context_middleware(app: FastAPI) -> None:
    @app.middleware("http")
    async def request_context_middleware(
        request: Request,
        call_next: Callable[[Request], Awaitable[Response]],
    ) -> Response:
        request_id = request.headers.get("X-Request-Id") or str(uuid4())
        trace_id = request_trace_id(request.headers.get("traceparent"))
        request.state.request_id = request_id
        request.state.trace_id = trace_id

        response = await call_next(request)
        response.headers["X-Request-Id"] = request_id
        return response
