from __future__ import annotations

from datetime import datetime, timezone
import logging

from sqlalchemy import select
from sqlalchemy.orm import Session

from core_agent_service.infrastructure.db.models import ConversationMapping, InboundEventRecord, OutboundStatusEventRecord


logger = logging.getLogger(__name__)


class MappingRepository:
    def __init__(self, session: Session):
        self._session = session

    def get_by_conversation_id(self, conversation_id: str) -> ConversationMapping | None:
        statement = select(ConversationMapping).where(ConversationMapping.conversation_id == conversation_id)
        mapping = self._session.scalars(statement).one_or_none()
        logger.info(
            "queried conversation mapping",
            extra={"conversation_id": conversation_id, "mapping_found": mapping is not None},
        )
        return mapping

    def upsert_mapping(self, conversation_id: str, decision_conversation_id: str, *, decision_agent_id: str | None = None) -> ConversationMapping:
        # conversation_id 是本地主键，decision_agent_id 与 decision_conversation_id 共同定位平台侧目标会话。
        mapping = self.get_by_conversation_id(conversation_id)
        if mapping is None:
            mapping = ConversationMapping(
                conversation_id=conversation_id,
                decision_agent_id=decision_agent_id,
                decision_conversation_id=decision_conversation_id,
                mapping_status="active",
            )
            self._session.add(mapping)
        else:
            if decision_agent_id is not None:
                mapping.decision_agent_id = decision_agent_id
            mapping.decision_conversation_id = decision_conversation_id
            mapping.mapping_status = "active"
            mapping.updated_at = datetime.now(timezone.utc)
        self._session.flush()
        logger.info(
            "mapping upserted",
            extra={
                "conversation_id": conversation_id,
                "decision_agent_id": decision_agent_id or getattr(mapping, "decision_agent_id", None),
                "decision_conversation_id": decision_conversation_id,
            },
        )
        return mapping


class InboundEventRepository:
    def __init__(self, session: Session):
        self._session = session

    def get_by_event_id(self, event_id: str) -> InboundEventRecord | None:
        # event_id 是入站消费的幂等键，业务处理前必须先查已有记录。
        record = self._session.get(InboundEventRecord, event_id)
        logger.info("queried inbound event", extra={"event_id": event_id, "event_found": record is not None})
        return record

    def save(self, *, event_id: str, event_type: str, occurred_at: datetime, conversation_id: str, content: str, processing_status: str, message_id: str | None = None, error_message: str | None = None) -> InboundEventRecord:
        # 入站事件使用 merge，保证同一 event_id 重放时不会生成重复主键记录。
        record = InboundEventRecord(
            event_id=event_id,
            event_type=event_type,
            occurred_at=occurred_at,
            conversation_id=conversation_id,
            message_id=message_id,
            content=content,
            processing_status=processing_status,
            error_message=error_message,
        )
        self._session.merge(record)
        self._session.flush()
        logger.info(
            "inbound event saved",
            extra={"event_id": event_id, "event_type": event_type, "conversation_id": conversation_id, "processing_status": processing_status},
        )
        return self._session.get(InboundEventRecord, event_id)


class OutboundStatusEventRepository:
    def __init__(self, session: Session):
        self._session = session

    def save(self, *, event_id: str, event_type: str, occurred_at: datetime, conversation_id: str, content: str | None = None) -> OutboundStatusEventRecord:
        # 出站状态事件单独落库，便于后续核对“服务是否已经向上游发过 accepted/completed/failed”。
        record = OutboundStatusEventRecord(
            event_id=event_id,
            event_type=event_type,
            occurred_at=occurred_at,
            conversation_id=conversation_id,
            content=content,
        )
        self._session.merge(record)
        self._session.flush()
        logger.info(
            "outbound status event saved",
            extra={"event_id": event_id, "event_type": event_type, "conversation_id": conversation_id},
        )
        return self._session.get(OutboundStatusEventRecord, event_id)
