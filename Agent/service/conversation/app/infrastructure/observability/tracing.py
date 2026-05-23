from uuid import uuid4

from opentelemetry import propagate, trace


def extract_trace_id_from_traceparent(traceparent: str | None) -> str | None:
    if traceparent is None:
        return None
    parts = traceparent.split("-")
    if len(parts) != 4:
        return None
    trace_id = parts[1].lower()
    if len(trace_id) != 32:
        return None
    if not all(character in "0123456789abcdef" for character in trace_id):
        return None
    if trace_id == "0" * 32:
        return None
    return trace_id


def request_trace_id(traceparent: str | None) -> str:
    return extract_trace_id_from_traceparent(traceparent) or _current_span_trace_id() or uuid4().hex


def inject_trace_context(carrier: dict[str, str], *, trace_id: str | None = None) -> None:
    if trace_id is None:
        propagate.inject(carrier)
        return
    carrier["traceparent"] = f"00-{trace_id}-0000000000000000-01"


def _current_span_trace_id() -> str | None:
    span_context = trace.get_current_span().get_span_context()
    if not span_context.is_valid:
        return None
    return f"{span_context.trace_id:032x}"
