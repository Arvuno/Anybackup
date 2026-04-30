from types import SimpleNamespace
import json
import logging

from core_agent_service.infrastructure.kweaver.adapter import KWeaverRelayAdapter
from core_agent_service.infrastructure.kweaver.client import create_client


TEST_KWEAVER_BASE_URL = "https://115.190.186.186"


def make_adapter(**kwargs):
    kwargs.setdefault("base_url", TEST_KWEAVER_BASE_URL)
    return KWeaverRelayAdapter(**kwargs)


class FakeTokenAuth:
    def __init__(self, token: str):
        self.token = token


class FakeConfigAuth:
    def __init__(self, platform: str | None = None):
        self.platform = platform


class FakeHttpSigninAuth:
    def __init__(
        self, base_url: str, *, username: str, password: str, tls_insecure: bool = False
    ):
        self.base_url = base_url
        self.username = username
        self.password = password
        self.tls_insecure = tls_insecure


class FakePasswordAuth:
    def __init__(self, base_url: str, username: str, password: str):
        self.base_url = base_url
        self.username = username
        self.password = password


class FakeKWeaverClient:
    def __init__(self, **kwargs):
        self.kwargs = kwargs
        self.agents = SimpleNamespace(get=self._get_agent)
        self.conversations = SimpleNamespace(send_message=self._send_message)
        self.sent_messages: list[dict] = []
        self.agent_gets: list[str] = []

    def _get_agent(self, agent_id: str):
        self.agent_gets.append(agent_id)
        return SimpleNamespace(
            id=agent_id,
            name="Decision Agent",
            status="published",
        )

    def _send_message(
        self,
        conversation_id: str,
        content: str,
        *,
        agent_id: str,
        agent_version: str = "latest",
        stream: bool = False,
        debug: bool = False,
        history=None,
    ):
        self.sent_messages.append(
            {
                "conversation_id": conversation_id,
                "content": content,
                "agent_id": agent_id,
                "agent_version": agent_version,
                "stream": stream,
            }
        )
        return SimpleNamespace(
            conversation_id=conversation_id or "decision_conv_created",
            id="msg_001",
            status="completed",
            content='{"status":"ok"}',
        )


class FakeHttp:
    def __init__(self):
        self.gets: list[dict] = []
        self.posts: list[dict] = []
        self.stream_posts: list[dict] = []
        self.stream_chunks = [
            {
                "delta": "业务处理中",
                "finished": False,
                "conversation_id": "decision_conv_created",
                "message_id": "msg_stream_001",
            },
            {
                "delta": "业务处理完成",
                "finished": True,
                "conversation_id": "decision_conv_created",
                "message_id": "msg_stream_002",
            },
            {
                "delta": "不应继续读取",
                "finished": False,
                "conversation_id": "decision_conv_created",
                "message_id": "msg_stream_003",
            },
        ]

    def get(self, path: str, *, params=None, headers=None):
        self.gets.append({"path": path, "params": params, "headers": headers})
        if "/agent-market/agent/" in path:
            agent_id = path.split("/agent-market/agent/", 1)[1].split("/version/", 1)[0]
            version = path.split("/version/", 1)[1].split("?", 1)[0]
            if version == "latest":
                version = "v9"
        elif "/agent/" in path:
            agent_id = path.rsplit("/agent/", 1)[1]
            version = None
        else:
            agent_id = "agent_001"
            version = "v9"
        response = {
            "id": agent_id,
            "key": f"{agent_id}_key",
        }
        if version is not None:
            response["version"] = version
            response["latest_version"] = version
        return response

    def post(self, path: str, *, json=None, timeout=None, headers=None):
        self.posts.append(
            {"path": path, "json": json, "timeout": timeout, "headers": headers}
        )
        return {
            "conversation_id": "decision_conv_created",
            "message_id": "msg_001",
            "content": "ok",
        }

    def stream_post(self, path: str, *, json=None, timeout=None, headers=None):
        self.stream_posts.append(
            {"path": path, "json": json, "timeout": timeout, "headers": headers}
        )
        yield from self.stream_chunks


class FakeConversationsWithHttp:
    def __init__(self):
        self._http = FakeHttp()


class FakeClientWithHttp:
    def __init__(self):
        self.agents = SimpleNamespace(
            get=lambda agent_id: SimpleNamespace(
                id=agent_id,
                key=f"{agent_id}_key",
                version="v0",
                name="Decision Agent",
                status="published",
            )
        )
        self.conversations = FakeConversationsWithHttp()


class FakeFailingHttp:
    def stream_post(self, path: str, *, json=None, timeout=None, headers=None):
        raise RuntimeError("Gateway Time-out")

    def post(self, path: str, *, json=None, timeout=None, headers=None):
        raise RuntimeError("Gateway Time-out")


class FakeClientWithFailingHttp:
    def __init__(self):
        self.agents = SimpleNamespace(
            get=lambda agent_id: SimpleNamespace(
                id=agent_id, version="v9", name="Decision Agent", status="published"
            )
        )
        self.conversations = SimpleNamespace(_http=FakeFailingHttp())


class FakePatchStreamHttp(FakeHttp):
    def __init__(self):
        super().__init__()
        self.stream_chunks = [
            {
                "seq_id": 0,
                "key": ["message_id"],
                "content": "msg_patch_001",
                "action": "upsert",
            },
            {
                "seq_id": 1,
                "key": ["conversation_id"],
                "content": "decision_conv_patch",
                "action": "upsert",
            },
            {
                "seq_id": 2,
                "key": ["message", "content", "final_answer", "answer", "text"],
                "content": "恢复",
                "action": "append",
            },
            {
                "seq_id": 3,
                "key": ["message", "content", "final_answer", "answer", "text"],
                "content": "完成",
                "action": "append",
            },
        ]


class FakeClientWithPatchStreamHttp:
    def __init__(self):
        self.agents = SimpleNamespace(
            get=lambda agent_id: SimpleNamespace(
                id=agent_id, version="v9", name="Decision Agent", status="published"
            )
        )
        self.conversations = SimpleNamespace(_http=FakePatchStreamHttp())


class FakeInspectingTraceHttp(FakeHttp):
    def __init__(self, trace_dir):
        super().__init__()
        self.trace_dir = trace_dir

    def stream_post(self, path: str, *, json=None, timeout=None, headers=None):
        self.stream_posts.append(
            {"path": path, "json": json, "timeout": timeout, "headers": headers}
        )
        yield {
            "delta": "第一片",
            "finished": False,
            "conversation_id": "decision_conv_created",
            "message_id": "msg_stream_001",
        }
        trace_files = list(self.trace_dir.glob("*.jsonl"))
        assert len(trace_files) == 1
        assert trace_files[0].read_text(encoding="utf-8").strip()
        yield {
            "delta": "完成",
            "finished": True,
            "conversation_id": "decision_conv_created",
            "message_id": "msg_stream_001",
        }


class FakeClientWithInspectingTraceHttp:
    def __init__(self, trace_dir):
        self.agents = SimpleNamespace(
            get=lambda agent_id: SimpleNamespace(
                id=agent_id, version="v9", name="Decision Agent", status="published"
            )
        )
        self.conversations = SimpleNamespace(_http=FakeInspectingTraceHttp(trace_dir))


def test_create_client_prefers_username_password_auth(monkeypatch):
    monkeypatch.setattr(
        "core_agent_service.infrastructure.kweaver.client._load_kweaver_types",
        lambda: (FakeKWeaverClient, FakeTokenAuth, FakeConfigAuth, FakeHttpSigninAuth),
    )

    client = create_client(
        base_url="https://kweaver.example.com",
        token="secret-token",
        username="demo_user",
        password="demo_pass",
        business_domain="bd_public",
        timeout=30.0,
        tls_insecure=True,
    )

    assert isinstance(client.kwargs["auth"], FakeHttpSigninAuth)
    assert client.kwargs["auth"].base_url == "https://kweaver.example.com"
    assert client.kwargs["auth"].username == "demo_user"
    assert client.kwargs["auth"].password == "demo_pass"
    assert client.kwargs["auth"].tls_insecure is True


def test_create_client_prefers_token_auth_when_username_password_missing(monkeypatch):
    monkeypatch.setattr(
        "core_agent_service.infrastructure.kweaver.client._load_kweaver_types",
        lambda: (FakeKWeaverClient, FakeTokenAuth, FakeConfigAuth, FakeHttpSigninAuth),
    )

    client = create_client(
        base_url="https://kweaver.example.com",
        token="secret-token",
        business_domain="bd_public",
        timeout=30.0,
        tls_insecure=True,
    )

    assert isinstance(client.kwargs["auth"], FakeTokenAuth)
    assert client.kwargs["auth"].token == "secret-token"


def test_create_client_falls_back_to_config_auth(monkeypatch):
    monkeypatch.setattr(
        "core_agent_service.infrastructure.kweaver.client._load_kweaver_types",
        lambda: (FakeKWeaverClient, FakeTokenAuth, FakeConfigAuth, FakeHttpSigninAuth),
    )

    client = create_client(
        base_url="https://kweaver.example.com",
        token=None,
        business_domain="bd_public",
        timeout=30.0,
        tls_insecure=True,
    )

    assert isinstance(client.kwargs["auth"], FakeConfigAuth)
    assert client.kwargs["auth"].platform == "https://kweaver.example.com"


def test_create_client_supports_legacy_password_auth_shape(monkeypatch):
    monkeypatch.setattr(
        "core_agent_service.infrastructure.kweaver.client._load_kweaver_types",
        lambda: (FakeKWeaverClient, FakeTokenAuth, FakeConfigAuth, FakePasswordAuth),
    )

    client = create_client(
        base_url="https://kweaver.example.com",
        token=None,
        username="demo_user",
        password="demo_pass",
        business_domain="bd_public",
        timeout=30.0,
        tls_insecure=True,
    )

    assert isinstance(client.kwargs["auth"], FakePasswordAuth)
    assert client.kwargs["auth"].base_url == "https://kweaver.example.com"
    assert client.kwargs["auth"].username == "demo_user"
    assert client.kwargs["auth"].password == "demo_pass"


def test_kweaver_adapter_sends_content_via_sdk_fallback():
    client = FakeKWeaverClient()
    adapter = make_adapter(client=client, agent_id="agent_001")

    decision_conversation_id = adapter.ensure_conversation(None)
    response = adapter.relay_message(
        decision_conversation_id,
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={
            "rabbitmq_url": "amqp://guest:guest@localhost:5672/",
        },
    )
    assert decision_conversation_id == ""
    assert client.sent_messages[0]["conversation_id"] == ""
    assert client.sent_messages[0]["agent_id"] == "agent_001"
    assert client.sent_messages[0]["stream"] is True
    assert client.sent_messages[0]["content"] == "hello kweaver"
    assert response["decision_conversation_id"] == "decision_conv_created"
    assert response["raw"]["status"] == "completed"


def test_kweaver_adapter_resolves_intent_with_raw_content_via_configured_agent():
    client = FakeKWeaverClient()
    adapter = make_adapter(client=client, agent_id="router_agent")

    def send_intent(
        conversation_id: str,
        content: str,
        *,
        agent_id: str,
        agent_version: str = "latest",
        stream: bool = False,
        debug: bool = False,
        history=None,
    ):
        client.sent_messages.append(
            {
                "conversation_id": conversation_id,
                "content": content,
                "agent_id": agent_id,
                "stream": stream,
            }
        )
        return SimpleNamespace(
            conversation_id="router_conv",
            id="msg_router",
            status="completed",
            content='{"agent_id":"target_agent","agent_name":"备份方案推荐Agent","is_intent_success":true}',
        )

    client.conversations.send_message = send_intent
    resolution = adapter.resolve_intent(
        "为 mysql 推荐备份方案",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        mq={},
        foundation=None,
    )

    assert client.sent_messages == [
        {
            "conversation_id": "",
            "content": "为 mysql 推荐备份方案",
            "agent_id": "router_agent",
            "stream": False,
        }
    ]
    assert resolution.is_intent_success is True
    assert resolution.agent_id == "target_agent"
    assert resolution.agent_name == "备份方案推荐Agent"


def test_kweaver_adapter_resolves_intent_through_agent_key_for_dip_history():
    client = FakeClientWithHttp()

    def post_intent(path: str, *, json=None, timeout=None, headers=None):
        client.conversations._http.posts.append(
            {"path": path, "json": json, "timeout": timeout, "headers": headers}
        )
        return {
            "conversation_id": "router_conv_created",
            "message_id": "msg_router",
            "content": '{"agent_id":"target_agent","agent_name":"备份方案推荐Agent","is_intent_success":true}',
        }

    client.conversations._http.post = post_intent
    adapter = make_adapter(client=client, agent_id="router_agent")

    resolution = adapter.resolve_intent(
        "为 mysql 推荐备份方案",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        mq={},
        foundation=None,
    )

    assert client.conversations._http.gets == [
        {
            "path": "/api/agent-factory/v3/agent-market/agent/router_agent/version/latest?is_visit=true",
            "params": None,
            "headers": {
                "accept": "application/json, text/plain, */*",
                "x-language": "zh-CN",
                "x-requested-with": "XMLHttpRequest",
            },
        }
    ]
    post = client.conversations._http.posts[0]
    assert post["path"] == "/api/agent-factory/v1/app/router_agent_key/chat/completion"
    assert post["json"]["agent_id"] == "router_agent"
    assert post["json"]["agent_version"] == "v9"
    assert "executor_version" not in post["json"]
    assert post["json"]["query"] == "为 mysql 推荐备份方案"
    assert post["json"]["stream"] is False
    assert post["json"]["inc_stream"] is False
    assert "conversation_id" not in post["json"]
    assert resolution.agent_id == "target_agent"
    assert resolution.agent_name == "备份方案推荐Agent"
    assert resolution.is_intent_success is True


def test_kweaver_adapter_logs_parsed_intent_json_payload(caplog):
    client = FakeKWeaverClient()
    adapter = make_adapter(client=client, agent_id="router_agent")

    def send_intent(
        conversation_id: str,
        content: str,
        *,
        agent_id: str,
        agent_version: str = "latest",
        stream: bool = False,
        debug: bool = False,
        history=None,
    ):
        return SimpleNamespace(
            conversation_id="router_conv",
            id="msg_router",
            status="completed",
            content='{"agent_id":"target_agent","agent_name":"订单恢复Agent","is_intent_success":true}',
        )

    client.conversations.send_message = send_intent

    with caplog.at_level(
        logging.INFO, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        adapter.resolve_intent(
            "订单库宕机了需要恢复",
            conversation_id="conv_001",
            turn_id="turn_001",
            source_message_id="msg_001",
            mq={},
            foundation=None,
        )

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "decision agent intent routing payload" in log_text
    assert '"agent_id": "target_agent"' in log_text
    assert '"agent_name": "订单恢复Agent"' in log_text
    assert '"is_intent_success": true' in log_text


def test_kweaver_adapter_sends_content_to_target_agent():
    client = FakeKWeaverClient()
    adapter = make_adapter(client=client, agent_id="router_agent")

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    assert client.sent_messages[0]["agent_id"] == "target_agent"
    assert client.sent_messages[0]["stream"] is True


def test_kweaver_adapter_streams_content_and_custom_querys_to_target_agent_until_finished():
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="router_agent")

    response = adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    assert client.conversations._http.posts == []
    assert client.conversations._http.gets == [
        {
            "path": "/api/agent-factory/v3/agent-market/agent/target_agent/version/latest?is_visit=true",
            "params": None,
            "headers": {
                "accept": "application/json, text/plain, */*",
                "x-language": "zh-CN",
                "x-requested-with": "XMLHttpRequest",
            },
        }
    ]
    stream_post = client.conversations._http.stream_posts[0]
    assert (
        stream_post["path"]
        == "/api/agent-factory/v1/app/target_agent_key/chat/completion"
    )
    assert stream_post["json"]["agent_id"] == "target_agent"
    assert "agent_key" not in stream_post["json"]
    assert stream_post["json"]["agent_version"] == "v9"
    assert "executor_version" not in stream_post["json"]
    assert stream_post["json"]["stream"] is True
    assert stream_post["json"]["inc_stream"] is True
    assert stream_post["json"]["chat_option"] == {
        "is_need_history": True,
        "is_need_doc_retrival_post_process": True,
        "is_need_progress": True,
        "enable_dependency_cache": True,
    }
    assert stream_post["json"]["query"] == "hello kweaver"
    assert stream_post["json"]["custom_querys"] == {
        "conversation_id": "conv_001",
        "turn_id": "turn_001",
        "source_message_id": "msg_001",
        "content": "hello kweaver",
        "mq": {},
    }
    assert response["decision_conversation_id"] == "decision_conv_created"
    assert response["reply_message_id"] == "msg_stream_002"
    assert response["content"] == "业务处理完成"
    assert response["raw"]["stream_chunk_count"] == 2


def test_kweaver_adapter_adds_foundation_runtime_info_to_custom_querys():
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="router_agent")

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        foundation={
            "endpoint": "https://115.190.150.254:9600",
            "ak": "foundation-ak",
            "sk": "foundation-sk",
        },
        agent_id="target_agent",
    )

    custom_querys = client.conversations._http.stream_posts[0]["json"]["custom_querys"]
    assert custom_querys["foundation"] == {
        "endpoint": "https://115.190.150.254:9600",
        "ak": "foundation-ak",
        "sk": "foundation-sk",
    }


def test_kweaver_adapter_sends_cli_compatible_headers_for_business_stream():
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="router_agent")

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    stream_post = client.conversations._http.stream_posts[0]
    assert stream_post["headers"] == {
        "Content-Type": "application/json; charset=utf-8",
        "accept": "text/event-stream",
        "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
        "origin": "https://115.190.186.186",
        "referer": "https://115.190.186.186/dip-hub/business-network/my-agents/usage?id=target_agent&version=v9&agentAppType=common&preRoute=%2F&filterParams=%257B%2522mode%2522%253A%2522myAgent%2522%257D",
        "responsetype": "text/event-stream",
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0",
        "x-Language": "zh-CN",
    }


def test_kweaver_adapter_uses_configured_base_url_for_business_stream_headers():
    client = FakeClientWithHttp()
    adapter = make_adapter(
        client=client,
        agent_id="router_agent",
        base_url="https://kweaver.example.com/",
    )

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    headers = client.conversations._http.stream_posts[0]["headers"]
    assert headers["origin"] == "https://kweaver.example.com"
    assert (
        headers["referer"]
        == "https://kweaver.example.com/dip-hub/business-network/my-agents/usage?id=target_agent&version=v9&agentAppType=common&preRoute=%2F&filterParams=%257B%2522mode%2522%253A%2522myAgent%2522%257D"
    )


def test_kweaver_adapter_does_not_log_stream_chunks_at_info(caplog):
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="router_agent")

    with caplog.at_level(
        logging.INFO, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        adapter.relay_message(
            "",
            conversation_id="conv_001",
            turn_id="turn_001",
            source_message_id="msg_001",
            content="hello kweaver",
            mq={},
            agent_id="target_agent",
        )

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "kweaver sdk conversation stream chunk" not in log_text
    assert "kweaver sdk conversation response" in log_text


def test_kweaver_adapter_logs_sampled_stream_progress_at_debug(caplog):
    client = FakeClientWithHttp()
    adapter = make_adapter(
        client=client, agent_id="router_agent", stream_progress_interval=2
    )

    with caplog.at_level(
        logging.DEBUG, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        adapter.relay_message(
            "",
            conversation_id="conv_001",
            turn_id="turn_001",
            source_message_id="msg_001",
            content="hello kweaver",
            mq={},
            agent_id="target_agent",
        )

    progress_logs = [
        record.getMessage()
        for record in caplog.records
        if "kweaver sdk conversation stream progress" in record.getMessage()
    ]
    assert len(progress_logs) == 2
    assert '"chunk_index": 1' in progress_logs[0]
    assert '"content_length":' in progress_logs[0]
    assert '"content":' not in "\n".join(progress_logs)
    assert '"chunk_index": 2' in progress_logs[1]
    assert '"finished": true' in progress_logs[1]


def test_kweaver_adapter_writes_sanitized_stream_trace_file(tmp_path):
    client = FakeClientWithHttp()
    client.conversations._http.stream_chunks = [
        {
            "delta": "endpoint=https://example.test，ak=abc，sk=def",
            "finished": True,
            "conversation_id": "decision_conv_created",
            "message_id": "msg_stream_001",
        }
    ]
    adapter = make_adapter(
        client=client,
        agent_id="router_agent",
        stream_trace_enabled=True,
        stream_trace_dir=str(tmp_path),
    )

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    trace_files = list(tmp_path.glob("*.jsonl"))
    assert len(trace_files) == 1
    trace_text = trace_files[0].read_text(encoding="utf-8")
    assert "ak=***" in trace_text
    assert "sk=***" in trace_text
    assert "abc" not in trace_text
    assert "def" not in trace_text


def test_kweaver_adapter_writes_stream_trace_while_stream_is_active(tmp_path):
    client = FakeClientWithInspectingTraceHttp(tmp_path)
    adapter = make_adapter(
        client=client,
        agent_id="router_agent",
        stream_trace_enabled=True,
        stream_trace_dir=str(tmp_path),
    )

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    trace_files = list(tmp_path.glob("*.jsonl"))
    assert len(trace_files) == 1
    trace_lines = trace_files[0].read_text(encoding="utf-8").splitlines()
    assert len(trace_lines) == 2


def test_kweaver_adapter_extracts_patch_style_stream_conversation_fields():
    client = FakeClientWithPatchStreamHttp()
    adapter = make_adapter(client=client, agent_id="router_agent")

    response = adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello kweaver",
        mq={},
        agent_id="target_agent",
    )

    assert response["decision_conversation_id"] == "decision_conv_patch"
    assert response["reply_message_id"] == "msg_patch_001"
    assert response["content"] == "恢复完成"
    assert response["raw"]["stream_chunk_count"] == 4


def test_kweaver_adapter_uses_configured_chat_timeout_with_sdk_http_resource():
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="agent_001", chat_timeout=600.0)

    response = adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello",
        mq={},
    )

    stream_post = client.conversations._http.stream_posts[0]
    assert (
        stream_post["path"] == "/api/agent-factory/v1/app/agent_001_key/chat/completion"
    )
    assert stream_post["timeout"] == 600.0
    assert stream_post["json"]["query"] == "hello"
    assert stream_post["json"]["custom_querys"] == {
        "conversation_id": "conv_001",
        "turn_id": "turn_001",
        "source_message_id": "msg_001",
        "content": "hello",
        "mq": {},
    }
    assert "conversation_id" not in stream_post["json"]
    assert response["decision_conversation_id"] == "decision_conv_created"


def test_kweaver_adapter_logs_sdk_conversation_request_and_response(caplog):
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="agent_001")

    with caplog.at_level(
        logging.INFO, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        adapter.relay_message(
            "",
            conversation_id="conv_001",
            turn_id="turn_001",
            source_message_id="msg_001",
            content="hello",
            mq={"rabbitmq_url": "amqp://user:secret@rabbitmq.middleware:5672/"},
        )

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "kweaver sdk conversation request" in log_text
    assert "kweaver sdk conversation response" in log_text
    assert '"query_length":' in log_text
    assert '"content_length":' in log_text
    assert '"content": "hello"' in log_text
    assert '"content": "业务处理完成"' in log_text
    assert '"conversation_id": "conv_001"' in log_text
    assert '"agent_id": "agent_001"' in log_text
    assert '"turn_id": "turn_001"' in log_text
    assert '"source_message_id": "msg_001"' in log_text
    assert "rabbitmq.middleware:5672" in log_text
    assert "secret" not in log_text


def test_kweaver_adapter_logs_sdk_conversation_error_response(caplog):
    client = FakeClientWithFailingHttp()
    adapter = make_adapter(client=client, agent_id="agent_001")

    with caplog.at_level(
        logging.INFO, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        try:
            adapter.relay_message(
                "",
                conversation_id="conv_001",
                turn_id="turn_001",
                source_message_id="msg_001",
                content="hello",
                mq={},
            )
        except RuntimeError:
            pass
        else:
            raise AssertionError("relay_message should raise SDK error")

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "kweaver sdk conversation request" in log_text
    assert "kweaver sdk conversation response" in log_text
    assert "Gateway Time-out" in log_text


def test_kweaver_adapter_redacts_foundation_credentials_in_logs(caplog):
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="agent_001")

    with caplog.at_level(
        logging.INFO, logger="core_agent_service.infrastructure.kweaver.adapter"
    ):
        adapter.relay_message(
            "",
            conversation_id="conv_001",
            turn_id="turn_001",
            source_message_id="msg_001",
            content="hello",
            mq={},
            foundation={
                "endpoint": "https://115.190.150.254:9600",
                "ak": "foundation-ak-secret",
                "sk": "foundation-sk-secret",
            },
        )

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "foundation-ak-secret" not in log_text
    assert "foundation-sk-secret" not in log_text
    assert '"ak": "***"' in log_text
    assert '"sk": "***"' in log_text


def test_kweaver_adapter_disables_sdk_default_chat_timeout_when_unconfigured():
    client = FakeClientWithHttp()
    adapter = make_adapter(client=client, agent_id="agent_001")

    adapter.relay_message(
        "",
        conversation_id="conv_001",
        turn_id="turn_001",
        source_message_id="msg_001",
        content="hello",
        mq={},
    )

    timeout = client.conversations._http.stream_posts[0]["timeout"]
    assert timeout.connect is None
    assert timeout.read is None
    assert timeout.write is None
    assert timeout.pool is None


def test_kweaver_adapter_probes_configured_agent_accessibility():
    client = FakeKWeaverClient()
    adapter = make_adapter(client=client, agent_id="agent_001")

    response = adapter.probe_connectivity()

    assert client.agent_gets == ["agent_001"]
    assert client.sent_messages == []
    assert response["agent_id"] == "agent_001"
    assert response["agent_name"] == "Decision Agent"
    assert response["agent_status"] == "published"


def test_kweaver_adapter_probe_fails_when_agent_is_not_accessible():
    client = FakeKWeaverClient()

    def raise_not_found(agent_id: str):
        raise RuntimeError(f"agent not found: {agent_id}")

    client.agents.get = raise_not_found
    adapter = make_adapter(client=client, agent_id="missing_agent")

    try:
        adapter.probe_connectivity()
    except RuntimeError as exc:
        assert str(exc) == "decision agent is not accessible: agent_id=missing_agent"
    else:
        raise AssertionError(
            "probe_connectivity should fail when agent is not accessible"
        )
    assert client.sent_messages == []
