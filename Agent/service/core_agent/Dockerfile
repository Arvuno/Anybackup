FROM docker.m.daocloud.io/library/python:3.12-slim

WORKDIR /app

ARG PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple
ARG PIP_TRUSTED_HOST=pypi.tuna.tsinghua.edu.cn

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_INDEX_URL=${PIP_INDEX_URL} \
    PIP_TRUSTED_HOST=${PIP_TRUSTED_HOST} \
    PIP_DEFAULT_TIMEOUT=120

COPY pyproject.toml README.md ./
COPY src ./src

RUN pip install --no-cache-dir setuptools wheel \
    && pip install --no-cache-dir --no-build-isolation .

CMD ["python", "-m", "core_agent_service"]
