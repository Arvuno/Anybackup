from typing import Protocol


class LockLease(Protocol):
    async def release(self) -> None:
        raise NotImplementedError


class GlobalLock(Protocol):
    async def acquire(self, key: str, ttl_ms: int) -> LockLease | None:
        raise NotImplementedError
