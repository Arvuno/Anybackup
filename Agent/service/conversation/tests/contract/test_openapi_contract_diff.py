from pathlib import Path

from app.bootstrap.app_factory import create_app
from app.interfaces.http.v1.contract_diff import find_missing_contract_paths

REPO_ROOT = Path(__file__).resolve().parents[5]
UPSTREAM_CONTRACT = REPO_ROOT / "docs/需求上下文/2026-04-会话管理/契约/接口契约.yaml"


def test_no_upstream_contract_paths_are_missing_from_app_openapi() -> None:
    schema = create_app().openapi()

    missing_paths = find_missing_contract_paths(
        upstream_contract_path=UPSTREAM_CONTRACT,
        generated_openapi=schema,
        service_prefix="/api/conversation_service/v1",
    )

    assert missing_paths == set()
