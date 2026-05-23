from fastapi.testclient import TestClient

from app.bootstrap.app_factory import create_app


def test_openapi_document_is_served_under_service_prefix() -> None:
    app = create_app()
    client = TestClient(app)

    response = client.get("/api/conversation_service/v1/openapi.json")

    assert response.status_code == 200
    body = response.json()
    assert body["info"]["title"] == "Conversation Service"
    assert body["info"]["version"] == "0.1.0"


def test_openapi_schema_contains_health_endpoints() -> None:
    app = create_app()

    schema = app.openapi()

    assert "/healthz" in schema["paths"]
    assert "/readyz" in schema["paths"]
