from fastapi.testclient import TestClient

from app.bootstrap.app_factory import create_app


def test_healthz_returns_process_liveness() -> None:
    app = create_app()
    client = TestClient(app)

    response = client.get("/healthz")

    assert response.status_code == 200
    assert response.json() == {
        "status": "ok",
        "service": "conversation_service",
    }
