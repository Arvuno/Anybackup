package meta_test

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func TestRegisterRemoteCommand_UsesUnifiedExecutionChain(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	called := []string{}

	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: stdout, Stderr: stderr},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			called = append(called, "execute")
			return []byte("{\"ok\":true}\n"), nil
		}),
	}, meta.CommandMeta{
		Domain:        "sla",
		CanonicalPath: []string{"sla", "list"},
		Use:           "list",
		Description:   "List sla templates",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			called = append(called, "bind")
		},
		Validate: func(rt *meta.Runtime) error {
			called = append(called, "validate")
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			called = append(called, "build")
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/api/sla/v1/templates",
				ReadOnly: true,
			}, nil
		},
	})

	root.SetArgs([]string{
		"sla", "list",
		"--tenant-id", "t",
		"--endpoint", "https://x",
		"--ak", "ak",
		"--sk", "sk",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	wantCalled := []string{"bind", "validate", "build", "execute"}
	if !reflect.DeepEqual(called, wantCalled) {
		t.Fatalf("called = %v, want %v", called, wantCalled)
	}
	if stdout.String() != "{\"ok\":true}\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRuntime_TypedFlagHelpers(t *testing.T) {
	rt := meta.NewRuntime()
	cmd := &cobra.Command{Use: "x"}
	rt.BindIntFlag(cmd.Flags(), "index", 0, "index")
	rt.BindBoolFlag(cmd.Flags(), "include-tenant-id", false, "include tenant")
	rt.BindStringSliceFlag(cmd.Flags(), "runner-types", nil, "runner types")

	_ = cmd.Flags().Set("index", "10")
	_ = cmd.Flags().Set("include-tenant-id", "true")
	_ = cmd.Flags().Set("runner-types", "VMBackup")
	_ = cmd.Flags().Set("runner-types", "NasBackup")

	if got := rt.Int("index"); got != 10 {
		t.Fatalf("index=%d", got)
	}
	if !rt.Bool("include-tenant-id") {
		t.Fatal("include-tenant-id false")
	}
	if got := rt.Strings("runner-types"); len(got) != 2 {
		t.Fatalf("runner-types=%v", got)
	}
}

func TestRegisterRemoteCommand_PropagatesCommandContext(t *testing.T) {
	type contextKey string

	const key contextKey = "trace"
	wantValue := "ctx-value"
	gotValue := ""

	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	root.SetContext(context.WithValue(context.Background(), key, wantValue))
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			if value, ok := ctx.Value(key).(string); ok {
				gotValue = value
			}
			return []byte("{}"), nil
		}),
	}, meta.CommandMeta{
		Domain:        "sla",
		CanonicalPath: []string{"sla", "list"},
		Use:           "list",
		Description:   "List sla templates",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{Method: http.MethodGet, Path: "/api/sla/v1/templates", ReadOnly: true}, nil
		},
	})

	root.SetArgs([]string{"sla", "list", "--tenant-id", "t", "--endpoint", "https://x", "--ak", "ak", "--sk", "sk"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotValue != wantValue {
		t.Fatalf("context value = %q, want %q", gotValue, wantValue)
	}
}

func TestRuntime_String_FallbackToStringMap(t *testing.T) {
	rt := meta.NewRuntime()
	rt.SetString("key", "from-map")
	if got := rt.String("key"); got != "from-map" {
		t.Fatalf("String() = %q, want %q", got, "from-map")
	}
}

func TestRuntime_Int_ReturnsZeroWhenUnbound(t *testing.T) {
	rt := meta.NewRuntime()
	if got := rt.Int("unknown"); got != 0 {
		t.Fatalf("Int(unknown) = %d, want 0", got)
	}
}

func TestRuntime_Bool_ReturnsFalseWhenUnbound(t *testing.T) {
	rt := meta.NewRuntime()
	if got := rt.Bool("unknown"); got != false {
		t.Fatalf("Bool(unknown) = %v, want false", got)
	}
}

func TestRuntime_Strings_ReturnsNilWhenUnbound(t *testing.T) {
	rt := meta.NewRuntime()
	if got := rt.Strings("unknown"); got != nil {
		t.Fatalf("Strings(unknown) = %v, want nil", got)
	}
}

func TestRuntime_SetIntAndGet(t *testing.T) {
	rt := meta.NewRuntime()
	rt.SetInt("count", 42)
	if got := rt.Int("count"); got != 42 {
		t.Fatalf("Int(count) = %d, want 42", got)
	}
}

func TestRuntime_SetBoolAndGet(t *testing.T) {
	rt := meta.NewRuntime()
	rt.SetBool("flag", true)
	if got := rt.Bool("flag"); !got {
		t.Fatal("Bool(flag) = false, want true")
	}
}

func TestRuntime_BindStringFlag_RoundTrip(t *testing.T) {
	rt := meta.NewRuntime()
	cmd := &cobra.Command{Use: "x"}
	rt.BindStringFlag(cmd.Flags(), "name", "default", "name flag")

	if got := rt.String("name"); got != "default" {
		t.Fatalf("String(name) = %q, want %q (default)", got, "default")
	}

	_ = cmd.Flags().Set("name", "custom")
	if got := rt.String("name"); got != "custom" {
		t.Fatalf("String(name) = %q, want %q (set)", got, "custom")
	}
}

func TestRuntime_Changed_AfterExplicitFlagSet(t *testing.T) {
	rt := meta.NewRuntime()
	cmd := &cobra.Command{Use: "x"}
	rt.BindBoolFlag(cmd.Flags(), "include-tenant-id", false, "include tenant")

	if rt.Changed("include-tenant-id") {
		t.Fatal("Changed(include-tenant-id) = true before parsing, want false")
	}

	_ = cmd.Flags().Set("include-tenant-id", "false")
	rt.CaptureChangedFlags(cmd.Flags())

	if !rt.Changed("include-tenant-id") {
		t.Fatal("Changed(include-tenant-id) = false after explicit set, want true")
	}
}

func TestRegisterRemoteCommand_CapturesChangedFlagsBeforeValidate(t *testing.T) {
	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			return []byte("{}"), nil
		}),
	}, meta.CommandMeta{
		Domain:        "test",
		CanonicalPath: []string{"test", "changed"},
		Use:           "changed",
		Description:   "Test changed state capture",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindBoolFlag(cmd.Flags(), "include-tenant-id", false, "include tenant")
		},
		Validate: func(rt *meta.Runtime) error {
			if !rt.Changed("include-tenant-id") {
				return &testError{msg: "changed state missing in validate"}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			if !rt.Changed("include-tenant-id") {
				return console.RequestSpec{}, &testError{msg: "changed state missing in build"}
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/api/test/v1/changed",
				ReadOnly: true,
			}, nil
		},
	})

	root.SetArgs([]string{
		"test", "changed",
		"--tenant-id", "t",
		"--endpoint", "https://x",
		"--ak", "ak",
		"--sk", "sk",
		"--include-tenant-id=false",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestRegisterRemoteCommand_ChangedStateIsIsolatedPerExecution(t *testing.T) {
	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}

	var runCount int
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			return []byte("{}"), nil
		}),
	}, meta.CommandMeta{
		Domain:        "test",
		CanonicalPath: []string{"test", "changed-isolation"},
		Use:           "changed-isolation",
		Description:   "Test changed state isolation across executions",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindBoolFlag(cmd.Flags(), "include-tenant-id", false, "include tenant")
		},
		Validate: func(rt *meta.Runtime) error {
			runCount++
			if runCount == 1 && !rt.Changed("include-tenant-id") {
				return &testError{msg: "first run: expected changed(include-tenant-id)=true"}
			}
			if runCount == 2 && rt.Changed("include-tenant-id") {
				return &testError{msg: "second run: expected changed(include-tenant-id)=false"}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/api/test/v1/changed-isolation",
				ReadOnly: true,
			}, nil
		},
	})

	cmd, _, err := root.Find([]string{"test", "changed-isolation"})
	if err != nil {
		t.Fatalf("Find(changed-isolation) error = %v", err)
	}
	if err := cmd.Flags().Set("tenant-id", "t"); err != nil {
		t.Fatalf("Set(tenant-id) error = %v", err)
	}
	if err := cmd.Flags().Set("endpoint", "https://x"); err != nil {
		t.Fatalf("Set(endpoint) error = %v", err)
	}
	if err := cmd.Flags().Set("ak", "ak"); err != nil {
		t.Fatalf("Set(ak) error = %v", err)
	}
	if err := cmd.Flags().Set("sk", "sk"); err != nil {
		t.Fatalf("Set(sk) error = %v", err)
	}
	if err := cmd.Flags().Set("include-tenant-id", "false"); err != nil {
		t.Fatalf("Set(include-tenant-id) error = %v", err)
	}
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("first RunE() error = %v", err)
	}

	secondRunCmd := &cobra.Command{Use: "changed-isolation-second-run"}
	if err := cmd.RunE(secondRunCmd, nil); err != nil {
		t.Fatalf("second RunE() error = %v", err)
	}
}

func TestRegisterRemoteCommand_ValidationErrorPropagates(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: stdout, Stderr: stderr},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			t.Fatal("Execute should not be called")
			return nil, nil
		}),
	}, meta.CommandMeta{
		Domain:        "test",
		CanonicalPath: []string{"test", "fail"},
		Use:           "fail",
		Description:   "Test validation error",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		Validate: func(rt *meta.Runtime) error {
			return &testError{msg: "validation failed"}
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{}, nil
		},
	})

	root.SetArgs([]string{"test", "fail", "--tenant-id", "t", "--endpoint", "https://x", "--ak", "ak", "--sk", "sk"})
	err := root.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "validation failed") {
			t.Fatalf("err = %q, want 'validation failed'", err.Error())
		}
		return
	}
	t.Fatal("Execute() error = nil, want validation error")
}

func TestRegisterRemoteCommand_BuildRequestErrorPropagates(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: stdout, Stderr: stderr},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			t.Fatal("Execute should not be called")
			return nil, nil
		}),
	}, meta.CommandMeta{
		Domain:        "test",
		CanonicalPath: []string{"test", "fail"},
		Use:           "fail",
		Description:   "Test build request error",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{}, &testError{msg: "build failed"}
		},
	})

	root.SetArgs([]string{"test", "fail", "--tenant-id", "t", "--endpoint", "https://x", "--ak", "ak", "--sk", "sk"})
	err := root.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want build error")
	}
	if !strings.Contains(err.Error(), "build failed") {
		t.Fatalf("err = %q, want 'build failed'", err.Error())
	}
}

func TestRegisterRemoteCommand_CustomExecuteOverridesDefaultFlow(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	called := []string{}

	root := &cobra.Command{Use: "foundation-cli", SilenceErrors: true, SilenceUsage: true}
	meta.Register(root, meta.Dependencies{
		Streams: output.Streams{Stdout: stdout, Stderr: stderr},
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			called = append(called, "execute")
			return []byte("{}"), nil
		}),
	}, meta.CommandMeta{
		Domain:        "test",
		CanonicalPath: []string{"test", "custom"},
		Use:           "custom",
		Description:   "Test custom executor",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		Validate: func(rt *meta.Runtime) error {
			called = append(called, "validate")
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			called = append(called, "custom")
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			called = append(called, "build")
			return console.RequestSpec{Method: http.MethodPost, Path: "/should/not/run"}, nil
		},
	})

	root.SetArgs([]string{"test", "custom", "--tenant-id", "t", "--endpoint", "https://x", "--ak", "ak", "--sk", "sk"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	wantCalled := []string{"validate", "custom"}
	if !reflect.DeepEqual(called, wantCalled) {
		t.Fatalf("called = %v, want %v", called, wantCalled)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
