from dataclasses import dataclass
from typing import Any

from app.application.models.conversation import (
    ConversationMessageRecord,
    ConversationRecord,
    ConversationStatusEventRecord,
)


@dataclass(frozen=True, slots=True)
class ConversationWritebackIdempotencyRecord:
    writeback_id: int
    idempotency_key: str
    conversation_id: int
    output_id: str | None
    request_hash: str
    result_status: str
    created_time: int
    updated_time: int
    result_message_id: int | None = None
    reject_code: str | None = None
    reject_reason: str | None = None
    trace_id: str | None = None
    correlation_id: str | None = None


@dataclass(frozen=True, slots=True)
class ConversationContextDeltaRecord:
    context_delta_id: int
    conversation_id: int
    delta_payload: dict[str, Any]
    merge_status: str
    created_time: int
    updated_time: int
    turn_id: int | None = None
    source_message_id: int | None = None
    base_snapshot_version: int | None = None
    created_by_agent: str | None = None
    trace_id: str | None = None


@dataclass(frozen=True, slots=True)
class ConversationContextSnapshotRecord:
    context_snapshot_id: int
    conversation_id: int
    snapshot_version: int
    short_summary: str
    structured_state: dict[str, Any]
    status: str
    created_by: str
    created_time: int
    updated_time: int
    last_message_id: int | None = None
    trace_id: str | None = None


@dataclass(frozen=True, slots=True)
class ConversationReasoningTraceRecord:
    reasoning_trace_id: str
    conversation_id: int
    trace_payload: dict[str, Any]
    created_time: int
    updated_time: int
    source_message_id: int | None = None
    core_agent_run_id: str | None = None
    created_by_agent: str | None = None
    trace_id: str | None = None


@dataclass(frozen=True, slots=True)
class ConversationCandidateSelectionRecord:
    selection_id: int
    conversation_id: int
    reasoning_trace_id: str
    candidate_option_id: str
    action: str
    idempotency_key: str
    created_by_user_id: str
    created_time: int
    updated_time: int
    comment: str | None = None
    trace_id: str | None = None
    correlation_id: str | None = None


@dataclass(frozen=True, slots=True)
class DecisionAgentAgUiEventResult:
    result_status: str
    idempotent: bool
    conversation: ConversationRecord | None = None
    message: ConversationMessageRecord | None = None
    status_event: ConversationStatusEventRecord | None = None
    reject_code: str | None = None
    reject_reason: str | None = None


@dataclass(frozen=True, slots=True)
class ConversationContextResult:
    conversation_id: int
    snapshot: ConversationContextSnapshotRecord | None
    pending_delta_count: int
    conversation: ConversationRecord | None = None


@dataclass(frozen=True, slots=True)
class ContextMergeRunResult:
    lock_acquired: bool
    scanned_delta_count: int = 0
    merged_delta_count: int = 0
    failed_delta_count: int = 0
