from __future__ import annotations

from dataclasses import dataclass
import logging

from core_agent_service.domain.events import InboundEvent


logger = logging.getLogger(__name__)


@dataclass(slots=True)
class RelayResult:
    # 返回值只保留最小中转版真正需要向上游反馈的结果。
    status: str
    decision_conversation_id: str


class RelayService:
    def __init__(
        self,
        *,
        mapping_repository,
        inbound_event_repository,
        adapter,
        mq_runtime_info: dict | None = None,
        foundation_runtime_info: dict | None = None,
    ):
        self._mapping_repository = mapping_repository
        self._inbound_event_repository = inbound_event_repository
        self._adapter = adapter
        # MQ 运行时信息随 query 传给 Decision Agent，用于下游按实际连接上下文处理回调或排障。
        self._mq_runtime_info = mq_runtime_info or {}
        # Foundation 连接信息由部署期注入，只在转发业务 Agent 时进入 custom_querys。
        self._foundation_runtime_info = foundation_runtime_info

    def handle_message(self, raw_event: dict) -> RelayResult:
        # 应用层先把原始字典校验成领域对象，后续逻辑始终围绕稳定字段展开。
        event = InboundEvent.model_validate(raw_event)
        logger.info(
            "handling inbound message event",
            extra={
                "event_id": event.event_id,
                "event_type": event.event_type,
                "conversation_id": event.payload.conversation_id,
            },
        )
        existing_event = self._inbound_event_repository.get_by_event_id(event.event_id)
        mapping = self._mapping_repository.get_by_conversation_id(
            event.payload.conversation_id
        )
        if existing_event is not None:
            # 幂等重放只返回已有处理结果，不能再次触发 Decision Agent 或重复推进业务状态。
            logger.info(
                "duplicate inbound event ignored",
                extra={
                    "event_id": event.event_id,
                    "conversation_id": event.payload.conversation_id,
                    "processing_status": existing_event.processing_status,
                },
            )
            return RelayResult(
                status=existing_event.processing_status,
                decision_conversation_id=mapping.decision_conversation_id
                if mapping
                else "",
            )
        # 首次对话没有映射时会返回空 conversation_id，让 KWeaver 在首条消息时创建平台会话。
        if mapping is None:
            logger.info(
                "resolving target decision agent",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "message_id": event.payload.message_id,
                },
            )
            intent = self._adapter.resolve_intent(
                event.payload.content,
                conversation_id=event.payload.conversation_id,
                turn_id=event.payload.turn_id or "",
                source_message_id=event.payload.message_id or "",
                mq=self._mq_runtime_info,
                foundation=self._foundation_runtime_info,
            )
            if intent.is_intent_success is not True or not intent.agent_id:
                logger.error(
                    "target decision agent intent resolution failed",
                    extra={
                        "conversation_id": event.payload.conversation_id,
                        "message_id": event.payload.message_id,
                    },
                )
                raise RuntimeError("target decision agent intent resolution failed")
            target_agent_id = intent.agent_id
            logger.info(
                "target decision agent resolved",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "decision_agent_id": target_agent_id,
                    "decision_agent_name": intent.agent_name,
                },
            )
        else:
            # 旧数据可能尚未包含 decision_agent_id；这种情况下只能回退配置 Agent，并通过日志暴露风险。
            target_agent_id = (
                getattr(mapping, "decision_agent_id", None)
                or self._adapter.default_agent_id
            )
            logger.info(
                "target decision agent reused",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "decision_agent_id": target_agent_id,
                    "mapping_has_decision_agent_id": getattr(
                        mapping, "decision_agent_id", None
                    )
                    is not None,
                },
            )

        # Agent 变化时必须创建新会话，不能复用其他 Agent 的上下文
        mapping_agent_id = (
            getattr(mapping, "decision_agent_id", None) if mapping else None
        )
        agent_changed = (
            mapping_agent_id is not None and mapping_agent_id != target_agent_id
        )
        if agent_changed:
            logger.warning(
                "target decision agent changed, creating new kweaver conversation",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "old_agent_id": mapping_agent_id,
                    "new_agent_id": target_agent_id,
                    "old_decision_conversation_id": mapping.decision_conversation_id
                    if mapping
                    else None,
                },
            )
            requested_conversation_id = self._adapter.ensure_conversation(None)
        else:
            requested_conversation_id = self._adapter.ensure_conversation(
                mapping.decision_conversation_id if mapping else None
            )
        response_payload = self._adapter.relay_message(
            requested_conversation_id,
            conversation_id=event.payload.conversation_id,
            turn_id=event.payload.turn_id or "",
            source_message_id=event.payload.message_id or "",
            content=event.payload.content,
            mq=self._mq_runtime_info,
            foundation=self._foundation_runtime_info,
            agent_id=target_agent_id,
        )
        decision_conversation_id = str(response_payload["decision_conversation_id"])
        # 映射以平台真实返回为准，这样首次创建和后续复用都走同一条更新路径。
        if (
            mapping is None
            or mapping.decision_conversation_id != decision_conversation_id
            or getattr(mapping, "decision_agent_id", None) != target_agent_id
        ):
            logger.info(
                "upserting conversation mapping",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "decision_conversation_id": decision_conversation_id,
                    "decision_agent_id": target_agent_id,
                    "previous_decision_conversation_id": mapping.decision_conversation_id
                    if mapping
                    else None,
                },
            )
            self._mapping_repository.upsert_mapping(
                event.payload.conversation_id,
                decision_conversation_id,
                decision_agent_id=target_agent_id,
            )
        else:
            logger.info(
                "reusing existing conversation mapping",
                extra={
                    "conversation_id": event.payload.conversation_id,
                    "decision_conversation_id": decision_conversation_id,
                    "decision_agent_id": target_agent_id,
                },
            )
        # 入站事件单独落表，方便排查“消息到了没、处理到哪一步”的问题。
        self._inbound_event_repository.save(
            event_id=event.event_id,
            event_type=event.event_type,
            occurred_at=event.occurred_at,
            conversation_id=event.payload.conversation_id,
            message_id=event.payload.message_id,
            content=event.payload.content,
            processing_status="completed",
        )
        logger.info(
            "inbound event persisted",
            extra={
                "event_id": event.event_id,
                "conversation_id": event.payload.conversation_id,
            },
        )
        logger.info(
            "message relayed to decision agent",
            extra={
                "conversation_id": event.payload.conversation_id,
                "decision_conversation_id": decision_conversation_id,
                "decision_agent_id": target_agent_id,
            },
        )
        return RelayResult(
            status="completed",
            decision_conversation_id=decision_conversation_id,
        )

    def handle_cancel(self, raw_event: dict) -> RelayResult:
        # 取消事件当前阶段只记录事实并复用已有映射，不额外向 Decision Agent 发控制指令。
        event = InboundEvent.model_validate(raw_event)
        logger.info(
            "handling cancel event",
            extra={
                "event_id": event.event_id,
                "conversation_id": event.payload.conversation_id,
            },
        )
        self._inbound_event_repository.save(
            event_id=event.event_id,
            event_type=event.event_type,
            occurred_at=event.occurred_at,
            conversation_id=event.payload.conversation_id,
            message_id=event.payload.message_id,
            content=event.payload.content,
            processing_status="cancelled",
        )
        logger.info(
            "cancel event persisted",
            extra={
                "event_id": event.event_id,
                "conversation_id": event.payload.conversation_id,
            },
        )
        mapping = self._mapping_repository.get_by_conversation_id(
            event.payload.conversation_id
        )
        return RelayResult(
            status="completed",
            decision_conversation_id=self._adapter.ensure_conversation(
                mapping.decision_conversation_id if mapping else None
            ),
        )
