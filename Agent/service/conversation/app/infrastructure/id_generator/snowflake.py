import time
from collections.abc import Callable

from app.application.ports.id_generator import IdGenerator


class SnowflakeIdGenerator(IdGenerator):
    _NODE_ID_BITS = 10
    _SEQUENCE_BITS = 12
    _MAX_NODE_ID = (1 << _NODE_ID_BITS) - 1
    _MAX_SEQUENCE = (1 << _SEQUENCE_BITS) - 1
    _NODE_SHIFT = _SEQUENCE_BITS
    _TIME_SHIFT = _NODE_ID_BITS + _SEQUENCE_BITS

    def __init__(
        self,
        *,
        node_id: int,
        epoch_ms: int,
        time_ms: Callable[[], int] | None = None,
    ) -> None:
        if node_id < 0 or node_id > self._MAX_NODE_ID:
            raise ValueError("node_id must be between 0 and 1023")
        self._node_id = node_id
        self._epoch_ms = epoch_ms
        self._time_ms = time_ms or (lambda: int(time.time() * 1000))
        self._last_timestamp_ms = -1
        self._sequence = 0

    def next_id(self) -> int:
        timestamp_ms = self._time_ms()
        if timestamp_ms < self._last_timestamp_ms:
            raise RuntimeError("clock moved backwards")

        if timestamp_ms == self._last_timestamp_ms:
            self._sequence = (self._sequence + 1) & self._MAX_SEQUENCE
            if self._sequence == 0:
                timestamp_ms = self._wait_next_millis(timestamp_ms)
        else:
            self._sequence = 0

        self._last_timestamp_ms = timestamp_ms
        return (
            ((timestamp_ms - self._epoch_ms) << self._TIME_SHIFT)
            | (self._node_id << self._NODE_SHIFT)
            | self._sequence
        )

    def _wait_next_millis(self, timestamp_ms: int) -> int:
        next_timestamp_ms = self._time_ms()
        while next_timestamp_ms <= timestamp_ms:
            next_timestamp_ms = self._time_ms()
        return next_timestamp_ms
