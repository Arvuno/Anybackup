import asyncio
import json
from collections.abc import Iterator
from pathlib import Path

import pytest
from dependency_injector import providers
from fastapi.testclient import TestClient
from sqlalchemy.ext.asyncio import create_async_engine

from app.application.ports.error_catalog import ErrorDescriptor
from app.bootstrap.app_factory import create_app
from app.bootstrap.settings import Settings
from app.infrastructure.persistence.sqlalchemy.models import Base

API_PREFIX = "/api/conversation_service/v1"
TRACEPARENT = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"


@pytest.fixture()
def client(tmp_path: Path) -> Iterator[TestClient]:
    database_path = tmp_path / "conversation.db"
    database_url = f"sqlite+aiosqlite:///{database_path.as_posix()}"

    async def create_schema() -> None:
        engine = create_async_engine(database_url)
        async with engine.begin() as connection:
            await connection.run_sync(Base.metadata.create_all)
        await engine.dispose()

    asyncio.run(create_schema())
    app = create_app(Settings(database_url=database_url, snowflake_node_id=15))

    with TestClient(app) as test_client:
        yield test_client

    asyncio.run(app.state.container.database_engine().dispose())


def test_missing_x_user_uses_localized_error_and_trace_headers(client: TestClient) -> None:
    response = client.get(
        f"{API_PREFIX}/conversations",
        headers={
            "Accept-Language": "zh-CN",
            "X-Request-Id": "req-sec-001",
            "traceparent": TRACEPARENT,
        },
    )

    assert response.status_code == 401
    assert response.headers["Content-Language"] == "zh-CN"
    assert response.headers["X-Request-Id"] == "req-sec-001"
    assert response.json() == {
        "error": {
            "code": "UNAUTHORIZED",
            "message": "身份信息缺失或无效。",
            "retryable": True,
            "details": {"scope": "page"},
        },
        "request_id": "req-sec-001",
        "trace_id": "0af7651916cd43dd8448eb211c80319c",
    }


def test_error_message_can_use_english_locale(client: TestClient) -> None:
    response = client.get(
        f"{API_PREFIX}/conversations",
        headers={"Accept-Language": "en-US", "X-Request-Id": "req-sec-002"},
    )

    assert response.status_code == 401
    assert response.headers["Content-Language"] == "en-US"
    assert (
        response.json()["error"]["message"]
        == "Authentication information is missing or invalid."
    )
    assert response.json()["error"]["code"] == "UNAUTHORIZED"


def test_error_response_does_not_echo_sensitive_header_value(client: TestClient) -> None:
    response = client.get(
        f"{API_PREFIX}/conversations",
        headers={
            "X-User": '{"sub":',
            "Authorization": "Bearer access_token=secret-value",
            "Accept-Language": "en-US",
        },
    )

    body = json.dumps(response.json())
    assert response.status_code == 401
    assert "secret-value" not in body
    assert "access_token" not in body


def test_platform_error_catalog_override_is_used() -> None:
    app = create_app()
    app.state.container.error_catalog.override(
        providers.Object(
            FakeErrorCatalog(
                ErrorDescriptor(
                    code="UNAUTHORIZED",
                    http_status=401,
                    retryable=False,
                    message_key="fake.unauthorized",
                    scope="page",
                    source="platform",
                )
            )
        )
    )
    app.state.container.message_catalog.override(providers.Object(FakeMessageCatalog()))
    client = TestClient(app)

    response = client.get(
        f"{API_PREFIX}/conversations",
        headers={"Accept-Language": "en-US"},
    )

    assert response.status_code == 401
    assert response.json()["error"]["message"] == "platform message"
    assert response.json()["error"]["retryable"] is False
    assert response.json()["error"]["details"] == {"scope": "page"}


class FakeErrorCatalog:
    def __init__(self, descriptor: ErrorDescriptor) -> None:
        self._descriptor = descriptor

    def resolve_code(self, code: str) -> ErrorDescriptor:
        return self._descriptor

    def resolve_reason(self, reason: object) -> ErrorDescriptor:
        return self._descriptor


class FakeMessageCatalog:
    def render(self, message_key: str, locale: str) -> str:
        return "platform message"
