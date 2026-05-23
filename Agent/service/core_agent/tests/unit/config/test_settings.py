import pytest

from core_agent_service.config.settings import Settings


def test_settings_reads_minimum_required_values(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")

    settings = Settings(_env_file=None)

    assert settings.database_url.startswith("postgresql+psycopg://")
    assert settings.rabbitmq_url.startswith("amqp://")
    assert settings.rabbitmq_exchange == "conversation.message.events"
    assert settings.rabbitmq_queue == "core_agent.message.events"
    assert settings.rabbitmq_status_exchange == "core_agent.run_status.events"
    assert settings.rabbitmq_status_routing_key is None
    assert settings.rabbitmq_consumer_count == 1
    assert settings.kweaver_base_url == "https://kweaver.example.com"
    assert settings.kweaver_token is None
    assert settings.kweaver_decision_agent_id == "agent_001"
    assert settings.kweaver_timeout == 30.0
    assert settings.kweaver_chat_timeout is None
    assert settings.kweaver_tls_insecure is True
    assert settings.kweaver_stream_trace_enabled is True
    assert settings.kweaver_stream_trace_dir == "/tmp/core-agent-service-trace-log"
    assert settings.core_agent_log_level == "DEBUG"
    assert settings.core_agent_log_to_stdout is False
    assert settings.core_agent_log_file.endswith("logs/core_agent_service.log")


def test_settings_reads_optional_token(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_TOKEN", "secret-token")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")

    settings = Settings(_env_file=None)

    assert settings.kweaver_token == "secret-token"


def test_settings_reads_optional_username_password(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_USERNAME", "demo_user")
    monkeypatch.setenv("KWEAVER_PASSWORD", "demo_pass")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")

    settings = Settings(_env_file=None)

    assert settings.kweaver_username == "demo_user"
    assert settings.kweaver_password == "demo_pass"


def test_settings_reads_optional_foundation_credentials(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("FOUNDATION_ENDPOINT", "https://115.190.150.254:9600")
    monkeypatch.setenv("FOUNDATION_AK", "foundation-ak")
    monkeypatch.setenv("FOUNDATION_SK", "foundation-sk")

    settings = Settings(_env_file=None)

    assert settings.foundation_endpoint == "https://115.190.150.254:9600"
    assert settings.foundation_ak == "foundation-ak"
    assert settings.foundation_sk == "foundation-sk"


def test_settings_rejects_partial_foundation_credentials(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("FOUNDATION_ENDPOINT", "https://115.190.150.254:9600")

    with pytest.raises(ValueError, match="FOUNDATION_ENDPOINT, FOUNDATION_AK and FOUNDATION_SK"):
        Settings(_env_file=None)


def test_settings_reads_optional_tls_insecure(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("KWEAVER_TLS_INSECURE", "true")

    settings = Settings(_env_file=None)

    assert settings.kweaver_tls_insecure is True


def test_settings_reads_optional_probe_on_startup(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("KWEAVER_PROBE_ON_STARTUP", "true")

    settings = Settings(_env_file=None)

    assert settings.kweaver_probe_on_startup is True


def test_settings_reads_optional_rabbitmq_consumer_count(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("RABBITMQ_CONSUMER_COUNT", "3")

    settings = Settings(_env_file=None)

    assert settings.rabbitmq_consumer_count == 3


def test_settings_rejects_invalid_rabbitmq_consumer_count(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("RABBITMQ_CONSUMER_COUNT", "0")

    with pytest.raises(ValueError, match="RABBITMQ_CONSUMER_COUNT"):
        Settings(_env_file=None)


def test_settings_reads_optional_chat_timeout(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("KWEAVER_CHAT_TIMEOUT", "900")

    settings = Settings(_env_file=None)

    assert settings.kweaver_chat_timeout == 900.0


def test_settings_treats_blank_chat_timeout_as_unlimited(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("KWEAVER_CHAT_TIMEOUT", "")

    settings = Settings(_env_file=None)

    assert settings.kweaver_chat_timeout is None


def test_settings_enables_probe_on_startup_by_default(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.delenv("KWEAVER_PROBE_ON_STARTUP", raising=False)

    settings = Settings(_env_file=None)

    assert settings.kweaver_probe_on_startup is True


def test_settings_reads_logging_flags(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+psycopg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")
    monkeypatch.setenv("CORE_AGENT_LOG_LEVEL", "DEBUG")
    monkeypatch.setenv("CORE_AGENT_LOG_TO_STDOUT", "true")
    monkeypatch.setenv("CORE_AGENT_LOG_FILE", "logs/custom.log")

    settings = Settings(_env_file=None)

    assert settings.core_agent_log_level == "DEBUG"
    assert settings.core_agent_log_to_stdout is True
    assert settings.core_agent_log_file == "logs/custom.log"


def test_settings_rejects_asyncpg_database_url(monkeypatch):
    monkeypatch.setenv("DATABASE_URL", "postgresql+asyncpg://user:pass@localhost:5432/core_agent")
    monkeypatch.setenv("KWEAVER_BASE_URL", "https://kweaver.example.com")
    monkeypatch.setenv("KWEAVER_DECISION_AGENT_ID", "agent_001")

    with pytest.raises(ValueError, match="postgresql\\+psycopg"):
        Settings(_env_file=None)
