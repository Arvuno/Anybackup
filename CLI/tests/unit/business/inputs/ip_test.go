package inputs_test

import (
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
)

func TestRequireIP(t *testing.T) {
	if err := inputs.RequireIP("ip", "10.4.111.55"); err != nil {
		t.Fatalf("RequireIP(valid ipv4) error = %v", err)
	}
	if err := inputs.RequireIP("ip", "2001:db8::1"); err != nil {
		t.Fatalf("RequireIP(valid ipv6) error = %v", err)
	}
	if err := inputs.RequireIP("ip", ""); err == nil {
		t.Fatal("RequireIP(empty) error = nil, want missing --ip")
	}
	if err := inputs.RequireIP("ip", "10.0.0.0/24"); err == nil {
		t.Fatal("RequireIP(cidr) error = nil, want invalid --ip")
	}
	if err := inputs.RequireIP("ip", "bad-ip"); err == nil {
		t.Fatal("RequireIP(invalid) error = nil, want invalid --ip")
	}
}
