from app.bootstrap.settings import Settings


def test_settings_load_default_service_baseline() -> None:
    settings = Settings()

    assert settings.service_name == "conversation_service"
    assert settings.api_prefix == "/api/conversation_service/v1"
    assert settings.auth_context_source == "nginx_ingress_x_user"
    assert settings.auth_user_header_name == "X-User"
    assert settings.redis_url == "redis://localhost:6379/0"
    assert settings.rabbitmq_exchange == "conversation.message.events"
    assert settings.core_agent_status_exchange == "core_agent.run_status.events"
    assert settings.core_agent_status_queue == "conversation.core_agent.run_status"
    assert settings.decision_agent_ag_ui_exchange == "decision_agent.ag_ui.events"
    assert settings.decision_agent_ag_ui_queue == "conversation.decision_agent.ag_ui"
    assert settings.decision_agent_ag_ui_routing_key == "decision_agent.session.ag_ui_event.v1"
    assert settings.snowflake_node_id == 1
    assert settings.snowflake_epoch_ms == 1_735_689_600_000


def test_settings_allow_environment_override(monkeypatch) -> None:
    monkeypatch.setenv("SNOWFLAKE_NODE_ID", "7")
    monkeypatch.setenv("LOG_LEVEL", "DEBUG")

    settings = Settings()

    assert settings.snowflake_node_id == 7
    assert settings.log_level == "DEBUG"
