package domains_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/client"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

func TestClientDeployValidateSupportsJobName(t *testing.T) {
	def := commandByPath(t, client.Commands(), "client", "deploy")
	rt := meta.NewRuntime()
	rt.SetString("os", "linux")
	rt.Body = []byte(`{"jobName":"deploy-a","hostList":["h-1"],"runners":[{"runnerType":"Gen"}]}`)

	if err := def.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestClientDeployValidateRequiresOS(t *testing.T) {
	def := commandByPath(t, client.Commands(), "client", "deploy")
	rt := meta.NewRuntime()
	rt.Body = []byte(`{"jobName":"deploy-a","hostList":["h-1"],"runners":[{"runnerType":"Gen"}]}`)

	if err := def.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want missing --os error")
	}
}

func TestClientDeployExecuteMergesHostListFromCreateResponse(t *testing.T) {
	def := commandByPath(t, client.Commands(), "client", "deploy")
	rt := meta.NewRuntime()
	rt.SetString("os", "linux")
	rt.Body = []byte(`{"jobName":"dsadadad","LinuxInstallPath":"/opt","WindowsInstallPath":"C:\\BackupAgent","concurrent":1,"uninstallType":false,"hostList":["79b53974313b11f199cd0050568943e2"],"runners":[{"runnerType":"Gen"}]}`)

	var deployReqBody []byte
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && req.Path == "/deploy/v1/hostConfig/nameExist?name=dsadadad":
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case req.Method == "POST" && req.Path == "/deploy/v1/hostConfig/Linux":
					return []byte(`{"error":null,"status":"success","responseData":["cfg-1","cfg-2"]}`), nil
				case req.Method == "GET" && req.Path == "/deploy/v1/job/config/nameExist?name=dsadadad":
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case req.Method == "POST" && req.Path == "/deploy/v1/job/config":
					deployReqBody = append([]byte(nil), req.Body...)
					return []byte(`{"error":null,"status":"success","responseData":{"jobId":"job-1"}}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if err := def.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(deployReqBody, &got); err != nil {
		t.Fatalf("json.Unmarshal(deployReqBody) error = %v, body=%s", err, string(deployReqBody))
	}

	hostList, ok := got["hostList"].([]any)
	if !ok || len(hostList) != 2 || hostList[0] != "cfg-1" || hostList[1] != "cfg-2" {
		t.Fatalf("hostList = %v, want [cfg-1 cfg-2]", got["hostList"])
	}
	if got["jobName"] != "dsadadad" {
		t.Fatalf("jobName = %v, want dsadadad", got["jobName"])
	}
	if got["LinuxInstallPath"] != "/opt" {
		t.Fatalf("LinuxInstallPath = %v, want /opt", got["LinuxInstallPath"])
	}
}

func TestClientDeployExecuteFailsWhenCreateResponseMissingIDs(t *testing.T) {
	def := commandByPath(t, client.Commands(), "client", "deploy")
	rt := meta.NewRuntime()
	rt.SetString("os", "linux")
	rt.Body = []byte(`{"jobName":"deploy-a","hostList":["h-1"],"runners":[{"runnerType":"Gen"}]}`)

	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && strings.HasPrefix(req.Path, "/deploy/v1/hostConfig/nameExist?name="):
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case req.Method == "POST" && req.Path == "/deploy/v1/hostConfig/Linux":
					return []byte(`{"error":null,"status":"success","responseData":[]}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if err := def.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	err := def.Execute(context.Background(), rt, deps)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing id array error")
	}
	if !strings.Contains(err.Error(), "missing id array") {
		t.Fatalf("error = %q, want missing id array", err.Error())
	}
}
