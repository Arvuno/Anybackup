from __future__ import annotations

import importlib
from functools import lru_cache
from typing import Any, Protocol


class KWeaverClientProtocol(Protocol):
    agents: Any
    conversations: Any


@lru_cache(maxsize=1)
def _load_kweaver_types() -> tuple[type[Any], type[Any], type[Any], type[Any] | None]:
    # 运行时依赖必须来自 pip 安装的 kweaver-sdk，这里显式校验关键类型是否存在。
    # 新版 SDK 把用户名密码登录能力暴露为 HttpSigninAuth；部分旧版则可能仍然叫 PasswordAuth。
    module = importlib.import_module("kweaver")
    client_class = getattr(module, "KWeaverClient", None)
    token_auth_class = getattr(module, "TokenAuth", None)
    config_auth_class = getattr(module, "ConfigAuth", None)
    password_auth_class = getattr(module, "HttpSigninAuth", None) or getattr(module, "PasswordAuth", None)
    if client_class is None or token_auth_class is None or config_auth_class is None:
        raise ImportError("kweaver-sdk does not expose KWeaverClient, TokenAuth and ConfigAuth")
    return client_class, token_auth_class, config_auth_class, password_auth_class


def create_client(
    *,
    base_url: str,
    token: str | None,
    username: str | None = None,
    password: str | None = None,
    business_domain: str = "bd_public",
    timeout: float = 30.0,
    tls_insecure: bool = False,
) -> KWeaverClientProtocol:
    client_class, token_auth_class, config_auth_class, password_auth_class = _load_kweaver_types()
    normalized_token = token.strip() if isinstance(token, str) else None
    normalized_username = username.strip() if isinstance(username, str) else None
    normalized_password = password.strip() if isinstance(password, str) else None

    # 优先使用显式用户名密码，避免服务进程依赖宿主机 ~/.kweaver 目录中的预置登录态。
    if normalized_username and normalized_password:
        if password_auth_class is None:
            raise ImportError("kweaver-sdk does not expose username/password authentication")
        try:
            auth = password_auth_class(
                base_url,
                username=normalized_username,
                password=normalized_password,
                tls_insecure=tls_insecure,
            )
        except TypeError:
            # 兼容仍然使用 PasswordAuth(base_url, username, password) 的旧版 SDK。
            auth = password_auth_class(base_url, normalized_username, normalized_password)
    elif normalized_token:
        # 没有用户名密码时继续支持显式 token，便于非交互式部署环境复用已有密钥配置。
        auth = token_auth_class(normalized_token)
    else:
        # 最后才回退到本机当前登录态，兼容开发机上已执行过 kweaver auth login 的场景。
        auth = config_auth_class(base_url)
    return client_class(
        base_url=base_url,
        auth=auth,
        business_domain=business_domain,
        timeout=timeout,
        tls_insecure=tls_insecure,
    )
