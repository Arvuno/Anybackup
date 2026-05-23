from app.domain.shared.redaction import RedactionAction, RedactionPolicy


def test_redaction_policy_rejects_high_risk_credentials() -> None:
    policy = RedactionPolicy()

    decision = policy.inspect_user_text("please use access_token=abc123")

    assert decision.action is RedactionAction.REJECT
    assert decision.reason == "credential"


def test_redaction_policy_allows_ordinary_user_text() -> None:
    policy = RedactionPolicy()

    decision = policy.inspect_user_text("restore the order database from last backup")

    assert decision.action is RedactionAction.ALLOW
    assert decision.reason is None
