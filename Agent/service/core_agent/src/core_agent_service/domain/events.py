from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, ConfigDict, model_validator


class InboundPayload(BaseModel):
    conversation_id: str
    message_id: str | None = None
    turn_id: str | None = None
    content: str

    model_config = ConfigDict(extra="forbid")

    @model_validator(mode="after")
    def _validate_content(self) -> "InboundPayload":
        if not self.content:
            raise ValueError("payload.content is required")
        return self


class InboundEvent(BaseModel):
    event_id: str
    event_type: str
    occurred_at: datetime
    source_service: str
    payload: InboundPayload

    model_config = ConfigDict(extra="forbid")

    @model_validator(mode="after")
    def _validate_message_id(self) -> "InboundEvent":
        # 发送消息必须带上游用户消息 ID；取消事件仍保持最小取消指令结构。
        if self.event_type == "conversation.message.sent":
            if not self.payload.message_id:
                raise ValueError("payload.message_id is required for conversation.message.sent")
            if not self.payload.turn_id:
                raise ValueError("payload.turn_id is required for conversation.message.sent")
        return self
