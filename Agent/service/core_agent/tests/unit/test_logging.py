import logging
from pathlib import Path

from core_agent_service.config.settings import Settings
from core_agent_service.infrastructure.logging import configure_logging


def test_configure_logging_writes_to_file(tmp_path):
    log_file = tmp_path / "service.log"
    settings = Settings(
        database_url="sqlite://",
        kweaver_base_url="https://example.com",
        kweaver_decision_agent_id="agent_001",
        core_agent_log_file=str(log_file),
    )

    configure_logging(settings)
    logger = logging.getLogger("core_agent_service.test")
    logger.info("file logging works")

    assert log_file.exists()
    assert "file logging works" in log_file.read_text(encoding="utf-8")


def test_configure_logging_adds_stdout_handler_when_enabled(tmp_path):
    log_file = tmp_path / "stdout.log"
    settings = Settings(
        database_url="sqlite://",
        kweaver_base_url="https://example.com",
        kweaver_decision_agent_id="agent_001",
        core_agent_log_file=str(log_file),
        core_agent_log_to_stdout=True,
    )

    configure_logging(settings)

    root_logger = logging.getLogger()
    assert any(isinstance(handler, logging.StreamHandler) for handler in root_logger.handlers)


def test_configure_logging_skips_file_handler_when_log_file_is_blank():
    settings = Settings(
        database_url="sqlite://",
        kweaver_base_url="https://example.com",
        kweaver_decision_agent_id="agent_001",
        core_agent_log_file="",
        core_agent_log_to_stdout=True,
    )

    configure_logging(settings)

    root_logger = logging.getLogger()
    assert any(isinstance(handler, logging.StreamHandler) for handler in root_logger.handlers)
    assert not any(isinstance(handler, logging.FileHandler) for handler in root_logger.handlers)
