package domains_test

import (
	"net/http"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/storage"
)

func TestStorageReadMappings(t *testing.T) {
	cmds := storage.Commands()
	serviceList := commandByPath(t, cmds, "storage", "service", "list")
	poolList := commandByPath(t, cmds, "storage", "pool", "list")
	nodeList := commandByPath(t, cmds, "storage", "pool", "node", "list")
	deviceList := commandByPath(t, cmds, "storage", "pool", "node", "device", "list")

	spec, err := serviceList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("service list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/mstoragesvcmgm/v1/storage-svc?onlyStorage=true" || !spec.ReadOnly {
		t.Fatalf("service list spec = %#v", spec)
	}

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	spec, err = poolList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("pool list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/storageresmgm/v1/svc-1/pool" {
		t.Fatalf("pool list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("pool list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("pool list count = %q, want %q", got, "30")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	spec, err = nodeList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("node list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/storageresmgm/v1/svc-1/pool/node" {
		t.Fatalf("node list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("node list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "100" {
		t.Fatalf("node list count = %q, want %q", got, "100")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("node-id", "node-1")
	spec, err = deviceList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("device list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/storageresmgm/v1/svc-1/device" {
		t.Fatalf("device list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("nodeId"); got != "node-1" {
		t.Fatalf("device list nodeId = %q, want %q", got, "node-1")
	}
	if got := query["types"]; len(got) != 6 || got[0] != "1" || got[1] != "2" || got[2] != "3" || got[3] != "4" || got[4] != "5" || got[5] != "7" {
		t.Fatalf("device list types = %v, want %v", got, []string{"1", "2", "3", "4", "5", "7"})
	}
	if got := query.Get("authorized"); got != "2" {
		t.Fatalf("device list authorized = %q, want %q", got, "2")
	}
}

func TestStorageWriteMappings(t *testing.T) {
	cmds := storage.Commands()
	poolCreate := commandByPath(t, cmds, "storage", "pool", "create")
	poolDelete := commandByPath(t, cmds, "storage", "pool", "delete")

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.Body = []byte(`{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1},"resource":[]}`)
	spec, err := poolCreate.BuildRequest(rt)
	if err != nil {
		t.Fatalf("pool create BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/storageresmgm/v1/svc-1/pool" || spec.ReadOnly {
		t.Fatalf("pool create spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("pool-id", "pool-1")
	spec, err = poolDelete.BuildRequest(rt)
	if err != nil {
		t.Fatalf("pool delete BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.ReadOnly {
		t.Fatalf("pool delete spec = %#v", spec)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/storageresmgm/v1/svc-1/pool" {
		t.Fatalf("pool delete path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("id"); got != "pool-1" {
		t.Fatalf("pool delete id = %q, want %q", got, "pool-1")
	}
}

func TestStorageValidate(t *testing.T) {
	cmds := storage.Commands()
	poolList := commandByPath(t, cmds, "storage", "pool", "list")
	poolCreate := commandByPath(t, cmds, "storage", "pool", "create")
	poolDelete := commandByPath(t, cmds, "storage", "pool", "delete")
	deviceList := commandByPath(t, cmds, "storage", "pool", "node", "device", "list")

	if err := poolList.Validate(meta.NewRuntime()); err == nil {
		t.Fatal("pool list Validate() error = nil, want missing --storage-svc-id")
	}

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetInt("count", 0)
	if err := poolList.Validate(rt); err == nil {
		t.Fatal("pool list Validate() error = nil, want invalid --count")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	if err := deviceList.Validate(rt); err == nil {
		t.Fatal("device list Validate() error = nil, want missing --node-id")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	if err := poolDelete.Validate(rt); err == nil {
		t.Fatal("pool delete Validate() error = nil, want missing --pool-id")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	if err := poolCreate.Validate(rt); err == nil {
		t.Fatal("pool create Validate() error = nil, want missing body")
	}
}

func TestStoragePoolCreateBodyValidation(t *testing.T) {
	poolCreate := commandByPath(t, storage.Commands(), "storage", "pool", "create")

	valid := `{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1},"resource":[]}`

	t.Run("invalid json", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid json")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"type":3,"redundancyPolicy":{"policy":1},"resource":[]}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing name")
		}
	})

	t.Run("missing type", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","redundancyPolicy":{"policy":1},"resource":[]}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing type")
		}
	})

	t.Run("missing redundancyPolicy", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","type":3,"resource":[]}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing redundancyPolicy")
		}
	})

	t.Run("missing redundancyPolicy.policy", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","type":3,"redundancyPolicy":{},"resource":[]}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing redundancyPolicy.policy")
		}
	})

	t.Run("missing resource", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1}}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing resource")
		}
	})

	t.Run("invalid warnThreshold", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1},"resource":[],"warnThreshold":99}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid warnThreshold")
		}
	})

	t.Run("invalid safeThreshold", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(`{"name":"pool-a","type":3,"redundancyPolicy":{"policy":1},"resource":[],"safeThreshold":1}`)
		if err := poolCreate.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid safeThreshold")
		}
	})

	t.Run("valid body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("storage-svc-id", "svc-1")
		rt.Body = []byte(valid)
		if err := poolCreate.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	})
}
