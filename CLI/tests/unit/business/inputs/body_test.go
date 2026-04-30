package inputs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
)

func TestLoadBody_RejectsMissingAndInvalidJSON(t *testing.T) {
	if _, err := inputs.LoadBody("", ""); err == nil {
		t.Fatal("LoadBody() error = nil, want missing body error")
	}

	if _, err := inputs.LoadBody("{bad json}", ""); err == nil {
		t.Fatal("LoadBody() error = nil, want invalid json")
	}
}

func TestLoadBody_ReadsBodyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "body.json")
	if err := os.WriteFile(path, []byte("{\"ok\":true}"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := inputs.LoadBody("", path)
	if err != nil {
		t.Fatalf("LoadBody() error = %v", err)
	}
	if string(got) != "{\"ok\":true}" {
		t.Fatalf("body = %q", string(got))
	}
}
