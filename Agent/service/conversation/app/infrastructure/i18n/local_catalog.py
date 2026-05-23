from app.application.ports.error_catalog import ErrorDescriptor
from app.domain.shared.errors import ErrorReason


class LocalErrorCatalog:
    _BY_CODE: dict[str, ErrorDescriptor] = {
        "BAD_REQUEST": ErrorDescriptor(
            code="BAD_REQUEST",
            http_status=400,
            retryable=False,
            message_key="common.error.bad_request",
            scope="page",
        ),
        "UNAUTHORIZED": ErrorDescriptor(
            code="UNAUTHORIZED",
            http_status=401,
            retryable=True,
            message_key="common.error.unauthorized",
            scope="page",
        ),
        "FORBIDDEN": ErrorDescriptor(
            code="FORBIDDEN",
            http_status=403,
            retryable=False,
            message_key="common.error.forbidden",
            scope="page",
        ),
        "NOT_FOUND": ErrorDescriptor(
            code="NOT_FOUND",
            http_status=404,
            retryable=False,
            message_key="common.error.not_found",
            scope="page",
        ),
        "VALIDATION_FAILED": ErrorDescriptor(
            code="VALIDATION_FAILED",
            http_status=422,
            retryable=False,
            message_key="common.error.validation_failed",
            scope="field",
        ),
        "IDEMPOTENCY_CONFLICT": ErrorDescriptor(
            code="IDEMPOTENCY_CONFLICT",
            http_status=409,
            retryable=False,
            message_key="common.error.idempotency_conflict",
            scope="message",
        ),
        "INTERNAL_ERROR": ErrorDescriptor(
            code="INTERNAL_ERROR",
            http_status=500,
            retryable=True,
            message_key="common.error.internal_error",
            scope="system",
        ),
        "CONVERSATION_ARCHIVED": ErrorDescriptor(
            code="CONVERSATION_ARCHIVED",
            http_status=409,
            retryable=False,
            message_key="conversation.error.conversation_archived",
            scope="conversation",
        ),
        "CONVERSATION_EXPIRED": ErrorDescriptor(
            code="CONVERSATION_EXPIRED",
            http_status=409,
            retryable=False,
            message_key="conversation.error.conversation_expired",
            scope="conversation",
        ),
        "CONVERSATION_BUSY": ErrorDescriptor(
            code="CONVERSATION_BUSY",
            http_status=409,
            retryable=True,
            message_key="conversation.error.conversation_busy",
            scope="conversation",
        ),
        "CONVERSATION_STATE_CONFLICT": ErrorDescriptor(
            code="CONVERSATION_STATE_CONFLICT",
            http_status=409,
            retryable=False,
            message_key="conversation.error.conversation_state_conflict",
            scope="conversation",
        ),
        "CONVERSATION_WRITEBACK_STALE": ErrorDescriptor(
            code="CONVERSATION_WRITEBACK_STALE",
            http_status=409,
            retryable=False,
            message_key="conversation.error.conversation_writeback_stale",
            scope="message",
        ),
        "LEGAL_HOLD_PERMISSION_DENIED": ErrorDescriptor(
            code="LEGAL_HOLD_PERMISSION_DENIED",
            http_status=403,
            retryable=False,
            message_key="conversation.error.legal_hold_permission_denied",
            scope="page",
        ),
    }
    _REASON_TO_CODE: dict[ErrorReason, str] = {
        ErrorReason.INVALID_STATUS_TRANSITION: "VALIDATION_FAILED",
        ErrorReason.SENSITIVE_INPUT_REJECTED: "VALIDATION_FAILED",
        ErrorReason.CONVERSATION_ARCHIVED: "CONVERSATION_ARCHIVED",
        ErrorReason.CONVERSATION_EXPIRED: "CONVERSATION_EXPIRED",
        ErrorReason.CONVERSATION_BUSY: "CONVERSATION_BUSY",
        ErrorReason.CONVERSATION_STATE_CONFLICT: "CONVERSATION_STATE_CONFLICT",
        ErrorReason.CONVERSATION_NOT_FOUND: "NOT_FOUND",
        ErrorReason.CONVERSATION_ACCESS_DENIED: "FORBIDDEN",
        ErrorReason.CONVERSATION_WRITEBACK_STALE: "CONVERSATION_WRITEBACK_STALE",
        ErrorReason.CHILD_CONVERSATION_MISMATCH: "VALIDATION_FAILED",
    }

    def resolve_code(self, code: str) -> ErrorDescriptor:
        return self._BY_CODE.get(code, self._BY_CODE["INTERNAL_ERROR"])

    def resolve_reason(self, reason: ErrorReason) -> ErrorDescriptor:
        return self.resolve_code(self._REASON_TO_CODE.get(reason, "BAD_REQUEST"))


class LocalMessageCatalog:
    _MESSAGES: dict[str, dict[str, str]] = {
        "zh-CN": {
            "common.error.bad_request": "请求参数不正确。",
            "common.error.unauthorized": "身份信息缺失或无效。",
            "common.error.forbidden": "没有权限执行该操作。",
            "common.error.not_found": "资源不存在。",
            "common.error.validation_failed": "输入内容不符合要求，请检查后重试。",
            "common.error.idempotency_conflict": "请求幂等键已被不同内容使用。",
            "common.error.internal_error": "服务暂时不可用，请稍后重试。",
            "conversation.error.conversation_archived": "会话已归档，恢复后才能继续。",
            "conversation.error.conversation_expired": "会话已过期，仅支持历史查看。",
            "conversation.error.conversation_busy": "上一轮对话仍在处理中，请稍后再试。",
            "conversation.error.conversation_state_conflict": "当前会话状态不允许执行该操作。",
            "conversation.error.conversation_writeback_stale": "Agent 写回已过期，请刷新会话。",
            "conversation.error.legal_hold_permission_denied": "没有权限修改保留标记。",
        },
        "en-US": {
            "common.error.bad_request": "The request is invalid.",
            "common.error.unauthorized": "Authentication information is missing or invalid.",
            "common.error.forbidden": "You do not have permission to perform this action.",
            "common.error.not_found": "The resource was not found.",
            "common.error.validation_failed": (
                "The input is invalid. Please review it and try again."
            ),
            "common.error.idempotency_conflict": (
                "The idempotency key has already been used with different content."
            ),
            "common.error.internal_error": "The service is temporarily unavailable.",
            "conversation.error.conversation_archived": (
                "This conversation is archived. Restore it before continuing."
            ),
            "conversation.error.conversation_expired": (
                "This conversation has expired and is read-only."
            ),
            "conversation.error.conversation_busy": (
                "The previous turn is still being processed. Please try again later."
            ),
            "conversation.error.conversation_state_conflict": (
                "The current conversation state does not allow this operation."
            ),
            "conversation.error.conversation_writeback_stale": (
                "The agent writeback is stale. Refresh the conversation."
            ),
            "conversation.error.legal_hold_permission_denied": (
                "You do not have permission to change legal hold."
            ),
        },
    }

    def render(self, message_key: str, locale: str) -> str:
        messages = self._MESSAGES.get(locale) or self._MESSAGES["zh-CN"]
        return messages.get(message_key) or self._MESSAGES["zh-CN"].get(
            message_key,
            "服务暂时不可用，请稍后重试。",
        )
