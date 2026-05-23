from typing import Protocol

from app.application.models.outbox import OutboxPublishMessage


class EventPublisher(Protocol):
    async def publish(self, message: OutboxPublishMessage) -> None:
        raise NotImplementedError
