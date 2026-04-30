from app.infrastructure.observability.logging import redact_event_dict


def test_log_redaction_masks_sensitive_values_recursively() -> None:
    event = {
        "event": "security.test",
        "authorization": "Bearer access_token=secret",
        "nested": {
            "password": "secret-password",
            "email": "user@example.com",
            "safe": "line1\r\nline2",
        },
    }

    redacted = redact_event_dict(event)

    assert redacted["authorization"] == "[REDACTED:credential]"
    assert redacted["nested"]["password"] == "[REDACTED:credential]"
    assert redacted["nested"]["email"] == "u***@example.com"
    assert redacted["nested"]["safe"] == "line1 line2"
