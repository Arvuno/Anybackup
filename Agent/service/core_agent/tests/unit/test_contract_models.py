from datetime import datetime, timezone

import pytest

from core_agent_service.domain.events import InboundEvent
from core_agent_service.domain.statuses import OutboundStatusPayload, build_status_event


def test_inbound_event_accepts_minimal_message():
    event = InboundEvent.model_validate(
        {
            "event_id": "evt_001",
            "event_type": "conversation.message.sent",
            "occurred_at": "2026-04-22T10:00:00Z",
            "source_service": "conversation_service",
            "payload": {
                "conversation_id": "conv_001",
                "message_id": "msg_001",
                "turn_id": "turn_001",
                "content": "hello",
            },
        }
    )

    assert event.payload.conversation_id == "conv_001"
    assert event.payload.message_id == "msg_001"
    assert event.payload.turn_id == "turn_001"


def test_sent_event_rejects_missing_message_id():
    with pytest.raises(ValueError):
        InboundEvent.model_validate(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": datetime.now(timezone.utc).isoformat(),
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_001", "content": "hello"},
            }
        )


def test_sent_event_rejects_missing_turn_id():
    with pytest.raises(ValueError):
        InboundEvent.model_validate(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": datetime.now(timezone.utc).isoformat(),
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_001", "message_id": "msg_001", "content": "hello"},
            }
        )


def test_cancel_event_allows_missing_message_id():
    event = InboundEvent.model_validate(
        {
            "event_id": "evt_002",
            "event_type": "conversation.message.cancel_requested",
            "occurred_at": datetime.now(timezone.utc).isoformat(),
            "source_service": "conversation_service",
            "payload": {"conversation_id": "conv_001", "content": "cancel"},
        }
    )

    assert event.payload.message_id is None


def test_inbound_event_rejects_missing_content():
    with pytest.raises(ValueError):
        InboundEvent.model_validate(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": datetime.now(timezone.utc).isoformat(),
                "source_service": "conversation_service",
                "payload": {
                    "conversation_id": "conv_001",
                    "message_id": "msg_001",
                    "turn_id": "turn_001",
                    "content": "",
                },
            }
        )


def test_outbound_status_payload_accepts_minimal_completed_state():
    payload = OutboundStatusPayload(
        conversation_id="conv_001",
        content="success",
    )

    assert payload.content == "success"


def test_outbound_status_event_uses_minimal_payload_without_status_or_error_fields():
    event = build_status_event(
        event_id="evt_101",
        event_type="core_agent.run.failed",
        conversation_id="conv_001",
        message_id="msg_001",
        content="redacted timeout",
    )

    assert event.model_dump(mode="json")["payload"] == {
        "conversation_id": "conv_001",
        "message_id": "msg_001",
        "content": "redacted timeout",
    }


def test_outbound_status_payload_rejects_missing_content():
    with pytest.raises(ValueError):
        OutboundStatusPayload(conversation_id="conv_001")
