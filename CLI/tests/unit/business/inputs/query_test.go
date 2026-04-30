package inputs_test

import (
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
)

func TestBuildQuery_RepeatedAndMappedKeys(t *testing.T) {
	q := inputs.NewQueryBuilder().
		Add("index", "1").
		AddMapped("production-system-id", "productionSystemId", "ps-1").
		AddAll("runner-types", []string{"VMBackup", "NasBackup"}).
		Encode()
	if !strings.Contains(q, "productionSystemId=ps-1") {
		t.Fatalf("q=%s", q)
	}
	if strings.Count(q, "runner-types=") != 2 {
		t.Fatalf("q=%s", q)
	}
}

func TestBuildQuery_EncodeEmpty(t *testing.T) {
	q := inputs.NewQueryBuilder().Encode()
	if q != "" {
		t.Fatalf("expected empty, got %s", q)
	}
}

func TestBuildQuery_Chaining(t *testing.T) {
	b := inputs.NewQueryBuilder()
	result := b.Add("a", "1").Add("b", "2")
	if result != b {
		t.Fatal("expected chaining to return same builder")
	}
}

func TestBuildQuery_EncodesSpecialCharacters(t *testing.T) {
	q := inputs.NewQueryBuilder().
		Add("name", "hello world").
		Add("filter", "a&b=c").
		Encode()
	if !strings.Contains(q, "hello+world") && !strings.Contains(q, "hello%20world") {
		t.Fatalf("expected URL-encoded space in q=%s", q)
	}
	if !strings.Contains(q, "a") {
		t.Fatalf("expected key 'a' in q=%s", q)
	}
}

func TestBuildQuery_AddAllWithEmptySlice(t *testing.T) {
	q := inputs.NewQueryBuilder().
		Add("a", "1").
		AddAll("empty", []string{}).
		Encode()
	if strings.Count(q, "empty=") != 0 {
		t.Fatalf("empty slice should not add params, got q=%s", q)
	}
}

func TestBuildQuery_MultipleAddSameKey(t *testing.T) {
	q := inputs.NewQueryBuilder().
		Add("key", "a").
		Add("key", "b").
		Encode()
	if strings.Count(q, "key=") != 2 {
		t.Fatalf("expected 2 repeated keys, got q=%s", q)
	}
}

func TestBuildQuery_AddMappedUsesKeyParam(t *testing.T) {
	q := inputs.NewQueryBuilder().
		AddMapped("flag-name", "backendKey", "value").
		Encode()
	if !strings.Contains(q, "backendKey=value") {
		t.Fatalf("expected backendKey=value in q=%s", q)
	}
	if strings.Contains(q, "flag-name") {
		t.Fatalf("flag-name should not appear in q=%s", q)
	}
}
