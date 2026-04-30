from typing import Protocol


class MessageCatalog(Protocol):
    def render(self, message_key: str, locale: str) -> str:
        raise NotImplementedError
