from __future__ import annotations

import asyncio
import json
import logging
import time
from dataclasses import replace
from typing import Any, Protocol

from conversation_agent_mq_mock.messages import (
    IncomingConversationMessage,
    OutgoingMqMessage,
    build_ag_ui_message,
    build_core_status_message,
)
from conversation_agent_mq_mock.scenario import (
    build_scenario,
    should_fail,
    should_reject,
    should_replay,
)
from conversation_agent_mq_mock.settings import Settings
from conversation_agent_mq_mock.state import SessionTracker, core_agent_run_id

logger = logging.getLogger(__name__)


class MessagePublisher(Protocol):
    async def publish(self, message: OutgoingMqMessage) -> None:
        raise NotImplementedError


class AgentMqMockRunner:
    def __init__(
        self,
        *,
        settings: Settings,
        publisher: MessagePublisher,
        tracker: SessionTracker | None = None,
    ) -> None:
        self._settings = settings
        self._publisher = publisher
        self._tracker = tracker or SessionTracker()

    async def handle_body(self, body: dict[str, Any]) -> None:
        incoming = IncomingConversationMessage.from_body(body)
        if not self._tracker.mark_seen(incoming.event_id):
            logger.info(
                "mock_duplicate_input_skipped",
                extra={
                    "event_id": incoming.event_id,
                    "conversation_id": incoming.conversation_id,
                    "message_id": incoming.message_id,
                },
            )
            return

        run_id = core_agent_run_id(incoming.conversation_id, incoming.message_id)
        logger.info(
            "mock_input_accepted",
            extra={
                "event_id": incoming.event_id,
                "conversation_id": incoming.conversation_id,
                "message_id": incoming.message_id,
                "core_agent_run_id": run_id,
            },
        )

        if should_reject(incoming):
            await self._publisher.publish(
                build_core_status_message(
                    incoming,
                    kind="rejected",
                    core_agent_run_id=run_id,
                    now_ms=_now_ms(),
                    status_queue=self._settings.core_status_queue,
                    error_code="UNSUPPORTED_TASK",
                    error_message="Mock core agent rejected this unsupported task.",
                    retryable=False,
                )
            )
            return

        await self._publisher.publish(
            build_core_status_message(
                incoming,
                kind="accepted",
                core_agent_run_id=run_id,
                now_ms=_now_ms(),
                status_queue=self._settings.core_status_queue,
            )
        )
        await _sleep_ms(self._settings.delay_ms)

        if should_fail(incoming):
            await self._publisher.publish(
                build_core_status_message(
                    incoming,
                    kind="failed",
                    core_agent_run_id=run_id,
                    now_ms=_now_ms(),
                    status_queue=self._settings.core_status_queue,
                    error_code="CORE_AGENT_MOCK_FAILURE",
                    error_message="Mock core agent generated a simulated failure.",
                    retryable=True,
                )
            )
            return

        scenario = build_scenario(incoming, core_agent_run_id=run_id)
        sequence_numbers = self._tracker.reserve_ag_ui_sequences(
            incoming.conversation_id,
            len(scenario.ag_ui_steps),
        )
        replay_output = should_replay(incoming)
        for step, sequence in zip(scenario.ag_ui_steps, sequence_numbers, strict=True):
            message = build_ag_ui_message(
                incoming,
                step=replace(step, sequence=sequence),
                now_ms=_now_ms(),
                exchange=self._settings.ag_ui_exchange,
                routing_key=self._settings.ag_ui_routing_key,
            )
            await self._publisher.publish(message)
            logger.info(
                "mock_ag_ui_published",
                extra={
                    "conversation_id": incoming.conversation_id,
                    "message_id": incoming.message_id,
                    "sequence": step.sequence,
                    "scenario_id": scenario.scenario_id,
                },
            )
            if replay_output:
                await self._publisher.publish(message)
                logger.info(
                    "mock_ag_ui_replayed",
                    extra={
                        "conversation_id": incoming.conversation_id,
                        "message_id": incoming.message_id,
                        "sequence": step.sequence,
                        "scenario_id": scenario.scenario_id,
                    },
                )
            await _sleep_ms(self._settings.delay_ms)


def decode_json_body(raw: bytes) -> dict[str, Any]:
    body = json.loads(raw.decode("utf-8"))
    if not isinstance(body, dict):
        raise ValueError("MQ body must be a JSON object")
    return body


async def _sleep_ms(delay_ms: int) -> None:
    if delay_ms > 0:
        await asyncio.sleep(delay_ms / 1000)


def _now_ms() -> int:
    return int(time.time() * 1000)
