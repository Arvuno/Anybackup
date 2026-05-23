from dataclasses import dataclass

from app.domain.shared.errors import ErrorReason


@dataclass(frozen=True, slots=True)
class ErrorDescriptor:
    code: str
    http_status: int
    retryable: bool
    message_key: str


LOCAL_DOMAIN_ERROR_CATALOG: dict[ErrorReason, ErrorDescriptor] = {
    ErrorReason.INVALID_STATUS_TRANSITION: ErrorDescriptor(
        code="VALIDATION_FAILED",
        http_status=422,
        retryable=False,
        message_key="error.invalid_status_transition",
    ),
    ErrorReason.SENSITIVE_INPUT_REJECTED: ErrorDescriptor(
        code="VALIDATION_FAILED",
        http_status=422,
        retryable=False,
        message_key="error.sensitive_input_rejected",
    ),
    ErrorReason.CONVERSATION_ARCHIVED: ErrorDescriptor(
        code="CONVERSATION_ARCHIVED",
        http_status=409,
        retryable=False,
        message_key="error.conversation_archived",
    ),
    ErrorReason.CONVERSATION_EXPIRED: ErrorDescriptor(
        code="CONVERSATION_EXPIRED",
        http_status=409,
        retryable=False,
        message_key="error.conversation_expired",
    ),
    ErrorReason.CONVERSATION_BUSY: ErrorDescriptor(
        code="CONVERSATION_BUSY",
        http_status=409,
        retryable=True,
        message_key="error.conversation_busy",
    ),
    ErrorReason.CONVERSATION_STATE_CONFLICT: ErrorDescriptor(
        code="CONVERSATION_STATE_CONFLICT",
        http_status=409,
        retryable=False,
        message_key="error.conversation_state_conflict",
    ),
    ErrorReason.CONVERSATION_NOT_FOUND: ErrorDescriptor(
        code="CONVERSATION_NOT_FOUND",
        http_status=404,
        retryable=False,
        message_key="error.conversation_not_found",
    ),
    ErrorReason.CONVERSATION_ACCESS_DENIED: ErrorDescriptor(
        code="FORBIDDEN",
        http_status=403,
        retryable=False,
        message_key="error.conversation_access_denied",
    ),
    ErrorReason.CONVERSATION_WRITEBACK_STALE: ErrorDescriptor(
        code="CONVERSATION_WRITEBACK_STALE",
        http_status=409,
        retryable=False,
        message_key="error.conversation_writeback_stale",
    ),
    ErrorReason.CHILD_CONVERSATION_MISMATCH: ErrorDescriptor(
        code="VALIDATION_FAILED",
        http_status=422,
        retryable=False,
        message_key="error.child_conversation_mismatch",
    ),
}
