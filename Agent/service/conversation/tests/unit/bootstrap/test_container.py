from dependency_injector import providers

from app.bootstrap.app_factory import create_app
from app.bootstrap.container import Container
from app.bootstrap.settings import Settings


def test_container_settings_provider_can_be_overridden() -> None:
    container = Container()

    with container.settings.override(
        providers.Object(
            Settings(
                service_name="conversation_service_test",
                snowflake_node_id=9,
            )
        )
    ):
        settings = container.settings()

    assert settings.service_name == "conversation_service_test"
    assert settings.snowflake_node_id == 9


def test_app_factory_attaches_container_to_app_state() -> None:
    app = create_app()

    assert hasattr(app.state.container, "settings")
    assert app.state.container.settings().service_name == "conversation_service"
