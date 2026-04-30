from pydantic import field_validator, model_validator
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    # 所有运行参数统一从环境变量或 .env 读取，避免在业务代码里散落配置判断。
    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

    database_url: str
    rabbitmq_url: str = "amqp://guest:guest@localhost:5672/"
    rabbitmq_exchange: str = "conversation.message.events"
    rabbitmq_exchange_type: str = "topic"
    rabbitmq_queue: str = "core_agent.message.events"
    rabbitmq_status_exchange: str = "core_agent.run_status.events"
    rabbitmq_status_routing_key: str | None = None
    rabbitmq_consumer_count: int = 1
    kweaver_base_url: str
    kweaver_token: str | None = None
    kweaver_username: str | None = None
    kweaver_password: str | None = None
    kweaver_decision_agent_id: str
    kweaver_business_domain: str = "bd_public"
    kweaver_timeout: float = 30.0
    kweaver_chat_timeout: float | None = None
    kweaver_stream_progress_interval: int = 100
    kweaver_stream_trace_enabled: bool = True
    kweaver_stream_trace_dir: str = "/tmp/core-agent-service-trace-log"
    kweaver_tls_insecure: bool = True
    kweaver_probe_on_startup: bool = True
    foundation_endpoint: str | None = None
    foundation_ak: str | None = None
    foundation_sk: str | None = None
    core_agent_log_level: str = "DEBUG"
    core_agent_log_to_stdout: bool = False
    core_agent_log_file: str = "logs/core_agent_service.log"

    @field_validator("database_url")
    @classmethod
    def validate_database_url(cls, value: str) -> str:
        # 当前服务数据库访问只支持同步 SQLAlchemy 链路，避免把异步驱动方言带入运行时。
        if value.startswith("postgresql+asyncpg://"):
            raise ValueError("DATABASE_URL must use postgresql+psycopg:// for PostgreSQL")
        if value.startswith("postgresql://") or value.startswith("postgres://"):
            raise ValueError("DATABASE_URL must use postgresql+psycopg:// for PostgreSQL")
        return value

    @field_validator("kweaver_chat_timeout", mode="before")
    @classmethod
    def normalize_kweaver_chat_timeout(cls, value):
        # 不配置或显式留空时表示不限制 Decision Agent 单次调用等待时间。
        if value is None:
            return None
        if isinstance(value, str) and value.strip().lower() in {"", "none", "null", "unlimited"}:
            return None
        return value

    @field_validator("rabbitmq_status_routing_key", mode="before")
    @classmethod
    def normalize_rabbitmq_status_routing_key(cls, value):
        # 默认按会话服务绑定的三类状态 routing key 投递；显式留空也表示不覆盖。
        if value is None:
            return None
        if isinstance(value, str) and value.strip().lower() in {"", "none", "null"}:
            return None
        return value

    @field_validator("foundation_endpoint", "foundation_ak", "foundation_sk", mode="before")
    @classmethod
    def normalize_optional_foundation_value(cls, value):
        # Foundation 参数由安装期注入；空字符串按未配置处理，避免把空凭据传给业务 Agent。
        if value is None:
            return None
        if isinstance(value, str) and value.strip() == "":
            return None
        return value

    @model_validator(mode="after")
    def validate_foundation_credentials(self):
        # Foundation 连接信息必须三项同时存在，否则业务 Agent 无法完整调用 foundation cli。
        values = [self.foundation_endpoint, self.foundation_ak, self.foundation_sk]
        if any(values) and not all(values):
            raise ValueError("FOUNDATION_ENDPOINT, FOUNDATION_AK and FOUNDATION_SK must be configured together")
        return self

    @field_validator("rabbitmq_consumer_count")
    @classmethod
    def validate_rabbitmq_consumer_count(cls, value: int) -> int:
        # 每个 consumer worker 独占 MQ 连接和数据库会话；数量至少为 1。
        if value < 1:
            raise ValueError("RABBITMQ_CONSUMER_COUNT must be greater than or equal to 1")
        return value
