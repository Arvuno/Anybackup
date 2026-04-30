from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class CandidateSelectionCommand:
    conversation_id: int
    reasoning_trace_id: str
    candidate_option_id: str
    action: str
    idempotency_key: str
    user_id: str
    comment: str | None = None
    request_id: str | None = None
