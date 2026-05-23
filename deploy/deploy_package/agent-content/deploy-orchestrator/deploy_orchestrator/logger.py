from __future__ import annotations

import logging
import sys
from datetime import datetime
from pathlib import Path


def configure_logging(command: str, config_file: str | Path, verbose: bool = False) -> Path:
    config_path = Path(config_file).resolve()
    log_dir = config_path.parent / "log"
    log_dir.mkdir(parents=True, exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    log_file = log_dir / f"{command}-{timestamp}.log"

    handlers: list[logging.Handler] = [logging.FileHandler(log_file, encoding="utf-8")]
    if verbose:
        handlers.append(logging.StreamHandler(sys.stderr))

    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(name)s - %(message)s",
        handlers=handlers,
        force=True,
    )
    return log_file


def get_logger(name: str) -> logging.Logger:
    return logging.getLogger(name)
