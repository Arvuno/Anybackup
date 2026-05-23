from core_agent_service.application.relay_service import RelayService


class FakeMapping:
    def __init__(
        self,
        decision_conversation_id: str,
        decision_agent_id: str | None = "target_agent_001",
    ):
        self.decision_conversation_id = decision_conversation_id
        self.decision_agent_id = decision_agent_id


class FakeMappingRepository:
    def __init__(self):
        self.mapping = None

    def get_by_conversation_id(self, conversation_id: str):
        return self.mapping

    def upsert_mapping(
        self,
        conversation_id: str,
        decision_conversation_id: str,
        *,
        decision_agent_id: str | None = None,
    ):
        self.mapping = FakeMapping(decision_conversation_id, decision_agent_id)
        return self.mapping


class FakeInboundEventRepository:
    def __init__(self):
        self.saved: list[dict] = []
        self.existing = None

    def get_by_event_id(self, event_id: str):
        return self.existing

    def save(self, **kwargs):
        self.saved.append(kwargs)
        return kwargs


class FakeAdapter:
    def __init__(self):
        self.sent: list[dict] = []
        self.intent_inputs: list[str] = []
        self.default_agent_id = "router_agent"

    def ensure_conversation(self, existing_conversation_id: str | None) -> str:
        return existing_conversation_id or ""

    def resolve_intent(
        self,
        content: str,
        *,
        conversation_id: str,
        turn_id: str,
        source_message_id: str,
        mq: dict,
        foundation: dict | None = None,
    ):
        self.intent_inputs.append(
            {
                "content": content,
                "conversation_id": conversation_id,
                "turn_id": turn_id,
                "source_message_id": source_message_id,
                "mq": mq,
                "foundation": foundation,
            }
        )
        return type(
            "IntentResolution",
            (),
            {
                "agent_id": "target_agent_001",
                "agent_name": "备份方案推荐Agent",
                "is_intent_success": True,
            },
        )()

    def relay_message(
        self,
        decision_conversation_id: str,
        *,
        conversation_id: str,
        turn_id: str,
        source_message_id: str,
        content: str,
        mq: dict,
        foundation: dict | None = None,
        agent_id: str | None = None,
    ) -> dict[str, object]:
        self.sent.append(
            {
                "decision_conversation_id": decision_conversation_id,
                "agent_id": agent_id,
                "conversation_id": conversation_id,
                "turn_id": turn_id,
                "source_message_id": source_message_id,
                "content": content,
                "mq": mq,
                "foundation": foundation,
            }
        )
        return {
            "decision_conversation_id": decision_conversation_id or "decision_conv_001",
            "raw": {"status": "completed"},
        }


def test_relay_service_creates_mapping_when_missing():
    mapping_repository = FakeMappingRepository()
    inbound_event_repository = FakeInboundEventRepository()
    adapter = FakeAdapter()
    service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
        mq_runtime_info={"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
    )

    result = service.handle_message(
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

    assert result.status == "completed"
    assert result.decision_conversation_id == "decision_conv_001"
    assert adapter.intent_inputs == [
        {
            "content": "hello",
            "conversation_id": "conv_001",
            "turn_id": "turn_001",
            "source_message_id": "msg_001",
            "mq": {"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
            "foundation": None,
        }
    ]
    assert adapter.sent == [
        {
            "decision_conversation_id": "",
            "agent_id": "target_agent_001",
            "conversation_id": "conv_001",
            "turn_id": "turn_001",
            "source_message_id": "msg_001",
            "content": "hello",
            "mq": {"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
            "foundation": None,
        }
    ]


def test_relay_service_reuses_existing_mapping():
    mapping_repository = FakeMappingRepository()
    mapping_repository.mapping = FakeMapping("decision_conv_existing")
    inbound_event_repository = FakeInboundEventRepository()
    adapter = FakeAdapter()
    service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
        mq_runtime_info={"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
    )

    result = service.handle_message(
        {
            "event_id": "evt_002",
            "event_type": "conversation.message.sent",
            "occurred_at": "2026-04-22T10:00:00Z",
            "source_service": "conversation_service",
            "payload": {
                "conversation_id": "conv_001",
                "message_id": "msg_002",
                "turn_id": "turn_002",
                "content": "world",
            },
        }
    )

    assert result.decision_conversation_id == "decision_conv_existing"
    assert adapter.intent_inputs == []
    assert adapter.sent == [
        {
            "decision_conversation_id": "decision_conv_existing",
            "agent_id": "target_agent_001",
            "conversation_id": "conv_001",
            "turn_id": "turn_002",
            "source_message_id": "msg_002",
            "content": "world",
            "mq": {"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
            "foundation": None,
        }
    ]


def test_relay_service_passes_foundation_runtime_info_to_adapter():
    mapping_repository = FakeMappingRepository()
    inbound_event_repository = FakeInboundEventRepository()
    adapter = FakeAdapter()
    service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
        mq_runtime_info={"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
        foundation_runtime_info={
            "endpoint": "https://115.190.150.254:9600",
            "ak": "foundation-ak",
            "sk": "foundation-sk",
        },
    )

    service.handle_message(
        {
            "event_id": "evt_foundation",
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

    assert adapter.sent[0]["foundation"] == {
        "endpoint": "https://115.190.150.254:9600",
        "ak": "foundation-ak",
        "sk": "foundation-sk",
    }


def test_relay_service_does_not_call_decision_agent_for_duplicate_event():
    mapping_repository = FakeMappingRepository()
    mapping_repository.mapping = FakeMapping("decision_conv_existing")
    inbound_event_repository = FakeInboundEventRepository()
    inbound_event_repository.existing = type(
        "InboundRecord",
        (),
        {
            "event_id": "evt_duplicate",
            "conversation_id": "conv_001",
            "processing_status": "completed",
        },
    )()
    adapter = FakeAdapter()
    service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
        mq_runtime_info={"rabbitmq_url": "amqp://guest:guest@localhost:5672/"},
    )

    result = service.handle_message(
        {
            "event_id": "evt_duplicate",
            "event_type": "conversation.message.sent",
            "occurred_at": "2026-04-23T10:00:00Z",
            "source_service": "conversation_service",
            "payload": {
                "conversation_id": "conv_001",
                "message_id": "msg_duplicate",
                "turn_id": "turn_duplicate",
                "content": "hello again",
            },
        }
    )

    assert result.status == "completed"
    assert result.decision_conversation_id == "decision_conv_existing"
    assert adapter.sent == []


def test_relay_service_fails_when_intent_resolution_misses_target_agent():
    mapping_repository = FakeMappingRepository()
    inbound_event_repository = FakeInboundEventRepository()
    adapter = FakeAdapter()

    def resolve_failed(
        content: str,
        *,
        conversation_id: str,
        turn_id: str,
        source_message_id: str,
        mq: dict,
        foundation: dict | None = None,
    ):
        return type(
            "IntentResolution",
            (),
            {"agent_id": None, "agent_name": None, "is_intent_success": False},
        )()

    adapter.resolve_intent = resolve_failed
    service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
    )

    try:
        service.handle_message(
            {
                "event_id": "evt_003",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {
                    "conversation_id": "conv_003",
                    "message_id": "msg_003",
                    "turn_id": "turn_003",
                    "content": "unknown intent",
                },
            }
        )
    except RuntimeError as exc:
        assert str(exc) == "target decision agent intent resolution failed"
    else:
        raise AssertionError(
            "handle_message should fail when intent resolution misses target agent"
        )

    assert adapter.sent == []
    assert inbound_event_repository.saved == []
