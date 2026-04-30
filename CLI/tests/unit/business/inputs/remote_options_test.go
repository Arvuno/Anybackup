package inputs_test

import (
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
)

func TestRemoteOptions_ValidateDefaultsAndRejectsUnsupportedVersion(t *testing.T) {
	opts := inputs.RemoteOptions{
		TenantID: "tenant-a",
		Endpoint: "https://10.10.1.8:9600",
		AK:       "ak",
		SK:       "sk",
	}

	if err := opts.Normalize(); err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if opts.TargetVersion != "9.0.9.0" {
		t.Fatalf("TargetVersion = %q, want 9.0.9.0", opts.TargetVersion)
	}

	opts.TargetVersion = "8.0.8.0"
	if err := opts.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want compatibility error")
	}
}

func TestRemoteOptions_RejectsNonURLLikeEndpoint(t *testing.T) {
	opts := inputs.RemoteOptions{
		TenantID:      "tenant-a",
		Endpoint:      "10.10.1.8:9600",
		AK:            "ak",
		SK:            "sk",
		TargetVersion: "9.0.9.0",
	}

	if err := opts.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want invalid endpoint")
	}
}

func TestRemoteOptions_NormalizeTenantIDFromEnvWhenFlagMissing(t *testing.T) {
	t.Setenv("FOUNDATION_TENANT_ID", "tenant-from-env")

	opts := inputs.RemoteOptions{
		Endpoint: "https://10.10.1.8:9600",
		AK:       "ak",
		SK:       "sk",
	}

	if err := opts.Normalize(); err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if opts.TenantID != "tenant-from-env" {
		t.Fatalf("TenantID = %q, want %q", opts.TenantID, "tenant-from-env")
	}
}

func TestRemoteOptions_AllowsEmptyTenantIDWhenNoEnv(t *testing.T) {
	t.Setenv("FOUNDATION_TENANT_ID", "")

	opts := inputs.RemoteOptions{
		Endpoint: "https://10.10.1.8:9600",
		AK:       "ak",
		SK:       "sk",
	}

	if err := opts.Normalize(); err != nil {
		t.Fatalf("Normalize() error = %v, want nil", err)
	}
	if opts.TenantID != "" {
		t.Fatalf("TenantID = %q, want empty", opts.TenantID)
	}
}
