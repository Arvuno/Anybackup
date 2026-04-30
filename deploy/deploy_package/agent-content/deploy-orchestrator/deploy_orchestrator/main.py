from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

from deploy_orchestrator.cli_runner import CliRunner, CommandExecutionError
from deploy_orchestrator.config_loader import load_config
from deploy_orchestrator.context import DeploymentContext
from deploy_orchestrator.discovery import DiscoveryService
from deploy_orchestrator.logger import configure_logging, get_logger
from deploy_orchestrator.planner import build_execution_plan
from deploy_orchestrator.resources import build_resource_executors
from deploy_orchestrator.schema import DeployConfig
from deploy_orchestrator.state_store import StateStore


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="deploy-orchestrator")
    parser.add_argument("--verbose", action="store_true")
    subparsers = parser.add_subparsers(dest="command", required=True)
    for name in ("plan", "apply", "status", "export-state"):
        subparser = subparsers.add_parser(name)
        subparser.add_argument("--verbose", action="store_true")
        subparser.add_argument("-f", "--file", required=True, dest="config_file")
        subparser.add_argument("--env", default=None)
        subparser.add_argument("--state-file", default=None)
    return parser


def _default_state_file(config: DeployConfig, config_file: str | Path) -> Path:
    root = Path(config_file).resolve().parent
    return root / ".deploy" / f"state.{config.env}.json"


def _initial_state(config: DeployConfig) -> dict[str, object]:
    return {
        "project": config.project,
        "env": config.env,
        "skills": {},
        "bkns": {},
        "agents": {},
        "bindings": {},
    }


def run_plan(config: DeployConfig, discovery: DiscoveryService, executors: dict[str, object]) -> list[dict[str, str]]:
    logger = get_logger(__name__)
    spec_lookup = {
        "skill": {item.name: item for item in config.skills},
        "bkn": {item.name: item for item in config.bkns},
        "agent": {item.name: item for item in config.agents},
        "binding": {item.name: item for item in config.bindings},
    }
    rows: list[dict[str, str]] = []
    base_path = Path(config.globals.get("__config_dir__", "."))
    for item in build_execution_plan(config, base_path=base_path):
        logger.info("Planning resource %s/%s", item.resource_kind, item.resource_name)
        executor = executors[item.resource_kind]
        spec = spec_lookup[item.resource_kind][item.resource_name]
        discovered = discovery.find_by_name(item.resource_kind, item.resource_name)
        rows.append(
            {
                "kind": item.resource_kind,
                "name": item.resource_name,
                "action": executor.plan(spec, discovered),
            }
        )
    return rows


def run_apply(
    config: DeployConfig,
    state_store: StateStore,
    discovery: DiscoveryService,
    executors: dict[str, object],
) -> dict[str, object]:
    logger = get_logger(__name__)
    state = _initial_state(config)
    context = DeploymentContext()
    spec_lookup = {
        "skill": {item.name: item for item in config.skills},
        "bkn": {item.name: item for item in config.bkns},
        "agent": {item.name: item for item in config.agents},
        "binding": {item.name: item for item in config.bindings},
    }
    base_path = Path(config.globals.get("__config_dir__", "."))
    for item in build_execution_plan(config, base_path=base_path):
        logger.info("Applying resource %s/%s", item.resource_kind, item.resource_name)
        executor = executors[item.resource_kind]
        spec = spec_lookup[item.resource_kind][item.resource_name]
        discovered = discovery.find_by_name(item.resource_kind, item.resource_name)
        executor.apply(spec, discovered, context, state)
    state_store.save(state)
    return state


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    log_file = configure_logging(args.command, args.config_file, verbose=args.verbose)
    logger = get_logger(__name__)
    logger.info("Starting command %s with config %s", args.command, Path(args.config_file).resolve())
    logger.info("Writing logs to %s", log_file)

    try:
        config = load_config(args.config_file)
        config.globals["__config_dir__"] = str(Path(args.config_file).resolve().parent)
        if args.env:
            config.env = args.env
        state_file = Path(args.state_file) if args.state_file else _default_state_file(config, args.config_file)
        logger.info("Using state file %s", state_file)
        state_store = StateStore(state_file)
        runner = CliRunner()
        business_domain = str(config.globals.get("business_domain", "") or "")
        workspace = str(Path(args.config_file).resolve().parent)
        discovery = DiscoveryService(runner=runner, business_domain=business_domain)
        executors = build_resource_executors(runner=runner, workspace=workspace, business_domain=business_domain)

        if args.command == "plan":
            rows = run_plan(config, discovery, executors)
            print(json.dumps(rows, ensure_ascii=False, indent=2))
            logger.info("Finished command %s", args.command)
            return 0
        if args.command == "apply":
            state = run_apply(config, state_store, discovery, executors)
            resource_count = sum(len(state[key]) for key in ("skills", "bkns", "agents", "bindings"))
            print("Command apply completed successfully", file=sys.stderr)
            print(f"Resources applied: {resource_count}", file=sys.stderr)
            print(f"Log file: {log_file}", file=sys.stderr)
            logger.info("Finished command %s", args.command)
            return 0
        if args.command == "status":
            print(json.dumps(state_store.load(), ensure_ascii=False, indent=2))
            logger.info("Finished command %s", args.command)
            return 0
        if args.command == "export-state":
            print(json.dumps(state_store.load(), ensure_ascii=False, indent=2))
            logger.info("Finished command %s", args.command)
            return 0
    except Exception:
        logger.exception("Command %s failed", args.command)
        print(f"Command {args.command} failed", file=sys.stderr)
        print(f"Reason: {_build_error_summary(sys.exc_info()[1])}", file=sys.stderr)
        print(f"Log file: {log_file}", file=sys.stderr)
        return 1
    parser.error(f"未知命令: {args.command}")
    return 2


def _build_error_summary(exc: BaseException | None) -> str:
    if isinstance(exc, CommandExecutionError):
        return "CLI command failed"
    if isinstance(exc, ValueError):
        return "Validation failed"
    if exc is None:
        return "Unknown error"
    return exc.__class__.__name__


if __name__ == "__main__":
    raise SystemExit(main())
