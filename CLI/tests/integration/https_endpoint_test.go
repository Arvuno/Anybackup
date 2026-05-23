package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anybackup-ai/Anybackup/CLI/cmd"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

func TestCLIJobList_AllowsHTTPSSelfSignedEndpoint(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/job_center/v1/jobs" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	streams := output.Streams{Stdout: stdout, Stderr: stderr}
	root := cmd.NewRootCommand(cmd.Dependencies{
		Streams: streams,
		Version: cmd.BuildVersionInfo("dev"),
		Console: cmd.NewDefaultExecutor(&http.Client{Timeout: 30 * time.Second}),
	})
	root.SetArgs([]string{
		"job", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
	})

	err := root.Execute()
	if err != nil {
		_ = output.WriteError(streams, err)
	}

	if exitCode := clierrors.ExitCodeOf(err); exitCode != clierrors.ExitSuccess {
		t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, clierrors.ExitSuccess, stderr.String())
	}
	if got := stdout.String(); got != `{"items":[]}` {
		t.Fatalf("stdout = %q, want %q", got, `{"items":[]}`)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}
