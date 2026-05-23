from dataclasses import dataclass
from enum import StrEnum


class RedactionAction(StrEnum):
    ALLOW = "allow"
    REJECT = "reject"


@dataclass(frozen=True, slots=True)
class RedactionDecision:
    action: RedactionAction
    reason: str | None = None


class RedactionPolicy:
    _HIGH_RISK_TOKENS = (
        "access_token=",
        "refresh_token=",
        "password=",
        "api_key=",
        "private key",
        "system prompt",
        "internal reasoning",
    )

    def inspect_user_text(self, text: str) -> RedactionDecision:
        normalized = text.lower()
        if any(token in normalized for token in self._HIGH_RISK_TOKENS):
            return RedactionDecision(action=RedactionAction.REJECT, reason="credential")
        return RedactionDecision(action=RedactionAction.ALLOW)
