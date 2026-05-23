package inputs_test

import (
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
)

func TestRequireNonEmpty(t *testing.T) {
	if err := inputs.RequireNonEmpty("name", ""); err == nil {
		t.Fatal("expected error for empty value")
	}
	if err := inputs.RequireNonEmpty("name", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatePaging(t *testing.T) {
	if err := inputs.ValidatePaging(-1, 10); err == nil {
		t.Fatal("expected error for negative index")
	}
	if err := inputs.ValidatePaging(0, 0); err == nil {
		t.Fatal("expected error for zero count")
	}
	if err := inputs.ValidatePaging(0, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireOneOfString_Strict(t *testing.T) {
	allow := map[string]struct{}{"1": {}, "2": {}, "4": {}}
	if err := inputs.RequireOneOfString("exec-status", "9", allow); err == nil {
		t.Fatal("want enum validation error")
	}
	if err := inputs.RequireOneOfString("exec-status", "1", allow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireAllInSet(t *testing.T) {
	allow := map[string]struct{}{"VMBackup": {}, "NasBackup": {}}
	if err := inputs.RequireAllInSet("runner-types", []string{"VMBackup", "Invalid"}, allow); err == nil {
		t.Fatal("expected error for invalid value")
	}
	if err := inputs.RequireAllInSet("runner-types", []string{"VMBackup", "NasBackup"}, allow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireNonEmpty_WhitespaceOnly(t *testing.T) {
	if err := inputs.RequireNonEmpty("name", "   "); err == nil {
		t.Fatal("expected error for whitespace-only value")
	}
	if err := inputs.RequireNonEmpty("name", "\t\n"); err == nil {
		t.Fatal("expected error for tab/newline-only value")
	}
}

func TestRequireNonEmpty_ErrorMessage(t *testing.T) {
	err := inputs.RequireNonEmpty("object-id", "")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if msg == "" {
		t.Fatal("error message should not be empty")
	}
}

func TestValidatePaging_EdgeCases(t *testing.T) {
	if err := inputs.ValidatePaging(0, -1); err == nil {
		t.Fatal("expected error for negative count")
	}
	if err := inputs.ValidatePaging(-1, -1); err == nil {
		t.Fatal("expected error for negative index (checked first)")
	}
	if err := inputs.ValidatePaging(1, 1); err != nil {
		t.Fatalf("unexpected error for valid paging: %v", err)
	}
}

func TestRequireOneOfString_EdgeCases(t *testing.T) {
	allow := map[string]struct{}{"a": {}, "b": {}}

	if err := inputs.RequireOneOfString("flag", "", allow); err == nil {
		t.Fatal("expected error for empty value")
	}

	emptyAllow := map[string]struct{}{}
	if err := inputs.RequireOneOfString("flag", "a", emptyAllow); err == nil {
		t.Fatal("expected error for empty allow map")
	}

	if err := inputs.RequireOneOfString("flag", "a", allow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireAllInSet_EdgeCases(t *testing.T) {
	allow := map[string]struct{}{"a": {}, "b": {}}

	if err := inputs.RequireAllInSet("flag", []string{}, allow); err != nil {
		t.Fatalf("unexpected error for empty values: %v", err)
	}

	emptyAllow := map[string]struct{}{}
	if err := inputs.RequireAllInSet("flag", []string{"a"}, emptyAllow); err == nil {
		t.Fatal("expected error for empty allow map")
	}
}
