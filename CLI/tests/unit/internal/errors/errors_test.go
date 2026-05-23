package errors_test

import (
	stderrors "errors"
	"testing"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func TestExitCodeOf_DefaultsGenericErrorsToInvalidArgs(t *testing.T) {
	got := clierrors.ExitCodeOf(stderrors.New("unexpected write failure"))
	if got != clierrors.ExitTransport {
		t.Fatalf("ExitCodeOf(generic) = %d, want %d", got, clierrors.ExitTransport)
	}
}

func TestExitCodeOf_UsesWrappedCLIErrorExitCode(t *testing.T) {
	err := clierrors.Wrap(clierrors.CodeConflict, clierrors.ExitConflict, "conflict", stderrors.New("backend"))
	got := clierrors.ExitCodeOf(err)
	if got != clierrors.ExitConflict {
		t.Fatalf("ExitCodeOf(cli error) = %d, want %d", got, clierrors.ExitConflict)
	}
}

func TestExitCodeOf_MapsCLIUsageErrorsToInvalidArgs(t *testing.T) {
	got := clierrors.ExitCodeOf(stderrors.New(`unknown command "nosuch" for "foundation-cli"`))
	if got != clierrors.ExitInvalidArgs {
		t.Fatalf("ExitCodeOf(usage error) = %d, want %d", got, clierrors.ExitInvalidArgs)
	}
}
