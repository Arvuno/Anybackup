package output_test

import (
	"bytes"
	"errors"
	"testing"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

func TestWriteRaw_UsesOnlyStdout(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	streams := output.Streams{Stdout: stdout, Stderr: stderr}

	if err := output.WriteRaw(streams, []byte("{\"ok\":true}\n")); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "{\"ok\":true}\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestWriteError_UsesOnlyStderr(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	streams := output.Streams{Stdout: stdout, Stderr: stderr}
	err := clierrors.Wrap(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --endpoint", errors.New("flag missing"))

	if writeErr := output.WriteError(streams, err); writeErr != nil {
		t.Fatal(writeErr)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if stderr.String() != "Cli.InvalidArgument: missing --endpoint\n" {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
