from __future__ import annotations

import logging
import sys
from pathlib import Path


def configure_logging(settings) -> None:
    formatter = logging.Formatter(
        fmt="%(asctime)s %(levelname)s [%(name)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    handlers: list[logging.Handler] = []

    # 空字符串表示禁用文件日志，避免在容器里把工作目录误当成日志文件路径。
    raw_log_file = (settings.core_agent_log_file or "").strip()
    if raw_log_file:
        log_file = Path(raw_log_file)
        if not log_file.is_absolute():
            log_file = Path.cwd() / log_file
        log_file.parent.mkdir(parents=True, exist_ok=True)

        file_handler = logging.FileHandler(log_file, encoding="utf-8")
        file_handler.setFormatter(formatter)
        handlers.append(file_handler)

    if settings.core_agent_log_to_stdout:
        stream_handler = logging.StreamHandler(sys.stdout)
        stream_handler.setFormatter(formatter)
        handlers.append(stream_handler)

    if not handlers:
        raise ValueError("logging requires at least one enabled handler")

    root_logger = logging.getLogger()
    root_logger.handlers.clear()
    root_logger.setLevel(getattr(logging, settings.core_agent_log_level.upper(), logging.INFO))
    for handler in handlers:
        root_logger.addHandler(handler)

    # 调试阶段服务自身使用 DEBUG，但第三方库的底层连接细节会淹没业务链路日志。
    logging.getLogger("pika").setLevel(logging.CRITICAL)
    logging.getLogger("httpx").setLevel(logging.INFO)
    logging.getLogger("httpcore").setLevel(logging.INFO)

    logging.getLogger("core_agent_service.bootstrap").info(
        "logging configured",
        extra={
            "log_file": raw_log_file or None,
            "log_to_stdout": settings.core_agent_log_to_stdout,
            "log_level": settings.core_agent_log_level,
        },
    )
