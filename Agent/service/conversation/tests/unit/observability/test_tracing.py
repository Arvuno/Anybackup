from app.infrastructure.observability.tracing import (
    extract_trace_id_from_traceparent,
    inject_trace_context,
)


def test_extract_trace_id_from_w3c_traceparent() -> None:
    trace_id = extract_trace_id_from_traceparent(
        "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
    )

    assert trace_id == "0af7651916cd43dd8448eb211c80319c"


def test_inject_trace_context_sets_traceparent_header() -> None:
    carrier: dict[str, str] = {}

    inject_trace_context(carrier, trace_id="0af7651916cd43dd8448eb211c80319c")

    assert carrier["traceparent"].startswith("00-0af7651916cd43dd8448eb211c80319c-")
