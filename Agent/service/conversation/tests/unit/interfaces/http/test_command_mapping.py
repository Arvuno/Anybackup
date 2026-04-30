from app.interfaces.http.v1.command_mapping import build_create_conversation_command
from app.interfaces.http.v1.schemas import CreateConversationRequest, UserMessageRequest


def test_create_conversation_request_maps_to_application_command() -> None:
    request = CreateConversationRequest(
        initial_message=UserMessageRequest(
            type="user_message",
            content="restore database",
            client_message_id="client-msg-001",
        ),
        title="DB restore",
        scenario_binding={"scenario_id": "scenario-001", "task_type": "restore_db"},
        tags=["database", "restore"],
        source="web",
        idempotency_key="body-key",
    )

    command = build_create_conversation_command(
        request=request,
        idempotency_key_header="header-key",
        request_id="req-001",
    )

    assert command.initial_message_content == "restore database"
    assert command.initial_message_client_id == "client-msg-001"
    assert command.title == "DB restore"
    assert command.scenario_binding == {
        "scenario_id": "scenario-001",
        "task_type": "restore_db",
        "asset_refs": [],
    }
    assert command.tags == ("database", "restore")
    assert command.source == "web"
    assert command.idempotency_key == "header-key"
    assert command.request_id == "req-001"
