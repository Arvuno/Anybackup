"""Pytest configuration for core_agent_service."""

from collections.abc import Iterator

import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session, sessionmaker

from core_agent_service.infrastructure.db.base import Base
from core_agent_service.infrastructure.db import models as _db_models  # noqa: F401


@pytest.fixture()
def db_engine(tmp_path):
    database_path = tmp_path / "core_agent_service_test.db"
    engine = create_engine(
        f"sqlite+pysqlite:///{database_path}",
        connect_args={"check_same_thread": False},
    )
    Base.metadata.create_all(engine)
    try:
        yield engine
    finally:
        Base.metadata.drop_all(engine)
        engine.dispose()


@pytest.fixture()
def db_session(db_engine) -> Iterator[Session]:
    session_factory = sessionmaker(bind=db_engine, expire_on_commit=False)
    session = session_factory()
    try:
        yield session
    finally:
        session.close()


@pytest.fixture()
def db_session_factory(db_engine):
    return sessionmaker(bind=db_engine, expire_on_commit=False)
