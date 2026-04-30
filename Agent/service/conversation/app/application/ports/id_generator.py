from typing import Protocol


class IdGenerator(Protocol):
    def next_id(self) -> int:
        raise NotImplementedError
