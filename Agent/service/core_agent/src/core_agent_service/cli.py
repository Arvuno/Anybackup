from __future__ import annotations

import argparse
import json
from datetime import datetime, timezone
from uuid import uuid4

import pika

from core_agent_service.bootstrap import build_app
from core_agent_service.config.settings import Settings


def _new_event_id() -> str:
    return f"evt_{uuid4().hex[:12]}"


def _new_message_id() -> str:
    return f"msg_{uuid4().hex[:12]}"


def _build_sample_event(conversation_id: str, content: str, event_type: str) -> dict:
    # CLI 测试事件保持和真实消费格式一致，避免出现“命令能跑、实际消息不兼容”的偏差。
    payload = {
        "conversation_id": conversation_id,
        "content": content,
    }
    if event_type == "conversation.message.sent":
        # 发送事件必须带用户消息 ID，取消事件仍保持最小取消指令结构。
        message_id = _new_message_id()
        payload["message_id"] = message_id
        payload["turn_id"] = message_id
    return {
        "event_id": _new_event_id(),
        "event_type": event_type,
        "occurred_at": datetime.now(timezone.utc).isoformat(),
        "source_service": "local_cli",
        "payload": payload,
    }


def publish_event_to_rabbitmq(settings, event: dict, *, routing_key: str, connection_factory=None) -> None:
    if connection_factory is not None:
        connection = connection_factory()
    else:
        connection = pika.BlockingConnection(pika.URLParameters(settings.rabbitmq_url))

    channel = connection.channel()
    channel.exchange_declare(
        exchange=settings.rabbitmq_exchange,
        exchange_type=settings.rabbitmq_exchange_type,
        durable=True,
    )
    channel.basic_publish(
        exchange=settings.rabbitmq_exchange,
        routing_key=routing_key,
        body=json.dumps(event, ensure_ascii=False),
        properties=pika.BasicProperties(
            content_type="application/json",
            delivery_mode=2,
        ),
    )


def _cmd_publish_sample(args) -> int:
    # publish-sample 直接走本地 consumer，适合不依赖 RabbitMQ 的快速冒烟验证。
    app = build_app()
    body = json.dumps(
        _build_sample_event(
            conversation_id=args.conversation_id,
            content=args.content,
            event_type="conversation.message.sent",
        )
    ).encode("utf-8")
    app["consumer"].handle_delivery(body)
    print(json.dumps({"status": "published", "conversationId": args.conversation_id}, ensure_ascii=False))
    return 0


def _cmd_publish_rabbitmq(args) -> int:
    # publish-rabbitmq 会真实连到 RabbitMQ，适合验证完整中转链路。
    settings = Settings()
    event_type = "conversation.message.cancel_requested" if args.cancel else "conversation.message.sent"
    routing_key = "conversation.message.cancel_requested.v1" if args.cancel else "conversation.message.sent.v1"
    event = _build_sample_event(
        conversation_id=args.conversation_id,
        content=args.content,
        event_type=event_type,
    )
    publish_event_to_rabbitmq(
        settings,
        event,
        routing_key=routing_key,
    )
    print(
        json.dumps(
            {
                "status": "published_to_rabbitmq",
                "conversationId": args.conversation_id,
                "routingKey": routing_key,
            },
            ensure_ascii=False,
        )
    )
    return 0


def _cmd_show_mapping(args) -> int:
    # 这里直接读本地映射表，用于确认某个 conversation_id 是否已经绑定到平台会话。
    app = build_app()
    mapping = app["relay_service"]._mapping_repository.get_by_conversation_id(args.conversation_id)
    if mapping is None:
        print(json.dumps({"conversationId": args.conversation_id, "mapping": None}, ensure_ascii=False))
        return 0

    print(
        json.dumps(
            {
                "conversationId": args.conversation_id,
                "decisionAgentId": mapping.decision_agent_id,
                "decisionConversationId": mapping.decision_conversation_id,
                "mappingStatus": mapping.mapping_status,
            },
            ensure_ascii=False,
        )
    )
    return 0


def build_parser() -> argparse.ArgumentParser:
    # CLI 只保留最小原型调试命令，不扩展成复杂运维入口。
    parser = argparse.ArgumentParser(prog="python -m core_agent_service.cli")
    subparsers = parser.add_subparsers(dest="command", required=True)

    publish_parser = subparsers.add_parser("publish-sample")
    publish_parser.add_argument("--conversation-id", required=True)
    publish_parser.add_argument("--content", required=True)
    publish_parser.set_defaults(func=_cmd_publish_sample)

    rabbitmq_publish_parser = subparsers.add_parser("publish-rabbitmq")
    rabbitmq_publish_parser.add_argument("--conversation-id", required=True)
    rabbitmq_publish_parser.add_argument("--content", required=True)
    rabbitmq_publish_parser.add_argument("--cancel", action="store_true")
    rabbitmq_publish_parser.set_defaults(func=_cmd_publish_rabbitmq)

    mapping_parser = subparsers.add_parser("show-mapping")
    mapping_parser.add_argument("--conversation-id", required=True)
    mapping_parser.set_defaults(func=_cmd_show_mapping)

    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    raise SystemExit(main())
