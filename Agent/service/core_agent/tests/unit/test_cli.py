from core_agent_service.cli import _build_sample_event, build_parser, publish_event_to_rabbitmq


def test_cli_has_publish_sample_and_show_mapping_commands():
    parser = build_parser()
    subcommands = parser._subparsers._group_actions[0].choices

    assert "publish-sample" in subcommands
    assert "publish-rabbitmq" in subcommands
    assert "show-mapping" in subcommands


class FakeChannel:
    def __init__(self):
        self.calls: list[tuple[str, tuple, dict]] = []

    def exchange_declare(self, *args, **kwargs):
        self.calls.append(("exchange_declare", args, kwargs))

    def basic_publish(self, *args, **kwargs):
        self.calls.append(("basic_publish", args, kwargs))


class FakeConnection:
    def __init__(self, channel: FakeChannel):
        self._channel = channel

    def channel(self):
        return self._channel


class FakeSettings:
    rabbitmq_url = "amqp://guest:guest@localhost:5672/"
    rabbitmq_exchange = "conversation.message.events"
    rabbitmq_exchange_type = "topic"


def test_publish_event_to_rabbitmq_uses_real_exchange_and_routing_key():
    fake_channel = FakeChannel()
    event = {
        "event_id": "evt_001",
        "event_type": "conversation.message.sent",
        "occurred_at": "2026-04-22T10:00:00Z",
        "source_service": "local_cli",
        "payload": {
            "conversation_id": "conv_001",
            "content": "hello",
        },
    }

    publish_event_to_rabbitmq(
        FakeSettings(),
        event,
        routing_key="conversation.message.sent.v1",
        connection_factory=lambda: FakeConnection(fake_channel),
    )

    call_names = [name for name, _, _ in fake_channel.calls]
    assert "exchange_declare" in call_names
    assert "basic_publish" in call_names


def test_build_sample_sent_event_includes_message_id():
    event = _build_sample_event("conv_001", "hello", "conversation.message.sent")

    assert event["payload"]["message_id"].startswith("msg_")
