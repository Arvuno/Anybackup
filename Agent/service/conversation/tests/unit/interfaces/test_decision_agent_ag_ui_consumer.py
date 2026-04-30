import pytest

from app.application.commands.agent_events import DecisionAgentAgUiEventCommand
from app.interfaces.mq import decision_agent_ag_ui_consumer
from app.interfaces.mq.decision_agent_ag_ui_consumer import (
    RabbitMqDecisionAgentAgUiConsumer,
    _command_from_body,
)


@pytest.mark.asyncio
async def test_ag_ui_consumer_binds_decision_agent_exchange(monkeypatch) -> None:
    connection = FakeConnection()

    async def fake_connect_robust(url: str) -> FakeConnection:
        connection.url = url
        return connection

    monkeypatch.setattr(
        decision_agent_ag_ui_consumer.aio_pika,
        "connect_robust",
        fake_connect_robust,
    )
    consumer = RabbitMqDecisionAgentAgUiConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange_name="decision_agent.ag_ui.events",
        queue_name="conversation.decision_agent.ag_ui",
        prefetch_count=10,
        handler=NoopHandler(),
    )

    await consumer.start()

    assert connection.url == "amqp://guest:guest@localhost:5672/"
    assert connection.channel_instance.qos_calls == [10]
    assert connection.channel_instance.declared_exchanges == [
        {"name": "decision_agent.ag_ui.events", "durable": True}
    ]
    assert connection.channel_instance.declared_queues == [
        {"name": "conversation.decision_agent.ag_ui", "durable": True}
    ]
    assert connection.channel_instance.queue.bind_calls == [
        {
            "exchange": "decision_agent.ag_ui.events",
            "routing_key": "decision_agent.session.ag_ui_event.v1",
        }
    ]
    assert connection.channel_instance.queue.consume_calls == [False]


def test_command_from_body_accepts_markdown_payload_ag_ui() -> None:
    command = _command_from_body(
        {
            "event_id": "evt-markdown-001",
            "event_type": "decision_agent.session.ag_ui_event",
            "source_service": "decision_agent_session",
            "occurred_at": "2026-04-28T10:00:00Z",
            "payload": {
                "conversation_id": "100",
                "turn_id": "200",
                "message_id": "901",
                "content": "方案设计已生成。",
                "sequence": 1,
                "ag_ui": "# 方案设计\n\n这里是 Markdown 内容。",
            },
        }
    )

    assert command.ag_ui == "# 方案设计\n\n这里是 Markdown 内容。"


@pytest.mark.parametrize("ag_ui", [None, "", "  ", {"version": "1.x", "events": []}])
def test_command_from_body_rejects_missing_empty_or_object_ag_ui(ag_ui: object) -> None:
    with pytest.raises(ValueError, match="ag_ui"):
        _command_from_body(
            {
                "event_id": "evt-markdown-invalid",
                "event_type": "decision_agent.session.ag_ui_event",
                "source_service": "decision_agent_session",
                "payload": {
                    "conversation_id": "100",
                    "turn_id": "200",
                    "message_id": "901",
                    "content": "方案设计已生成。",
                    "sequence": 1,
                    "ag_ui": ag_ui,
                },
            }
        )


class NoopHandler:
    async def handle(self, command: DecisionAgentAgUiEventCommand) -> object:
        return command


class FakeConnection:
    def __init__(self) -> None:
        self.url = ""
        self.channel_instance = FakeChannel()

    async def channel(self) -> "FakeChannel":
        return self.channel_instance

    async def close(self) -> None:
        pass


class FakeChannel:
    def __init__(self) -> None:
        self.qos_calls: list[int] = []
        self.declared_exchanges: list[dict[str, object]] = []
        self.declared_queues: list[dict[str, object]] = []
        self.queue = FakeQueue()

    async def set_qos(self, *, prefetch_count: int) -> None:
        self.qos_calls.append(prefetch_count)

    async def declare_exchange(self, name: str, exchange_type, *, durable: bool):
        del exchange_type
        self.declared_exchanges.append({"name": name, "durable": durable})
        return FakeExchange(name=name)

    async def declare_queue(self, name: str, *, durable: bool) -> "FakeQueue":
        self.declared_queues.append({"name": name, "durable": durable})
        return self.queue


class FakeExchange:
    def __init__(self, *, name: str) -> None:
        self.name = name


class FakeQueue:
    def __init__(self) -> None:
        self.bind_calls: list[dict[str, str]] = []
        self.consume_calls: list[bool] = []

    async def bind(self, exchange: FakeExchange, *, routing_key: str) -> None:
        self.bind_calls.append({"exchange": exchange.name, "routing_key": routing_key})

    async def consume(self, callback, *, no_ack: bool) -> None:
        del callback
        self.consume_calls.append(no_ack)
