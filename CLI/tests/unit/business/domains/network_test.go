package domains_test

import (
	"net/http"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/network"
)

func TestNetworkReadMappings(t *testing.T) {
	cmds := network.Commands()
	subnetList := commandByPath(t, cmds, "network", "subnet", "list")
	nodeList := commandByPath(t, cmds, "network", "subnet", "node", "list")
	ipList := commandByPath(t, cmds, "network", "subnet", "node", "ip", "list")

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	spec, err := subnetList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("subnet list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || !spec.ReadOnly {
		t.Fatalf("subnet list spec = %#v", spec)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/clusters/v1/svc-1/subnet" {
		t.Fatalf("subnet list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("planeType"); got != "3" {
		t.Fatalf("subnet list planeType = %q, want %q", got, "3")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	spec, err = nodeList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("node list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || !spec.ReadOnly {
		t.Fatalf("node list spec = %#v", spec)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/clusters/v1/svc-1/subnet/nodes" {
		t.Fatalf("node list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("node list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("node list count = %q, want %q", got, "30")
	}
	if got := query.Get("planeType"); got != "3" {
		t.Fatalf("node list planeType = %q, want %q", got, "3")
	}
	if got := query.Get("subnetId"); got != "subnet-1" {
		t.Fatalf("node list subnetId = %q, want %q", got, "subnet-1")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	spec, err = ipList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("ip list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || !spec.ReadOnly {
		t.Fatalf("ip list spec = %#v", spec)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/clusters/v1/svc-1/subnet/nodes/node_ip_addresses" {
		t.Fatalf("ip list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("ip list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("ip list count = %q, want %q", got, "30")
	}
	if got := query.Get("subnetId"); got != "subnet-1" {
		t.Fatalf("ip list subnetId = %q, want %q", got, "subnet-1")
	}
	if got := query.Get("nodeId"); got != "node-1" {
		t.Fatalf("ip list nodeId = %q, want %q", got, "node-1")
	}
}

func TestNetworkWriteMappings(t *testing.T) {
	cmds := network.Commands()
	ipSet := commandByPath(t, cmds, "network", "subnet", "node", "ip", "set")
	ipRemove := commandByPath(t, cmds, "network", "subnet", "node", "ip", "remove")

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "10.4.111.55")
	spec, err := ipSet.BuildRequest(rt)
	if err != nil {
		t.Fatalf("ip set BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/clusters/v1/svc-1/subnet/nodes/node_ip_addresses" || spec.ReadOnly {
		t.Fatalf("ip set spec = %#v", spec)
	}
	if got := string(spec.Body); got != `{"nodeId":"node-1","ip":"10.4.111.55","subnetId":"subnet-1"}` {
		t.Fatalf("ip set body = %q", got)
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "10.4.111.55")
	spec, err = ipRemove.BuildRequest(rt)
	if err != nil {
		t.Fatalf("ip remove BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/clusters/v1/svc-1/subnet/nodes/subnet-1/node-1" || spec.ReadOnly {
		t.Fatalf("ip remove spec = %#v", spec)
	}
	if got := string(spec.Body); got != `{"nodeId":"node-1","ip":"10.4.111.55","subnetId":"subnet-1"}` {
		t.Fatalf("ip remove body = %q", got)
	}
}

func TestNetworkValidate(t *testing.T) {
	cmds := network.Commands()
	subnetList := commandByPath(t, cmds, "network", "subnet", "list")
	nodeList := commandByPath(t, cmds, "network", "subnet", "node", "list")
	ipList := commandByPath(t, cmds, "network", "subnet", "node", "ip", "list")
	ipSet := commandByPath(t, cmds, "network", "subnet", "node", "ip", "set")
	ipRemove := commandByPath(t, cmds, "network", "subnet", "node", "ip", "remove")

	if err := subnetList.Validate(meta.NewRuntime()); err == nil {
		t.Fatal("subnet list Validate() error = nil, want missing --storage-svc-id")
	}

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetInt("plane-type", 5)
	if err := subnetList.Validate(rt); err == nil {
		t.Fatal("subnet list Validate() error = nil, want invalid --plane-type")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetInt("index", -1)
	if err := nodeList.Validate(rt); err == nil {
		t.Fatal("node list Validate() error = nil, want invalid --index")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	if err := ipList.Validate(rt); err == nil {
		t.Fatal("ip list Validate() error = nil, want missing --node-id")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "bad-ip")
	if err := ipSet.Validate(rt); err == nil {
		t.Fatal("ip set Validate() error = nil, want invalid --ip")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "  ")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "10.4.111.55")
	if err := ipRemove.Validate(rt); err == nil {
		t.Fatal("ip remove Validate() error = nil, want missing --subnet-id")
	}
}
