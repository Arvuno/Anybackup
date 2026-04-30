from __future__ import annotations

from core_agent_service.application.relay_service import RelayService
from core_agent_service.config.settings import Settings
from core_agent_service.infrastructure.logging import configure_logging
from core_agent_service.infrastructure.db.repositories import (
    InboundEventRepository,
    MappingRepository,
    OutboundStatusEventRepository,
)
from core_agent_service.infrastructure.db.session import create_session_factory, init_db
from core_agent_service.infrastructure.kweaver.adapter import KWeaverRelayAdapter
from core_agent_service.infrastructure.kweaver.client import create_client
from core_agent_service.infrastructure.mq.consumer import RabbitMQConsumer
from core_agent_service.infrastructure.mq.consumer import RabbitMQConsumerPool
from core_agent_service.infrastructure.mq.publisher import StatusPublisher
import logging


logger = logging.getLogger(__name__)


def _create_adapter(current_settings: Settings) -> KWeaverRelayAdapter:
    return KWeaverRelayAdapter(
        create_client(
            base_url=current_settings.kweaver_base_url,
            token=current_settings.kweaver_token,
            username=current_settings.kweaver_username,
            password=current_settings.kweaver_password,
            business_domain=current_settings.kweaver_business_domain,
            timeout=current_settings.kweaver_timeout,
            tls_insecure=current_settings.kweaver_tls_insecure,
        ),
        agent_id=current_settings.kweaver_decision_agent_id,
        base_url=current_settings.kweaver_base_url,
        chat_timeout=current_settings.kweaver_chat_timeout,
        stream_progress_interval=current_settings.kweaver_stream_progress_interval,
        stream_trace_enabled=current_settings.kweaver_stream_trace_enabled,
        stream_trace_dir=current_settings.kweaver_stream_trace_dir,
    )


def _foundation_runtime_info(current_settings: Settings) -> dict[str, str] | None:
    if not (
        current_settings.foundation_endpoint
        and current_settings.foundation_ak
        and current_settings.foundation_sk
    ):
        return None
    return {
        "endpoint": current_settings.foundation_endpoint,
        "ak": current_settings.foundation_ak,
        "sk": current_settings.foundation_sk,
    }


def _build_worker_components(current_settings: Settings, session_factory, *, worker_index: int, probe_on_startup: bool) -> dict[str, object]:
    # 每个 worker 独占数据库会话、KWeaver client、MQ consumer 和状态 publisher，避免并发线程共享非线程安全对象。
    session = session_factory()
    logger.info("database session created", extra={"worker_index": worker_index})
    mapping_repository = MappingRepository(session)
    inbound_event_repository = InboundEventRepository(session)
    outbound_event_repository = OutboundStatusEventRepository(session)
    adapter = _create_adapter(current_settings)
    logger.info(
        "kweaver adapter created",
        extra={"decision_agent_id": current_settings.kweaver_decision_agent_id, "worker_index": worker_index},
    )
    if probe_on_startup:
        adapter.probe_connectivity()
        logger.info("decision agent connectivity probe completed", extra={"worker_index": worker_index})
    elif worker_index == 1:
        logger.info("decision agent connectivity probe skipped", extra={"probe_on_startup": False})
    relay_service = RelayService(
        mapping_repository=mapping_repository,
        inbound_event_repository=inbound_event_repository,
        adapter=adapter,
        mq_runtime_info={
            "rabbitmq_url": current_settings.rabbitmq_url,
        },
        foundation_runtime_info=_foundation_runtime_info(current_settings),
    )
    publisher = StatusPublisher(
        rabbitmq_url=current_settings.rabbitmq_url,
        exchange=current_settings.rabbitmq_status_exchange,
        exchange_type=current_settings.rabbitmq_exchange_type,
        status_routing_key=current_settings.rabbitmq_status_routing_key,
        outbound_event_repository=outbound_event_repository,
    )
    logger.info(
        "publisher created",
        extra={
            "rabbitmq_url": current_settings.rabbitmq_url,
            "exchange": current_settings.rabbitmq_status_exchange,
            "worker_index": worker_index,
        },
    )
    consumer = RabbitMQConsumer(
        rabbitmq_url=current_settings.rabbitmq_url,
        exchange=current_settings.rabbitmq_exchange,
        exchange_type=current_settings.rabbitmq_exchange_type,
        queue=current_settings.rabbitmq_queue,
        relay_service=relay_service,
        publisher=publisher,
        session=session,
    )
    logger.info(
        "consumer created",
        extra={
            "rabbitmq_url": current_settings.rabbitmq_url,
            "exchange": current_settings.rabbitmq_exchange,
            "queue": current_settings.rabbitmq_queue,
            "worker_index": worker_index,
        },
    )
    return {
        "session": session,
        "relay_service": relay_service,
        "publisher": publisher,
        "consumer": consumer,
    }


def build_app(settings: Settings | None = None) -> dict[str, object]:
    # 这里统一完成运行时装配，业务层只感知抽象依赖，不直接关心 SDK、数据库和 MQ 细节。
    current_settings = settings or Settings()
    configure_logging(current_settings)
    logger.info("building core agent service app")
    init_db(current_settings.database_url)
    logger.info("database initialized", extra={"database_url": current_settings.database_url})
    session_factory = create_session_factory(current_settings.database_url)
    worker_components = [
        _build_worker_components(
            current_settings,
            session_factory,
            worker_index=index,
            probe_on_startup=current_settings.kweaver_probe_on_startup and index == 1,
        )
        for index in range(1, current_settings.rabbitmq_consumer_count + 1)
    ]
    consumers = [components["consumer"] for components in worker_components]
    consumer_runner = consumers[0] if len(consumers) == 1 else RabbitMQConsumerPool(consumers)
    first_worker = worker_components[0]

    # 运行时返回结构保持简单，便于 CLI、main 和测试复用同一套装配结果。
    logger.info(
        "rabbitmq runtime wired",
        extra={
            "rabbitmq_url": current_settings.rabbitmq_url,
            "exchange": current_settings.rabbitmq_exchange,
            "status_exchange": current_settings.rabbitmq_status_exchange,
            "status_routing_key": current_settings.rabbitmq_status_routing_key,
            "queue": current_settings.rabbitmq_queue,
            "consumer_count": current_settings.rabbitmq_consumer_count,
        },
    )
    return {
        "settings": current_settings,
        "session": first_worker["session"],
        "relay_service": first_worker["relay_service"],
        "publisher": first_worker["publisher"],
        "consumer": first_worker["consumer"],
        "consumers": consumers,
        "consumer_runner": consumer_runner,
    }
