from __future__ import annotations

from datetime import datetime, timezone

from pydantic import BaseModel, ConfigDict, model_validator


class OutboundStatusPayload(BaseModel):
    conversation_id: str
    message_id: str | None = None
    content: str

    model_config = ConfigDict(extra="forbid")

    @model_validator(mode="after")
    def _validate_content(self) -> "OutboundStatusPayload":
        if not self.content:
            raise ValueError("payload.content is required")
        return self


class OutboundStatusEvent(BaseModel):
    event_id: str
    event_type: str
    occurred_at: datetime
    source_service: str = "core_agent_service"
    payload: OutboundStatusPayload

    model_config = ConfigDict(extra="forbid")


def build_status_event(
    *,
    event_id: str,
    event_type: str,
    conversation_id: str,
    content: str,
    message_id: str | None = None,
) -> OutboundStatusEvent:
    # 下行 MQ 契约只允许状态链路字段，message_id 用于会话服务精确定位用户消息。
    return OutboundStatusEvent(
        event_id=event_id,
        event_type=event_type,
        occurred_at=datetime.now(timezone.utc),
        payload=OutboundStatusPayload(
            conversation_id=conversation_id,
            message_id=message_id,
            content=content,
        ),
    )
