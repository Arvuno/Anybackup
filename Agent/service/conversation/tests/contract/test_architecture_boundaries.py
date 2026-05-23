from pathlib import Path

import yaml

APP_ROOT = Path(__file__).resolve().parents[2] / "app"
SERVICE_ROOT = Path(__file__).resolve().parents[2]
REPO_ROOT = Path(__file__).resolve().parents[5]

FORBIDDEN_DOMAIN_IMPORTS = (
    "fastapi",
    "pydantic",
    "sqlalchemy",
    "aio_pika",
    "dependency_injector",
    "redis",
)

FORBIDDEN_APPLICATION_IMPORTS = (
    "fastapi",
    "sqlalchemy",
    "aio_pika",
    "redis",
)

FORBIDDEN_HTTP_ROUTE_IMPORTS = (
    "app.infrastructure.persistence",
    "app.infrastructure.messaging",
    "sqlalchemy",
    "aio_pika",
    "redis",
)


def test_required_ddd_layer_packages_exist() -> None:
    for package in (
        "domain",
        "application",
        "infrastructure",
        "interfaces",
        "bootstrap",
    ):
        package_dir = APP_ROOT / package
        assert package_dir.is_dir(), f"missing DDD package: {package}"
        assert (package_dir / "__init__.py").is_file(), f"missing package marker: {package}"


def test_domain_layer_has_no_forbidden_framework_imports() -> None:
    domain_dir = APP_ROOT / "domain"

    for source_file in domain_dir.rglob("*.py"):
        source = source_file.read_text(encoding="utf-8")
        for module_name in FORBIDDEN_DOMAIN_IMPORTS:
            assert f"import {module_name}" not in source
            assert f"from {module_name}" not in source


def test_application_layer_has_no_forbidden_infrastructure_imports() -> None:
    application_dir = APP_ROOT / "application"

    for source_file in application_dir.rglob("*.py"):
        source = source_file.read_text(encoding="utf-8")
        for module_name in FORBIDDEN_APPLICATION_IMPORTS:
            assert f"import {module_name}" not in source
            assert f"from {module_name}" not in source


def test_http_routes_do_not_depend_on_infrastructure_adapters() -> None:
    http_dir = APP_ROOT / "interfaces" / "http"

    for source_file in http_dir.rglob("*.py"):
        source = source_file.read_text(encoding="utf-8")
        for module_name in FORBIDDEN_HTTP_ROUTE_IMPORTS:
            assert f"import {module_name}" not in source
            assert f"from {module_name}" not in source


def test_conversation_service_does_not_implement_auth_middleware() -> None:
    assert not (APP_ROOT / "auth_middleware.py").exists()
    assert not (
        SERVICE_ROOT / "tests" / "unit" / "interfaces" / "http" / "test_auth_middleware.py"
    ).exists()


def test_conversation_service_chart_does_not_deploy_auth_middleware() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "conversation_service_chart"

    assert not (chart_root / "templates" / "auth-middleware-deployment.yaml").exists()
    assert not (chart_root / "templates" / "auth-middleware-service.yaml").exists()

    for source_file in chart_root.rglob("*"):
        if source_file.is_file() and source_file.suffix in {".yaml", ".yml", ".tpl", ".md", ".txt"}:
            source = source_file.read_text(encoding="utf-8")
            assert "authMiddleware" not in source
            assert "auth-middleware" not in source


def test_api_gateway_auth_middleware_does_not_point_to_conversation_service() -> None:
    values_file = REPO_ROOT / "src" / "helm" / "api_gateway_service_chart" / "values.yaml"

    values_text = values_file.read_text(encoding="utf-8")
    values = yaml.safe_load(values_text)

    assert "conversation-service-conversation-service-auth-middleware" not in values_text
    assert values["auth"]["tokenToXUser"]["userinfoUrl"] == (
        "http://auth-service-auth-service.anybackup-ai.svc.cluster.local"
        "/api/auth_service/v1/realms/master/protocol/openid-connect/userinfo"
    )


def test_service_charts_own_ingress_without_entrypoints() -> None:
    for chart_name in ("auth_service_chart", "conversation_service_chart", "web_service_chart"):
        chart_root = REPO_ROOT / "src" / "helm" / chart_name
        values = yaml.safe_load((chart_root / "values.yaml").read_text(encoding="utf-8"))

        assert (chart_root / "templates" / "ingressroute.yaml").exists()
        assert "ingressRoute" in values
        assert "entryPoints" not in values["ingressRoute"]

        for source_file in chart_root.rglob("*"):
            if source_file.is_file() and source_file.suffix in {
                ".yaml",
                ".yml",
                ".tpl",
                ".md",
                ".txt",
            }:
                source = source_file.read_text(encoding="utf-8")
                assert "entryPoints" not in source


def test_web_service_chart_does_not_configure_business_api_routes() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "web_service_chart"
    chart_text = "\n".join(
        source_file.read_text(encoding="utf-8")
        for source_file in chart_root.rglob("*")
        if source_file.is_file()
        and source_file.suffix in {".yaml", ".yml", ".tpl", ".md", ".txt"}
    )

    forbidden_terms = (
        "AUTH_SERVICE_UPSTREAM",
        "CONVERSATION_SERVICE_UPSTREAM",
        "authService",
        "conversationService",
        "/api/auth_service",
        "/api/conversation_service",
        "auth-service-auth-service",
        "conversation-service-conversation-service",
    )
    for term in forbidden_terms:
        assert term not in chart_text


def test_web_service_chart_overrides_nginx_with_static_spa_config() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "web_service_chart"
    deployment_template = (chart_root / "templates" / "deployment.yaml").read_text(
        encoding="utf-8"
    )
    nginx_config_template = (
        chart_root / "templates" / "nginx-configmap.yaml"
    ).read_text(encoding="utf-8")

    assert "/etc/nginx/templates/default.conf.template" in deployment_template
    assert "default.conf.template" in nginx_config_template
    assert "try_files $uri $uri/ /index.html;" in nginx_config_template
    assert "proxy_pass" not in nginx_config_template
    assert "/api/" not in nginx_config_template


def test_web_service_ingressroute_is_public_and_uses_no_auth_middleware() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "web_service_chart"
    values = yaml.safe_load((chart_root / "values.yaml").read_text(encoding="utf-8"))
    ingressroute_template = (
        chart_root / "templates" / "ingressroute.yaml"
    ).read_text(encoding="utf-8")

    assert values["ingressRoute"]["enabled"] is True
    assert values["ingressRoute"]["pathPrefix"] == "/"
    assert values["ingressRoute"]["priority"] < 10
    assert values["ingressRoute"]["middlewares"] == []
    assert "api-gateway-service-api-gateway-service-auth" not in ingressroute_template


def test_auth_service_chart_does_not_implement_auth_middleware() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "auth_service_chart"

    for source_file in chart_root.rglob("*"):
        if source_file.is_file() and source_file.suffix in {".yaml", ".yml", ".tpl", ".md", ".txt"}:
            source = source_file.read_text(encoding="utf-8")
            assert "authMiddleware" not in source
            assert "auth-middleware" not in source


def test_auth_service_keycloak_routes_do_not_use_x_user_gateway_middleware() -> None:
    auth_values = yaml.safe_load(
        (
            REPO_ROOT / "src" / "helm" / "auth_service_chart" / "values.yaml"
        ).read_text(encoding="utf-8")
    )

    assert auth_values["ingressRoute"]["protectedMiddlewares"] == []


def test_auth_service_keycloak_health_probes_follow_relative_path() -> None:
    deployment_template = (
        REPO_ROOT / "src" / "helm" / "auth_service_chart" / "templates" / "deployment.yaml"
    ).read_text(encoding="utf-8")

    assert "managementPathPrefix := trimSuffix" in deployment_template
    assert 'printf "%s/health/started" $managementPathPrefix' in deployment_template
    assert 'printf "%s/health/ready" $managementPathPrefix' in deployment_template
    assert 'printf "%s/health/live" $managementPathPrefix' in deployment_template
    assert "path: /health/started" not in deployment_template
    assert "path: /health/ready" not in deployment_template
    assert "path: /health/live" not in deployment_template


def test_api_gateway_does_not_own_service_routes() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "api_gateway_service_chart"
    values_file = REPO_ROOT / "src" / "helm" / "api_gateway_service_chart" / "values.yaml"

    assert not (chart_root / "templates" / "ingressroutes.yaml").exists()
    values = yaml.safe_load(values_file.read_text(encoding="utf-8"))
    assert "routes" not in values


def test_api_gateway_does_not_deploy_forward_auth_backend() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "api_gateway_service_chart"
    values_file = chart_root / "values.yaml"

    values = yaml.safe_load(values_file.read_text(encoding="utf-8"))
    assert "gatewayAuth" not in values
    assert not (chart_root / "templates" / "gateway-auth-configmap.yaml").exists()
    assert not (chart_root / "templates" / "gateway-auth-deployment.yaml").exists()
    assert not (chart_root / "templates" / "gateway-auth-service.yaml").exists()


def test_api_gateway_auth_middleware_uses_token_to_x_user_plugin() -> None:
    chart_root = REPO_ROOT / "src" / "helm" / "api_gateway_service_chart"
    values = yaml.safe_load((chart_root / "values.yaml").read_text(encoding="utf-8"))
    middleware_template = (
        chart_root / "templates" / "middleware-forward-auth.yaml"
    ).read_text(encoding="utf-8")

    local_plugins = values["traefik"]["experimental"]["localPlugins"]
    assert local_plugins["tokenToXUser"]["moduleName"] == (
        "github.com/anybackup/api-gateway-token-to-x-user"
    )
    assert local_plugins["tokenToXUser"]["type"] == "localPath"
    assert values["auth"]["tokenToXUser"]["enabled"] is True
    assert values["auth"]["tokenToXUser"]["userinfoUrl"] == (
        "http://auth-service-auth-service.anybackup-ai.svc.cluster.local"
        "/api/auth_service/v1/realms/master/protocol/openid-connect/userinfo"
    )
    assert "forwardAuth" not in values["auth"]
    assert (chart_root / "templates" / "token-to-x-user-plugin-configmap.yaml").exists()
    assert (chart_root / "plugins" / "token-to-x-user" / "token_to_x_user.go").exists()

    assert "kind: Middleware" in middleware_template
    assert "plugin:" in middleware_template
    assert "tokenToXUser:" in middleware_template
    assert "userinfoUrl:" in middleware_template
    assert "forwardAuth:" not in middleware_template
    assert "chain:" not in middleware_template


def test_helm_namespace_defaults_are_anybackup_ai() -> None:
    auth_values = yaml.safe_load(
        (
            REPO_ROOT / "src" / "helm" / "auth_service_chart" / "values.yaml"
        ).read_text(encoding="utf-8")
    )
    conversation_values = yaml.safe_load(
        (
            REPO_ROOT / "src" / "helm" / "conversation_service_chart" / "values.yaml"
        ).read_text(encoding="utf-8")
    )
    gateway_readme = (
        REPO_ROOT / "src" / "helm" / "api_gateway_service_chart" / "README.md"
    ).read_text(encoding="utf-8")

    assert auth_values["ingressRoute"]["protectedMiddlewares"] == []
    assert conversation_values["ingressRoute"]["middlewares"] == [
        {
            "name": "api-gateway-service-api-gateway-service-auth",
            "namespace": "anybackup-ai",
        }
    ]
    assert "--namespace api-gateway" not in gateway_readme
    assert "--namespace anybackup-ai" in gateway_readme
