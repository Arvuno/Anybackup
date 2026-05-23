package domains_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/protect"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

type fakeExecutor struct {
	exec func(req console.RequestSpec) ([]byte, error)
}

func (f fakeExecutor) Execute(_ context.Context, _ *meta.Runtime, req console.RequestSpec) ([]byte, error) {
	return f.exec(req)
}

func TestProtectPolicyCreateBind_Success(t *testing.T) {
	def := commandByPath(t, protect.Commands(), "protect", "policy", "create-bind")
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"name":"demo-policy","backupConfig":{"dayEnable":true}}`)

	var stdout bytes.Buffer
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &stdout},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && strings.HasPrefix(req.Path, "/api/sla/v1/common/name/exists"):
					return []byte(`{"status":"success","responseData":{"isExists":false}}`), nil
				case req.Method == "POST" && req.Path == "/api/sla/v1/group/backup_info":
					return []byte(`{"status":"success","responseData":{"id":"sla-1"}}`), nil
				case req.Method == "POST" && req.Path == "/backupmgm/v1/protect_object/obj-1/slas":
					var bindReq map[string]any
					if err := json.Unmarshal(req.Body, &bindReq); err != nil {
						return nil, errors.New("invalid bind request body json")
					}
					if got, _ := bindReq["objectId"].(string); got != "obj-1" {
						return nil, errors.New("bind request missing objectId")
					}
					return []byte(`{"status":"success","responseData":null}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if def.Validate != nil {
		if err := def.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	}
	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, `"policyId":"sla-1"`) {
		t.Fatalf("stdout missing policyId, got: %s", out)
	}
	if !strings.Contains(out, `"status":"success"`) {
		t.Fatalf("stdout missing success status, got: %s", out)
	}
}

func TestProtectPolicyCreateBind_BindFailRollbackSuccess(t *testing.T) {
	def := commandByPath(t, protect.Commands(), "protect", "policy", "create-bind")
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"name":"demo-policy","backupConfig":{"dayEnable":true}}`)

	rollbackCalled := false
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && strings.HasPrefix(req.Path, "/api/sla/v1/common/name/exists"):
					return []byte(`{"status":"success","responseData":{"isExists":false}}`), nil
				case req.Method == "POST" && req.Path == "/api/sla/v1/group/backup_info":
					return []byte(`{"status":"success","responseData":{"id":"sla-rollback"}}`), nil
				case req.Method == "POST" && req.Path == "/backupmgm/v1/protect_object/obj-1/slas":
					return nil, errors.New("bind failed")
				case req.Method == "DELETE" && req.Path == "/api/sla/v1/groups":
					rollbackCalled = true
					return []byte(`{"status":"success","responseData":null}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if def.Validate != nil {
		if err := def.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	}
	err := def.Execute(context.Background(), rt, deps)
	if err == nil {
		t.Fatal("Execute() error = nil, want bind-failed error")
	}
	if !rollbackCalled {
		t.Fatal("rollback delete was not called")
	}
	if !strings.Contains(err.Error(), "rollback delete succeeded") {
		t.Fatalf("error = %q, want rollback success message", err.Error())
	}
	if !strings.Contains(err.Error(), "policyId=sla-rollback") {
		t.Fatalf("error = %q, want created policy id", err.Error())
	}
}

func TestProtectPolicyCreateBind_CreateResponseDataStringID(t *testing.T) {
	def := commandByPath(t, protect.Commands(), "protect", "policy", "create-bind")
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"name":"demo-policy","backupConfig":{"dayEnable":true}}`)

	var stdout bytes.Buffer
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &stdout},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && strings.HasPrefix(req.Path, "/api/sla/v1/common/name/exists"):
					return []byte(`{"status":"success","responseData":{"isExists":false}}`), nil
				case req.Method == "POST" && req.Path == "/api/sla/v1/group/backup_info":
					return []byte(`{"status":"success","responseData":"sla-2"}`), nil
				case req.Method == "POST" && req.Path == "/backupmgm/v1/protect_object/obj-1/slas":
					return []byte(`{"status":"success","responseData":null}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if def.Validate != nil {
		if err := def.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	}
	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout.String(), `"policyId":"sla-2"`) {
		t.Fatalf("stdout missing policyId, got: %s", stdout.String())
	}
}

func TestProtectPolicyCreateBind_CreateResponseGroupIdsArray(t *testing.T) {
	def := commandByPath(t, protect.Commands(), "protect", "policy", "create-bind")
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"name":"demo-policy","backupConfig":{"dayEnable":true}}`)

	var stdout bytes.Buffer
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &stdout},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && strings.HasPrefix(req.Path, "/api/sla/v1/common/name/exists"):
					return []byte(`{"status":"success","responseData":{"isExists":false}}`), nil
				case req.Method == "POST" && req.Path == "/api/sla/v1/group/backup_info":
					return []byte(`{"status":"success","responseData":{"groupIds":["sla-3"]}}`), nil
				case req.Method == "POST" && req.Path == "/backupmgm/v1/protect_object/obj-1/slas":
					return []byte(`{"status":"success","responseData":null}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if def.Validate != nil {
		if err := def.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	}
	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout.String(), `"policyId":"sla-3"`) {
		t.Fatalf("stdout missing policyId, got: %s", stdout.String())
	}
}
