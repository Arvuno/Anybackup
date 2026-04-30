package domains_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/api"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/client"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/host"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/job"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/mysql"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/network"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/policy"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/protect"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/storage"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/timepoint"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/vmware"
)

func TestDomainCommandCounts(t *testing.T) {
	if got := len(policy.Commands()); got != 6 {
		t.Fatalf("policy commands count = %d, want 6", got)
	}
	if got := len(protect.Commands()); got != 12 {
		t.Fatalf("protect commands count = %d, want 12", got)
	}
	if got := len(timepoint.Commands()); got != 2 {
		t.Fatalf("timepoint commands count = %d, want 2", got)
	}
	if got := len(vmware.Commands()); got != 9 {
		t.Fatalf("vmware commands count = %d, want 9", got)
	}
	if got := len(mysql.Commands()); got != 14 {
		t.Fatalf("mysql commands count = %d, want 14", got)
	}
	if got := len(network.Commands()); got != 5 {
		t.Fatalf("network commands count = %d, want 5", got)
	}
	if got := len(storage.Commands()); got != 6 {
		t.Fatalf("storage commands count = %d, want 6", got)
	}
	if got := len(host.Commands()); got != 7 {
		t.Fatalf("host commands count = %d, want 7", got)
	}
	if got := len(job.Commands()); got != 9 {
		t.Fatalf("job commands count = %d, want 9", got)
	}
	if got := len(client.Commands()); got != 9 {
		t.Fatalf("client commands count = %d, want 9", got)
	}
	if got := len(api.Commands()); got != 1 {
		t.Fatalf("api commands count = %d, want 1", got)
	}
}

func TestPolicyWriteMappings(t *testing.T) {
	policyList := commandByPath(t, policy.Commands(), "policy", "list")
	backupDetail := commandByPath(t, policy.Commands(), "policy", "backup", "detail")
	backupCreate := commandByPath(t, policy.Commands(), "policy", "backup", "create")
	copyDetail := commandByPath(t, policy.Commands(), "policy", "copy", "detail")
	copyCreate := commandByPath(t, policy.Commands(), "policy", "copy", "create")
	deletePolicy := commandByPath(t, policy.Commands(), "policy", "delete")

	spec, err := policyList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("policy list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/templates" || !spec.ReadOnly {
		t.Fatalf("policy list spec = %#v", spec)
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
	spec, err = policyList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/api/sla/v1/templates" {
		t.Fatalf("policy list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("policy list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("policy list count = %q, want %q", got, "20")
	}
	if got := query.Get("filter"); got != "oracle" {
		t.Fatalf("policy list filter = %q, want %q", got, "oracle")
	}
	if got := query.Get("validateStatus"); got != "3" {
		t.Fatalf("policy list validateStatus = %q, want %q", got, "3")
	}
	if got := query.Get("disableMark"); got != "0" {
		t.Fatalf("policy list disableMark = %q, want %q", got, "0")
	}
	if got := query.Get("type"); got != "Database" {
		t.Fatalf("policy list type = %q, want %q", got, "Database")
	}
	if got := query.Get("copyMode"); got != "2" {
		t.Fatalf("policy list copyMode = %q, want %q", got, "2")
	}
	if got := query.Get("backupMode"); got != "5" {
		t.Fatalf("policy list backupMode = %q, want %q", got, "5")
	}
	if got := query["types"]; len(got) != 2 || got[0] != "Database" || got[1] != "Fileset" {
		t.Fatalf("policy list types = %v, want %v", got, []string{"Database", "Fileset"})
	}

	rt = meta.NewRuntime()
	if err := backupDetail.Validate(rt); err == nil {
		t.Fatal("policy backup detail Validate() error = nil, want missing --group-id")
	}
	rt.SetString("group-id", "group-1")
	spec, err = backupDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy backup detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/group/group-1/backup_detail" || !spec.ReadOnly {
		t.Fatalf("policy backup detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"name":"p1"}`)
	spec, err = backupCreate.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy backup create BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/api/sla/v1/group/backup_info" || spec.ReadOnly {
		t.Fatalf("policy backup create spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	if err := copyDetail.Validate(rt); err == nil {
		t.Fatal("policy copy detail Validate() error = nil, want missing --group-id")
	}
	rt.SetString("group-id", "group-1")
	spec, err = copyDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy copy detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/group/group-1/copy_detail" || !spec.ReadOnly {
		t.Fatalf("policy copy detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"name":"copy-1"}`)
	spec, err = copyCreate.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy copy create BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/api/sla/v1/group/copy_info" || spec.ReadOnly {
		t.Fatalf("policy copy create spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"groupIds":["a"]}`)
	spec, err = deletePolicy.BuildRequest(rt)
	if err != nil {
		t.Fatalf("policy delete BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/api/sla/v1/groups" || spec.ReadOnly {
		t.Fatalf("policy delete spec = %#v", spec)
	}
}

func TestProtectMappings(t *testing.T) {
	cmds := protect.Commands()

	policyBind := commandByPath(t, cmds, "protect", "policy", "bind")
	policyBindBatch := commandByPath(t, cmds, "protect", "policy", "bind-batch")
	policyBindList := commandByPath(t, cmds, "protect", "policy", "bind-list")
	policyBatchUnbind := commandByPath(t, cmds, "protect", "policy", "batch-unbind")
	backupStart := commandByPath(t, cmds, "protect", "backup", "start")
	backupBatchStart := commandByPath(t, cmds, "protect", "backup", "batch-start")

	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"slaId":"s1"}`)
	spec, err := policyBind.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect policy bind BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/obj-1/slas" || spec.ReadOnly {
		t.Fatalf("protect policy bind spec = %#v", spec)
	}
	var bindBody map[string]any
	if err := json.Unmarshal(spec.Body, &bindBody); err != nil {
		t.Fatalf("protect policy bind body json error = %v", err)
	}
	if got, _ := bindBody["objectId"].(string); got != "obj-1" {
		t.Fatalf("protect policy bind body objectId = %q, want %q", got, "obj-1")
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"objectId":"obj-2","slaIds":["s1"]}`)
	if _, err := policyBind.BuildRequest(rt); err == nil {
		t.Fatal("protect policy bind BuildRequest() error = nil, want objectId mismatch")
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectIds":["obj-1"],"slaId":"s1"}`)
	spec, err = policyBindBatch.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect policy bind-batch BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/slas" || spec.ReadOnly {
		t.Fatalf("protect policy bind-batch spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = policyBindList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect policy bind-list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/api/sla/v1/group/object/obj-1/templates" || !spec.ReadOnly {
		t.Fatalf("protect policy bind-list spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.SetInt("index", 2)
	rt.SetInt("count", 20)
	rt.SetString("filter", "oracle")
	rt.SetString("validate-status", "3")
	rt.SetString("disable-mark", "0")
	rt.SetString("type", "Database")
	spec, err = policyBindList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect policy bind-list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/api/sla/v1/group/object/obj-1/templates" {
		t.Fatalf("protect policy bind-list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("protect policy bind-list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("protect policy bind-list count = %q, want %q", got, "20")
	}
	if got := query.Get("filter"); got != "oracle" {
		t.Fatalf("protect policy bind-list filter = %q, want %q", got, "oracle")
	}
	if got := query.Get("validateStatus"); got != "3" {
		t.Fatalf("protect policy bind-list validateStatus = %q, want %q", got, "3")
	}
	if got := query.Get("disableMark"); got != "0" {
		t.Fatalf("protect policy bind-list disableMark = %q, want %q", got, "0")
	}
	if got := query.Get("type"); got != "Database" {
		t.Fatalf("protect policy bind-list type = %q, want %q", got, "Database")
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectIds":["obj-1"]}`)
	spec, err = policyBatchUnbind.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect policy batch-unbind BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/backupmgm/v1/protect_object/slas" || spec.ReadOnly {
		t.Fatalf("protect policy batch-unbind spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"force":true}`)
	spec, err = backupStart.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect backup start BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/obj-1/backup_task/start" || spec.ReadOnly {
		t.Fatalf("protect backup start spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectIds":["obj-1"]}`)
	spec, err = backupBatchStart.BuildRequest(rt)
	if err != nil {
		t.Fatalf("protect backup batch-start BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/batch/backup_tasks/start" || spec.ReadOnly {
		t.Fatalf("protect backup batch-start spec = %#v", spec)
	}

}

func TestTimepointMappings(t *testing.T) {
	protectCmds := protect.Commands()
	assertCommandMissing(t, protectCmds, "protect", "timepoint", "list")
	assertCommandMissing(t, protectCmds, "protect", "clean", "start")

	cmds := timepoint.Commands()
	timepointList := commandByPath(t, cmds, "timepoint", "list")
	cleanStart := commandByPath(t, cmds, "timepoint", "clean", "start")

	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.SetString("business", "4")
	rt.SetString("start-time", "100")
	rt.SetString("end-time", "200")
	rt.BindStringSliceFlag(newTestStringSliceBinder("backup-types", []string{"1", "4"}), "backup-types", nil, "backup types")
	spec, err := timepointList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("timepoint list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/protect_object/obj-1/time_points" {
		t.Fatalf("timepoint list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("business"); got != "4" {
		t.Fatalf("timepoint list business = %q, want %q", got, "4")
	}
	if got := query.Get("startTime"); got != "100" {
		t.Fatalf("timepoint list startTime = %q, want %q", got, "100")
	}
	if got := query.Get("endTime"); got != "200" {
		t.Fatalf("timepoint list endTime = %q, want %q", got, "200")
	}
	if got := query["backupTypes"]; len(got) != 2 || got[0] != "1" || got[1] != "4" {
		t.Fatalf("timepoint list backupTypes = %v, want %v", got, []string{"1", "4"})
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"force":true}`)
	spec, err = cleanStart.BuildRequest(rt)
	if err != nil {
		t.Fatalf("timepoint clean start BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/obj-1/clean_task/start" || spec.ReadOnly {
		t.Fatalf("timepoint clean start spec = %#v", spec)
	}
}

func TestTimepointValidate_MissingRequiredFields(t *testing.T) {
	cmds := timepoint.Commands()
	list := commandByPath(t, cmds, "timepoint", "list")
	cleanStart := commandByPath(t, cmds, "timepoint", "clean", "start")

	t.Run("list missing object-id", func(t *testing.T) {
		if err := list.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("list rejects invalid business", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		rt.SetString("business", "2")
		if err := list.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid --business")
		}
	})

	t.Run("clean start missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		if err := cleanStart.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}

func TestHostMySQLVMwareReadMappings(t *testing.T) {
	hostList := commandByPath(t, host.Commands(), "host", "list")
	hostObjectList := commandByPath(t, host.Commands(), "host", "object", "list")
	hostObjectDetail := commandByPath(t, host.Commands(), "host", "object", "detail")
	hostBackupConfigDetail := commandByPath(t, host.Commands(), "host", "backup-config", "detail")
	hostRestoreStart := commandByPath(t, host.Commands(), "host", "restore", "start")
	mysqlObjectList := commandByPath(t, mysql.Commands(), "mysql", "object", "list")
	mysqlObjectGet := commandByPath(t, mysql.Commands(), "mysql", "object", "get")
	mysqlAppBackupConfigDetail := commandByPath(t, mysql.Commands(), "mysql", "backup-config", "detail")
	mysqlAppBackupConfigSet := commandByPath(t, mysql.Commands(), "mysql", "backup-config", "set")
	mysqlDatasourceBackup := commandByPath(t, mysql.Commands(), "mysql", "datasource", "backup")
	mysqlRecoveryRange := commandByPath(t, mysql.Commands(), "mysql", "recovery", "range")
	mysqlRecoveryConfigDetail := commandByPath(t, mysql.Commands(), "mysql", "recovery-config", "detail")
	mysqlDatasourceRecovery := commandByPath(t, mysql.Commands(), "mysql", "datasource", "recovery")
	mysqlTimepointList := commandByPath(t, mysql.Commands(), "mysql", "recovery", "timepoint", "list")
	mysqlBackupDetail := commandByPath(t, mysql.Commands(), "mysql", "backup-detail")
	mysqlRecoveryDetail := commandByPath(t, mysql.Commands(), "mysql", "recovery-detail")
	mysqlAuthorize := commandByPath(t, mysql.Commands(), "mysql", "authorize")
	vmwareObjectList := commandByPath(t, vmware.Commands(), "vmware", "object", "list")
	vmwareObjectInfo := commandByPath(t, vmware.Commands(), "vmware", "object", "info")
	vmwareDatasourceGet := commandByPath(t, vmware.Commands(), "vmware", "datasource", "get")
	vmwareBackupConfigDetail := commandByPath(t, vmware.Commands(), "vmware", "backup-config", "detail")
	vmwareTimepointMetadata := commandByPath(t, vmware.Commands(), "vmware", "timepoint", "metadata")
	vmwareBackupDetail := commandByPath(t, vmware.Commands(), "vmware", "backup-detail")
	vmwareRecoveryDetail := commandByPath(t, vmware.Commands(), "vmware", "recovery-detail")

	spec, err := hostList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("host list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/resource_center/v1/host/group/host_list" || !spec.ReadOnly {
		t.Fatalf("host list spec = %#v", spec)
	}

	hostListRt := meta.NewRuntime()
	hostListRt.SetInt("index", 0)
	hostListRt.SetInt("count", 30)
	hostListRt.MarkChanged("index")
	hostListRt.MarkChanged("count")
	spec, err = hostList.BuildRequest(hostListRt)
	if err != nil {
		t.Fatalf("host list BuildRequest() with paging error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/resource_center/v1/host/group/host_list" {
		t.Fatalf("host list path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("index"); got != "0" {
		t.Fatalf("host list index = %q, want %q", got, "0")
	}
	if got := parsed.Query().Get("count"); got != "30" {
		t.Fatalf("host list count = %q, want %q", got, "30")
	}

	spec, err = hostObjectList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("host object list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/file/object_list" {
		t.Fatalf("host object list path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("count"); got != "10" {
		t.Fatalf("host object list default count = %q, want %q", got, "10")
	}
	var query url.Values

	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = hostObjectDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("host object detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/file/fileset/obj-1" || !spec.ReadOnly {
		t.Fatalf("host object detail spec = %#v", spec)
	}

	spec, err = hostBackupConfigDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("host backup-config detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/file/backup_config/obj-1" || !spec.ReadOnly {
		t.Fatalf("host backup-config detail spec = %#v", spec)
	}
	rt.Body = []byte(`{"target":"f1"}`)
	spec, err = hostRestoreStart.BuildRequest(rt)
	if err != nil {
		t.Fatalf("host restore start BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/file/recovery" || spec.ReadOnly {
		t.Fatalf("host restore start spec = %#v", spec)
	}

	spec, err = mysqlObjectList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("mysql object list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/database/object_list" {
		t.Fatalf("mysql object list path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("count"); got != "20" {
		t.Fatalf("mysql object list default count = %q, want %q", got, "20")
	}
	if got := parsed.Query().Get("appType"); got != "202" {
		t.Fatalf("mysql object list default appType = %q, want %q", got, "202")
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = mysqlObjectGet.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql object get BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/object" || parsed.Query().Get("objectId") != "obj-1" {
		t.Fatalf("mysql object get spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("system-id", "sys-1")
	rt.SetString("object-id", "obj-1")
	spec, err = mysqlAppBackupConfigDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql backup-config detail BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/app_backup_config" || parsed.Query().Get("systemId") != "sys-1" || parsed.Query().Get("objectId") != "obj-1" {
		t.Fatalf("mysql backup-config detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectId":"obj-1","systemId":"sys-1","commonConfigParams":{"backupDestination":{"regionParams":[{"storagePoolParams":[{"storagePoolId":"sp-1"}]}]}},"backupType":2,"backupGran":1}`)
	spec, err = mysqlAppBackupConfigSet.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql backup-config set BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPut || spec.Path != "/backupmgm/v1/mysql/app_backup_config" || spec.ReadOnly {
		t.Fatalf("mysql backup-config set spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetInt("index", 1)
	rt.SetInt("count", 50)
	rt.SetString("system-id", "sys-1")
	rt.SetString("object-id", "obj-1")
	rt.SetString("full-path", "/data/mysql")
	rt.SetString("request-id", "req-1")
	spec, err = mysqlDatasourceBackup.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql datasource backup BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/backup_datasource" {
		t.Fatalf("mysql datasource backup path = %q", parsed.Path)
	}
	query = parsed.Query()
	if query.Get("index") != "1" || query.Get("count") != "50" || query.Get("systemId") != "sys-1" || query.Get("objectId") != "obj-1" || query.Get("fullPath") != "/data/mysql" || query.Get("requestId") != "req-1" {
		t.Fatalf("mysql datasource backup query = %#v", query)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.SetString("storage-service-id", "ss-1")
	rt.SetString("restore-object-type", "2")
	rt.SetString("restore-gran", "3")
	rt.SetString("timestamp", "1710000000000")
	rt.SetString("backup-task-type", "1")
	rt.SetString("storage-pool-id", "sp-1")
	spec, err = mysqlRecoveryRange.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql recovery range BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/get_recovery_range" {
		t.Fatalf("mysql recovery range path = %q", parsed.Path)
	}
	query = parsed.Query()
	if query.Get("objectId") != "obj-1" || query.Get("storageServiceId") != "ss-1" || query.Get("restoreObjectType") != "2" || query.Get("restoreGran") != "3" || query.Get("timestamp") != "1710000000000" || query.Get("backupTaskType") != "1" || query.Get("storagePoolId") != "sp-1" {
		t.Fatalf("mysql recovery range query = %#v", query)
	}

	rt = meta.NewRuntime()
	rt.SetString("system-id", "sys-1")
	rt.SetString("object-id", "obj-1")
	spec, err = mysqlRecoveryConfigDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql recovery-config detail BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/recovery_config" || parsed.Query().Get("systemId") != "sys-1" || parsed.Query().Get("objectId") != "obj-1" {
		t.Fatalf("mysql recovery-config detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("data-set-id", "ds-1")
	rt.SetString("storage-service-id", "ss-1")
	rt.SetString("timestamp", "1710000000000")
	rt.SetString("restore-gran", "2")
	rt.SetString("request-id", "req-1")
	rt.SetString("full-path", "/")
	rt.SetInt("index", 0)
	rt.SetInt("count", 100)
	spec, err = mysqlDatasourceRecovery.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql datasource recovery BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/recovery_datasource" {
		t.Fatalf("mysql datasource recovery path = %q", parsed.Path)
	}
	query = parsed.Query()
	if query.Get("dataSetId") != "ds-1" || query.Get("storageServiceId") != "ss-1" || query.Get("timestamp") != "1710000000000" || query.Get("restoreGran") != "2" || query.Get("requestId") != "req-1" || query.Get("fullPath") != "/" || query.Get("index") != "0" || query.Get("count") != "100" {
		t.Fatalf("mysql datasource recovery query = %#v", query)
	}

	rt = meta.NewRuntime()
	rt.SetInt("index", 2)
	rt.SetInt("count", 30)
	rt.SetString("object-id", "obj-1")
	rt.SetString("storage-service-id", "ss-1")
	rt.SetString("storage-pool-id", "sp-1")
	rt.SetString("start-time", "100")
	rt.SetString("end-time", "200")
	rt.SetString("backup-task-type", "1")
	rt.SetString("restore-gran", "3")
	spec, err = mysqlTimepointList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql timepoint list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/time_points" {
		t.Fatalf("mysql timepoint list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if query.Get("index") != "2" || query.Get("count") != "30" || query.Get("objectId") != "obj-1" || query.Get("storageServiceId") != "ss-1" || query.Get("storagePoolId") != "sp-1" || query.Get("startTime") != "100" || query.Get("endTime") != "200" || query.Get("backupTaskType") != "1" || query.Get("restoreGran") != "3" {
		t.Fatalf("mysql timepoint list query = %#v", query)
	}

	rt = meta.NewRuntime()
	rt.SetString("task-id", "task-1")
	spec, err = mysqlBackupDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql backup-detail BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/get_backup_task_detail" || parsed.Query().Get("taskId") != "task-1" {
		t.Fatalf("mysql backup-detail spec = %#v", spec)
	}

	spec, err = mysqlRecoveryDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql recovery-detail BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/get_recovery_task_detail" || parsed.Query().Get("taskId") != "task-1" {
		t.Fatalf("mysql recovery-detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"instanceName":"inst-1","clientId":"c-1","username":"u-1","password":"p-1","systemId":"s-1","resourceId":"r-1","osUserName":"os-1","type":1}`)
	spec, err = mysqlAuthorize.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql authorize BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/resource_center/v1/databasealone/mysql/authorize" || spec.ReadOnly {
		t.Fatalf("mysql authorize spec = %#v", spec)
	}
	if got := string(spec.Body); got != string(rt.Body) {
		t.Fatalf("mysql authorize body = %q, want %q", got, string(rt.Body))
	}

	rt = meta.NewRuntime()
	rt.SetString("production-system-id", "ps-1")
	spec, err = vmwareObjectList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware object list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/ps-1/get_vms" || !spec.ReadOnly {
		t.Fatalf("vmware object list spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = vmwareObjectInfo.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware object info BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/obj-1/info" || !spec.ReadOnly {
		t.Fatalf("vmware object info spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("production-system-id", "ps-1")
	spec, err = vmwareDatasourceGet.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware datasource get BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/ps-1/sub_objects" || !spec.ReadOnly {
		t.Fatalf("vmware datasource get spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	spec, err = vmwareBackupConfigDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware backup-config detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/obj-1/backup_config" || !spec.ReadOnly {
		t.Fatalf("vmware backup-config detail spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetString("timestamp", "1710000000000")
	rt.SetString("data-set-id", "ds-1")
	rt.SetString("business", "1")
	rt.SetString("request-id", "req-1")
	rt.SetString("storage-service-id", "ss-1")
	spec, err = vmwareTimepointMetadata.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware timepoint metadata BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/virtual/vmware/time_point/get_metadata" {
		t.Fatalf("vmware timepoint metadata path = %q", parsed.Path)
	}
	query = parsed.Query()
	if query.Get("timestamp") != "1710000000000" || query.Get("dataSetId") != "ds-1" || query.Get("business") != "1" || query.Get("requestId") != "req-1" || query.Get("storageServiceId") != "ss-1" {
		t.Fatalf("vmware timepoint metadata query = %#v", query)
	}

	rt = meta.NewRuntime()
	rt.SetString("task-id", "task-1")
	spec, err = vmwareBackupDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware backup-detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/backup_task/task-1/detail" || !spec.ReadOnly {
		t.Fatalf("vmware backup-detail spec = %#v", spec)
	}

	spec, err = vmwareRecoveryDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("vmware recovery-detail BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/backupmgm/v1/virtual/vmware/recovery_task/task-1/detail" || !spec.ReadOnly {
		t.Fatalf("vmware recovery-detail spec = %#v", spec)
	}
}

func TestMySQLObjectListDefaultsAppType202(t *testing.T) {
	mysqlObjectList := commandByPath(t, mysql.Commands(), "mysql", "object", "list")

	spec, err := mysqlObjectList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("mysql object list BuildRequest() error = %v", err)
	}

	parsed := parseSpecURL(t, spec.Path)
	if got := parsed.Query().Get("appType"); got != "202" {
		t.Fatalf("mysql object list default appType = %q, want %q", got, "202")
	}
}

func TestMySQLDatasourceRecoveryBuildsDatasetStyleQuery(t *testing.T) {
	mysqlDatasourceRecovery := commandByPath(t, mysql.Commands(), "mysql", "datasource", "recovery")

	rt := meta.NewRuntime()
	rt.SetString("data-set-id", "ds-1")
	rt.SetString("storage-service-id", "ss-1")
	rt.SetString("timestamp", "1710000000000")
	rt.SetString("restore-gran", "2")
	rt.SetString("request-id", "req-1")
	rt.SetString("full-path", "/")
	rt.SetInt("index", 0)
	rt.SetInt("count", 100)

	spec, err := mysqlDatasourceRecovery.BuildRequest(rt)
	if err != nil {
		t.Fatalf("mysql datasource recovery BuildRequest() error = %v", err)
	}

	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/mysql/recovery_datasource" {
		t.Fatalf("mysql datasource recovery path = %q", parsed.Path)
	}

	query := parsed.Query()
	if query.Get("dataSetId") != "ds-1" || query.Get("storageServiceId") != "ss-1" || query.Get("timestamp") != "1710000000000" || query.Get("restoreGran") != "2" || query.Get("requestId") != "req-1" || query.Get("fullPath") != "/" || query.Get("index") != "0" || query.Get("count") != "100" {
		t.Fatalf("mysql datasource recovery query = %#v", query)
	}
}

func TestMySQLDatasourceRecoveryValidateRequiresDatasetAndTimestamp(t *testing.T) {
	mysqlDatasourceRecovery := commandByPath(t, mysql.Commands(), "mysql", "datasource", "recovery")

	rt := meta.NewRuntime()
	rt.SetString("storage-service-id", "ss-1")
	rt.SetString("timestamp", "1710000000000")
	if err := mysqlDatasourceRecovery.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want missing --data-set-id")
	}

	rt = meta.NewRuntime()
	rt.SetString("data-set-id", "ds-1")
	rt.SetString("storage-service-id", "ss-1")
	if err := mysqlDatasourceRecovery.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want missing --timestamp")
	}
}

func TestMySQLAuthorizeValidate(t *testing.T) {
	authorize := commandByPath(t, mysql.Commands(), "mysql", "authorize")

	newBody := func(m map[string]any) []byte {
		t.Helper()
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}
		return b
	}

	validBody := map[string]any{
		"instanceName": "inst-1",
		"clientId":     "c-1",
		"username":     "u-1",
		"password":     "p-1",
		"systemId":     "s-1",
		"resourceId":   "r-1",
		"osUserName":   "os-1",
		"type":         1,
	}

	cloneBody := func() map[string]any {
		t.Helper()
		body := make(map[string]any, len(validBody))
		for k, v := range validBody {
			body[k] = v
		}
		return body
	}

	t.Run("missing body", func(t *testing.T) {
		if err := authorize.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte("{")
		if err := authorize.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid JSON")
		}
	})

	requiredStringFields := []string{
		"instanceName",
		"clientId",
		"username",
		"password",
		"systemId",
		"resourceId",
		"osUserName",
	}

	for _, field := range requiredStringFields {
		field := field

		t.Run(field+" missing", func(t *testing.T) {
			rt := meta.NewRuntime()
			body := cloneBody()
			delete(body, field)
			rt.Body = newBody(body)
			if err := authorize.Validate(rt); err == nil {
				t.Fatalf("Validate() error = nil, want missing %q", field)
			}
		})

		t.Run(field+" null", func(t *testing.T) {
			rt := meta.NewRuntime()
			body := cloneBody()
			body[field] = nil
			rt.Body = newBody(body)
			if err := authorize.Validate(rt); err == nil {
				t.Fatalf("Validate() error = nil, want null %q", field)
			}
		})

		t.Run(field+" non-string", func(t *testing.T) {
			rt := meta.NewRuntime()
			body := cloneBody()
			body[field] = 1
			rt.Body = newBody(body)
			if err := authorize.Validate(rt); err == nil {
				t.Fatalf("Validate() error = nil, want non-string %q", field)
			}
		})

		t.Run(field+" blank string", func(t *testing.T) {
			rt := meta.NewRuntime()
			body := cloneBody()
			body[field] = "   "
			rt.Body = newBody(body)
			if err := authorize.Validate(rt); err == nil {
				t.Fatalf("Validate() error = nil, want blank %q", field)
			}
		})
	}

	t.Run("type missing", func(t *testing.T) {
		rt := meta.NewRuntime()
		body := cloneBody()
		delete(body, "type")
		rt.Body = newBody(body)
		if err := authorize.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing type")
		}
	})

	t.Run("type null", func(t *testing.T) {
		rt := meta.NewRuntime()
		body := cloneBody()
		body["type"] = nil
		rt.Body = newBody(body)
		if err := authorize.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want null type")
		}
	})

	t.Run("valid body passes", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = newBody(validBody)
		if err := authorize.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v, want nil", err)
		}
	})
}

func TestJobReadMappings(t *testing.T) {
	cmds := job.Commands()
	jobList := commandByPath(t, cmds, "job", "list")
	childList := commandByPath(t, cmds, "job", "child", "list")
	backupDetail := commandByPath(t, cmds, "job", "backup-detail")
	syncDetail := commandByPath(t, cmds, "job", "sync-detail")
	logs := commandByPath(t, cmds, "job", "logs")
	businessTypes := commandByPath(t, cmds, "job", "business-types")
	appTypes := commandByPath(t, cmds, "job", "app-types")

	spec, err := jobList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/jobs" || !spec.ReadOnly {
		t.Fatalf("job list spec = %#v", spec)
	}

	rt := meta.NewRuntime()
	rt.SetString("job-id", "job-1")
	spec, err = childList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job child list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/backupmgm/v1/task/job-1/child" {
		t.Fatalf("job child list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("job child list default index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("job child list default count = %q, want %q", got, "30")
	}
	if got := query.Get("taskId"); got != "job-1" {
		t.Fatalf("job child list taskId = %q, want %q", got, "job-1")
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-1")
	rt.SetInt("index", 2)
	rt.SetInt("count", 50)
	spec, err = childList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job child list BuildRequest() with paging error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	query = parsed.Query()
	if got := query.Get("index"); got != "2" {
		t.Fatalf("job child list index = %q, want %q", got, "2")
	}
	if got := query.Get("count"); got != "50" {
		t.Fatalf("job child list count = %q, want %q", got, "50")
	}
	if got := query.Get("taskId"); got != "job-1" {
		t.Fatalf("job child list taskId with paging = %q, want %q", got, "job-1")
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-1")
	spec, err = backupDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job backup-detail BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/task/job-1/detail" {
		t.Fatalf("job backup-detail path = %q", spec.Path)
	}

	spec, err = syncDetail.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job sync-detail BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/sync_task/job-1/detail" {
		t.Fatalf("job sync-detail path = %q", spec.Path)
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-1")
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	spec, err = logs.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job logs BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/job_center/v1/activity/job-1/logs" {
		t.Fatalf("job logs path = %q", parsed.Path)
	}
	if got := parsed.Query().Get("count"); got != "30" {
		t.Fatalf("job logs count = %q, want %q", got, "30")
	}

	spec, err = businessTypes.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job business-types BuildRequest() error = %v", err)
	}
	if spec.Path != "/job_center/v1/business_types" {
		t.Fatalf("job business-types path = %q", spec.Path)
	}

	spec, err = appTypes.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("job app-types BuildRequest() error = %v", err)
	}
	if spec.Path != "/job_center/v1/app_types" {
		t.Fatalf("job app-types path = %q", spec.Path)
	}
}

func TestJobWriteMappings(t *testing.T) {
	cmds := job.Commands()
	stop := commandByPath(t, cmds, "job", "stop")
	deleteCmd := commandByPath(t, cmds, "job", "delete")

	rt := meta.NewRuntime()
	rt.Body = []byte(`{"jobIds":["job-1"]}`)
	spec, err := stop.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job stop BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/job_center/v1/jobs/stop" || spec.ReadOnly {
		t.Fatalf("job stop spec = %#v", spec)
	}
	if got := string(spec.Body); got != `{"jobIds":["job-1"]}` {
		t.Fatalf("job stop body = %q, want %q", got, `{"jobIds":["job-1"]}`)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"jobIds":["job-2"]}`)
	spec, err = deleteCmd.BuildRequest(rt)
	if err != nil {
		t.Fatalf("job delete BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/job_center/v1/jobs" || spec.ReadOnly {
		t.Fatalf("job delete spec = %#v", spec)
	}
	if got := string(spec.Body); got != `{"jobIds":["job-2"]}` {
		t.Fatalf("job delete body = %q, want %q", got, `{"jobIds":["job-2"]}`)
	}
}

func TestClientReadMappings(t *testing.T) {
	cmds := client.Commands()
	list := commandByPath(t, cmds, "client", "list")
	datasourceList := commandByPath(t, cmds, "client", "datasource", "list")
	runnerList := commandByPath(t, cmds, "client", "runner", "list")
	runnerTypes := commandByPath(t, cmds, "client", "runner-types")
	deployHistory := commandByPath(t, cmds, "client", "deploy-history")
	deployConfigList := commandByPath(t, cmds, "client", "deploy-config", "list")
	deployLogList := commandByPath(t, cmds, "client", "deploy-log", "list")

	spec, err := list.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("client list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || !spec.ReadOnly {
		t.Fatalf("client list spec = %#v", spec)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/commons/all_clients" {
		t.Fatalf("client list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("client list default index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("client list default count = %q, want %q", got, "30")
	}
	if got := query.Get("type"); got != "0" {
		t.Fatalf("client list default type = %q, want %q", got, "0")
	}
	if got := query.Get("status"); got != "2" {
		t.Fatalf("client list default status = %q, want %q", got, "2")
	}
	if got := query.Get("clientType"); got != "0" {
		t.Fatalf("client list default clientType = %q, want %q", got, "0")
	}

	rt := meta.NewRuntime()
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	rt.SetString("type", "0")
	rt.SetString("status", "2")
	rt.SetString("client-type", "0")
	rt.MarkChanged("index")
	rt.MarkChanged("count")
	rt.MarkChanged("type")
	rt.MarkChanged("status")
	rt.MarkChanged("client-type")
	spec, err = list.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client list BuildRequest() with query error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/commons/all_clients" {
		t.Fatalf("client list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("client list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("client list count = %q, want %q", got, "30")
	}
	if got := query.Get("type"); got != "0" {
		t.Fatalf("client list type = %q, want %q", got, "0")
	}
	if got := query.Get("status"); got != "2" {
		t.Fatalf("client list status = %q, want %q", got, "2")
	}
	if got := query.Get("clientType"); got != "0" {
		t.Fatalf("client list clientType = %q, want %q", got, "0")
	}

	rt = meta.NewRuntime()
	rt.SetString("client-id", "13caa58fd84bab1beb0b7c9d76b7efae")
	rt.SetString("full-path", "")
	rt.SetInt("business-type", 1)
	rt.SetString("runner-type", "PostgreSQL")
	rt.SetString("runner-user", "root")
	rt.SetString("request-id", "")
	rt.SetInt("index", 0)
	rt.SetInt("count", 100)
	spec, err = datasourceList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client datasource list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/commons/v1/datasources" {
		t.Fatalf("client datasource list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("clientId"); got != "13caa58fd84bab1beb0b7c9d76b7efae" {
		t.Fatalf("client datasource list clientId = %q", got)
	}
	if got := query.Get("businessType"); got != "1" {
		t.Fatalf("client datasource list businessType = %q, want %q", got, "1")
	}
	if got := query.Get("runnerType"); got != "PostgreSQL" {
		t.Fatalf("client datasource list runnerType = %q, want %q", got, "PostgreSQL")
	}
	if got := query.Get("runnerUser"); got != "root" {
		t.Fatalf("client datasource list runnerUser = %q, want %q", got, "root")
	}
	if got := query.Get("index"); got != "0" {
		t.Fatalf("client datasource list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "100" {
		t.Fatalf("client datasource list count = %q, want %q", got, "100")
	}

	spec, err = runnerList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("client runner list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/commons/clientRunner" || !spec.ReadOnly {
		t.Fatalf("client runner list spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetInt("index", 1)
	rt.SetInt("count", 20)
	rt.MarkChanged("index")
	rt.MarkChanged("count")
	spec, err = runnerList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client runner list BuildRequest() with query error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/commons/clientRunner" {
		t.Fatalf("client runner list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "1" {
		t.Fatalf("client runner list index = %q, want %q", got, "1")
	}
	if got := query.Get("count"); got != "20" {
		t.Fatalf("client runner list count = %q, want %q", got, "20")
	}

	spec, err = runnerTypes.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("client runner-types BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/commons/client/runnerTypes" || !spec.ReadOnly {
		t.Fatalf("client runner-types spec = %#v", spec)
	}

	spec, err = deployHistory.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("client deploy-history BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/deploy/v1/job/history" || !spec.ReadOnly {
		t.Fatalf("client deploy-history spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	rt.MarkChanged("index")
	rt.MarkChanged("count")
	spec, err = deployHistory.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client deploy-history BuildRequest() with query error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/deploy/v1/job/history" {
		t.Fatalf("client deploy-history path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("client deploy-history index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("client deploy-history count = %q, want %q", got, "30")
	}

	spec, err = deployConfigList.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("client deploy-config list BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/deploy/v1/hostConfig" || !spec.ReadOnly {
		t.Fatalf("client deploy-config list spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	if err := deployLogList.Validate(rt); err == nil {
		t.Fatal("client deploy-log list Validate() error = nil, want missing --job-id")
	}

	rt.SetString("job-id", "job-1")
	spec, err = deployLogList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client deploy-log list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/deploy/v1/job/jobLog" {
		t.Fatalf("client deploy-log list path = %q", parsed.Path)
	}
	query = parsed.Query()
	if got := query.Get("jobId"); got != "job-1" {
		t.Fatalf("client deploy-log list jobId = %q, want %q", got, "job-1")
	}

	rt = meta.NewRuntime()
	rt.SetString("job-id", "job-1")
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	rt.MarkChanged("index")
	rt.MarkChanged("count")
	spec, err = deployLogList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("client deploy-log list BuildRequest() with query error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	query = parsed.Query()
	if got := query.Get("index"); got != "0" {
		t.Fatalf("client deploy-log list index = %q, want %q", got, "0")
	}
	if got := query.Get("count"); got != "30" {
		t.Fatalf("client deploy-log list count = %q, want %q", got, "30")
	}
}

func TestClientDatasourceListValidateBusinessTypeEnum(t *testing.T) {
	cmds := client.Commands()
	datasourceList := commandByPath(t, cmds, "client", "datasource", "list")

	rt := meta.NewRuntime()
	rt.SetString("client-id", "13caa58fd84bab1beb0b7c9d76b7efae")
	rt.SetInt("index", 0)
	rt.SetInt("count", 100)
	rt.SetInt("business-type", 4)
	if err := datasourceList.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want invalid --business-type")
	}

	for _, value := range []int{1, 2, 3} {
		rt = meta.NewRuntime()
		rt.SetString("client-id", "13caa58fd84bab1beb0b7c9d76b7efae")
		rt.SetInt("index", 0)
		rt.SetInt("count", 100)
		rt.SetInt("business-type", value)
		if err := datasourceList.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v, want nil for business-type=%d", err, value)
		}
	}
}

func TestStoragePoolNodeDeviceListBuildsRepeatableTypesQuery(t *testing.T) {
	cmds := storage.Commands()
	poolList := commandByPath(t, cmds, "storage", "pool", "list")
	poolCreate := commandByPath(t, cmds, "storage", "pool", "create")
	poolDelete := commandByPath(t, cmds, "storage", "pool", "delete")
	poolNodeList := commandByPath(t, cmds, "storage", "pool", "node", "list")
	nodeDeviceList := commandByPath(t, cmds, "storage", "pool", "node", "device", "list")

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	spec, err := poolList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool list BuildRequest() error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/pool?index=0&count=30" {
		t.Fatalf("storage pool list path = %q, want %q", spec.Path, "/storageresmgm/v1/svc-1/pool?index=0&count=30")
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/pool" {
		t.Fatalf("storage pool list sign path = %q, want %q", got, "/storageresmgm/v1/pool")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.Body = []byte(`{"name":"pool-1"}`)
	spec, err = poolCreate.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool create BuildRequest() error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/pool" {
		t.Fatalf("storage pool create path = %q, want %q", spec.Path, "/storageresmgm/v1/svc-1/pool")
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/pool" {
		t.Fatalf("storage pool create sign path = %q, want %q", got, "/storageresmgm/v1/pool")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("pool-id", "pool-1")
	spec, err = poolDelete.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool delete BuildRequest() error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/pool?id=pool-1" {
		t.Fatalf("storage pool delete path = %q, want %q", spec.Path, "/storageresmgm/v1/svc-1/pool?id=pool-1")
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/pool" {
		t.Fatalf("storage pool delete sign path = %q, want %q", got, "/storageresmgm/v1/pool")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	spec, err = poolNodeList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool node list BuildRequest() error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/pool/node?index=0&count=100" {
		t.Fatalf("storage pool node list path = %q, want %q", spec.Path, "/storageresmgm/v1/svc-1/pool/node?index=0&count=100")
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/pool/node" {
		t.Fatalf("storage pool node list sign path = %q, want %q", got, "/storageresmgm/v1/pool/node")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("node-id", "node-1")
	spec, err = nodeDeviceList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool node device list BuildRequest() error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/device?nodeId=node-1&types=1&types=2&types=3&types=4&types=5&types=7&authorized=2" {
		t.Fatalf("storage pool node device list path = %q", spec.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/device" {
		t.Fatalf("storage pool node device list sign path = %q, want %q", got, "/storageresmgm/v1/device")
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/storageresmgm/v1/svc-1/device" {
		t.Fatalf("storage pool node device list path = %q", parsed.Path)
	}
	query := parsed.Query()
	if got := query.Get("nodeId"); got != "node-1" {
		t.Fatalf("storage pool node device list nodeId = %q, want %q", got, "node-1")
	}
	if got := query["types"]; len(got) != 6 || got[0] != "1" || got[1] != "2" || got[2] != "3" || got[3] != "4" || got[4] != "5" || got[5] != "7" {
		t.Fatalf("storage pool node device list default types = %v, want %v", got, []string{"1", "2", "3", "4", "5", "7"})
	}
	if got := query.Get("authorized"); got != "2" {
		t.Fatalf("storage pool node device list authorized = %q, want %q", got, "2")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("node-id", "node-1")
	rt.BindStringSliceFlag(newTestStringSliceBinder("types", []string{"1", "4", "7"}), "types", nil, "types")
	rt.SetInt("authorized", 2)
	spec, err = nodeDeviceList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("storage pool node device list BuildRequest() with explicit types error = %v", err)
	}
	if spec.Path != "/storageresmgm/v1/svc-1/device?nodeId=node-1&types=1&types=4&types=7&authorized=2" {
		t.Fatalf("storage pool node device list explicit path = %q", spec.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/storageresmgm/v1/device" {
		t.Fatalf("storage pool node device list explicit sign path = %q, want %q", got, "/storageresmgm/v1/device")
	}
	parsed = parseSpecURL(t, spec.Path)
	query = parsed.Query()
	if got := query["types"]; len(got) != 3 || got[0] != "1" || got[1] != "4" || got[2] != "7" {
		t.Fatalf("storage pool node device list explicit types = %v, want %v", got, []string{"1", "4", "7"})
	}
}

func TestNetworkSignPathOverrides(t *testing.T) {
	cmds := network.Commands()
	subnetList := commandByPath(t, cmds, "network", "subnet", "list")
	subnetNodeList := commandByPath(t, cmds, "network", "subnet", "node", "list")
	subnetNodeIPList := commandByPath(t, cmds, "network", "subnet", "node", "ip", "list")
	subnetNodeIPSet := commandByPath(t, cmds, "network", "subnet", "node", "ip", "set")
	subnetNodeIPRemove := commandByPath(t, cmds, "network", "subnet", "node", "ip", "remove")

	rt := meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetInt("plane-type", 1)
	spec, err := subnetList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("network subnet list BuildRequest() error = %v", err)
	}
	if spec.Path != "/clusters/v1/svc-1/subnet?planeType=1" {
		t.Fatalf("network subnet list path = %q", spec.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/clusters/v1/subnet" {
		t.Fatalf("network subnet list sign path = %q, want %q", got, "/clusters/v1/subnet")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	rt.SetInt("plane-type", 3)
	spec, err = subnetNodeList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("network subnet node list BuildRequest() error = %v", err)
	}
	parsed := parseSpecURL(t, spec.Path)
	if parsed.Path != "/clusters/v1/svc-1/subnet/nodes" {
		t.Fatalf("network subnet node list path = %q", parsed.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/clusters/v1/subnet/nodes" {
		t.Fatalf("network subnet node list sign path = %q, want %q", got, "/clusters/v1/subnet/nodes")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetInt("index", 0)
	rt.SetInt("count", 30)
	spec, err = subnetNodeIPList.BuildRequest(rt)
	if err != nil {
		t.Fatalf("network subnet node ip list BuildRequest() error = %v", err)
	}
	parsed = parseSpecURL(t, spec.Path)
	if parsed.Path != "/clusters/v1/svc-1/subnet/nodes/node_ip_addresses" {
		t.Fatalf("network subnet node ip list path = %q", parsed.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/clusters/v1/subnet/nodes/node_ip_addresses" {
		t.Fatalf("network subnet node ip list sign path = %q, want %q", got, "/clusters/v1/subnet/nodes/node_ip_addresses")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "10.0.0.1")
	spec, err = subnetNodeIPSet.BuildRequest(rt)
	if err != nil {
		t.Fatalf("network subnet node ip set BuildRequest() error = %v", err)
	}
	if spec.Path != "/clusters/v1/svc-1/subnet/nodes/node_ip_addresses" {
		t.Fatalf("network subnet node ip set path = %q", spec.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/clusters/v1/subnet/nodes/node_ip_addresses" {
		t.Fatalf("network subnet node ip set sign path = %q, want %q", got, "/clusters/v1/subnet/nodes/node_ip_addresses")
	}

	rt = meta.NewRuntime()
	rt.SetString("storage-svc-id", "svc-1")
	rt.SetString("subnet-id", "subnet-1")
	rt.SetString("node-id", "node-1")
	rt.SetString("ip", "10.0.0.1")
	spec, err = subnetNodeIPRemove.BuildRequest(rt)
	if err != nil {
		t.Fatalf("network subnet node ip remove BuildRequest() error = %v", err)
	}
	if spec.Path != "/clusters/v1/svc-1/subnet/nodes/subnet-1/node-1" {
		t.Fatalf("network subnet node ip remove path = %q", spec.Path)
	}
	if got := spec.Headers.Get("X-Foundation-Cli-Sign-Path"); got != "/clusters/v1/subnet/nodes/subnet-1/node-1" {
		t.Fatalf("network subnet node ip remove sign path = %q, want %q", got, "/clusters/v1/subnet/nodes/subnet-1/node-1")
	}
}

func TestAPICommandValidation(t *testing.T) {
	def := api.Commands()[0]

	rt := meta.NewRuntime()
	rt.SetString("method", http.MethodGet)
	rt.SetString("path", "https://backend.example.com/api")
	if err := def.Validate(rt); err == nil {
		t.Fatal("Validate() error = nil, want invalid path")
	}

	rt = meta.NewRuntime()
	rt.SetString("method", "BAD METHOD")
	rt.SetString("path", "/api/safe/path")
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

func assertCommandMissing(t *testing.T, defs []meta.CommandMeta, path ...string) {
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
			t.Fatalf("command path %v should be absent", path)
		}
	}
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
