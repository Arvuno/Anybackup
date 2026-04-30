package policy

import (
	"encoding/json"
	"testing"
)

func TestNormalizeBackupCreateBodyDefaultsMissingFields(t *testing.T) {
	input := []byte(`{"name":"backup-policy","backupConfig":{"dayEnable":true}}`)

	body, err := normalizeBackupCreateBody(input)
	if err != nil {
		t.Fatalf("normalizeBackupCreateBody returned error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal normalized body: %v", err)
	}

	if got["validatePeriod"] != float64(1) {
		t.Fatalf("validatePeriod = %v, want 1", got["validatePeriod"])
	}
	if got["effectiveType"] != float64(1) {
		t.Fatalf("effectiveType = %v, want 1", got["effectiveType"])
	}
}

func TestNormalizeBackupCreateBodyPreservesExplicitValues(t *testing.T) {
	input := []byte(`{"name":"backup-policy","validatePeriod":2,"effectiveType":3,"backupConfig":{"dayEnable":true}}`)

	body, err := normalizeBackupCreateBody(input)
	if err != nil {
		t.Fatalf("normalizeBackupCreateBody returned error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal normalized body: %v", err)
	}

	if got["validatePeriod"] != float64(2) {
		t.Fatalf("validatePeriod = %v, want 2", got["validatePeriod"])
	}
	if got["effectiveType"] != float64(3) {
		t.Fatalf("effectiveType = %v, want 3", got["effectiveType"])
	}
}
