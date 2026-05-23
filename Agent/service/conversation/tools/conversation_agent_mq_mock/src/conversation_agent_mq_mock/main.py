from __future__ import annotations

import asyncio
import logging
import signal

from conversation_agent_mq_mock.amqp import RabbitMqMockService
from conversation_agent_mq_mock.settings import Settings


async def _main() -> None:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(name)s %(message)s",
    )
    settings = Settings.from_env()
    service = RabbitMqMockService(settings)

    loop = asyncio.get_running_loop()
    stop = asyncio.Event()
    for signame in ("SIGINT", "SIGTERM"):
        signum = getattr(signal, signame, None)
        if signum is not None:
            try:
                loop.add_signal_handler(signum, stop.set)
            except NotImplementedError:
                pass

    task = asyncio.create_task(service.run_forever())
    await stop.wait()
    task.cancel()
    await service.close()


def main() -> None:
    asyncio.run(_main())


if __name__ == "__main__":
    main()
