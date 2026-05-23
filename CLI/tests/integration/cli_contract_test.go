//go:build legacy_cli_contract
// +build legacy_cli_contract

package integration_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/cmd"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/signer"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/transport"
)

func TestCLIContract(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/sla/v1/templates":
			_, _ = w.Write([]byte(`[{"id":"sla-1"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/sla/v1/group/object/object-1/templates":
			_, _ = w.Write([]byte(`[{"id":"sla-bind-1"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/protect_object/object-1/time_points":
			_, _ = w.Write([]byte(`[{"id":"tp-1"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/resource_center/v1/host/group/host_list":
			_, _ = w.Write([]byte(`[{"id":"host-1"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/file/object_list":
			_, _ = w.Write([]byte(`[{"id":"host-object-1"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/database/object_list":
			_, _ = w.Write([]byte(`[{"id":"mysql-1"}]`))
		case r.Method == http.MethodPost && r.URL.Path == "/backupmgm/v1/protect_object/object-1/backup_task/start":
			_, _ = w.Write([]byte(`{"jobId":"backup-1"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/backupmgm/v1/protect_object/conflict/backup_task/start":
			http.Error(w, "conflict", http.StatusConflict)
		case r.Method == http.MethodPost && r.URL.Path == "/backupmgm/v1/file/recovery":
			_, _ = w.Write([]byte(`{"jobId":"restore-1"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/backupmgm/v1/protect_object/object-1/clean_task/start":
			_, _ = w.Write([]byte(`{"jobId":"clean-1"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/task/job-1/detail":
			_, _ = w.Write([]byte(`{"jobId":"job-1"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/job/job-2":
			_, _ = w.Write([]byte(`{"jobId":"job-2"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/sync_task/job-3/detail":
			_, _ = w.Write([]byte(`{"jobId":"job-3"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/jobs":
			_, _ = w.Write([]byte(`{"items":[]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/business_types":
			_, _ = w.Write([]byte(`{"items":["business-1"]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/app_types":
			_, _ = w.Write([]byte(`{"items":["app-1"]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/job/unauthorized":
			http.Error(w, "no auth", http.StatusUnauthorized)
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/job/disabled":
			http.Error(w, "disabled", http.StatusForbidden)
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/task/unauthorized/detail":
			http.Error(w, "no auth", http.StatusUnauthorized)
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/task/disabled/detail":
			http.Error(w, "disabled", http.StatusForbidden)
		case r.Method == http.MethodGet && r.URL.Path == "/custom/path":
			_, _ = w.Write([]byte(`{"raw":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	cases := []struct {
		name               string
		args               []string
		wantExit           int
		wantStdout         string
		wantStderrContains string
	}{
		{
			name:       "version",
			args:       []string{"version"},
			wantExit:   0,
			wantStdout: "{\"cliName\":\"foundation-cli\",\"cliVersion\":\"dev\",\"supportedTargetVersions\":[\"9.0.9.0\"],\"defaultTargetVersion\":\"9.0.9.0\"}\n",
		},
		{
			name:       "sla list",
			args:       []string{"sla", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `[{"id":"sla-1"}]`,
		},
		{
			name:       "object sla bind-list",
			args:       []string{"object", "sla", "bind-list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--object-id", "object-1"},
			wantExit:   0,
			wantStdout: `[{"id":"sla-bind-1"}]`,
		},
		{
			name:       "object timepoint list",
			args:       []string{"object", "timepoint", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--object-id", "object-1"},
			wantExit:   0,
			wantStdout: `[{"id":"tp-1"}]`,
		},
		{
			name:       "host list",
			args:       []string{"host", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `[{"id":"host-1"}]`,
		},
		{
			name:       "host object list",
			args:       []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `[{"id":"host-object-1"}]`,
		},
		{
			name:       "mysql list",
			args:       []string{"mysql", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `[{"id":"mysql-1"}]`,
		},
		{
			name:       "backup start",
			args:       []string{"backup", "start", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--object-id", "object-1", "-d", `{"force":true}`},
			wantExit:   0,
			wantStdout: `{"jobId":"backup-1"}`,
		},
		{
			name:       "restore file",
			args:       []string{"restore", "file", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "-d", `{"files":["/tmp/a.txt"]}`},
			wantExit:   0,
			wantStdout: `{"jobId":"restore-1"}`,
		},
		{
			name:       "clean start",
			args:       []string{"clean", "start", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--object-id", "object-1", "-d", `{"force":true}`},
			wantExit:   0,
			wantStdout: `{"jobId":"clean-1"}`,
		},
		{
			name:       "job backup-detail",
			args:       []string{"job", "backup-detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "job-1"},
			wantExit:   0,
			wantStdout: `{"jobId":"job-1"}`,
		},
		{
			name:       "job detail",
			args:       []string{"job", "detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "job-2"},
			wantExit:   0,
			wantStdout: `{"jobId":"job-2"}`,
		},
		{
			name:       "job sync-detail",
			args:       []string{"job", "sync-detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "job-3"},
			wantExit:   0,
			wantStdout: `{"jobId":"job-3"}`,
		},
		{
			name:       "job business-types",
			args:       []string{"job", "business-types", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `{"items":["business-1"]}`,
		},
		{
			name:       "job app-types",
			args:       []string{"job", "app-types", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `{"items":["app-1"]}`,
		},
		{
			name:       "job list",
			args:       []string{"job", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk"},
			wantExit:   0,
			wantStdout: `{"items":[]}`,
		},
		{
			name:       "api passthrough",
			args:       []string{"api", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--method", "GET", "--path", "/custom/path"},
			wantExit:   0,
			wantStdout: `{"raw":true}`,
		},
		{
			name:               "401 unauthorized",
			args:               []string{"job", "backup-detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "unauthorized"},
			wantExit:           201,
			wantStderrContains: "Cli.Unauthorized",
		},
		{
			name:               "403 disabled user",
			args:               []string{"job", "backup-detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "disabled"},
			wantExit:           201,
			wantStderrContains: "Cli.UserMissingOrDisabled",
		},
		{
			name:               "401 unauthorized job detail",
			args:               []string{"job", "detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "unauthorized"},
			wantExit:           201,
			wantStderrContains: "Cli.Unauthorized",
		},
		{
			name:               "403 disabled user job detail",
			args:               []string{"job", "detail", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "disabled"},
			wantExit:           201,
			wantStderrContains: "Cli.UserMissingOrDisabled",
		},
		{
			name:               "409 conflict",
			args:               []string{"backup", "start", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--object-id", "conflict", "-d", `{"force":true}`},
			wantExit:           409,
			wantStderrContains: "Cli.Conflict",
		},
		{
			name:               "unsupported target version",
			args:               []string{"sla", "list", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--target-version", "8.0.8.0"},
			wantExit:           551,
			wantStderrContains: "Cli.UnsupportedBackendVersion",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := runCLI(t, srv, tc.args)
			if exitCode != tc.wantExit {
				t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, tc.wantExit, stderr)
			}
			if stdout != tc.wantStdout {
				t.Fatalf("stdout = %q, want %q", stdout, tc.wantStdout)
			}
			if tc.wantStderrContains == "" {
				if stderr != "" {
					t.Fatalf("stderr = %q, want empty", stderr)
				}
				return
			}
			if !strings.Contains(stderr, tc.wantStderrContains) {
				t.Fatalf("stderr = %q, want substring %q", stderr, tc.wantStderrContains)
			}
		})
	}
}

func TestCLIHostList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path           string
		groupID        string
		filter         string
		isChild        string
		runnerTypes    []string
		clientOsFilter []string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/resource_center/v1/host/group/host_list":
			requests <- observedRequest{
				path:           r.URL.Path,
				groupID:        r.URL.Query().Get("groupId"),
				filter:         r.URL.Query().Get("filter"),
				isChild:        r.URL.Query().Get("isChild"),
				runnerTypes:    r.URL.Query()["runnerTypes"],
				clientOsFilter: r.URL.Query()["clientOsFilter"],
			}
			_, _ = w.Write([]byte(`[{"id":"host-1"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"host", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--runner-types", "VMBackup",
		"--runner-types", "NasBackup",
		"--group-id", "group-1",
		"--client-os-filter", "1",
		"--client-os-filter", "7",
		"--filter", "oracle",
		"--is-child", "false",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `[{"id":"host-1"}]` {
		t.Fatalf("stdout = %q, want %q", stdout, `[{"id":"host-1"}]`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/resource_center/v1/host/group/host_list" {
			t.Fatalf("path = %q, want %q", req.path, "/resource_center/v1/host/group/host_list")
		}
		if req.groupID != "group-1" {
			t.Fatalf("groupId = %q, want %q", req.groupID, "group-1")
		}
		if req.filter != "oracle" {
			t.Fatalf("filter = %q, want %q", req.filter, "oracle")
		}
		if req.isChild != "false" {
			t.Fatalf("isChild = %q, want %q", req.isChild, "false")
		}
		if len(req.runnerTypes) != 2 || req.runnerTypes[0] != "VMBackup" || req.runnerTypes[1] != "NasBackup" {
			t.Fatalf("runnerTypes = %v, want %v", req.runnerTypes, []string{"VMBackup", "NasBackup"})
		}
		if len(req.clientOsFilter) != 2 || req.clientOsFilter[0] != "1" || req.clientOsFilter[1] != "7" {
			t.Fatalf("clientOsFilter = %v, want %v", req.clientOsFilter, []string{"1", "7"})
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLIHostObjectList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path               string
		index              string
		count              string
		productionSystemID string
		name               string
		datasourceType     string
		intelligentArchive string
		execStatus         string
		isIncludeTenantID  string
		tenantID           string
		groupID            string
		hasTp              string
		isBackupConfig     string
		objectID           string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/file/object_list":
			requests <- observedRequest{
				path:               r.URL.Path,
				index:              r.URL.Query().Get("index"),
				count:              r.URL.Query().Get("count"),
				productionSystemID: r.URL.Query().Get("productionSystemId"),
				name:               r.URL.Query().Get("name"),
				datasourceType:     r.URL.Query().Get("datasourceType"),
				intelligentArchive: r.URL.Query().Get("intelligentArchive"),
				execStatus:         r.URL.Query().Get("execStatus"),
				isIncludeTenantID:  r.URL.Query().Get("isIncludeTenantId"),
				tenantID:           r.URL.Query().Get("tenantId"),
				groupID:            r.URL.Query().Get("groupId"),
				hasTp:              r.URL.Query().Get("hasTp"),
				isBackupConfig:     r.URL.Query().Get("isBackupConfig"),
				objectID:           r.URL.Query().Get("objectId"),
			}
			_, _ = w.Write([]byte(`[{"id":"host-object-1"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"host", "object", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--index", "2",
		"--count", "20",
		"--production-system-id", "12345678901234567890123456789012",
		"--name", "fileset-a",
		"--datasource-type", "5",
		"--intelligent-archive", "1",
		"--exec-status", "4",
		"--is-include-tenant-id", "false",
		"--query-tenant-id", "tenant-1",
		"--group-id", "group-1",
		"--has-tp", "1",
		"--is-backup-config", "2",
		"--object-id", "abcdefghijklmnopqrstuvwx12345678",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `[{"id":"host-object-1"}]` {
		t.Fatalf("stdout = %q, want %q", stdout, `[{"id":"host-object-1"}]`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/backupmgm/v1/file/object_list" {
			t.Fatalf("path = %q, want %q", req.path, "/backupmgm/v1/file/object_list")
		}
		if req.index != "2" || req.count != "20" {
			t.Fatalf("paging = (%q,%q), want (%q,%q)", req.index, req.count, "2", "20")
		}
		if req.productionSystemID != "12345678901234567890123456789012" {
			t.Fatalf("productionSystemId = %q", req.productionSystemID)
		}
		if req.name != "fileset-a" {
			t.Fatalf("name = %q, want %q", req.name, "fileset-a")
		}
		if req.datasourceType != "5" {
			t.Fatalf("datasourceType = %q, want %q", req.datasourceType, "5")
		}
		if req.intelligentArchive != "1" {
			t.Fatalf("intelligentArchive = %q, want %q", req.intelligentArchive, "1")
		}
		if req.execStatus != "4" {
			t.Fatalf("execStatus = %q, want %q", req.execStatus, "4")
		}
		if req.isIncludeTenantID != "false" {
			t.Fatalf("isIncludeTenantId = %q, want %q", req.isIncludeTenantID, "false")
		}
		if req.tenantID != "tenant-1" {
			t.Fatalf("tenantId = %q, want %q", req.tenantID, "tenant-1")
		}
		if req.groupID != "group-1" {
			t.Fatalf("groupId = %q, want %q", req.groupID, "group-1")
		}
		if req.hasTp != "1" {
			t.Fatalf("hasTp = %q, want %q", req.hasTp, "1")
		}
		if req.isBackupConfig != "2" {
			t.Fatalf("isBackupConfig = %q, want %q", req.isBackupConfig, "2")
		}
		if req.objectID != "abcdefghijklmnopqrstuvwx12345678" {
			t.Fatalf("objectId = %q", req.objectID)
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLIMySQLList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path               string
		isAscSort          string
		sortField          string
		index              string
		count              string
		filter             string
		includedPath       string
		groupID            string
		hostID             string
		emptyHostID        string
		hostTag            string
		objectStatus       string
		objectMode         string
		canBackup          string
		execStatus         string
		allChild           string
		nameFilter         string
		appType            string
		isolation          string
		objectID           string
		productionSystemID []string
		objectIDs          []string
		objectTypes        []string
		appTypes           []string
		isConfig           []string
		bindSlaStatus      []string
		includedBindSla    []string
		excludeBindSla     []string
		bindSlaIDs         []string
		protectStatus      []string
		lastBackupStatus   []string
		lastSnapshotStatus []string
		parentIDs          []string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/database/object_list":
			requests <- observedRequest{
				path:               r.URL.Path,
				isAscSort:          r.URL.Query().Get("isAscSort"),
				sortField:          r.URL.Query().Get("sortField"),
				index:              r.URL.Query().Get("index"),
				count:              r.URL.Query().Get("count"),
				filter:             r.URL.Query().Get("filter"),
				includedPath:       r.URL.Query().Get("includedPath"),
				groupID:            r.URL.Query().Get("groupId"),
				hostID:             r.URL.Query().Get("hostId"),
				emptyHostID:        r.URL.Query().Get("emptyHostId"),
				hostTag:            r.URL.Query().Get("hostTag"),
				objectStatus:       r.URL.Query().Get("objectStatus"),
				objectMode:         r.URL.Query().Get("objectMode"),
				canBackup:          r.URL.Query().Get("canBackup"),
				execStatus:         r.URL.Query().Get("execStatus"),
				allChild:           r.URL.Query().Get("allChild"),
				nameFilter:         r.URL.Query().Get("nameFilter"),
				appType:            r.URL.Query().Get("appType"),
				isolation:          r.URL.Query().Get("isolation"),
				objectID:           r.URL.Query().Get("objectId"),
				productionSystemID: r.URL.Query()["productionSystemIds"],
				objectIDs:          r.URL.Query()["objectIds"],
				objectTypes:        r.URL.Query()["objectTypes"],
				appTypes:           r.URL.Query()["appTypes"],
				isConfig:           r.URL.Query()["isConfig"],
				bindSlaStatus:      r.URL.Query()["bindSlaStatus"],
				includedBindSla:    r.URL.Query()["includedBindSla"],
				excludeBindSla:     r.URL.Query()["excludeBindSla"],
				bindSlaIDs:         r.URL.Query()["bindSlaIds"],
				protectStatus:      r.URL.Query()["protectStatus"],
				lastBackupStatus:   r.URL.Query()["lastBackupStatus"],
				lastSnapshotStatus: r.URL.Query()["lastSnapshotStatus"],
				parentIDs:          r.URL.Query()["parentIds"],
			}
			_, _ = w.Write([]byte(`[{"id":"mysql-1"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"mysql", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--is-asc-sort", "false",
		"--sort-field", "4",
		"--index", "2",
		"--count", "30",
		"--filter", "prod",
		"--included-path", "true",
		"--group-id", "default",
		"--host-id", "12345678901234567890123456789012",
		"--empty-host-id", "false",
		"--host-tag", "abcdefghijklmnopqrstuvwx12345678",
		"--production-system-ids", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"--production-system-ids", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"--object-ids", "11111111111111111111111111111111",
		"--object-ids", "22222222222222222222222222222222",
		"--object-types", "instance",
		"--object-types", "cluster",
		"--app-types", "MySQL",
		"--app-types", "MariaDB",
		"--is-config", "1",
		"--is-config", "3",
		"--object-status", "2",
		"--object-mode", "1",
		"--can-backup", "2",
		"--bind-sla-status", "1",
		"--bind-sla-status", "6",
		"--included-bind-sla", "2",
		"--included-bind-sla", "5",
		"--exclude-bind-sla", "1",
		"--exclude-bind-sla", "4",
		"--bind-sla-ids", "cccccccccccccccccccccccccccccccc",
		"--bind-sla-ids", "dddddddddddddddddddddddddddddddd",
		"--protect-status", "1",
		"--protect-status", "4",
		"--last-backup-status", "600",
		"--last-backup-status", "900",
		"--last-snapshot-status", "800",
		"--last-snapshot-status", "1100",
		"--exec-status", "4",
		"--parent-ids", "parent-a",
		"--parent-ids", "parent-b",
		"--all-child", "true",
		"--name-filter", "mysql-a",
		"--app-type", "42",
		"--isolation", "false",
		"--object-id", "object-single",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `[{"id":"mysql-1"}]` {
		t.Fatalf("stdout = %q, want %q", stdout, `[{"id":"mysql-1"}]`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/backupmgm/v1/database/object_list" {
			t.Fatalf("path = %q, want %q", req.path, "/backupmgm/v1/database/object_list")
		}
		if req.isAscSort != "false" || req.sortField != "4" {
			t.Fatalf("sort = (%q,%q)", req.isAscSort, req.sortField)
		}
		if req.index != "2" || req.count != "30" {
			t.Fatalf("paging = (%q,%q), want (%q,%q)", req.index, req.count, "2", "30")
		}
		if req.filter != "prod" || req.includedPath != "true" {
			t.Fatalf("filter/includedPath = (%q,%q)", req.filter, req.includedPath)
		}
		if req.groupID != "default" || req.hostID != "12345678901234567890123456789012" {
			t.Fatalf("groupId/hostId = (%q,%q)", req.groupID, req.hostID)
		}
		if req.emptyHostID != "false" || req.hostTag != "abcdefghijklmnopqrstuvwx12345678" {
			t.Fatalf("emptyHostId/hostTag = (%q,%q)", req.emptyHostID, req.hostTag)
		}
		if req.objectStatus != "2" || req.objectMode != "1" || req.canBackup != "2" || req.execStatus != "4" {
			t.Fatalf("status filters = (%q,%q,%q,%q)", req.objectStatus, req.objectMode, req.canBackup, req.execStatus)
		}
		if req.allChild != "true" || req.nameFilter != "mysql-a" || req.appType != "42" || req.isolation != "false" || req.objectID != "object-single" {
			t.Fatalf("tail filters = (%q,%q,%q,%q,%q)", req.allChild, req.nameFilter, req.appType, req.isolation, req.objectID)
		}
		if len(req.productionSystemID) != 2 || req.productionSystemID[0] != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" || req.productionSystemID[1] != "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
			t.Fatalf("productionSystemIds = %v", req.productionSystemID)
		}
		if len(req.objectIDs) != 2 || req.objectIDs[0] != "11111111111111111111111111111111" || req.objectIDs[1] != "22222222222222222222222222222222" {
			t.Fatalf("objectIds = %v", req.objectIDs)
		}
		if len(req.objectTypes) != 2 || req.objectTypes[0] != "instance" || req.objectTypes[1] != "cluster" {
			t.Fatalf("objectTypes = %v", req.objectTypes)
		}
		if len(req.appTypes) != 2 || req.appTypes[0] != "MySQL" || req.appTypes[1] != "MariaDB" {
			t.Fatalf("appTypes = %v", req.appTypes)
		}
		if len(req.isConfig) != 2 || req.isConfig[0] != "1" || req.isConfig[1] != "3" {
			t.Fatalf("isConfig = %v", req.isConfig)
		}
		if len(req.bindSlaStatus) != 2 || req.bindSlaStatus[0] != "1" || req.bindSlaStatus[1] != "6" {
			t.Fatalf("bindSlaStatus = %v", req.bindSlaStatus)
		}
		if len(req.includedBindSla) != 2 || req.includedBindSla[0] != "2" || req.includedBindSla[1] != "5" {
			t.Fatalf("includedBindSla = %v", req.includedBindSla)
		}
		if len(req.excludeBindSla) != 2 || req.excludeBindSla[0] != "1" || req.excludeBindSla[1] != "4" {
			t.Fatalf("excludeBindSla = %v", req.excludeBindSla)
		}
		if len(req.bindSlaIDs) != 2 || req.bindSlaIDs[0] != "cccccccccccccccccccccccccccccccc" || req.bindSlaIDs[1] != "dddddddddddddddddddddddddddddddd" {
			t.Fatalf("bindSlaIds = %v", req.bindSlaIDs)
		}
		if len(req.protectStatus) != 2 || req.protectStatus[0] != "1" || req.protectStatus[1] != "4" {
			t.Fatalf("protectStatus = %v", req.protectStatus)
		}
		if len(req.lastBackupStatus) != 2 || req.lastBackupStatus[0] != "600" || req.lastBackupStatus[1] != "900" {
			t.Fatalf("lastBackupStatus = %v", req.lastBackupStatus)
		}
		if len(req.lastSnapshotStatus) != 2 || req.lastSnapshotStatus[0] != "800" || req.lastSnapshotStatus[1] != "1100" {
			t.Fatalf("lastSnapshotStatus = %v", req.lastSnapshotStatus)
		}
		if len(req.parentIDs) != 2 || req.parentIDs[0] != "parent-a" || req.parentIDs[1] != "parent-b" {
			t.Fatalf("parentIds = %v", req.parentIDs)
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLIJobLogs_AlwaysSendsPagingQuery(t *testing.T) {
	type observedRequest struct {
		path      string
		index     string
		count     string
		startTime string
		endTime   string
		level     string
	}

	requests := make(chan observedRequest, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/activity/job-default/logs":
			requests <- observedRequest{
				path:      r.URL.Path,
				index:     r.URL.Query().Get("index"),
				count:     r.URL.Query().Get("count"),
				startTime: r.URL.Query().Get("startTime"),
				endTime:   r.URL.Query().Get("endTime"),
				level:     r.URL.Query().Get("level"),
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/activity/job-explicit/logs":
			requests <- observedRequest{
				path:      r.URL.Path,
				index:     r.URL.Query().Get("index"),
				count:     r.URL.Query().Get("count"),
				startTime: r.URL.Query().Get("startTime"),
				endTime:   r.URL.Query().Get("endTime"),
				level:     r.URL.Query().Get("level"),
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	cases := []struct {
		name      string
		args      []string
		wantPath  string
		wantIndex string
		wantCount string
		wantStart string
		wantEnd   string
		wantLevel string
	}{
		{
			name:      "default paging",
			args:      []string{"job", "logs", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "job-default"},
			wantPath:  "/job_center/v1/activity/job-default/logs",
			wantIndex: "0",
			wantCount: "30",
			wantStart: "",
			wantEnd:   "",
			wantLevel: "",
		},
		{
			name:      "explicit paging",
			args:      []string{"job", "logs", "--tenant-id", "t", "--endpoint", srv.URL, "--ak", "ak", "--sk", "sk", "--job-id", "job-explicit", "--index", "4", "--count", "7", "--start-time", "100", "--end-time", "200", "--level", "3"},
			wantPath:  "/job_center/v1/activity/job-explicit/logs",
			wantIndex: "4",
			wantCount: "7",
			wantStart: "100",
			wantEnd:   "200",
			wantLevel: "3",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := runCLI(t, srv, tc.args)
			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
			}
			if stdout != `{"items":[]}` {
				t.Fatalf("stdout = %q, want %q", stdout, `{"items":[]}`)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			select {
			case req := <-requests:
				if req.path != tc.wantPath {
					t.Fatalf("path = %q, want %q", req.path, tc.wantPath)
				}
				if req.index != tc.wantIndex {
					t.Fatalf("index = %q, want %q", req.index, tc.wantIndex)
				}
				if req.count != tc.wantCount {
					t.Fatalf("count = %q, want %q", req.count, tc.wantCount)
				}
				if req.startTime != tc.wantStart {
					t.Fatalf("startTime = %q, want %q", req.startTime, tc.wantStart)
				}
				if req.endTime != tc.wantEnd {
					t.Fatalf("endTime = %q, want %q", req.endTime, tc.wantEnd)
				}
				if req.level != tc.wantLevel {
					t.Fatalf("level = %q, want %q", req.level, tc.wantLevel)
				}
			default:
				t.Fatal("backend did not receive request")
			}
		})
	}
}

func TestCLIJobList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path             string
		index            string
		count            string
		objectID         string
		objectName       string
		startTime        string
		endTime          string
		sort             string
		remark           string
		clusterID        string
		storageServiceID string
		storagePoolID    string
		strategyName     string
		clientName       string
		taskContinueType string
		statuses         []string
		appTypes         []string
		operationTypes   []string
		businessTypes    []string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/job_center/v1/jobs":
			requests <- observedRequest{
				path:             r.URL.Path,
				index:            r.URL.Query().Get("index"),
				count:            r.URL.Query().Get("count"),
				objectID:         r.URL.Query().Get("objectId"),
				objectName:       r.URL.Query().Get("objectName"),
				startTime:        r.URL.Query().Get("startTime"),
				endTime:          r.URL.Query().Get("endTime"),
				sort:             r.URL.Query().Get("sort"),
				remark:           r.URL.Query().Get("remark"),
				clusterID:        r.URL.Query().Get("clusterId"),
				storageServiceID: r.URL.Query().Get("storageServiceId"),
				storagePoolID:    r.URL.Query().Get("storagePoolId"),
				strategyName:     r.URL.Query().Get("strategyName"),
				clientName:       r.URL.Query().Get("clientName"),
				taskContinueType: r.URL.Query().Get("taskContinueType"),
				statuses:         r.URL.Query()["statuses"],
				appTypes:         r.URL.Query()["appTypes"],
				operationTypes:   r.URL.Query()["operationTypes"],
				businessTypes:    r.URL.Query()["businessTypes"],
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"job", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--index", "2",
		"--count", "20",
		"--statuses", "100",
		"--statuses", "400",
		"--object-id", "object-1",
		"--object-name", "prod-db",
		"--app-types", "Oracle",
		"--app-types", "VMware",
		"--operation-types", "1",
		"--operation-types", "2",
		"--business-types", "10",
		"--business-types", "20",
		"--start-time", "100",
		"--end-time", "200",
		"--sort", "-startTime",
		"--remark", "nightly backup",
		"--cluster-id", "cluster-1",
		"--storage-service-id", "storage-service-1",
		"--storage-pool-id", "storage-pool-1",
		"--strategy-name", "gold",
		"--client-name", "client-a",
		"--task-continue-type", "1",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `{"items":[]}` {
		t.Fatalf("stdout = %q, want %q", stdout, `{"items":[]}`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/job_center/v1/jobs" {
			t.Fatalf("path = %q, want %q", req.path, "/job_center/v1/jobs")
		}
		if req.index != "2" || req.count != "20" {
			t.Fatalf("paging = (%q,%q), want (%q,%q)", req.index, req.count, "2", "20")
		}
		if req.objectID != "object-1" || req.objectName != "prod-db" {
			t.Fatalf("object = (%q,%q)", req.objectID, req.objectName)
		}
		if req.startTime != "100" || req.endTime != "200" {
			t.Fatalf("time range = (%q,%q)", req.startTime, req.endTime)
		}
		if req.sort != "-startTime" {
			t.Fatalf("sort = %q, want %q", req.sort, "-startTime")
		}
		if req.remark != "nightly backup" {
			t.Fatalf("remark = %q, want %q", req.remark, "nightly backup")
		}
		if req.clusterID != "cluster-1" {
			t.Fatalf("clusterId = %q, want %q", req.clusterID, "cluster-1")
		}
		if req.storageServiceID != "storage-service-1" {
			t.Fatalf("storageServiceId = %q, want %q", req.storageServiceID, "storage-service-1")
		}
		if req.storagePoolID != "storage-pool-1" {
			t.Fatalf("storagePoolId = %q, want %q", req.storagePoolID, "storage-pool-1")
		}
		if req.strategyName != "gold" {
			t.Fatalf("strategyName = %q, want %q", req.strategyName, "gold")
		}
		if req.clientName != "client-a" {
			t.Fatalf("clientName = %q, want %q", req.clientName, "client-a")
		}
		if req.taskContinueType != "1" {
			t.Fatalf("taskContinueType = %q, want %q", req.taskContinueType, "1")
		}
		if len(req.statuses) != 2 || req.statuses[0] != "100" || req.statuses[1] != "400" {
			t.Fatalf("statuses = %v, want %v", req.statuses, []string{"100", "400"})
		}
		if len(req.appTypes) != 2 || req.appTypes[0] != "Oracle" || req.appTypes[1] != "VMware" {
			t.Fatalf("appTypes = %v, want %v", req.appTypes, []string{"Oracle", "VMware"})
		}
		if len(req.operationTypes) != 2 || req.operationTypes[0] != "1" || req.operationTypes[1] != "2" {
			t.Fatalf("operationTypes = %v, want %v", req.operationTypes, []string{"1", "2"})
		}
		if len(req.businessTypes) != 2 || req.businessTypes[0] != "10" || req.businessTypes[1] != "20" {
			t.Fatalf("businessTypes = %v, want %v", req.businessTypes, []string{"10", "20"})
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLISLAList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path          string
		index         string
		count         string
		filter        string
		validateState string
		disableMark   string
		businessType  string
		copyMode      string
		backupMode    string
		types         []string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/sla/v1/templates":
			requests <- observedRequest{
				path:          r.URL.Path,
				index:         r.URL.Query().Get("index"),
				count:         r.URL.Query().Get("count"),
				filter:        r.URL.Query().Get("filter"),
				validateState: r.URL.Query().Get("validateStatus"),
				disableMark:   r.URL.Query().Get("disableMark"),
				businessType:  r.URL.Query().Get("type"),
				copyMode:      r.URL.Query().Get("copyMode"),
				backupMode:    r.URL.Query().Get("backupMode"),
				types:         r.URL.Query()["types"],
			}
			_, _ = w.Write([]byte(`[{"id":"sla-1"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"sla", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--index", "2",
		"--count", "20",
		"--filter", "oracle",
		"--validate-status", "3",
		"--disable-mark", "0",
		"--type", "Database",
		"--types", "Database",
		"--types", "Fileset",
		"--copy-mode", "2",
		"--backup-mode", "5",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `[{"id":"sla-1"}]` {
		t.Fatalf("stdout = %q, want %q", stdout, `[{"id":"sla-1"}]`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/api/sla/v1/templates" {
			t.Fatalf("path = %q, want %q", req.path, "/api/sla/v1/templates")
		}
		if req.index != "2" || req.count != "20" {
			t.Fatalf("paging = (%q,%q), want (%q,%q)", req.index, req.count, "2", "20")
		}
		if req.filter != "oracle" {
			t.Fatalf("filter = %q, want %q", req.filter, "oracle")
		}
		if req.validateState != "3" {
			t.Fatalf("validateStatus = %q, want %q", req.validateState, "3")
		}
		if req.disableMark != "0" {
			t.Fatalf("disableMark = %q, want %q", req.disableMark, "0")
		}
		if req.businessType != "Database" {
			t.Fatalf("type = %q, want %q", req.businessType, "Database")
		}
		if req.copyMode != "2" {
			t.Fatalf("copyMode = %q, want %q", req.copyMode, "2")
		}
		if req.backupMode != "5" {
			t.Fatalf("backupMode = %q, want %q", req.backupMode, "5")
		}
		if len(req.types) != 2 || req.types[0] != "Database" || req.types[1] != "Fileset" {
			t.Fatalf("types = %v, want %v", req.types, []string{"Database", "Fileset"})
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLIObjectTimepointList_ForwardsOptionalQueryParams(t *testing.T) {
	type observedRequest struct {
		path                string
		business            string
		startTime           string
		endTime             string
		storagePoolID       string
		isDuplication       string
		storageServiceID    string
		dataSetID           string
		expirationStartTime string
		expirationEndTime   string
		usable              string
		timePointType       string
		businesses          []string
		backupTypes         []string
		includeStorageTypes []string
		excludeStorageTypes []string
	}

	requests := make(chan observedRequest, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/backupmgm/v1/protect_object/object-1/time_points":
			requests <- observedRequest{
				path:                r.URL.Path,
				business:            r.URL.Query().Get("business"),
				startTime:           r.URL.Query().Get("startTime"),
				endTime:             r.URL.Query().Get("endTime"),
				storagePoolID:       r.URL.Query().Get("storagePoolId"),
				isDuplication:       r.URL.Query().Get("isDuplication"),
				storageServiceID:    r.URL.Query().Get("storageServiceId"),
				dataSetID:           r.URL.Query().Get("dataSetId"),
				expirationStartTime: r.URL.Query().Get("expirationStartTime"),
				expirationEndTime:   r.URL.Query().Get("expirationEndTime"),
				usable:              r.URL.Query().Get("usable"),
				timePointType:       r.URL.Query().Get("timePointType"),
				businesses:          r.URL.Query()["businesses"],
				backupTypes:         r.URL.Query()["backupTypes"],
				includeStorageTypes: r.URL.Query()["includeStorageTypes"],
				excludeStorageTypes: r.URL.Query()["excludeStorageTypes"],
			}
			_, _ = w.Write([]byte(`[{"id":"tp-1"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	args := []string{
		"object", "timepoint", "list",
		"--tenant-id", "t",
		"--endpoint", srv.URL,
		"--ak", "ak",
		"--sk", "sk",
		"--object-id", "object-1",
		"--business", "4",
		"--start-time", "100",
		"--end-time", "200",
		"--storage-pool-id", "12345678901234567890123456789012",
		"--is-duplication", "false",
		"--storage-service-id", "abcdefghijabcdefghijabcdefghijab",
		"--data-set-id", "fedcbafedcbafedcbafedcbafedcbafe",
		"--businesses", "1",
		"--businesses", "6",
		"--expiration-start-time", "300",
		"--expiration-end-time", "400",
		"--usable", "2",
		"--backup-types", "1",
		"--backup-types", "4",
		"--include-storage-types", "2",
		"--include-storage-types", "3",
		"--exclude-storage-types", "5",
		"--time-point-type", "3",
	}

	stdout, stderr, exitCode := runCLI(t, srv, args)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr = %q", exitCode, stderr)
	}
	if stdout != `[{"id":"tp-1"}]` {
		t.Fatalf("stdout = %q, want %q", stdout, `[{"id":"tp-1"}]`)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	select {
	case req := <-requests:
		if req.path != "/backupmgm/v1/protect_object/object-1/time_points" {
			t.Fatalf("path = %q, want %q", req.path, "/backupmgm/v1/protect_object/object-1/time_points")
		}
		if req.business != "4" || req.startTime != "100" || req.endTime != "200" {
			t.Fatalf("business/start/end = (%q,%q,%q)", req.business, req.startTime, req.endTime)
		}
		if req.storagePoolID != "12345678901234567890123456789012" {
			t.Fatalf("storagePoolId = %q", req.storagePoolID)
		}
		if req.isDuplication != "false" {
			t.Fatalf("isDuplication = %q, want %q", req.isDuplication, "false")
		}
		if req.storageServiceID != "abcdefghijabcdefghijabcdefghijab" {
			t.Fatalf("storageServiceId = %q", req.storageServiceID)
		}
		if req.dataSetID != "fedcbafedcbafedcbafedcbafedcbafe" {
			t.Fatalf("dataSetId = %q", req.dataSetID)
		}
		if req.expirationStartTime != "300" || req.expirationEndTime != "400" {
			t.Fatalf("expiration range = (%q,%q)", req.expirationStartTime, req.expirationEndTime)
		}
		if req.usable != "2" {
			t.Fatalf("usable = %q, want %q", req.usable, "2")
		}
		if req.timePointType != "3" {
			t.Fatalf("timePointType = %q, want %q", req.timePointType, "3")
		}
		if len(req.businesses) != 2 || req.businesses[0] != "1" || req.businesses[1] != "6" {
			t.Fatalf("businesses = %v, want %v", req.businesses, []string{"1", "6"})
		}
		if len(req.backupTypes) != 2 || req.backupTypes[0] != "1" || req.backupTypes[1] != "4" {
			t.Fatalf("backupTypes = %v, want %v", req.backupTypes, []string{"1", "4"})
		}
		if len(req.includeStorageTypes) != 2 || req.includeStorageTypes[0] != "2" || req.includeStorageTypes[1] != "3" {
			t.Fatalf("includeStorageTypes = %v, want %v", req.includeStorageTypes, []string{"2", "3"})
		}
		if len(req.excludeStorageTypes) != 1 || req.excludeStorageTypes[0] != "5" {
			t.Fatalf("excludeStorageTypes = %v, want %v", req.excludeStorageTypes, []string{"5"})
		}
	default:
		t.Fatal("backend did not receive request")
	}
}

func TestCLIJobCommands_RejectInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "job backup-detail missing job id",
			args: []string{"job", "backup-detail", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk"},
		},
		{
			name: "job detail missing job id",
			args: []string{"job", "detail", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk"},
		},
		{
			name: "job sync-detail missing job id",
			args: []string{"job", "sync-detail", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk"},
		},
		{
			name: "job logs missing job id",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk"},
		},
		{
			name: "job logs rejects negative index",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--job-id", "job-1", "--index", "-1"},
		},
		{
			name: "job logs rejects non-positive count",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--job-id", "job-1", "--count", "0"},
		},
		{
			name: "job logs rejects negative start time",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--job-id", "job-1", "--start-time", "-1"},
		},
		{
			name: "job logs rejects negative end time",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--job-id", "job-1", "--end-time", "-1"},
		},
		{
			name: "job logs rejects invalid level",
			args: []string{"job", "logs", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--job-id", "job-1", "--level", "5"},
		},
		{
			name: "job list rejects invalid status",
			args: []string{"job", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--statuses", "999"},
		},
		{
			name: "job list rejects negative start time",
			args: []string{"job", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--start-time", "-1"},
		},
		{
			name: "job list rejects invalid sort",
			args: []string{"job", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--sort", "createdAt"},
		},
		{
			name: "job list rejects invalid task continue type",
			args: []string{"job", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--task-continue-type", "2"},
		},
		{
			name: "job list rejects invalid operation type",
			args: []string{"job", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--operation-types", "abc"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requestCount int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				http.Error(w, "backend should not be called", http.StatusInternalServerError)
			}))
			defer srv.Close()

			args := append([]string{}, tc.args...)
			for i := range args {
				if args[i] == "http://example.invalid" {
					args[i] = srv.URL
				}
			}

			stdout, stderr, exitCode := runCLI(t, srv, args)
			if exitCode != clierrors.ExitInvalidArgs {
				t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, clierrors.ExitInvalidArgs, stderr)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, "Cli.InvalidArgument") {
				t.Fatalf("stderr = %q, want substring %q", stderr, "Cli.InvalidArgument")
			}
			if got := atomic.LoadInt32(&requestCount); got != 0 {
				t.Fatalf("backend request count = %d, want 0", got)
			}
		})
	}
}

func TestCLIHostCommands_RejectInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "host list rejects invalid client os filter",
			args: []string{"host", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--client-os-filter", "9"},
		},
		{
			name: "host list rejects invalid is-child",
			args: []string{"host", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--is-child", "maybe"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requestCount int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				http.Error(w, "backend should not be called", http.StatusInternalServerError)
			}))
			defer srv.Close()

			args := append([]string{}, tc.args...)
			for i := range args {
				if args[i] == "http://example.invalid" {
					args[i] = srv.URL
				}
			}

			stdout, stderr, exitCode := runCLI(t, srv, args)
			if exitCode != clierrors.ExitInvalidArgs {
				t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, clierrors.ExitInvalidArgs, stderr)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, "Cli.InvalidArgument") {
				t.Fatalf("stderr = %q, want substring %q", stderr, "Cli.InvalidArgument")
			}
			if got := atomic.LoadInt32(&requestCount); got != 0 {
				t.Fatalf("backend request count = %d, want 0", got)
			}
		})
	}
}

func TestCLIHostObjectCommands_RejectInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "host object list rejects count below range",
			args: []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--count", "0"},
		},
		{
			name: "host object list rejects index above range",
			args: []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--index", "20001"},
		},
		{
			name: "host object list rejects invalid datasource type",
			args: []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--datasource-type", "9"},
		},
		{
			name: "host object list rejects invalid include tenant id",
			args: []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--is-include-tenant-id", "maybe"},
		},
		{
			name: "host object list rejects invalid object id length",
			args: []string{"host", "object", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--object-id", "short"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requestCount int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				http.Error(w, "backend should not be called", http.StatusInternalServerError)
			}))
			defer srv.Close()

			args := append([]string{}, tc.args...)
			for i := range args {
				if args[i] == "http://example.invalid" {
					args[i] = srv.URL
				}
			}

			stdout, stderr, exitCode := runCLI(t, srv, args)
			if exitCode != clierrors.ExitInvalidArgs {
				t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, clierrors.ExitInvalidArgs, stderr)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, "Cli.InvalidArgument") {
				t.Fatalf("stderr = %q, want substring %q", stderr, "Cli.InvalidArgument")
			}
			if got := atomic.LoadInt32(&requestCount); got != 0 {
				t.Fatalf("backend request count = %d, want 0", got)
			}
		})
	}
}

func TestCLIMySQLCommands_RejectInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "mysql list rejects invalid sort field",
			args: []string{"mysql", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--sort-field", "9"},
		},
		{
			name: "mysql list rejects invalid bool",
			args: []string{"mysql", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--is-asc-sort", "maybe"},
		},
		{
			name: "mysql list rejects count above range",
			args: []string{"mysql", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--count", "101"},
		},
		{
			name: "mysql list rejects invalid object ids length",
			args: []string{"mysql", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--object-ids", "short"},
		},
		{
			name: "mysql list rejects invalid bind sla status",
			args: []string{"mysql", "list", "--tenant-id", "t", "--endpoint", "http://example.invalid", "--ak", "ak", "--sk", "sk", "--bind-sla-status", "11"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requestCount int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)
				http.Error(w, "backend should not be called", http.StatusInternalServerError)
			}))
			defer srv.Close()

			args := append([]string{}, tc.args...)
			for i := range args {
				if args[i] == "http://example.invalid" {
					args[i] = srv.URL
				}
			}

			stdout, stderr, exitCode := runCLI(t, srv, args)
			if exitCode != clierrors.ExitInvalidArgs {
				t.Fatalf("exitCode = %d, want %d, stderr = %q", exitCode, clierrors.ExitInvalidArgs, stderr)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, "Cli.InvalidArgument") {
				t.Fatalf("stderr = %q, want substring %q", stderr, "Cli.InvalidArgument")
			}
			if got := atomic.LoadInt32(&requestCount); got != 0 {
				t.Fatalf("backend request count = %d, want 0", got)
			}
		})
	}
}

func runCLI(t *testing.T, srv *httptest.Server, args []string) (string, string, int) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	streams := output.Streams{Stdout: stdout, Stderr: stderr}
	root := cmd.NewRootCommand(cmd.Dependencies{
		Streams: streams,
		Version: cmd.BuildVersionInfo("dev"),
		Console: meta.ExecutorFunc(func(ctx context.Context, rt *meta.Runtime, req console.RequestSpec) ([]byte, error) {
			exec := console.Executor{
				Signer: signer.Signer{},
				Client: transport.Client{Doer: srv.Client()},
			}
			return exec.Execute(ctx, rt.Remote, req)
		}),
	})
	root.SetArgs(args)

	err := root.Execute()
	if err != nil {
		_ = output.WriteError(streams, err)
	}
	return stdout.String(), stderr.String(), clierrors.ExitCodeOf(err)
}
