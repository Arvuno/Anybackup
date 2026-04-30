from app.domain.shared.error_catalog import LOCAL_DOMAIN_ERROR_CATALOG
from app.domain.shared.errors import ErrorReason


def test_local_domain_error_catalog_covers_all_error_reasons() -> None:
    assert set(LOCAL_DOMAIN_ERROR_CATALOG) == set(ErrorReason)


def test_conversation_archived_error_descriptor_is_stable() -> None:
    descriptor = LOCAL_DOMAIN_ERROR_CATALOG[ErrorReason.CONVERSATION_ARCHIVED]

    assert descriptor.code == "CONVERSATION_ARCHIVED"
    assert descriptor.http_status == 409
    assert descriptor.retryable is False


def test_conversation_state_conflict_error_descriptor_is_stable() -> None:
    descriptor = LOCAL_DOMAIN_ERROR_CATALOG[ErrorReason.CONVERSATION_STATE_CONFLICT]

    assert descriptor.code == "CONVERSATION_STATE_CONFLICT"
    assert descriptor.http_status == 409
    assert descriptor.retryable is False
