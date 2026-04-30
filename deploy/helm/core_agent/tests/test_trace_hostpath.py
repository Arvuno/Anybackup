from pathlib import Path


CHART_ROOT = Path(__file__).resolve().parents[1]
TRACE_DIR = "/tmp/core-agent-service-trace-log"


def test_values_default_trace_dir_matches_host_path_contract():
    values = (CHART_ROOT / "values.yaml").read_text(encoding="utf-8")

    assert "kweaverStreamTraceEnabled: true" in values
    assert f'kweaverStreamTraceDir: "{TRACE_DIR}"' in values


def test_deployment_mounts_trace_hostpath_when_trace_enabled():
    deployment = (CHART_ROOT / "templates" / "deployment.yaml").read_text(encoding="utf-8")

    assert "KWEAVER_STREAM_TRACE_ENABLED" in deployment
    assert "KWEAVER_STREAM_TRACE_DIR" in deployment
    assert ".Values.config.kweaverStreamTraceEnabled" in deployment
    assert "- name: kweaver-stream-trace" in deployment
    assert "mountPath: {{ .Values.config.kweaverStreamTraceDir | quote }}" in deployment
    assert "hostPath:" in deployment
    assert "path: {{ .Values.config.kweaverStreamTraceDir | quote }}" in deployment
    assert "type: DirectoryOrCreate" in deployment


def test_values_define_foundation_connection_defaults():
    values = (CHART_ROOT / "values.yaml").read_text(encoding="utf-8")

    assert 'foundationEndpoint: ""' in values
    assert "foundationAkKey: foundation-ak" in values
    assert "foundationSkKey: foundation-sk" in values
    assert 'foundationAk: ""' in values
    assert 'foundationSk: ""' in values


def test_secret_and_deployment_expose_foundation_env():
    secret = (CHART_ROOT / "templates" / "secret.yaml").read_text(encoding="utf-8")
    deployment = (CHART_ROOT / "templates" / "deployment.yaml").read_text(encoding="utf-8")

    assert "{{ .Values.secrets.foundationAkKey }}" in secret
    assert "{{ .Values.secrets.foundationSkKey }}" in secret
    assert "FOUNDATION_ENDPOINT" in deployment
    assert ".Values.config.foundationEndpoint" in deployment
    assert "FOUNDATION_AK" in deployment
    assert ".Values.secrets.foundationAkKey" in deployment
    assert "FOUNDATION_SK" in deployment
    assert ".Values.secrets.foundationSkKey" in deployment


def test_deploy_script_accepts_foundation_connection_arguments():
    deploy_script = (CHART_ROOT / "scripts" / "deploy.sh").read_text(encoding="utf-8")

    assert 'FOUNDATION_ENDPOINT="https://115.190.150.254:9600"' in deploy_script
    assert 'FOUNDATION_AK="8422beff9f1eeadc2ea96f0ed13d03849953748c"' in deploy_script
    assert 'FOUNDATION_SK="bbe1121d876126f98b049106e4e43b4778485631b8dc23fe9119fbc99bfb3a33"' in deploy_script
    assert "--foundation-endpoint" in deploy_script
    assert "--foundation-ak" in deploy_script
    assert "--foundation-sk" in deploy_script
    assert 'config.foundationEndpoint=${FOUNDATION_ENDPOINT}' in deploy_script
    assert 'secrets.foundationAk=${FOUNDATION_AK}' in deploy_script
    assert 'secrets.foundationSk=${FOUNDATION_SK}' in deploy_script
