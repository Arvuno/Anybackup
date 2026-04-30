import json

from fastapi.testclient import TestClient

from app.bootstrap.app_factory import create_app


def test_missing_user_context_uses_contract_error_response() -> None:
    client = TestClient(create_app())

    response = client.post(
        "/api/conversation_service/v1/conversations",
        headers={
            "Idempotency-Key": "header-key",
            "X-Request-Id": "req-001",
        },
        json={
            "initial_message": {
                "type": "user_message",
                "content": "restore database",
            }
        },
    )

    assert response.status_code == 401
    assert response.json() == {
        "error": {
            "code": "UNAUTHORIZED",
            "message": "身份信息缺失或无效。",
            "retryable": True,
            "details": {"scope": "page"},
        },
        "request_id": "req-001",
        "trace_id": response.json()["trace_id"],
    }
    assert response.json()["trace_id"] is not None
    assert response.headers["Content-Language"] == "zh-CN"
    assert response.headers["X-Request-Id"] == "req-001"


def test_validation_error_uses_contract_error_response() -> None:
    client = TestClient(create_app())

    response = client.post(
        "/api/conversation_service/v1/conversations",
        headers={
            "X-Request-Id": "req-test-001",
            "X-User": json.dumps({"sub": "user-001"}),
        },
        json={},
    )

    assert response.status_code == 422
    assert response.json() == {
        "error": {
            "code": "VALIDATION_FAILED",
            "message": "输入内容不符合要求，请检查后重试。",
            "retryable": False,
            "details": {"scope": "field"},
        },
        "request_id": "req-test-001",
        "trace_id": response.json()["trace_id"],
    }
    assert response.json()["trace_id"] is not None
    assert response.headers["Content-Language"] == "zh-CN"
    assert response.headers["X-Request-Id"] == "req-test-001"
