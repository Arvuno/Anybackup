from core_agent_service.main import main
from core_agent_service.bootstrap import build_app
from core_agent_service.config.settings import Settings
from core_agent_service.infrastructure.kweaver.adapter import KWeaverRelayAdapter


def test_build_app_exposes_minimal_runtime_components(tmp_path):
    original_probe = KWeaverRelayAdapter.probe_connectivity
    KWeaverRelayAdapter.probe_connectivity = lambda self: {"decision_conversation_id": "decision_conv_created"}
    try:
        app = build_app(
            Settings(
                database_url="sqlite://",
                kweaver_base_url="https://example.com",
                kweaver_decision_agent_id="agent_001",
                foundation_endpoint="https://115.190.150.254:9600",
                foundation_ak="foundation-ak",
                foundation_sk="foundation-sk",
                core_agent_log_file=str(tmp_path / "core_agent_service.log"),
            )
        )

        assert "settings" in app
        assert "consumer" in app
        assert "consumer_runner" in app
        assert "consumers" in app
        assert "publisher" in app
        assert "relay_service" in app
        assert len(app["consumers"]) == 1
        assert app["consumer_runner"] is app["consumer"]
        assert app["relay_service"]._adapter._base_url == "https://example.com"
        assert app["relay_service"]._adapter._chat_timeout is None
        assert app["relay_service"]._adapter._stream_progress_interval == 100
        assert app["relay_service"]._adapter._stream_trace_enabled is True
        assert app["relay_service"]._mq_runtime_info == {
            "rabbitmq_url": "amqp://guest:guest@localhost:5672/"
        }
        assert app["relay_service"]._foundation_runtime_info == {
            "endpoint": "https://115.190.150.254:9600",
            "ak": "foundation-ak",
            "sk": "foundation-sk",
        }
    finally:
        KWeaverRelayAdapter.probe_connectivity = original_probe


def test_main_runs_consumer_once(monkeypatch):
    calls: list[str] = []

    class FakeConsumer:
        def run(self) -> None:
            calls.append("run")

    monkeypatch.setattr(
        "core_agent_service.main.build_app",
        lambda: {
            "settings": object(),
            "relay_service": object(),
            "publisher": object(),
            "consumer": FakeConsumer(),
        },
    )

    assert main() == 0
    assert calls == ["run"]


def test_main_runs_consumer_runner_when_present(monkeypatch):
    calls: list[str] = []

    class FakeConsumer:
        def run(self) -> None:
            calls.append("consumer")

    class FakeRunner:
        def run(self) -> None:
            calls.append("runner")

    monkeypatch.setattr(
        "core_agent_service.main.build_app",
        lambda: {
            "settings": object(),
            "relay_service": object(),
            "publisher": object(),
            "consumer": FakeConsumer(),
            "consumer_runner": FakeRunner(),
        },
    )

    assert main() == 0
    assert calls == ["runner"]


def test_build_app_creates_multiple_independent_consumers(monkeypatch, tmp_path):
    monkeypatch.setattr(KWeaverRelayAdapter, "probe_connectivity", lambda self: {"decision_conversation_id": "decision_conv_created"})

    app = build_app(
        Settings(
            database_url=f"sqlite+pysqlite:///{tmp_path / 'runtime.db'}",
            kweaver_base_url="https://example.com",
            kweaver_decision_agent_id="agent_001",
            kweaver_probe_on_startup=False,
            rabbitmq_consumer_count=3,
            core_agent_log_file=str(tmp_path / "core_agent_service.log"),
        )
    )

    assert len(app["consumers"]) == 3
    assert app["consumer"] is app["consumers"][0]
    assert app["consumer_runner"].consumers == tuple(app["consumers"])
    assert len({id(consumer.session) for consumer in app["consumers"]}) == 3


def test_build_app_probes_decision_agent_by_default(monkeypatch, tmp_path):
    calls: list[str] = []

    def fake_probe(self):
        calls.append("hello")
        return {"decision_conversation_id": "decision_conv_created"}

    monkeypatch.setattr(KWeaverRelayAdapter, "probe_connectivity", fake_probe)

    build_app(
        Settings(
            database_url="sqlite://",
            kweaver_base_url="https://example.com",
            kweaver_decision_agent_id="agent_001",
            core_agent_log_file=str(tmp_path / "core_agent_service.log"),
        )
    )

    assert calls == ["hello"]


def test_build_app_probes_decision_agent_when_enabled(monkeypatch, tmp_path):
    calls: list[str] = []

    def fake_probe(self):
        calls.append("hello")
        return {"decision_conversation_id": "decision_conv_created"}

    monkeypatch.setattr(KWeaverRelayAdapter, "probe_connectivity", fake_probe)

    build_app(
        Settings(
            database_url="sqlite://",
            kweaver_base_url="https://example.com",
            kweaver_decision_agent_id="agent_001",
            kweaver_probe_on_startup=True,
            core_agent_log_file=str(tmp_path / "core_agent_service.log"),
        )
    )

    assert calls == ["hello"]
