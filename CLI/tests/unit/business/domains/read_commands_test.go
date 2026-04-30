//go:build legacy_cli_contract
// +build legacy_cli_contract

package domains_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/api"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/host"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/job"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/mysql"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/object"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/sla"
)

func TestReadCommandMappings_AreFixed(t *testing.T) {
	slaCmds := sla.Commands()
	if len(slaCmds) != 4 {
		t.Fatalf("sla commands count = %d, want 4", len(slaCmds))
	}

	slaListMeta := slaCmds[1]
	objCmds := object.Commands()
	if len(objCmds) != 5 {
		t.Fatalf("object commands count = %d, want 5", len(objCmds))
	}
	objSlaBindListMeta := commandByPath(t, objCmds, "object", "sla", "bind-list")
	objTimepointMeta := commandByPath(t, objCmds, "object", "timepoint", "list")

	spec, err := slaListMeta.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/templates" || !spec.ReadOnly {
		t.Fatalf("sla list spec = %#v", spec)
	}

	rt := meta.NewRuntime()
	rt.SetInt("index", 2)
	rt.SetInt("count", 20)
	rt.SetString("filter", "oracle")
	rt.SetString("validate-status", "3")
	rt.SetString("disable-mark", "0")
	rt.SetString("type", "Database")
	rt.SetString("copy-mode", "2")
	rt.SetString("backup-mode", "5")
	rt.BindStringSliceFlag(newTestStringSliceBinder("types", []string{"Database", "Fileset"}), "types", nil, "types")
	spec, err = slaListMeta.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/api/sla/v1/templates" {
		t.Fatalf("sla list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("sla list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("sla list count = %q, want %q", got, "20")
	}
	if got := query.Get("filter"); got != "oracle" {
		t.Fatalf("sla list filter = %q, want %q", got, "oracle")
	}
	if got := query.Get("validateStatus"); got != "3" {
		t.Fatalf("sla list validateStatus = %q, want %q", got, "3")
	}
	if got := query.Get("disableMark"); got != "0" {
		t.Fatalf("sla list disableMark = %q, want %q", got, "0")
	}
	if got := query.Get("type"); got != "Database" {
		t.Fatalf("sla list type = %q, want %q", got, "Database")
	}
	if got := query.Get("copyMode"); got != "2" {
		t.Fatalf("sla list copyMode = %q, want %q", got, "2")
	}
	if got := query.Get("backupMode"); got != "5" {
		t.Fatalf("sla list backupMode = %q, want %q", got, "5")
	}
	if got := query["types"]; len(got) != 2 || got[0] != "Database" || got[1] != "Fileset" {
		t.Fatalf("sla list types = %v, want %v", got, []string{"Database", "Fileset"})
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = objSlaBindListMeta.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/group/object/obj-1/templates" || !spec.ReadOnly {
		t.Fatalf("object sla bind-list spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = objTimepointMeta.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/protect_object/obj-1/time_points" {
		t.Fatalf("object spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.SetString("business", "4")
	rt.SetString("start-time", "100")
	rt.SetString("end-time", "200")
	rt.SetString("storage-pool-id", "12345678901234567890123456789012")
	rt.SetString("is-duplication", "false")
	rt.SetString("storage-service-id", "abcdefghijabcdefghijabcdefghijab")
	rt.SetString("data-set-id", "fedcbafedcbafedcbafedcbafedcbafe")
	rt.SetString("expiration-start-time", "300")
	rt.SetString("expiration-end-time", "400")
	rt.SetString("usable", "2")
	rt.SetString("time-point-type", "3")
	rt.BindStringSliceFlag(newTestStringSliceBinder("businesses", []string{"1", "6"}), "businesses", nil, "businesses")
	rt.BindStringSliceFlag(newTestStringSliceBinder("backup-types", []string{"1", "4"}), "backup-types", nil, "backup types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("include-storage-types", []string{"2", "3"}), "include-storage-types", nil, "include storage types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("exclude-storage-types", []string{"5"}), "exclude-storage-types", nil, "exclude storage types")
	spec, err = objTimepointMeta.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/protect_object/obj-1/time_points" {
		t.Fatalf("object timepoint path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("business"); got != "4" {
		t.Fatalf("object timepoint business = %q, want %q", got, "4")
	}
	if got := query.Get("startTime"); got != "100" {
		t.Fatalf("object timepoint startTime = %q, want %q", got, "100")
	}
	if got := query.Get("endTime"); got != "200" {
		t.Fatalf("object timepoint endTime = %q, want %q", got, "200")
	}
	if got := query.Get("storagePoolId"); got != "12345678901234567890123456789012" {
		t.Fatalf("object timepoint storagePoolId = %q", got)
	}
	if got := query.Get("isDuplication"); got != "false" {
		t.Fatalf("object timepoint isDuplication = %q, want %q", got, "false")
	}
	if got := query.Get("storageServiceId"); got != "abcdefghijabcdefghijabcdefghijab" {
		t.Fatalf("object timepoint storageServiceId = %q", got)
	}
	if got := query.Get("dataSetId"); got != "fedcbafedcbafedcbafedcbafedcbafe" {
		t.Fatalf("object timepoint dataSetId = %q", got)
	}
	if got := query.Get("expirationStartTime"); got != "300" {
		t.Fatalf("object timepoint expirationStartTime = %q, want %q", got, "300")
	}
	if got := query.Get("expirationEndTime"); got != "400" {
		t.Fatalf("object timepoint expirationEndTime = %q, want %q", got, "400")
	}
	if got := query.Get("usable"); got != "2" {
		t.Fatalf("object timepoint usable = %q, want %q", got, "2")
	}
	if got := query.Get("timePointType"); got != "3" {
		t.Fatalf("object timepoint timePointType = %q, want %q", got, "3")
	}
	if got := query["businesses"]; len(got) != 2 || got[0] != "1" || got[1] != "6" {
		t.Fatalf("object timepoint businesses = %v, want %v", got, []string{"1", "6"})
	}
	if got := query["backupTypes"]; len(got) != 2 || got[0] != "1" || got[1] != "4" {
		t.Fatalf("object timepoint backupTypes = %v, want %v", got, []string{"1", "4"})
	}
	if got := query["includeStorageTypes"]; len(got) != 2 || got[0] != "2" || got[1] != "3" {
		t.Fatalf("object timepoint includeStorageTypes = %v, want %v", got, []string{"2", "3"})
	}
	if got := query["excludeStorageTypes"]; len(got) != 1 || got[0] != "5" {
		t.Fatalf("object timepoint excludeStorageTypes = %v, want %v", got, []string{"5"})
	}
}

func TestHostReadCommandMappings_AreFixed(t *testing.T) {
	hostCmds := host.Commands()
	if len(hostCmds) != 4 {
		t.Fatalf("host commands count = %d, want 4", len(hostCmds))
	}

	hostList := hostCmds[0]
	spec, err := hostList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("host list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/resource_center/v1/host/group/host_list" || !spec.ReadOnly {
		t.Fatalf("host list spec = %#v", spec)
	}

	rt := meta.NewRuntime()
	rt.SetString("group-id", "group-1")
	rt.SetString("filter", "oracle")
	rt.SetString("is-child", "false")
	rt.BindStringSliceFlag(newTestStringSliceBinder("runner-types", []string{"VMBackup", "NasBackup"}), "runner-types", nil, "runner types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("client-os-filter", []string{"1", "7"}), "client-os-filter", nil, "client os filter")
	spec, err = hostList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("host list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/resource_center/v1/host/group/host_list" {
		t.Fatalf("host list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("groupId"); got != "group-1" {
		t.Fatalf("host list groupId = %q, want %q", got, "group-1")
	}
	if got := query.Get("filter"); got != "oracle" {
		t.Fatalf("host list filter = %q, want %q", got, "oracle")
	}
	if got := query.Get("isChild"); got != "false" {
		t.Fatalf("host list isChild = %q, want %q", got, "false")
	}
	if got := query["runnerTypes"]; len(got) != 2 || got[0] != "VMBackup" || got[1] != "NasBackup" {
		t.Fatalf("host list runnerTypes = %v, want %v", got, []string{"VMBackup", "NasBackup"})
	}
	if got := query["clientOsFilter"]; len(got) != 2 || got[0] != "1" || got[1] != "7" {
		t.Fatalf("host list clientOsFilter = %v, want %v", got, []string{"1", "7"})
	}

	hostObjectList := hostCmds[2]
	spec, err = hostObjectList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("host object list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/file/object_list" {
		t.Fatalf("host object list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("host object list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "10" {
		t.Fatalf("host object list count = %q, want %q", got, "10")
	}

	rt = meta.NewRuntime()
	rt.SetInt("index", 2)
	rt.SetInt("count", 20)
	rt.SetString("production-system-id", "12345678901234567890123456789012")
	rt.SetString("name", "fileset-a")
	rt.SetString("datasource-type", "5")
	rt.SetString("intelligent-archive", "1")
	rt.SetString("exec-status", "4")
	rt.SetString("is-include-tenant-id", "false")
	rt.SetString("query-tenant-id", "tenant-1")
	rt.SetString("group-id", "group-1")
	rt.SetString("has-tp", "1")
	rt.SetString("is-backup-config", "2")
	rt.SetString("object-id", "abcdefghijklmnopqrstuvwx12345678")
	spec, err = hostObjectList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("host object list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/file/object_list" {
		t.Fatalf("host object list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("host object list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("host object list count = %q, want %q", got, "20")
	}
	if got := query.Get("productionSystemId"); got != "12345678901234567890123456789012" {
		t.Fatalf("host object list productionSystemId = %q", got)
	}
	if got := query.Get("name"); got != "fileset-a" {
		t.Fatalf("host object list name = %q, want %q", got, "fileset-a")
	}
	if got := query.Get("datasourceType"); got != "5" {
		t.Fatalf("host object list datasourceType = %q, want %q", got, "5")
	}
	if got := query.Get("intelligentArchive"); got != "1" {
		t.Fatalf("host object list intelligentArchive = %q, want %q", got, "1")
	}
	if got := query.Get("execStatus"); got != "4" {
		t.Fatalf("host object list execStatus = %q, want %q", got, "4")
	}
	if got := query.Get("isIncludeTenantId"); got != "false" {
		t.Fatalf("host object list isIncludeTenantId = %q, want %q", got, "false")
	}
	if got := query.Get("tenantId"); got != "tenant-1" {
		t.Fatalf("host object list tenantId = %q, want %q", got, "tenant-1")
	}
	if got := query.Get("groupId"); got != "group-1" {
		t.Fatalf("host object list groupId = %q, want %q", got, "group-1")
	}
	if got := query.Get("hasTp"); got != "1" {
		t.Fatalf("host object list hasTp = %q, want %q", got, "1")
	}
	if got := query.Get("isBackupConfig"); got != "2" {
		t.Fatalf("host object list isBackupConfig = %q, want %q", got, "2")
	}
	if got := query.Get("objectId"); got != "abcdefghijklmnopqrstuvwx12345678" {
		t.Fatalf("host object list objectId = %q", got)
	}
}

func TestMySQLReadCommandMappings_AreFixed(t *testing.T) {
	mysqlCmds := mysql.Commands()
	if len(mysqlCmds) != 3 {
		t.Fatalf("mysql commands count = %d, want 3", len(mysqlCmds))
	}

	mysqlList := mysqlCmds[0]
	spec, err := mysqlList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("mysql list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/database/object_list" {
		t.Fatalf("mysql list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("mysql list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("mysql list count = %q, want %q", got, "20")
	}

	rt := meta.NewRuntime()
	rt.SetString("is-asc-sort", "false")
	rt.SetString("sort-field", "4")
	rt.SetInt("index", 2)
	rt.SetInt("count", 30)
	rt.SetString("filter", "prod")
	rt.SetString("included-path", "true")
	rt.SetString("group-id", "default")
	rt.SetString("host-id", "12345678901234567890123456789012")
	rt.SetString("empty-host-id", "false")
	rt.SetString("host-tag", "abcdefghijklmnopqrstuvwx12345678")
	rt.SetString("object-status", "2")
	rt.SetString("object-mode", "1")
	rt.SetString("can-backup", "2")
	rt.SetString("exec-status", "4")
	rt.SetString("all-child", "true")
	rt.SetString("name-filter", "mysql-a")
	rt.SetString("app-type", "42")
	rt.SetString("isolation", "false")
	rt.SetString("object-id", "object-single")
	rt.BindStringSliceFlag(newTestStringSliceBinder("production-system-ids", []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}), "production-system-ids", nil, "production system ids")
	rt.BindStringSliceFlag(newTestStringSliceBinder("object-ids", []string{"11111111111111111111111111111111", "22222222222222222222222222222222"}), "object-ids", nil, "object ids")
	rt.BindStringSliceFlag(newTestStringSliceBinder("object-types", []string{"instance", "cluster"}), "object-types", nil, "object types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("app-types", []string{"MySQL", "MariaDB"}), "app-types", nil, "app types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("is-config", []string{"1", "3"}), "is-config", nil, "is config")
	rt.BindStringSliceFlag(newTestStringSliceBinder("bind-sla-status", []string{"1", "6"}), "bind-sla-status", nil, "bind sla status")
	rt.BindStringSliceFlag(newTestStringSliceBinder("included-bind-sla", []string{"2", "5"}), "included-bind-sla", nil, "included bind sla")
	rt.BindStringSliceFlag(newTestStringSliceBinder("exclude-bind-sla", []string{"1", "4"}), "exclude-bind-sla", nil, "exclude bind sla")
	rt.BindStringSliceFlag(newTestStringSliceBinder("bind-sla-ids", []string{"cccccccccccccccccccccccccccccccc", "dddddddddddddddddddddddddddddddd"}), "bind-sla-ids", nil, "bind sla ids")
	rt.BindStringSliceFlag(newTestStringSliceBinder("protect-status", []string{"1", "4"}), "protect-status", nil, "protect status")
	rt.BindStringSliceFlag(newTestStringSliceBinder("last-backup-status", []string{"600", "900"}), "last-backup-status", nil, "last backup status")
	rt.BindStringSliceFlag(newTestStringSliceBinder("last-snapshot-status", []string{"800", "1100"}), "last-snapshot-status", nil, "last snapshot status")
	rt.BindStringSliceFlag(newTestStringSliceBinder("parent-ids", []string{"parent-a", "parent-b"}), "parent-ids", nil, "parent ids")
	spec, err = mysqlList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/database/object_list" {
		t.Fatalf("mysql list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("isAscSort"); got != "false" {
		t.Fatalf("mysql list isAscSort = %q, want %q", got, "false")
	}
	if got := query.Get("sortField"); got != "4" {
		t.Fatalf("mysql list sortField = %q, want %q", got, "4")
	}
	if got := query.Get("index"); got != "2" {
		t.Fatalf("mysql list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("mysql list count = %q, want %q", got, "30")
	}
	if got := query.Get("filter"); got != "prod" {
		t.Fatalf("mysql list filter = %q, want %q", got, "prod")
	}
	if got := query.Get("includedPath"); got != "true" {
		t.Fatalf("mysql list includedPath = %q, want %q", got, "true")
	}
	if got := query.Get("groupId"); got != "default" {
		t.Fatalf("mysql list groupId = %q, want %q", got, "default")
	}
	if got := query.Get("hostId"); got != "12345678901234567890123456789012" {
		t.Fatalf("mysql list hostId = %q", got)
	}
	if got := query.Get("emptyHostId"); got != "false" {
		t.Fatalf("mysql list emptyHostId = %q, want %q", got, "false")
	}
	if got := query.Get("hostTag"); got != "abcdefghijklmnopqrstuvwx12345678" {
		t.Fatalf("mysql list hostTag = %q", got)
	}
	if got := query.Get("objectStatus"); got != "2" {
		t.Fatalf("mysql list objectStatus = %q, want %q", got, "2")
	}
	if got := query.Get("objectMode"); got != "1" {
		t.Fatalf("mysql list objectMode = %q, want %q", got, "1")
	}
	if got := query.Get("canBackup"); got != "2" {
		t.Fatalf("mysql list canBackup = %q, want %q", got, "2")
	}
	if got := query.Get("execStatus"); got != "4" {
		t.Fatalf("mysql list execStatus = %q, want %q", got, "4")
	}
	if got := query.Get("allChild"); got != "true" {
		t.Fatalf("mysql list allChild = %q, want %q", got, "true")
	}
	if got := query.Get("nameFilter"); got != "mysql-a" {
		t.Fatalf("mysql list nameFilter = %q, want %q", got, "mysql-a")
	}
	if got := query.Get("appType"); got != "42" {
		t.Fatalf("mysql list appType = %q, want %q", got, "42")
	}
	if got := query.Get("isolation"); got != "false" {
		t.Fatalf("mysql list isolation = %q, want %q", got, "false")
	}
	if got := query.Get("objectId"); got != "object-single" {
		t.Fatalf("mysql list objectId = %q, want %q", got, "object-single")
	}
	if got := query["productionSystemIds"]; len(got) != 2 || got[0] != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" || got[1] != "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Fatalf("mysql list productionSystemIds = %v", got)
	}
	if got := query["objectIds"]; len(got) != 2 || got[0] != "11111111111111111111111111111111" || got[1] != "22222222222222222222222222222222" {
		t.Fatalf("mysql list objectIds = %v", got)
	}
	if got := query["objectTypes"]; len(got) != 2 || got[0] != "instance" || got[1] != "cluster" {
		t.Fatalf("mysql list objectTypes = %v", got)
	}
	if got := query["appTypes"]; len(got) != 2 || got[0] != "MySQL" || got[1] != "MariaDB" {
		t.Fatalf("mysql list appTypes = %v", got)
	}
	if got := query["isConfig"]; len(got) != 2 || got[0] != "1" || got[1] != "3" {
		t.Fatalf("mysql list isConfig = %v", got)
	}
	if got := query["bindSlaStatus"]; len(got) != 2 || got[0] != "1" || got[1] != "6" {
		t.Fatalf("mysql list bindSlaStatus = %v", got)
	}
	if got := query["includedBindSla"]; len(got) != 2 || got[0] != "2" || got[1] != "5" {
		t.Fatalf("mysql list includedBindSla = %v", got)
	}
	if got := query["excludeBindSla"]; len(got) != 2 || got[0] != "1" || got[1] != "4" {
		t.Fatalf("mysql list excludeBindSla = %v", got)
	}
	if got := query["bindSlaIds"]; len(got) != 2 || got[0] != "cccccccccccccccccccccccccccccccc" || got[1] != "dddddddddddddddddddddddddddddddd" {
		t.Fatalf("mysql list bindSlaIds = %v", got)
	}
	if got := query["protectStatus"]; len(got) != 2 || got[0] != "1" || got[1] != "4" {
		t.Fatalf("mysql list protectStatus = %v", got)
	}
	if got := query["lastBackupStatus"]; len(got) != 2 || got[0] != "600" || got[1] != "900" {
		t.Fatalf("mysql list lastBackupStatus = %v", got)
	}
	if got := query["lastSnapshotStatus"]; len(got) != 2 || got[0] != "800" || got[1] != "1100" {
		t.Fatalf("mysql list lastSnapshotStatus = %v", got)
	}
	if got := query["parentIds"]; len(got) != 2 || got[0] != "parent-a" || got[1] != "parent-b" {
		t.Fatalf("mysql list parentIds = %v", got)
	}
}

func TestJobCommandMappings_AreFixed(t *testing.T) {
	jobCmds := job.Commands()
	if len(jobCmds) != 7 {
		t.Fatalf("job commands count = %d, want 7", len(jobCmds))
	}

	jobList := commandByPath(t, jobCmds, "job", "list")
	jobBackupDetail := commandByPath(t, jobCmds, "job", "backup-detail")
	jobDetail := commandByPath(t, jobCmds, "job", "detail")
	jobSyncDetail := commandByPath(t, jobCmds, "job", "sync-detail")
	jobLogs := commandByPath(t, jobCmds, "job", "logs")
	jobBusinessTypes := commandByPath(t, jobCmds, "job", "business-types")
	jobAppTypes := commandByPath(t, jobCmds, "job", "app-types")

	spec, err := jobList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/jobs" || !spec.ReadOnly {
		t.Fatalf("job list spec = %#v", spec)
	}

	rt := meta.NewRuntime()
	rt.SetInt("index", 2)
	rt.SetInt("count", 20)
	rt.SetString("object-id", "object-1")
	rt.SetString("object-name", "prod-db")
	rt.SetString("start-time", "100")
	rt.SetString("end-time", "200")
	rt.SetString("sort", "-startTime")
	rt.SetString("remark", "nightly backup")
	rt.SetString("cluster-id", "cluster-1")
	rt.SetString("storage-service-id", "storage-service-1")
	rt.SetString("storage-pool-id", "storage-pool-1")
	rt.SetString("strategy-name", "gold")
	rt.SetString("client-name", "client-a")
	rt.SetString("task-continue-type", "1")
	rt.BindStringSliceFlag(newTestStringSliceBinder("statuses", []string{"100", "400"}), "statuses", nil, "statuses")
	rt.BindStringSliceFlag(newTestStringSliceBinder("app-types", []string{"Oracle", "VMware"}), "app-types", nil, "app types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("operation-types", []string{"1", "2"}), "operation-types", nil, "operation types")
	rt.BindStringSliceFlag(newTestStringSliceBinder("business-types", []string{"10", "20"}), "business-types", nil, "business types")
	spec, err = jobList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/job_center/v1/jobs" {
		t.Fatalf("job list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("job list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("job list count = %q, want %q", got, "20")
	}
	if got := query.Get("objectId"); got != "object-1" {
		t.Fatalf("job list objectId = %q, want %q", got, "object-1")
	}
	if got := query.Get("objectName"); got != "prod-db" {
		t.Fatalf("job list objectName = %q, want %q", got, "prod-db")
	}
	if got := query.Get("startTime"); got != "100" {
		t.Fatalf("job list startTime = %q, want %q", got, "100")
	}
	if got := query.Get("endTime"); got != "200" {
		t.Fatalf("job list endTime = %q, want %q", got, "200")
	}
	if got := query.Get("sort"); got != "-startTime" {
		t.Fatalf("job list sort = %q, want %q", got, "-startTime")
	}
	if got := query.Get("remark"); got != "nightly backup" {
		t.Fatalf("job list remark = %q, want %q", got, "nightly backup")
	}
	if got := query.Get("clusterId"); got != "cluster-1" {
		t.Fatalf("job list clusterId = %q, want %q", got, "cluster-1")
	}
	if got := query.Get("storageServiceId"); got != "storage-service-1" {
		t.Fatalf("job list storageServiceId = %q, want %q", got, "storage-service-1")
	}
	if got := query.Get("storagePoolId"); got != "storage-pool-1" {
		t.Fatalf("job list storagePoolId = %q, want %q", got, "storage-pool-1")
	}
	if got := query.Get("strategyName"); got != "gold" {
		t.Fatalf("job list strategyName = %q, want %q", got, "gold")
	}
	if got := query.Get("clientName"); got != "client-a" {
		t.Fatalf("job list clientName = %q, want %q", got, "client-a")
	}
	if got := query.Get("taskContinueType"); got != "1" {
		t.Fatalf("job list taskContinueType = %q, want %q", got, "1")
	}
	if got := query["statuses"]; len(got) != 2 || got[0] != "100" || got[1] != "400" {
		t.Fatalf("job list statuses = %v, want %v", got, []string{"100", "400"})
	}
	if got := query["appTypes"]; len(got) != 2 || got[0] != "Oracle" || got[1] != "VMware" {
		t.Fatalf("job list appTypes = %v, want %v", got, []string{"Oracle", "VMware"})
	}
	if got := query["operationTypes"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("job list operationTypes = %v, want %v", got, []string{"1", "2"})
	}
	if got := query["businessTypes"]; len(got) != 2 || got[0] != "10" || got[1] != "20" {
		t.Fatalf("job list businessTypes = %v, want %v", got, []string{"10", "20"})
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-9")
	spec, err = jobBackupDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job backup-detail BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/task/job-9/detail" {
		t.Fatalf("job backup-detail path = %q", spec.Path)
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-9")
	spec, err = jobDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/job/job-9" || !spec.ReadOnly {
		t.Fatalf("job detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-9")
	spec, err = jobSyncDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job sync-detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/sync_task/job-9/detail" || !spec.ReadOnly {
		t.Fatalf("job sync-detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-9")
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	rt.SetString("start-time", "100")
	rt.SetString("end-time", "200")
	rt.SetString("level", "3")
	spec, err = jobLogs.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job logs BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || !spec.ReadOnly {
		t.Fatalf("job logs spec = %#v", spec)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/job_center/v1/activity/job-9/logs" {
		t.Fatalf("job logs path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("index"); got != "0" {
		t.Fatalf("job logs index = %q, want %q", got, "0")
	}
	if got := parsed.Query().Get("count"); got != "30" {
		t.Fatalf("job logs count = %q, want %q", got, "30")
	}
	if got := parsed.Query().Get("startTime"); got != "100" {
		t.Fatalf("job logs startTime = %q, want %q", got, "100")
	}
	if got := parsed.Query().Get("endTime"); got != "200" {
		t.Fatalf("job logs endTime = %q, want %q", got, "200")
	}
	if got := parsed.Query().Get("level"); got != "3" {
		t.Fatalf("job logs level = %q, want %q", got, "3")
	}

	spec, err = jobBusinessTypes.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job business-types BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/business_types" || !spec.ReadOnly {
		t.Fatalf("job business-types spec = %#v", spec)
	}

	spec, err = jobAppTypes.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job app-types BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/app_types" || !spec.ReadOnly {
		t.Fatalf("job app-types spec = %#v", spec)
	}
}

func TestAPICommand_RejectsAbsoluteURL(t *testing.T) {
	def := api.Commands()[0]
	rt := meta.NewRuntime()
	rt.SetString("method", http.MethodGet)
	rt.SetString("path", "https://backend.example.com/api")

	if err := def.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want invalid path")
	}
}

func TestAPICommand_RejectsInvalidMethod(t *testing.T) {
	def := api.Commands()[0]
	rt := meta.NewRuntime()
	rt.SetString("method", "BAD METHOD")
	rt.SetString("path", "/api/sla/v1/templates")

	if err := def.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want invalid method")
	}
}

func commandByPath(t *testing.T, defs []meta.CommandMeta, path ...string) meta.CommandMeta {
	t.Helper()

	for _, def := range defs {
		if len(def.CanonicalPath) != len(path) {
			continue
		}
		match := true
		for i := range path {
			if def.CanonicalPath[i] != path[i] {
				match = false
				break
			}
		}
		if match {
			return def
		}
	}
	t.Fatalf("command path %v not found", path)
	return meta.CommandMeta{}
}

func parseSpecURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", raw, err)
	}
	return parsed
}

type testStringSliceBinder struct {
	name   string
	values []string
}

func newTestStringSliceBinder(name string, values []string) *testStringSliceBinder {
	return &testStringSliceBinder{name: name, values: values}
}

func (b *testStringSliceBinder) StringSliceVar(p *[]string, name string, value []string, usage string) {
	if name != b.name {
		*p = value
		return
	}
	*p = append([]string(nil), b.values...)
}
