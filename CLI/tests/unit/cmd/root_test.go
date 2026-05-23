package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	foundationcmd "github.com/anybackup-ai/Anybackup/CLI/cmd"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

func TestVersionCommand_PrintsStableJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
		Streams: output.Streams{
			Stdout: &stdout,
			Stderr: &stderr,
		},
		Version: foundationcmd.BuildVersionInfo("dev"),
	})

	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	const wantStdout = "{\"cliName\":\"foundation-cli\",\"cliVersion\":\"dev\",\"supportedTargetVersions\":[\"9.0.9.0\"],\"defaultTargetVersion\":\"9.0.9.0\"}\n"
	if got := stdout.String(); got != wantStdout {
		t.Fatalf("stdout = %q, want %q", got, wantStdout)
	}

	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestRootCommand_RegistersRemoteCommands(t *testing.T) {
	root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
		Streams: output.Streams{
			Stdout: &bytes.Buffer{},
			Stderr: &bytes.Buffer{},
		},
		Version: foundationcmd.BuildVersionInfo("dev"),
	})

	tests := [][]string{
		{"api"},
		{"policy", "list"},
		{"policy", "backup", "detail"},
		{"policy", "backup", "create"},
		{"policy", "copy", "detail"},
		{"policy", "copy", "create"},
		{"policy", "delete"},
		{"protect", "policy", "bind"},
		{"protect", "policy", "bind-batch"},
		{"protect", "policy", "bind-list"},
		{"protect", "policy", "batch-unbind"},
		{"protect", "backup", "start"},
		{"protect", "backup", "batch-start"},
		{"timepoint", "list"},
		{"timepoint", "clean", "start"},
		{"job", "list"},
		{"job", "backup-detail"},
		{"job", "child", "list"},
		{"job", "sync-detail"},
		{"job", "logs"},
		{"job", "business-types"},
		{"job", "app-types"},
		{"job", "stop"},
		{"job", "delete"},
		{"host", "list"},
		{"host", "object", "create"},
		{"host", "object", "list"},
		{"host", "backup-config"},
		{"host", "object", "detail"},
		{"host", "backup-config", "detail"},
		{"host", "restore", "start"},
		{"client", "deploy"},
		{"client", "deploy-history"},
		{"client", "list"},
		{"client", "runner", "list"},
		{"client", "runner-types"},
		{"client", "deploy-config", "list"},
		{"client", "deploy-log", "list"},
		{"vmware", "object", "list"},
		{"vmware", "object", "info"},
		{"vmware", "datasource", "get"},
		{"vmware", "backup-detail"},
		{"vmware", "recovery-detail"},
		{"vmware", "backup-config", "create"},
		{"vmware", "backup-config", "detail"},
		{"vmware", "restore-config", "create"},
		{"vmware", "timepoint", "metadata"},
		{"mysql", "object", "list"},
		{"mysql", "object", "get"},
		{"mysql", "backup-config", "set"},
		{"mysql", "backup-config", "detail"},
		{"mysql", "target-instance", "list"},
		{"mysql", "datasource", "backup"},
		{"mysql", "datasource", "recovery"},
		{"mysql", "recovery", "range"},
		{"mysql", "recovery-config", "detail"},
		{"mysql", "restore-config", "create"},
		{"mysql", "recovery", "timepoint", "list"},
		{"mysql", "backup-detail"},
		{"mysql", "recovery-detail"},
		{"mysql", "authorize"},
		{"network", "subnet", "list"},
		{"network", "subnet", "node", "list"},
		{"network", "subnet", "node", "ip", "list"},
		{"network", "subnet", "node", "ip", "set"},
		{"network", "subnet", "node", "ip", "remove"},
		{"storage", "service", "list"},
		{"storage", "pool", "list"},
		{"storage", "pool", "create"},
		{"storage", "pool", "delete"},
		{"storage", "pool", "node", "list"},
		{"storage", "pool", "node", "device", "list"},
	}

	for _, path := range tests {
		if _, _, err := root.Find(path); err != nil {
			t.Fatalf("Find(%v) error = %v", path, err)
		}
	}

	if cmd, remaining, err := root.Find([]string{"protect", "timepoint", "list"}); err != nil {
		t.Fatalf("Find(protect timepoint list) error = %v, want unresolved under protect", err)
	} else if cmd.Name() != "protect" || len(remaining) != 2 {
		t.Fatalf("Find(protect timepoint list) = (%s, %v), want unresolved under protect", cmd.Name(), remaining)
	}

	if cmd, remaining, err := root.Find([]string{"protect", "clean", "start"}); err != nil {
		t.Fatalf("Find(protect clean start) error = %v, want unresolved under protect", err)
	} else if cmd.Name() != "protect" || len(remaining) != 2 {
		t.Fatalf("Find(protect clean start) = (%s, %v), want unresolved under protect", cmd.Name(), remaining)
	}

	root.SetArgs([]string{"sla", "binding", "list"})
	if err := root.Execute(); err == nil {
		t.Fatal("Execute(sla binding list) error = nil, want unknown command")
	} else if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("Execute(sla binding list) error = %v, want unknown command", err)
	}

	if cmd, remaining, err := root.Find([]string{"job", "vmware", "backup-detail"}); err != nil {
		t.Fatalf("Find(job vmware backup-detail) error = %v", err)
	} else if cmd.Name() != "job" || len(remaining) != 2 {
		t.Fatalf("Find(job vmware backup-detail) = (%s, %v), want unresolved under job", cmd.Name(), remaining)
	}

	if cmd, remaining, err := root.Find([]string{"job", "mysql", "recovery-detail"}); err != nil {
		t.Fatalf("Find(job mysql recovery-detail) error = %v", err)
	} else if cmd.Name() != "job" || len(remaining) != 2 {
		t.Fatalf("Find(job mysql recovery-detail) = (%s, %v), want unresolved under job", cmd.Name(), remaining)
	}
}

func TestRootCommand_ExecutesStorageCommands(t *testing.T) {
	t.Run("service list", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"storage", "service", "list",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(storage service list) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodGet || !req.ReadOnly {
			t.Fatalf("storage service list request = %#v", req)
		}
		if req.Path != "/mstoragesvcmgm/v1/storage-svc?onlyStorage=true" {
			t.Fatalf("storage service list path = %q", req.Path)
		}
	})

	t.Run("pool create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		body := `{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1},"resource":[]}`
		root.SetArgs([]string{
			"storage", "pool", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--storage-svc-id", "svc-1",
			"--data", body,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(storage pool create) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodPost || req.ReadOnly {
			t.Fatalf("storage pool create request = %#v", req)
		}
		if req.Path != "/storageresmgm/v1/svc-1/pool" {
			t.Fatalf("storage pool create path = %q", req.Path)
		}
		if got := string(req.Body); got != body {
			t.Fatalf("storage pool create body = %q, want %q", got, body)
		}
	})
}

func TestRootCommand_ExecutesTimepointCommands(t *testing.T) {
	t.Run("timepoint list", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"timepoint", "list",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--object-id", "obj-1",
			"--business", "4",
			"--start-time", "100",
			"--end-time", "200",
			"--storage-pool-id", "12345678901234567890123456789012",
			"--is-duplication", "true",
			"--storage-service-id", "12345678901234567890123456789012",
			"--data-set-id", "12345678901234567890123456789012",
			"--businesses", "1",
			"--businesses", "6",
			"--expiration-start-time", "300",
			"--expiration-end-time", "400",
			"--usable", "2",
			"--backup-types", "1",
			"--backup-types", "4",
			"--include-storage-types", "NAS",
			"--include-storage-types", "OBJECT",
			"--exclude-storage-types", "TAPE",
			"--exclude-storage-types", "LOCAL",
			"--time-point-type", "3",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(timepoint list) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}

		req := reqs[0]
		if req.Method != http.MethodGet || !req.ReadOnly {
			t.Fatalf("timepoint list request = %#v", req)
		}

		parsed := mustParseURL(t, req.Path)
		if parsed.Path != "/backupmgm/v1/protect_object/obj-1/time_points" {
			t.Fatalf("timepoint list path = %q", parsed.Path)
		}

		query := parsed.Query()
		if got := query.Get("business"); got != "4" {
			t.Fatalf("timepoint list business = %q, want %q", got, "4")
		}
		if got := query.Get("startTime"); got != "100" {
			t.Fatalf("timepoint list startTime = %q, want %q", got, "100")
		}
		if got := query.Get("endTime"); got != "200" {
			t.Fatalf("timepoint list endTime = %q, want %q", got, "200")
		}
		if got := query.Get("storagePoolId"); got != "12345678901234567890123456789012" {
			t.Fatalf("timepoint list storagePoolId = %q", got)
		}
		if got := query.Get("isDuplication"); got != "true" {
			t.Fatalf("timepoint list isDuplication = %q, want %q", got, "true")
		}
		if got := query.Get("storageServiceId"); got != "12345678901234567890123456789012" {
			t.Fatalf("timepoint list storageServiceId = %q", got)
		}
		if got := query.Get("dataSetId"); got != "12345678901234567890123456789012" {
			t.Fatalf("timepoint list dataSetId = %q", got)
		}
		if got := query["businesses"]; len(got) != 2 || got[0] != "1" || got[1] != "6" {
			t.Fatalf("timepoint list businesses = %v, want %v", got, []string{"1", "6"})
		}
		if got := query.Get("expirationStartTime"); got != "300" {
			t.Fatalf("timepoint list expirationStartTime = %q, want %q", got, "300")
		}
		if got := query.Get("expirationEndTime"); got != "400" {
			t.Fatalf("timepoint list expirationEndTime = %q, want %q", got, "400")
		}
		if got := query.Get("usable"); got != "2" {
			t.Fatalf("timepoint list usable = %q, want %q", got, "2")
		}
		if got := query["backupTypes"]; len(got) != 2 || got[0] != "1" || got[1] != "4" {
			t.Fatalf("timepoint list backupTypes = %v, want %v", got, []string{"1", "4"})
		}
		if got := query["includeStorageTypes"]; len(got) != 2 || got[0] != "NAS" || got[1] != "OBJECT" {
			t.Fatalf("timepoint list includeStorageTypes = %v, want %v", got, []string{"NAS", "OBJECT"})
		}
		if got := query["excludeStorageTypes"]; len(got) != 2 || got[0] != "TAPE" || got[1] != "LOCAL" {
			t.Fatalf("timepoint list excludeStorageTypes = %v, want %v", got, []string{"TAPE", "LOCAL"})
		}
		if got := query.Get("timePointType"); got != "3" {
			t.Fatalf("timepoint list timePointType = %q, want %q", got, "3")
		}
	})

	t.Run("timepoint clean start", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"timepoint", "clean", "start",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--object-id", "obj-1",
			"--data", `{"force":true}`,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(timepoint clean start) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}

		req := reqs[0]
		if req.Method != http.MethodPost || req.ReadOnly {
			t.Fatalf("timepoint clean start request = %#v", req)
		}
		if req.Path != "/backupmgm/v1/protect_object/obj-1/clean_task/start" {
			t.Fatalf("timepoint clean start path = %q", req.Path)
		}
		if got := string(req.Body); got != `{"force":true}` {
			t.Fatalf("timepoint clean start body = %q, want %q", got, `{"force":true}`)
		}
	})
}

func TestRootCommand_ExecutesMySQLAuthorize(t *testing.T) {
	executeMySQLAuthorize := func(t *testing.T, bodyFlag string, bodyValue string) []console.RequestSpec {
		t.Helper()

		var reqs []console.RequestSpec
		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"mysql", "authorize",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			bodyFlag, bodyValue,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(mysql authorize %s) error = %v", bodyFlag, err)
		}
		return reqs
	}

	assertMySQLAuthorizeRequest := func(t *testing.T, reqs []console.RequestSpec, rawBody []byte) {
		t.Helper()

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}

		req := reqs[0]
		if req.Method != http.MethodPost || req.ReadOnly {
			t.Fatalf("mysql authorize request = %#v", req)
		}
		if req.Path != "/resource_center/v1/databasealone/mysql/authorize" {
			t.Fatalf("mysql authorize path = %q", req.Path)
		}
		if !bytes.Equal(req.Body, rawBody) {
			t.Fatalf("mysql authorize body = %q, want %q", string(req.Body), string(rawBody))
		}
	}

	t.Run("mysql authorize with --data", func(t *testing.T) {
		rawBody := []byte("{\"instanceName\":\"inst-1\",\"clientId\":\"c-1\",\"username\":\"u-1\",\"password\":\"p-1\",\"systemId\":\"s-1\",\"resourceId\":\"r-1\",\"osUserName\":\"os-1\",\"type\":1}")
		reqs := executeMySQLAuthorize(t, "--data", string(rawBody))
		assertMySQLAuthorizeRequest(t, reqs, rawBody)
	})

	t.Run("mysql authorize with --body-file", func(t *testing.T) {
		bodyFile := filepath.Join(t.TempDir(), "mysql-authorize.json")
		rawBody := []byte("{\n  \"instanceName\": \"inst-2\",\n  \"clientId\": \"c-2\",\n  \"username\": \"u-2\",\n  \"password\": \"p-2\",\n  \"systemId\": \"s-2\",\n  \"resourceId\": \"r-2\",\n  \"osUserName\": \"os-2\",\n  \"type\": 1\n}")
		if err := os.WriteFile(bodyFile, rawBody, 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q) error = %v", bodyFile, err)
		}

		reqs := executeMySQLAuthorize(t, "--body-file", bodyFile)
		assertMySQLAuthorizeRequest(t, reqs, rawBody)
	})
}

func TestRootCommand_ExecutesClientDeployConfigCreate_WithNameExistPrecheck(t *testing.T) {
	t.Run("name not exists then create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				switch req.Path {
				case "/deploy/v1/hostConfig/nameExist?name=linux-config":
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case "/deploy/v1/hostConfig/Linux":
					return []byte(`{}`), nil
				default:
					t.Fatalf("unexpected request path = %q", req.Path)
					return nil, nil
				}
			}),
		})

		body := `{"name":"linux-config","hostList":["10.0.0.10"],"port":22,"administrator":true,"account":{"rootPassword":"<encrypted>"}}`
		root.SetArgs([]string{
			"client", "deploy-config", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--os", "linux",
			"--data", body,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(client deploy-config create) error = %v", err)
		}

		if len(reqs) != 2 {
			t.Fatalf("executed requests = %d, want 2", len(reqs))
		}

		precheck := reqs[0]
		if precheck.Method != http.MethodGet || !precheck.ReadOnly {
			t.Fatalf("precheck request = %#v", precheck)
		}
		if precheck.Path != "/deploy/v1/hostConfig/nameExist?name=linux-config" {
			t.Fatalf("precheck path = %q", precheck.Path)
		}

		createReq := reqs[1]
		if createReq.Method != http.MethodPost || createReq.ReadOnly {
			t.Fatalf("create request = %#v", createReq)
		}
		if createReq.Path != "/deploy/v1/hostConfig/Linux" {
			t.Fatalf("create path = %q", createReq.Path)
		}
		assertJSONEq(t, createReq.Body, []byte(body))
	})

	t.Run("name exists returns error and skips create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				if req.Path == "/deploy/v1/hostConfig/nameExist?name=linux-config" {
					return []byte(`{"error":null,"status":"success","responseData":{"exist":true}}`), nil
				}
				t.Fatalf("unexpected request path = %q", req.Path)
				return nil, nil
			}),
		})

		root.SetArgs([]string{
			"client", "deploy-config", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--os", "linux",
			"--data", `{"name":"linux-config","hostList":["10.0.0.10"],"port":22,"administrator":true,"account":{"rootPassword":"<encrypted>"}}`,
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("Execute(client deploy-config create) error = nil, want duplicate name error")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("err = %q, want duplicate name error", err.Error())
		}
		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		if reqs[0].Path != "/deploy/v1/hostConfig/nameExist?name=linux-config" {
			t.Fatalf("precheck path = %q", reqs[0].Path)
		}
	})
}

func TestRootCommand_ExecutesClientDeploy_WithNameExistPrecheck(t *testing.T) {
	t.Run("name not exists then create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				switch req.Path {
				case "/deploy/v1/hostConfig/nameExist?name=sad":
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case "/deploy/v1/hostConfig/Linux":
					return []byte(`{"error":null,"status":"success","responseData":["cfg-1","cfg-2"]}`), nil
				case "/deploy/v1/job/config/nameExist?name=sad":
					return []byte(`{"error":null,"status":"success","responseData":{"exist":false}}`), nil
				case "/deploy/v1/job/config":
					return []byte(`{}`), nil
				default:
					t.Fatalf("unexpected request path = %q", req.Path)
					return nil, nil
				}
			}),
		})

		body := `{"name":"sad","clients":["c-1"],"runnerType":"Basic"}`
		root.SetArgs([]string{
			"client", "deploy",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--os", "linux",
			"--data", body,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(client deploy) error = %v", err)
		}

		if len(reqs) != 4 {
			t.Fatalf("executed requests = %d, want 4", len(reqs))
		}

		hostConfigPrecheck := reqs[0]
		if hostConfigPrecheck.Method != http.MethodGet || !hostConfigPrecheck.ReadOnly {
			t.Fatalf("hostConfig precheck request = %#v", hostConfigPrecheck)
		}
		if hostConfigPrecheck.Path != "/deploy/v1/hostConfig/nameExist?name=sad" {
			t.Fatalf("hostConfig precheck path = %q", hostConfigPrecheck.Path)
		}

		hostConfigCreateReq := reqs[1]
		if hostConfigCreateReq.Method != http.MethodPost || hostConfigCreateReq.ReadOnly {
			t.Fatalf("hostConfig create request = %#v", hostConfigCreateReq)
		}
		if hostConfigCreateReq.Path != "/deploy/v1/hostConfig/Linux" {
			t.Fatalf("hostConfig create path = %q", hostConfigCreateReq.Path)
		}
		assertJSONEq(t, hostConfigCreateReq.Body, []byte(body))

		jobPrecheck := reqs[2]
		if jobPrecheck.Method != http.MethodGet || !jobPrecheck.ReadOnly {
			t.Fatalf("job precheck request = %#v", jobPrecheck)
		}
		if jobPrecheck.Path != "/deploy/v1/job/config/nameExist?name=sad" {
			t.Fatalf("job precheck path = %q", jobPrecheck.Path)
		}

		createReq := reqs[3]
		if createReq.Method != http.MethodPost || createReq.ReadOnly {
			t.Fatalf("job create request = %#v", createReq)
		}
		if createReq.Path != "/deploy/v1/job/config" {
			t.Fatalf("job create path = %q", createReq.Path)
		}
		assertJSONEq(t, createReq.Body, []byte(`{"name":"sad","clients":["c-1"],"runnerType":"Basic","hostList":["cfg-1","cfg-2"]}`))
	})

	t.Run("name exists returns error and skips create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				if req.Path == "/deploy/v1/hostConfig/nameExist?name=sad" {
					return []byte(`{"error":null,"status":"success","responseData":{"exist":true}}`), nil
				}
				t.Fatalf("unexpected request path = %q", req.Path)
				return nil, nil
			}),
		})

		root.SetArgs([]string{
			"client", "deploy",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--os", "linux",
			"--data", `{"name":"sad","clients":["c-1"],"runnerType":"Basic"}`,
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("Execute(client deploy) error = nil, want duplicate name error")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("err = %q, want duplicate name error", err.Error())
		}
		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		if reqs[0].Path != "/deploy/v1/hostConfig/nameExist?name=sad" {
			t.Fatalf("precheck path = %q", reqs[0].Path)
		}
	})
}

func TestRootCommand_ExecutesPolicyBackupCreate_WithNameExistPrecheck(t *testing.T) {
	t.Run("name not exists then create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				switch req.Path {
				case "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5":
					return []byte(`{"status":"success","error":null,"responseData":{"isExists":false}}`), nil
				case "/api/sla/v1/group/backup_info":
					return []byte(`{}`), nil
				default:
					t.Fatalf("unexpected request path = %q", req.Path)
					return nil, nil
				}
			}),
		})

		body := `{"name":"备份策略","validatePeriod":1,"effectiveType":1,"backupConfig":{"dayEnable":true}}`
		root.SetArgs([]string{
			"policy", "backup", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--data", body,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(policy backup create) error = %v", err)
		}

		if len(reqs) != 2 {
			t.Fatalf("executed requests = %d, want 2", len(reqs))
		}

		precheck := reqs[0]
		if precheck.Method != http.MethodGet || !precheck.ReadOnly {
			t.Fatalf("precheck request = %#v", precheck)
		}
		if precheck.Path != "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
			t.Fatalf("precheck path = %q", precheck.Path)
		}

		createReq := reqs[1]
		if createReq.Method != http.MethodPost || createReq.ReadOnly {
			t.Fatalf("create request = %#v", createReq)
		}
		if createReq.Path != "/api/sla/v1/group/backup_info" {
			t.Fatalf("create path = %q", createReq.Path)
		}
		assertJSONEq(t, createReq.Body, []byte(body))
	})

	t.Run("name exists returns error and skips create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				if req.Path == "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
					return []byte(`{"status":"success","error":null,"responseData":{"isExists":true}}`), nil
				}
				t.Fatalf("unexpected request path = %q", req.Path)
				return nil, nil
			}),
		})

		root.SetArgs([]string{
			"policy", "backup", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--data", `{"name":"备份策略","validatePeriod":1,"effectiveType":1,"backupConfig":{"dayEnable":true}}`,
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("Execute(policy backup create) error = nil, want duplicate name error")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("err = %q, want duplicate name error", err.Error())
		}
		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		if reqs[0].Path != "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
			t.Fatalf("precheck path = %q", reqs[0].Path)
		}
	})
}

func TestRootCommand_ExecutesPolicyCopyCreate_WithNameExistPrecheck(t *testing.T) {
	t.Run("name not exists then create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				switch req.Path {
				case "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5":
					return []byte(`{"status":"success","error":null,"responseData":{"isExists":false}}`), nil
				case "/api/sla/v1/group/copy_info":
					return []byte(`{}`), nil
				default:
					t.Fatalf("unexpected request path = %q", req.Path)
					return nil, nil
				}
			}),
		})

		body := `{"name":"备份策略","validatePeriod":1,"effectiveType":1,"copyType":1,"copyConfigs":[{"copyMode":1,"copyConfig":{"dayEnable":true}}]}`
		root.SetArgs([]string{
			"policy", "copy", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--data", body,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(policy copy create) error = %v", err)
		}

		if len(reqs) != 2 {
			t.Fatalf("executed requests = %d, want 2", len(reqs))
		}

		precheck := reqs[0]
		if precheck.Method != http.MethodGet || !precheck.ReadOnly {
			t.Fatalf("precheck request = %#v", precheck)
		}
		if precheck.Path != "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
			t.Fatalf("precheck path = %q", precheck.Path)
		}

		createReq := reqs[1]
		if createReq.Method != http.MethodPost || createReq.ReadOnly {
			t.Fatalf("create request = %#v", createReq)
		}
		if createReq.Path != "/api/sla/v1/group/copy_info" {
			t.Fatalf("create path = %q", createReq.Path)
		}
		assertJSONEq(t, createReq.Body, []byte(body))
	})

	t.Run("name exists returns error and skips create", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				if req.Path == "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
					return []byte(`{"status":"success","error":null,"responseData":{"isExists":true}}`), nil
				}
				t.Fatalf("unexpected request path = %q", req.Path)
				return nil, nil
			}),
		})

		root.SetArgs([]string{
			"policy", "copy", "create",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--data", `{"name":"备份策略","validatePeriod":1,"effectiveType":1,"copyType":1,"copyConfigs":[{"copyMode":1,"copyConfig":{"dayEnable":true}}]}`,
		})

		err := root.Execute()
		if err == nil {
			t.Fatal("Execute(policy copy create) error = nil, want duplicate name error")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("err = %q, want duplicate name error", err.Error())
		}
		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		if reqs[0].Path != "/api/sla/v1/common/name/exists?name=%E5%A4%87%E4%BB%BD%E7%AD%96%E7%95%A5" {
			t.Fatalf("precheck path = %q", reqs[0].Path)
		}
	})
}

func TestRootCommand_ExecutesNetworkCommands(t *testing.T) {
	t.Run("network subnet list", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"network", "subnet", "list",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--storage-svc-id", "svc-1",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(network subnet list) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodGet || !req.ReadOnly {
			t.Fatalf("network subnet list request = %#v", req)
		}
		if req.Path != "/clusters/v1/svc-1/subnet?planeType=3" {
			t.Fatalf("network subnet list path = %q", req.Path)
		}
	})

	t.Run("network subnet node list", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"network", "subnet", "node", "list",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--storage-svc-id", "svc-1",
			"--subnet-id", "subnet-1",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(network subnet node list) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodGet || !req.ReadOnly {
			t.Fatalf("network subnet node list request = %#v", req)
		}
		parsed := mustParseURL(t, req.Path)
		if parsed.Path != "/clusters/v1/svc-1/subnet/nodes" {
			t.Fatalf("network subnet node list path = %q", parsed.Path)
		}
		if got := parsed.Query().Get("planeType"); got != "3" {
			t.Fatalf("network subnet node list planeType = %q, want %q", got, "3")
		}
	})

	t.Run("network subnet node ip set", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"network", "subnet", "node", "ip", "set",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--storage-svc-id", "svc-1",
			"--subnet-id", "subnet-1",
			"--node-id", "node-1",
			"--ip", "10.4.111.55",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(network subnet node ip set) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodPost || req.ReadOnly {
			t.Fatalf("network subnet node ip set request = %#v", req)
		}
		if req.Path != "/clusters/v1/svc-1/subnet/nodes/node_ip_addresses" {
			t.Fatalf("network subnet node ip set path = %q", req.Path)
		}
		if got := string(req.Body); got != `{"nodeId":"node-1","ip":"10.4.111.55","subnetId":"subnet-1"}` {
			t.Fatalf("network subnet node ip set body = %q", got)
		}
	})

	t.Run("network subnet node ip remove", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"network", "subnet", "node", "ip", "remove",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--storage-svc-id", "svc-1",
			"--subnet-id", "subnet-1",
			"--node-id", "node-1",
			"--ip", "10.4.111.55",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(network subnet node ip remove) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodDelete || req.ReadOnly {
			t.Fatalf("network subnet node ip remove request = %#v", req)
		}
		if req.Path != "/clusters/v1/svc-1/subnet/nodes/subnet-1/node-1" {
			t.Fatalf("network subnet node ip remove path = %q", req.Path)
		}
		if got := string(req.Body); got != `{"nodeId":"node-1","ip":"10.4.111.55","subnetId":"subnet-1"}` {
			t.Fatalf("network subnet node ip remove body = %q", got)
		}
	})
}

func TestRootCommand_ExecutesJobWriteCommands(t *testing.T) {
	t.Run("job stop with --job-id", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"job", "stop",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--job-id", "job-1",
			"--job-id", "job-2",
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(job stop) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodPost || req.ReadOnly {
			t.Fatalf("job stop request = %#v", req)
		}
		if req.Path != "/job_center/v1/jobs/stop" {
			t.Fatalf("job stop path = %q", req.Path)
		}
		assertJSONEq(t, req.Body, []byte(`{"jobIds":["job-1","job-2"]}`))
	})

	t.Run("job delete with --data", func(t *testing.T) {
		var reqs []console.RequestSpec

		root := foundationcmd.NewRootCommand(foundationcmd.Dependencies{
			Streams: output.Streams{
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Version: foundationcmd.BuildVersionInfo("dev"),
			Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
				reqs = append(reqs, req)
				return []byte("{}"), nil
			}),
		})

		root.SetArgs([]string{
			"job", "delete",
			"--tenant-id", "tenant-1",
			"--endpoint", "https://foundation.example",
			"--ak", "ak",
			"--sk", "sk",
			"--data", `{"jobIds":["job-9"]}`,
		})

		if err := root.Execute(); err != nil {
			t.Fatalf("Execute(job delete) error = %v", err)
		}

		if len(reqs) != 1 {
			t.Fatalf("executed requests = %d, want 1", len(reqs))
		}
		req := reqs[0]
		if req.Method != http.MethodDelete || req.ReadOnly {
			t.Fatalf("job delete request = %#v", req)
		}
		if req.Path != "/job_center/v1/jobs" {
			t.Fatalf("job delete path = %q", req.Path)
		}
		assertJSONEq(t, req.Body, []byte(`{"jobIds":["job-9"]}`))
	})
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", raw, err)
	}
	return parsed
}

func assertJSONEq(t *testing.T, got, want []byte) {
	t.Helper()

	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("json.Unmarshal(got) error = %v; got=%q", err, string(got))
	}

	var wantValue any
	if err := json.Unmarshal(want, &wantValue); err != nil {
		t.Fatalf("json.Unmarshal(want) error = %v; want=%q", err, string(want))
	}

	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("json body mismatch\ngot:  %s\nwant: %s", string(got), string(want))
	}
}
