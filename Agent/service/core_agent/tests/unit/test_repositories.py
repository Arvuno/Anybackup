from datetime import datetime, timezone

from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from core_agent_service.infrastructure.db.base import Base
from core_agent_service.infrastructure.db.models import ConversationMapping
from core_agent_service.infrastructure.db.repositories import MappingRepository, InboundEventRepository, OutboundStatusEventRepository


def _make_session() -> Session:
    engine = create_engine("sqlite+pysqlite:///:memory:", future=True)
    Base.metadata.create_all(engine)
    return Session(engine)


def test_conversation_mapping_model_has_unique_conversation_id():
    assert ConversationMapping.__tablename__ == "conversation_mappings"
    assert ConversationMapping.__table__.c.conversation_id.unique is True


def test_mapping_repository_upserts_by_conversation_id():
    session = _make_session()
    repository = MappingRepository(session)

    repository.upsert_mapping("conv_001", "decision_conv_001", decision_agent_id="agent_001")
    repository.upsert_mapping("conv_001", "decision_conv_002", decision_agent_id="agent_002")

    mapping = repository.get_by_conversation_id("conv_001")
    assert mapping is not None
    assert mapping.decision_agent_id == "agent_002"
    assert mapping.decision_conversation_id == "decision_conv_002"


def test_event_repositories_persist_minimal_records():
    session = _make_session()
    inbound_repository = InboundEventRepository(session)
    outbound_repository = OutboundStatusEventRepository(session)
    occurred_at = datetime(2026, 4, 22, 10, 0, tzinfo=timezone.utc)

    inbound = inbound_repository.save(
        event_id="evt_001",
        event_type="conversation.message.sent",
        occurred_at=occurred_at,
        conversation_id="conv_001",
        message_id="msg_001",
        content="hello",
        processing_status="processing",
    )
    outbound = outbound_repository.save(
        event_id="evt_status_001",
        event_type="core_agent.run.completed",
        occurred_at=occurred_at,
        conversation_id="conv_001",
        content="completed",
    )

    assert inbound.event_id == "evt_001"
    assert inbound.message_id == "msg_001"
    assert outbound.event_id == "evt_status_001"
