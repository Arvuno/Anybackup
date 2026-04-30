from __future__ import annotations

import json
import os
import shutil
import subprocess
from dataclasses import dataclass
from pathlib import Path

from deploy_orchestrator.logger import get_logger


@dataclass
class CommandResult:
    returncode: int
    stdout: str
    stderr: str


class CommandExecutionError(RuntimeError):
    def __init__(self, args: list[str], result: CommandResult):
        normalized_args = [str(arg) for arg in args]
        message = f"命令执行失败: {' '.join(normalized_args)}"
        if result.stderr.strip():
            message = f"{message}\n{result.stderr.strip()}"
        super().__init__(message)
        self.args_list = normalized_args
        self.result = result


class CliRunner:
    def __init__(self):
        self.logger = get_logger(__name__)

    def run(self, args: list[str], cwd: str | Path | None = None) -> CommandResult:
        normalized_args = [str(arg) for arg in args]
        resolved_args = self._resolve_args(normalized_args)
        self.logger.info("Executing command: %s", " ".join(normalized_args))
        completed = subprocess.run(
            resolved_args,
            cwd=cwd,
            capture_output=True,
            text=True,
            check=False,
        )
        self.logger.info("Command finished with returncode=%s", completed.returncode)
        if completed.stderr.strip():
            self.logger.warning("Command stderr: %s", completed.stderr.strip())
        return CommandResult(
            returncode=completed.returncode,
            stdout=completed.stdout,
            stderr=completed.stderr,
        )

    def run_checked(self, args: list[str], cwd: str | Path | None = None) -> CommandResult:
        result = self.run(args, cwd=cwd)
        if result.returncode != 0 and self._is_known_windows_node_teardown_bug(result):
            self.logger.warning("Ignoring known Windows Node teardown assertion because command returned JSON output")
            return result
        if result.returncode != 0:
            raise CommandExecutionError(args, result)
        return result

    def _resolve_args(self, args: list[str]) -> list[str]:
        if not args:
            return args
        executable = args[0]
        resolved = shutil.which(executable)
        if resolved:
            return [resolved, *args[1:]]
        fallback = self._resolve_windows_npm_global_command(executable)
        if fallback is not None:
            return [str(fallback), *args[1:]]
        return args

    @staticmethod
    def _resolve_windows_npm_global_command(executable: str) -> Path | None:
        if os.name != "nt" or executable != "kweaver":
            return None
        appdata = os.environ.get("APPDATA")
        if not appdata:
            return None
        candidate = Path(appdata) / "npm" / "kweaver.cmd"
        try:
            return candidate if candidate.exists() else None
        except OSError:
            return candidate

    @staticmethod
    def _is_known_windows_node_teardown_bug(result: CommandResult) -> bool:
        if "Assertion failed: !(handle->flags & UV_HANDLE_CLOSING)" not in result.stderr:
            return False
        text = result.stdout.strip()
        if not text:
            return False
        try:
            payload = json.loads(text)
        except json.JSONDecodeError:
            return False
        return isinstance(payload, (dict, list))
