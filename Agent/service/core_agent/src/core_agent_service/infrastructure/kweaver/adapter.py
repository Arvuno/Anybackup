from __future__ import annotations

import json
import logging
from dataclasses import dataclass
from pathlib import Path
from types import SimpleNamespace
from typing import Any
from uuid import uuid4

import httpx
from kweaver.resources.conversations import _parse_message

from core_agent_service.infrastructure.kweaver.client import KWeaverClientProtocol


logger = logging.getLogger(__name__)

KWEAVER_WEB_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0"
BUSINESS_AGENT_CHAT_OPTION = {
    "is_need_history": True,
    "is_need_doc_retrival_post_process": True,
    "is_need_progress": True,
    "enable_dependency_cache": True,
}


@dataclass(slots=True)
class IntentResolution:
    # 意图识别结果来自配置文件中的路由 Agent，后续会话必须使用命中的目标 Agent。
    agent_id: str | None
    agent_name: str | None
    is_intent_success: bool


@dataclass(slots=True)
class AgentChatTarget:
    # 业务会话请求沿用页面调用版本；agent_key 仅用于日志排障，不写入请求体。
    agent_id: str
    agent_key: str
    agent_version: str


class KWeaverRelayAdapter:
    def __init__(
        self,
        client: KWeaverClientProtocol,
        agent_id: str,
        base_url: str,
        chat_timeout: float | None = None,
        stream_progress_interval: int = 100,
        stream_trace_enabled: bool = False,
        stream_trace_dir: str = "/tmp/core-agent-service-trace-log",
    ):
        self._client = client
        self._agent_id = agent_id
        self._base_url = self._normalize_base_url(base_url)
        self._chat_timeout = chat_timeout
        self._stream_progress_interval = max(1, int(stream_progress_interval or 100))
        self._stream_trace_enabled = stream_trace_enabled
        self._stream_trace_dir = stream_trace_dir

    @property
    def default_agent_id(self) -> str:
        return self._agent_id

    def ensure_conversation(self, existing_conversation_id: str | None) -> str:
        if existing_conversation_id:
            logger.info(
                "decision agent conversation reuse: %s",
                json.dumps(
                    {"decision_conversation_id": existing_conversation_id},
                    ensure_ascii=False,
                ),
            )
            return existing_conversation_id
        # 首次对话不预先建会话，而是把空 conversation_id 交给 KWeaver 在 send_message 时创建。
        logger.info("decision agent conversation create on first message")
        return ""

    def probe_connectivity(self) -> dict[str, object]:
        logger.info(
            "decision agent connectivity probe start: %s",
            json.dumps({"agent_id": self._agent_id}, ensure_ascii=False),
        )
        response = self._verify_agent_accessible()
        logger.info(
            "decision agent connectivity probe success: %s",
            json.dumps(
                {
                    "agent_id": response.get("agent_id"),
                    "agent_name": response.get("agent_name"),
                    "agent_status": response.get("agent_status"),
                },
                ensure_ascii=False,
            ),
        )
        return response

    def _verify_agent_accessible(self) -> dict[str, object]:
        # 启动探针先校验智能体可见性，避免把地址、业务域或权限问题误判成 query 协议问题。
        try:
            agent = self._client.agents.get(self._agent_id)
        except Exception as exc:
            logger.exception(
                "decision agent accessibility probe failed: %s",
                json.dumps({"agent_id": self._agent_id}, ensure_ascii=False),
            )
            raise RuntimeError(
                f"decision agent is not accessible: agent_id={self._agent_id}"
            ) from exc

        logger.info(
            "decision agent accessibility probe success: %s",
            json.dumps(
                {
                    "agent_id": self._agent_id,
                    "agent_name": self._extract_string(agent, "name"),
                    "agent_status": self._extract_string(agent, "status"),
                },
                ensure_ascii=False,
            ),
        )
        return {
            "agent_id": self._agent_id,
            "agent_name": self._extract_string(agent, "name"),
            "agent_status": self._extract_string(agent, "status"),
        }

    def resolve_intent(
        self,
        content: str,
        *,
        conversation_id: str,
        turn_id: str,
        source_message_id: str,
        mq: dict,
        foundation: dict[str, Any] | None = None,
    ) -> IntentResolution:
        # 配置文件中的 Agent 只负责意图路由，请求必须直接透传用户原文，不能构造 query JSON。
        # 上游消息标识与 MQ 信息放入 custom_querys 供主控 Agent 按需读取。
        custom_querys: dict[str, Any] = {
            "conversation_id": conversation_id,
            "turn_id": turn_id,
            "source_message_id": source_message_id,
            "content": content,
            "mq": mq,
        }
        if foundation:
            custom_querys["foundation"] = foundation
        logger.info(
            "decision agent intent routing request: %s",
            json.dumps(
                {
                    "agent_id": self._agent_id,
                    "content_length": len(content),
                    "custom_querys_length": self._log_content_length(custom_querys),
                },
                ensure_ascii=False,
            ),
        )
        message = self._send_message(
            "",
            content,
            agent_id=self._agent_id,
            stream=True,
            custom_querys=custom_querys,
        )
        payload = self._parse_intent_payload(getattr(message, "content", None))
        logger.info(
            "decision agent intent routing payload: %s",
            json.dumps(self._sanitize_for_log(payload), ensure_ascii=False),
        )
        resolution = IntentResolution(
            agent_id=self._extract_payload_string(payload, "agent_id"),
            agent_name=self._extract_payload_string(payload, "agent_name"),
            is_intent_success=payload.get("is_intent_success") is True,
        )
        logger.info(
            "decision agent intent routing response: %s",
            json.dumps(
                {
                    "agent_id": resolution.agent_id,
                    "agent_name": resolution.agent_name,
                    "is_intent_success": resolution.is_intent_success,
                },
                ensure_ascii=False,
            ),
        )
        return resolution

    def relay_message(
        self,
        decision_conversation_id: str,
        *,
        conversation_id: str,
        turn_id: str,
        source_message_id: str,
        content: str,
        mq: dict,
        foundation: dict[str, Any] | None = None,
        agent_id: str | None = None,
    ) -> dict[str, object]:
        # 业务 Agent 的 query 只接收用户原文；上游消息标识与 MQ 信息放入 custom_querys 供业务侧按需读取。
        target_agent_id = agent_id or self._agent_id
        custom_querys: dict[str, Any] = {
            "conversation_id": conversation_id,
            "turn_id": turn_id,
            "source_message_id": source_message_id,
            "content": content,
            "mq": mq,
        }
        if foundation:
            custom_querys["foundation"] = foundation
        request_payload = {
            "agent_id": target_agent_id,
            "decision_conversation_id": decision_conversation_id,
            "query_length": len(content),
            "custom_querys_length": self._log_content_length(custom_querys),
            "stream": True,
            "chat_timeout": self._chat_timeout,
        }
        logger.info(
            "decision agent request payload: %s",
            json.dumps(request_payload, ensure_ascii=False),
        )
        # 真实请求只通过 SDK 发送，业务层不感知 conversations.send_message 的参数细节。
        message = self._send_message(
            decision_conversation_id,
            content,
            agent_id=target_agent_id,
            stream=True,
            custom_querys=custom_querys,
        )
        # 这里把 SDK 返回对象规整成稳定字典，避免上层代码依赖第三方对象结构。
        response_payload = {
            "decision_conversation_id": self._extract_conversation_id(message),
            "content": self._normalize_content(getattr(message, "content", None)),
            "reply_message_id": self._extract_string(message, "id"),
            "status": self._extract_string(message, "status"),
            "raw": self._serialize_message(message),
        }
        logger.info(
            "decision agent response payload: %s",
            json.dumps(
                self._sanitize_for_log(
                    self._include_content_lengths_for_log(response_payload)
                ),
                ensure_ascii=False,
            ),
        )
        return response_payload

    def _send_message(
        self,
        decision_conversation_id: str,
        query: str,
        *,
        agent_id: str,
        stream: bool = False,
        custom_querys: dict[str, Any] | None = None,
    ):
        # SDK 当前公开方法固定 chat/completion timeout=120 秒；这里在适配层复用 SDK resource 的 HTTP 客户端，
        # 只覆盖本次长耗时 Decision Agent 调用的等待时间，避免业务层感知 SDK 细节。
        conversations = self._client.conversations
        http_client = getattr(conversations, "_http", None)
        if http_client is None or not hasattr(http_client, "post"):
            request_payload = {
                "conversation_id": decision_conversation_id,
                "content": query,
                "content_length": len(query),
                "custom_querys_length": self._log_content_length(custom_querys),
                "agent_id": agent_id,
                "stream": stream,
            }
            logger.info(
                "kweaver sdk conversation request: %s",
                json.dumps(self._sanitize_for_log(request_payload), ensure_ascii=False),
            )
            message = conversations.send_message(
                decision_conversation_id,
                content=query,
                agent_id=agent_id,
                stream=stream,
            )
            if stream:
                message = self._collect_stream_chunks(
                    message, requested_conversation_id=decision_conversation_id
                )
            logger.info(
                "kweaver sdk conversation response: %s",
                json.dumps(
                    self._sanitize_for_log(
                        self._include_content_lengths_for_log(
                            self._serialize_message(message)
                        )
                    ),
                    ensure_ascii=False,
                ),
            )
            return message

        if stream:
            target = self._resolve_agent_chat_target(agent_id, http_client=http_client)
            path = f"/api/agent-factory/v1/app/{target.agent_key}/chat/completion"
            body: dict[str, Any] = {
                "agent_id": target.agent_id,
                "agent_version": target.agent_version,
                "query": query,
                "conversation_id": decision_conversation_id,
                "stream": True,
                "inc_stream": True,
                "chat_option": dict(BUSINESS_AGENT_CHAT_OPTION),
            }
            if custom_querys is not None:
                body["custom_querys"] = custom_querys
        else:
            target = self._resolve_agent_chat_target(agent_id, http_client=http_client)
            path = f"/api/agent-factory/v1/app/{target.agent_key}/chat/completion"
            body: dict[str, Any] = {
                "agent_id": target.agent_id,
                "agent_version": target.agent_version,
                "query": query,
                "conversation_id": decision_conversation_id,
                "stream": False,
                "inc_stream": False,
            }
            if custom_querys is not None:
                body["custom_querys"] = custom_querys
        if not body.get("conversation_id"):
            body.pop("conversation_id", None)

        timeout = self._effective_http_timeout()
        request_log_payload = {
            "path": path,
            "body": body,
            "timeout": self._serialize_timeout(timeout),
        }
        logger.info(
            "kweaver sdk conversation request: %s",
            json.dumps(
                self._sanitize_for_log(
                    self._include_content_lengths_for_log(request_log_payload)
                ),
                ensure_ascii=False,
            ),
        )

        if stream:
            try:
                message = self._stream_message(
                    http_client,
                    path,
                    body,
                    timeout,
                    headers=self._stream_headers(target.agent_id, target.agent_version),
                    requested_conversation_id=decision_conversation_id,
                )
            except Exception as exc:
                logger.error(
                    "kweaver sdk conversation response: %s",
                    json.dumps(
                        self._sanitize_for_log(self._serialize_exception(exc)),
                        ensure_ascii=False,
                    ),
                )
                raise
            logger.info(
                "kweaver sdk conversation response: %s",
                json.dumps(
                    self._sanitize_for_log(
                        self._include_content_lengths_for_log(
                            self._serialize_message(message)
                        )
                    ),
                    ensure_ascii=False,
                ),
            )
            return message

        try:
            data = http_client.post(path, json=body, timeout=timeout)
        except Exception as exc:
            logger.error(
                "kweaver sdk conversation response: %s",
                json.dumps(
                    self._sanitize_for_log(self._serialize_exception(exc)),
                    ensure_ascii=False,
                ),
            )
            raise
        logger.info(
            "kweaver sdk conversation response: %s",
            json.dumps(
                self._sanitize_for_log(self._include_content_lengths_for_log(data)),
                ensure_ascii=False,
            ),
        )
        message = _parse_message(data)
        logger.info(
            "kweaver sdk conversation parsed response: %s",
            json.dumps(
                self._sanitize_for_log(
                    self._include_content_lengths_for_log(
                        self._serialize_message(message)
                    )
                ),
                ensure_ascii=False,
            ),
        )
        return message

    def _stream_message(
        self,
        http_client: Any,
        path: str,
        body: dict[str, Any],
        timeout: Any,
        *,
        headers: dict[str, str] | None = None,
        requested_conversation_id: str,
    ) -> Any:
        # 流式调用参考 kweaver-sdk 的 ConversationsResource._stream_message；
        # 这里消费原始 event 是为了保留 conversation_id、message_id 等映射所需字段。
        stream_post = getattr(http_client, "stream_post", None)
        if stream_post is None:
            message = self._client.conversations.send_message(
                requested_conversation_id,
                content=str(body.get("query", "")),
                agent_id=str(body.get("agent_id", "")),
                stream=True,
            )
            return self._collect_stream_chunks(
                message, requested_conversation_id=requested_conversation_id
            )

        chunks: list[dict[str, Any]] = []
        trace_path = self._open_stream_trace(
            requested_conversation_id=requested_conversation_id
        )
        for event in stream_post(path, json=body, timeout=timeout, headers=headers):
            if not isinstance(event, dict):
                continue
            chunks.append(event)
            self._log_stream_progress(chunks, event)
            trace_path = self._append_stream_trace_event(trace_path, len(chunks), event)
            if event.get("finished") is True:
                break
        message = self._build_message_from_stream_events(
            chunks, requested_conversation_id=requested_conversation_id
        )
        self._log_stream_trace_written(trace_path, len(chunks), message)
        return message

    def _resolve_agent_chat_target(
        self, agent_id: str, *, http_client: Any | None = None
    ) -> AgentChatTarget:
        # 对齐 KWeaver Web 页面：直接读取最新发布版本详情，避免普通 Agent 详情缺少版本字段。
        try:
            if http_client is not None and hasattr(http_client, "get"):
                agent = http_client.get(
                    f"/api/agent-factory/v3/agent-market/agent/{agent_id}/version/latest?is_visit=true",
                    headers=self._agent_info_headers(),
                )
                agent_version = self._extract_agent_version(agent)
            else:
                agent = self._client.agents.get(agent_id)
                agent_version = self._extract_agent_version(agent)
        except Exception:
            logger.exception(
                "kweaver agent chat target resolve failed: %s",
                json.dumps({"agent_id": agent_id}, ensure_ascii=False),
            )
            raise
        agent_key = self._extract_string(agent, "key") or agent_id
        target = AgentChatTarget(
            agent_id=agent_id, agent_key=agent_key, agent_version=agent_version
        )
        logger.info(
            "kweaver agent chat target resolved: %s",
            json.dumps(
                {
                    "agent_id": target.agent_id,
                    "agent_key": target.agent_key,
                    "agent_version": target.agent_version,
                },
                ensure_ascii=False,
            ),
        )
        return target

    def _stream_headers(
        self, agent_id: str | None = None, agent_version: str | None = None
    ) -> dict[str, str]:
        # 对齐 KWeaver Web 页面可由服务稳定设置的请求头；认证、业务域和 Cookie 不在这里固化。
        referer_agent_id = agent_id or ""
        referer_agent_version = agent_version or ""
        return {
            "Content-Type": "application/json; charset=utf-8",
            "accept": "text/event-stream",
            "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
            "origin": self._base_url,
            "referer": f"{self._base_url}/dip-hub/business-network/my-agents/usage?id={referer_agent_id}&version={referer_agent_version}&agentAppType=common&preRoute=%2F&filterParams=%257B%2522mode%2522%253A%2522myAgent%2522%257D",
            "responsetype": "text/event-stream",
            "user-agent": KWEAVER_WEB_USER_AGENT,
            "x-Language": "zh-CN",
        }

    @staticmethod
    def _normalize_base_url(base_url: str) -> str:
        # 请求头中的 origin/referer 只需要站点根地址，统一去掉尾部斜杠，避免配置差异影响比对。
        value = str(base_url or "").strip()
        if not value:
            raise ValueError("kweaver base url is required")
        return value.rstrip("/")

    @staticmethod
    def _agent_info_headers() -> dict[str, str]:
        # 对齐 CLI 获取 Agent 信息的请求头；认证和业务域仍由 SDK HTTP 层统一注入。
        return {
            "accept": "application/json, text/plain, */*",
            "x-language": "zh-CN",
            "x-requested-with": "XMLHttpRequest",
        }

    def _collect_stream_chunks(
        self, chunks: Any, *, requested_conversation_id: str
    ) -> Any:
        events: list[dict[str, Any]] = []
        trace_path = self._open_stream_trace(
            requested_conversation_id=requested_conversation_id
        )
        try:
            iterator = iter(chunks)
        except TypeError:
            # 部分测试替身或旧 SDK 可能在 stream=True 时仍返回普通 Message，直接兼容该结果。
            return chunks
        for chunk in iterator:
            event = {
                "delta": getattr(chunk, "delta", ""),
                "finished": getattr(chunk, "finished", False),
                "references": getattr(chunk, "references", []),
            }
            events.append(event)
            self._log_stream_progress(events, event)
            trace_path = self._append_stream_trace_event(trace_path, len(events), event)
            if event["finished"] is True:
                break
        message = self._build_message_from_stream_events(
            events, requested_conversation_id=requested_conversation_id
        )
        self._log_stream_trace_written(trace_path, len(events), message)
        return message

    def _log_stream_progress(
        self, events: list[dict[str, Any]], event: dict[str, Any]
    ) -> None:
        # DEBUG 只记录采样进度，不输出完整 chunk，避免长流式响应污染主日志。
        chunk_index = len(events)
        finished = event.get("finished") is True
        if (
            chunk_index != 1
            and chunk_index % self._stream_progress_interval != 0
            and not finished
        ):
            return
        progress_payload = {
            "chunk_index": chunk_index,
            "finished": finished,
            "conversation_id": self._extract_stream_field(events, "conversation_id"),
            "message_id": (
                self._extract_stream_field(events, "message_id")
                or self._extract_stream_field(events, "id")
            ),
            "content_length": self._log_content_length(
                self._extract_stream_content(event)
            ),
        }
        logger.debug(
            "kweaver sdk conversation stream progress: %s",
            json.dumps(self._sanitize_for_log(progress_payload), ensure_ascii=False),
        )

    def _open_stream_trace(self, *, requested_conversation_id: str) -> Path | None:
        if not self._stream_trace_enabled:
            return None
        try:
            trace_dir = Path(self._stream_trace_dir)
            trace_dir.mkdir(parents=True, exist_ok=True)
            # 长耗时流式调用可能很久没有 finished 片段，trace 文件必须先落地，便于运行中排障。
            conversation_id = self._safe_trace_name(
                requested_conversation_id or "active-conversation"
            )
            trace_path = trace_dir / f"{conversation_id}_{uuid4().hex[:12]}.jsonl"
            trace_path.touch(exist_ok=False)
            return trace_path
        except OSError as exc:
            logger.warning(
                "kweaver sdk conversation stream trace open failed: %s",
                json.dumps(
                    {"error_type": type(exc).__name__, "error_message": str(exc)},
                    ensure_ascii=False,
                ),
            )
            return None

    def _append_stream_trace_event(
        self, trace_path: Path | None, chunk_index: int, event: dict[str, Any]
    ) -> Path | None:
        if trace_path is None:
            return None
        try:
            trace_event = {
                "chunk_index": chunk_index,
                "event": self._sanitize_for_log(event),
            }
            with trace_path.open("a", encoding="utf-8") as stream:
                stream.write(json.dumps(trace_event, ensure_ascii=False) + "\n")
            return trace_path
        except OSError as exc:
            logger.warning(
                "kweaver sdk conversation stream trace write failed: %s",
                json.dumps(
                    {
                        "error_type": type(exc).__name__,
                        "error_message": str(exc),
                        "path": str(trace_path),
                    },
                    ensure_ascii=False,
                ),
            )
            return None

    @staticmethod
    def _log_stream_trace_written(
        trace_path: Path | None, stream_chunk_count: int, message: Any
    ) -> None:
        if trace_path is None:
            return
        logger.info(
            "kweaver sdk conversation stream trace written: %s",
            json.dumps(
                {
                    "conversation_id": getattr(message, "conversation_id", None),
                    "message_id": getattr(message, "id", None),
                    "stream_chunk_count": stream_chunk_count,
                    "path": str(trace_path),
                },
                ensure_ascii=False,
            ),
        )

    @staticmethod
    def _safe_trace_name(value: Any) -> str:
        text = str(value or "").strip() or "unknown"
        return "".join(
            char if char.isalnum() or char in {"-", "_"} else "_" for char in text
        )[:120]

    def _build_message_from_stream_events(
        self, events: list[dict[str, Any]], *, requested_conversation_id: str
    ) -> Any:
        if not events:
            raise ValueError("decision agent stream response is empty")
        final_event = next(
            (event for event in reversed(events) if event.get("finished") is True),
            events[-1],
        )
        content = self._collect_patch_stream_content(events)
        if not content:
            content = self._extract_stream_content(final_event)
        if not content:
            content = "".join(self._extract_stream_content(event) for event in events)
        conversation_id = (
            self._extract_stream_field(events, "conversation_id")
            or requested_conversation_id
        )
        if not conversation_id:
            raise ValueError("decision agent stream response missing conversation_id")
        message_id = self._extract_stream_field(
            events, "message_id"
        ) or self._extract_stream_field(events, "id")
        status = self._extract_stream_field(events, "status") or (
            "completed" if final_event.get("finished") is True else None
        )
        return SimpleNamespace(
            conversation_id=conversation_id,
            id=message_id,
            status=status,
            content=content,
            stream_chunk_count=len(events),
            raw_stream_events=events,
        )

    @classmethod
    def _extract_stream_content(cls, event: dict[str, Any]) -> str:
        key_path = event.get("key")
        if (
            isinstance(key_path, list)
            and key_path[-5:]
            == ["message", "content", "final_answer", "answer", "text"]
            and isinstance(event.get("content"), str)
        ):
            return event["content"]
        value = event.get("delta", event.get("answer", ""))
        if isinstance(value, str):
            return value
        if value is None:
            return ""
        return json.dumps(cls._sanitize_for_log(value), ensure_ascii=False)

    @classmethod
    def _collect_patch_stream_content(cls, events: list[dict[str, Any]]) -> str:
        parts: list[str] = []
        for event in events:
            key_path = event.get("key")
            if (
                isinstance(key_path, list)
                and key_path[-5:]
                == ["message", "content", "final_answer", "answer", "text"]
                and isinstance(event.get("content"), str)
            ):
                if event.get("action") == "upsert":
                    parts = [event["content"]]
                else:
                    parts.append(event["content"])
        return "".join(parts)

    @staticmethod
    def _extract_stream_field(
        events: list[dict[str, Any]], field_name: str
    ) -> str | None:
        for event in reversed(events):
            key_path = event.get("key")
            if isinstance(key_path, list) and key_path == [field_name]:
                value = event.get("content")
                if isinstance(value, str) and value.strip():
                    return value
            if isinstance(key_path, list) and key_path == ["message", field_name]:
                value = event.get("content")
                if isinstance(value, str) and value.strip():
                    return value
            value = event.get(field_name)
            if isinstance(value, str) and value.strip():
                return value
            message = event.get("message")
            if isinstance(message, dict):
                value = message.get(field_name)
                if isinstance(value, str) and value.strip():
                    return value
        return None

    def _effective_http_timeout(self) -> float | httpx.Timeout:
        # KWeaver SDK 的 HTTP resource 在 timeout=None 时会回落到底层客户端默认超时；
        # 因此这里显式传入 httpx.Timeout(None)，让“不限制超时”的配置真正落到 chat/completion 请求上。
        if self._chat_timeout is None:
            return httpx.Timeout(None)
        return self._chat_timeout

    @staticmethod
    def _extract_conversation_id(message: Any) -> str:
        value = getattr(message, "conversation_id", None)
        if not isinstance(value, str) or not value.strip():
            raise ValueError("decision agent response missing conversation_id")
        return value

    @staticmethod
    def _extract_string(message: Any, field_name: str) -> str | None:
        if isinstance(message, dict):
            value = message.get(field_name)
        else:
            value = getattr(message, field_name, None)
        return value if isinstance(value, str) and value.strip() else None

    @classmethod
    def _extract_agent_version(cls, agent: Any) -> str:
        # Agent 详情接口返回的是当前最新可用版本；会话调用必须沿用该版本，不能在代码里写死。
        for field_name in (
            "version",
            "agent_version",
            "current_version",
            "latest_version",
        ):
            value = cls._extract_string(agent, field_name)
            if value:
                return value
        raise ValueError("kweaver agent response missing version")

    @staticmethod
    def _normalize_content(content: Any) -> Any:
        if isinstance(content, bytes):
            return content.decode("utf-8", errors="replace")
        return content

    @classmethod
    def _parse_intent_payload(cls, content: Any) -> dict[str, Any]:
        # 路由 Agent 的稳定契约是 JSON；这里仅做轻量容错，避免 SDK 返回 bytes 或前后附带说明文本导致解析失败。
        normalized = cls._normalize_content(content)
        if isinstance(normalized, dict):
            return normalized
        if not isinstance(normalized, str) or not normalized.strip():
            raise ValueError("decision agent intent response is empty")
        text = normalized.strip()
        try:
            parsed = json.loads(text)
        except json.JSONDecodeError:
            start = text.find("{")
            end = text.rfind("}")
            if start < 0 or end <= start:
                raise ValueError("decision agent intent response is not json") from None
            parsed = json.loads(text[start : end + 1])
        if not isinstance(parsed, dict):
            raise ValueError("decision agent intent response must be a json object")
        return parsed

    @staticmethod
    def _extract_payload_string(payload: dict[str, Any], field_name: str) -> str | None:
        value = payload.get(field_name)
        return value if isinstance(value, str) and value.strip() else None

    @classmethod
    def _summarize_request_body_for_log(cls, body: dict[str, Any]) -> dict[str, Any]:
        # chat/completion 的 query/custom_querys 可能包含用户原文和凭据，只记录长度和协议字段。
        summarized = {
            key: value
            for key, value in body.items()
            if key not in {"query", "custom_querys"}
        }
        summarized["query_length"] = len(str(body.get("query", "")))
        summarized["custom_querys_length"] = len(str(body.get("custom_querys", "")))
        return summarized

    @classmethod
    def _summarize_query_for_log(cls, query: Any) -> Any:
        # 业务上下文中的 content 不进通用 SDK 日志，保留长度即可定位消息规模。
        if isinstance(query, dict):
            summarized: dict[str, Any] = {}
            for key, value in query.items():
                if key == "content":
                    summarized["content_length"] = len(str(value))
                else:
                    summarized[key] = cls._summarize_query_for_log(value)
            return summarized
        if isinstance(query, list):
            return [cls._summarize_query_for_log(item) for item in query]
        return query

    @classmethod
    def _summarize_message_for_log(cls, message: Any) -> dict[str, Any]:
        return cls._summarize_mapping_for_log(cls._serialize_message(message))

    @classmethod
    def _summarize_mapping_for_log(cls, value: Any) -> Any:
        # 平台返回内容可能是自然语言正文或下游敏感文本，日志只保留 content_length。
        if isinstance(value, dict):
            summarized: dict[str, Any] = {}
            for key, item in value.items():
                if key == "content":
                    summarized["content_length"] = cls._log_content_length(item)
                else:
                    summarized[key] = cls._summarize_mapping_for_log(item)
            return cls._sanitize_for_log(summarized)
        if isinstance(value, list):
            return [cls._summarize_mapping_for_log(item) for item in value]
        return cls._sanitize_for_log(value)

    @classmethod
    def _summarize_response_for_log(
        cls, response_payload: dict[str, Any]
    ) -> dict[str, Any]:
        return cls._summarize_mapping_for_log(response_payload)

    @staticmethod
    def _log_content_length(value: Any) -> int:
        if value is None:
            return 0
        if isinstance(value, bytes):
            return len(value)
        if isinstance(value, str):
            return len(value)
        return len(json.dumps(value, ensure_ascii=False))

    @classmethod
    def _include_content_lengths_for_log(cls, value: Any) -> Any:
        # 日志保留正文，同时补充长度，方便排查内容是否被截断或错误转换。
        if isinstance(value, dict):
            enriched: dict[str, Any] = {}
            for key, item in value.items():
                enriched[key] = cls._include_content_lengths_for_log(item)
                if key == "content":
                    enriched["content_length"] = cls._log_content_length(item)
            return enriched
        if isinstance(value, list):
            return [cls._include_content_lengths_for_log(item) for item in value]
        return value

    @classmethod
    def _sanitize_for_log(cls, value: Any) -> Any:
        # 日志要保留排障需要的请求和响应内容，但连接串中的凭据必须脱敏。
        if isinstance(value, dict):
            sanitized: dict[Any, Any] = {}
            for key, item in value.items():
                if isinstance(key, str) and key.lower() in {
                    "ak",
                    "sk",
                    "token",
                    "password",
                    "secret",
                }:
                    sanitized[key] = "***"
                else:
                    sanitized[key] = cls._sanitize_for_log(item)
            return sanitized
        if isinstance(value, list):
            return [cls._sanitize_for_log(item) for item in value]
        if isinstance(value, tuple):
            return [cls._sanitize_for_log(item) for item in value]
        if isinstance(value, str):
            parsed = cls._try_parse_json_text(value)
            if parsed is not None:
                return json.dumps(cls._sanitize_for_log(parsed), ensure_ascii=False)
            return cls._redact_inline_credentials(cls._redact_url_credentials(value))
        return value

    @staticmethod
    def _try_parse_json_text(value: Any) -> Any | None:
        if not isinstance(value, str) or not value.strip().startswith(("{", "[")):
            return None
        try:
            return json.loads(value)
        except json.JSONDecodeError:
            return None

    @staticmethod
    def _redact_url_credentials(value: str) -> str:
        # 只处理形如 scheme://user:password@host 的连接串，普通文本保持原样。
        import re

        return re.sub(
            r"([A-Za-z][A-Za-z0-9+.-]*://)([^/@:\s]+):([^/@\s]+)@", r"\1\2:***@", value
        )

    @staticmethod
    def _redact_inline_credentials(value: str) -> str:
        # 用户原文有时会携带 ak/sk/tenant-id，进入日志前统一替换为占位符。
        import re

        return re.sub(
            r"(?i)\b(ak|sk|tenant-id|tenant_id|token|password|secret)\s*=\s*[^，,\s]+",
            lambda match: f"{match.group(1)}=***",
            value,
        )

    @staticmethod
    def _serialize_timeout(timeout: Any) -> Any:
        if isinstance(timeout, httpx.Timeout):
            return {
                "connect": timeout.connect,
                "read": timeout.read,
                "write": timeout.write,
                "pool": timeout.pool,
            }
        return timeout

    @staticmethod
    def _serialize_exception(exc: Exception) -> dict[str, Any]:
        return {
            "error_type": type(exc).__name__,
            "error_message": str(exc),
            "status_code": getattr(exc, "status_code", None),
            "error_code": getattr(exc, "error_code", None),
            "trace_id": getattr(exc, "trace_id", None),
        }

    @classmethod
    def _serialize_message(cls, message: Any) -> dict[str, object]:
        # raw 字段用于日志排查，尽量保留平台原始关键字段，方便比对 SDK 实际返回。
        return {
            "conversation_id": cls._extract_string(message, "conversation_id"),
            "id": cls._extract_string(message, "id"),
            "status": cls._extract_string(message, "status"),
            "content": cls._normalize_content(getattr(message, "content", None)),
            "stream_chunk_count": getattr(message, "stream_chunk_count", None),
        }
