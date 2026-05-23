from collections.abc import Callable
from typing import Protocol
from uuid import uuid4

RELEASE_LOCK_SCRIPT = """
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
"""


class RedisLockClient(Protocol):
    async def set(self, name: str, value: str, *, nx: bool, px: int) -> bool:
        raise NotImplementedError

    async def eval(self, script: str, numkeys: int, *keys_and_args: str) -> int:
        raise NotImplementedError


class RedisLockLease:
    def __init__(self, redis_client: RedisLockClient, key: str, owner_token: str) -> None:
        self._redis_client = redis_client
        self._key = key
        self._owner_token = owner_token
        self._released = False

    async def release(self) -> None:
        if self._released:
            return
        await self._redis_client.eval(
            RELEASE_LOCK_SCRIPT,
            1,
            self._key,
            self._owner_token,
        )
        self._released = True


class RedisGlobalLock:
    def __init__(
        self,
        redis_client: RedisLockClient,
        *,
        owner_token_factory: Callable[[], str] | None = None,
    ) -> None:
        self._redis_client = redis_client
        self._owner_token_factory = owner_token_factory or (lambda: uuid4().hex)

    async def acquire(self, key: str, ttl_ms: int) -> RedisLockLease | None:
        owner_token = self._owner_token_factory()
        acquired = await self._redis_client.set(key, owner_token, nx=True, px=ttl_ms)
        if not acquired:
            return None
        return RedisLockLease(self._redis_client, key, owner_token)
