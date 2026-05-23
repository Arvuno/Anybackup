import logging
from collections.abc import Callable
from typing import Any

from app.application.models.conversation import AuthenticatedUser
from app.application.models.writeback import (
    ContextMergeRunResult,
    ConversationContextResult,
    ConversationContextSnapshotRecord,
)
from app.application.ports.id_generator import IdGenerator
from app.application.ports.locking import GlobalLock
from app.application.ports.unit_of_work import UnitOfWork
from app.application.use_cases.access import ensure_conversation_owner

logger = logging.getLogger(__name__)


class GetConversationContextHandler:
    def __init__(self, *, unit_of_work_factory: Callable[[], UnitOfWork]) -> None:
        self._unit_of_work_factory = unit_of_work_factory

    async def handle(
        self,
        conversation_id: int,
        user: AuthenticatedUser,
    ) -> ConversationContextResult | None:
        logger.info(
            "conversation_context_get_enter",
            extra={"conversation_id": conversation_id, "user_id": user.user_id},
        )
        async with self._unit_of_work_factory() as unit_of_work:
            conversation = await unit_of_work.conversations.get_record_by_id(conversation_id)
            if conversation is None:
                return None
            ensure_conversation_owner(conversation, user)
            snapshot = await unit_of_work.context_snapshots.get_current_by_conversation(
                conversation_id
            )
            pending_count = await unit_of_work.context_deltas.count_pending_by_conversation(
                conversation_id
            )
            return ConversationContextResult(
                conversation_id=conversation_id,
                snapshot=snapshot,
                pending_delta_count=pending_count,
                conversation=conversation,
            )


class ContextMergeWorker:
    def __init__(
        self,
        *,
        unit_of_work_factory: Callable[[], UnitOfWork],
        id_generator: IdGenerator,
        lock: GlobalLock,
        lock_ttl_ms: int,
        batch_size: int = 100,
        fail_merge: bool = False,
    ) -> None:
        self._unit_of_work_factory = unit_of_work_factory
        self._id_generator = id_generator
        self._lock = lock
        self._lock_ttl_ms = lock_ttl_ms
        self._batch_size = batch_size
        self._fail_merge = fail_merge

    async def run_once(self, *, worker_id: str, now_ms: int) -> ContextMergeRunResult:
        logger.info("context_merge_worker_enter", extra={"worker_id": worker_id})
        lease = await self._lock.acquire(
            "conversation_service:lock:context_merge:worker",
            ttl_ms=self._lock_ttl_ms,
        )
        if lease is None:
            return ContextMergeRunResult(lock_acquired=False)
        merged = 0
        failed = 0
        scanned = 0
        try:
            async with self._unit_of_work_factory() as unit_of_work:
                deltas = await unit_of_work.context_deltas.list_pending(limit=self._batch_size)
                scanned = len(deltas)
                for delta in deltas:
                    try:
                        if self._fail_merge:
                            raise RuntimeError("forced context merge failure")
                        current = await unit_of_work.context_snapshots.get_current_by_conversation(
                            delta.conversation_id
                        )
                        next_version = await unit_of_work.context_snapshots.next_version(
                            delta.conversation_id
                        )
                        snapshot = _merge_snapshot(
                            context_snapshot_id=self._id_generator.next_id(),
                            delta_payload=delta.delta_payload,
                            current=current,
                            conversation_id=delta.conversation_id,
                            source_message_id=delta.source_message_id,
                            next_version=next_version,
                            now_ms=now_ms,
                            trace_id=delta.trace_id,
                        )
                        await unit_of_work.context_snapshots.add(snapshot)
                        await unit_of_work.context_deltas.mark_status(
                            context_delta_id=delta.context_delta_id,
                            merge_status="merged",
                            now_ms=now_ms,
                        )
                        merged += 1
                    except Exception:
                        logger.exception(
                            "context_merge_delta_failed",
                            extra={"context_delta_id": delta.context_delta_id},
                        )
                        await unit_of_work.context_deltas.mark_status(
                            context_delta_id=delta.context_delta_id,
                            merge_status="failed",
                            now_ms=now_ms,
                        )
                        failed += 1
            return ContextMergeRunResult(
                lock_acquired=True,
                scanned_delta_count=scanned,
                merged_delta_count=merged,
                failed_delta_count=failed,
            )
        finally:
            await lease.release()


def _merge_snapshot(
    *,
    context_snapshot_id: int,
    delta_payload: dict[str, Any],
    current: ConversationContextSnapshotRecord | None,
    conversation_id: int,
    source_message_id: int | None,
    next_version: int,
    now_ms: int,
    trace_id: str | None,
) -> ConversationContextSnapshotRecord:
    summary_delta = str(delta_payload.get("summary_delta") or "").strip()
    previous_summary = current.short_summary if current is not None else ""
    short_summary = " ".join(value for value in (previous_summary, summary_delta) if value).strip()
    structured_state = dict(current.structured_state if current is not None else {})
    key_variables = delta_payload.get("key_variables")
    if isinstance(key_variables, dict):
        structured_state.update(key_variables)
    return ConversationContextSnapshotRecord(
        context_snapshot_id=context_snapshot_id,
        conversation_id=conversation_id,
        snapshot_version=next_version,
        short_summary=short_summary or "No summary",
        structured_state=structured_state,
        last_message_id=source_message_id,
        status="current",
        created_by="summary_updater",
        trace_id=trace_id,
        created_time=now_ms,
        updated_time=now_ms,
    )
