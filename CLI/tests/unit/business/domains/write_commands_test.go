//go:build legacy_cli_contract
// +build legacy_cli_contract

package domains_test

import (
	"net/http"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/backup"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/clean"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/client"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/host"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/job"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/mysql"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/object"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/restore"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/sla"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/virtual"
)

func TestWriteCommandMappings_RequireBodyAndDisableRetry(t *testing.T) {
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"force":true}`)

	spec, err := backup.Commands()[0].BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/protect_object/obj-1/backup_task/start" || spec.ReadOnly {
		t.Fatalf("backup spec = %#v", spec)
	}

	spec, err = clean.Commands()[0].BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/protect_object/obj-1/clean_task/start" || spec.ReadOnly {
		t.Fatalf("clean spec = %#v", spec)
	}

	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectId":"obj-1","files":["/tmp/a.txt"]}`)
	spec, err = restore.Commands()[0].BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Path != "/backupmgm/v1/file/recovery" || spec.ReadOnly {
		t.Fatalf("restore spec = %#v", spec)
	}
}

func TestHostWriteCommandMappings(t *testing.T) {
	hostCmds := host.Commands()

	// host object create -> POST /backupmgm/v1/file/fileset
	objCreate := hostCmds[1]
	if objCreate.CanonicalPath[1] != "object" || objCreate.CanonicalPath[2] != "create" {
		t.Fatalf("host object create path = %v", objCreate.CanonicalPath)
	}

	// host backup-config -> POST /backupmgm/v1/batch/backup_tasks/start
	backupConfig := hostCmds[3]
	if backupConfig.CanonicalPath[1] != "backup-config" {
		t.Fatalf("host backup-config path = %v", backupConfig.CanonicalPath)
	}
	spec, err := backupConfig.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/batch/backup_tasks/start" || spec.ReadOnly {
		t.Fatalf("host backup-config spec = %#v", spec)
	}
}

func TestClientWriteCommandMappings(t *testing.T) {
	clientCmds := client.Commands()

	// client deploy -> POST /deploy/v1/job/config
	deploy := clientCmds[0]
	if deploy.CanonicalPath[1] != "deploy" {
		t.Fatalf("client deploy path = %v", deploy.CanonicalPath)
	}
	spec, err := deploy.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/deploy/v1/job/config" || spec.ReadOnly {
		t.Fatalf("client deploy spec = %#v", spec)
	}
}

func TestVirtualWriteCommandMappings(t *testing.T) {
	virtualCmds := virtual.Commands()

	// virtual vmware backup-config -> POST /backupmgm/v1/virtual/vmware/{objectId}/backup_config
	rt := meta.NewRuntime()
	rt.SetString("object-id", "vm-1")
	rt.Body = []byte(`{"key":"value"}`)
	spec, err := virtualCmds[3].BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/virtual/vmware/vm-1/backup_config" || spec.ReadOnly {
		t.Fatalf("virtual backup-config spec = %#v", spec)
	}

	// virtual vmware recover -> POST /backupmgm/v1/virtual/vmware/{objectId}/recovery
	rt = meta.NewRuntime()
	rt.SetString("object-id", "vm-2")
	rt.Body = []byte(`{"key":"value"}`)
	spec, err = virtualCmds[4].BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/virtual/vmware/vm-2/recovery" || spec.ReadOnly {
		t.Fatalf("virtual recover spec = %#v", spec)
	}
}

func TestMySQLWriteCommandMappings(t *testing.T) {
	mysqlCmds := mysql.Commands()

	// mysql backup-config -> POST /backupmgm/v1/mysql/app_backup_config
	spec, err := mysqlCmds[1].BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/mysql/app_backup_config" || spec.ReadOnly {
		t.Fatalf("mysql backup-config spec = %#v", spec)
	}

	// mysql recover -> POST /backupmgm/v1/mysql/recovery
	spec, err = mysqlCmds[2].BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/mysql/recovery" || spec.ReadOnly {
		t.Fatalf("mysql recover spec = %#v", spec)
	}
}

func TestObjectWriteCommandMappings(t *testing.T) {
	objCmds := object.Commands()
	slaBind := commandByPath(t, objCmds, "object", "sla", "bind")
	slaBindBatch := commandByPath(t, objCmds, "object", "sla", "bind-batch")
	slaUnbindBatch := commandByPath(t, objCmds, "object", "sla", "unbind-batch")

	// object sla bind -> POST /backupmgm/v1/protect_object/{objectId}/slas
	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	rt.Body = []byte(`{"slaId":"sla-1"}`)
	spec, err := slaBind.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/obj-1/slas" || spec.ReadOnly {
		t.Fatalf("object sla bind spec = %#v", spec)
	}

	// object sla bind-batch -> POST /backupmgm/v1/protect_object/slas
	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectIds":["obj-1"],"slaId":"sla-1"}`)
	spec, err = slaBindBatch.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/protect_object/slas" || spec.ReadOnly {
		t.Fatalf("object sla bind-batch spec = %#v", spec)
	}

	// object sla unbind-batch -> DELETE /backupmgm/v1/protect_object/slas
	rt = meta.NewRuntime()
	rt.Body = []byte(`{"objectIds":["obj-1"]}`)
	spec, err = slaUnbindBatch.BuildRequest(rt)
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/backupmgm/v1/protect_object/slas" || spec.ReadOnly {
		t.Fatalf("object sla unbind-batch spec = %#v", spec)
	}
}

func TestBackupStartBatchCommandMapping(t *testing.T) {
	backupCmds := backup.Commands()
	if len(backupCmds) != 2 {
		t.Fatalf("backup commands count = %d, want 2", len(backupCmds))
	}

	// backup start-batch -> POST /backupmgm/v1/batch/backup_tasks/start
	spec, err := backupCmds[1].BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/backupmgm/v1/batch/backup_tasks/start" || spec.ReadOnly {
		t.Fatalf("backup start-batch spec = %#v", spec)
	}
}

func TestJobListCommandMapping(t *testing.T) {
	jobCmds := job.Commands()
	if len(jobCmds) != 7 {
		t.Fatalf("job commands count = %d, want 7", len(jobCmds))
	}

	// job list -> GET /job_center/v1/jobs
	spec, err := commandByPath(t, jobCmds, "job", "list").BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodGet || spec.Path != "/job_center/v1/jobs" || !spec.ReadOnly {
		t.Fatalf("job list spec = %#v", spec)
	}
}

func TestSLAWriteCommandMappings(t *testing.T) {
	slaCmds := sla.Commands()
	backupCreate := commandByPath(t, slaCmds, "sla", "backup", "create")
	copyCreate := commandByPath(t, slaCmds, "sla", "copy", "create")
	backupDelete := commandByPath(t, slaCmds, "sla", "backup", "delete")

	spec, err := backupCreate.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/api/sla/v1/group/backup_info" || spec.ReadOnly {
		t.Fatalf("sla backup create spec = %+v", spec)
	}

	spec, err = copyCreate.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodPost || spec.Path != "/api/sla/v1/group/copy_info" || spec.ReadOnly {
		t.Fatalf("sla copy create spec = %+v", spec)
	}

	spec, err = backupDelete.BuildRequest(meta.NewRuntime())
	if err != nil {
		t.Fatalf("BuildRequest() error = %v", err)
	}
	if spec.Method != http.MethodDelete || spec.Path != "/api/sla/v1/groups" || spec.ReadOnly {
		t.Fatalf("sla backup delete spec = %+v", spec)
	}
}

func TestSLAValidate_MissingBody(t *testing.T) {
	slaCmds := sla.Commands()
	backupCreate := commandByPath(t, slaCmds, "sla", "backup", "create")
	copyCreate := commandByPath(t, slaCmds, "sla", "copy", "create")
	backupDelete := commandByPath(t, slaCmds, "sla", "backup", "delete")

	for _, cmd := range []struct {
		name string
		meta meta.CommandMeta
	}{
		{"sla backup create", backupCreate},
		{"sla copy create", copyCreate},
		{"sla backup delete", backupDelete},
	} {
		t.Run(cmd.name, func(t *testing.T) {
			if cmd.meta.Validate == nil {
				t.Skip("no Validate function")
			}
			if err := cmd.meta.Validate(meta.NewRuntime()); err == nil {
				t.Fatal("Validate() error = nil, want missing body error")
			}
		})
	}
}

func TestObjectSLABindList_Validate(t *testing.T) {
	objCmds := object.Commands()
	bindingList := commandByPath(t, objCmds, "object", "sla", "bind-list")

	if err := bindingList.Validate(meta.NewRuntime()); err == nil {
		t.Fatal("Validate() error = nil, want missing --object-id error")
	}

	rt := meta.NewRuntime()
	rt.SetString("object-id", "obj-1")
	if err := bindingList.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObjectValidate_MissingRequiredFields(t *testing.T) {
	objCmds := object.Commands()
	timepointList := commandByPath(t, objCmds, "object", "timepoint", "list")
	slaBind := commandByPath(t, objCmds, "object", "sla", "bind")
	slaBindBatch := commandByPath(t, objCmds, "object", "sla", "bind-batch")
	slaUnbindBatch := commandByPath(t, objCmds, "object", "sla", "unbind-batch")

	t.Run("timepoint list missing object-id", func(t *testing.T) {
		if err := timepointList.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("sla bind missing object-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte(`{"slaId":"sla-1"}`)
		if err := slaBind.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("sla bind missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		if err := slaBind.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("sla bind-batch missing body", func(t *testing.T) {
		if err := slaBindBatch.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("sla unbind-batch missing body", func(t *testing.T) {
		if err := slaUnbindBatch.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("timepoint list rejects invalid business", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		rt.SetString("business", "2")
		if err := timepointList.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid --business")
		}
	})
}

func TestBackupValidate_MissingRequiredFields(t *testing.T) {
	backupCmds := backup.Commands()

	t.Run("start missing object-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte(`{"force":true}`)
		if err := backupCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("start missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		if err := backupCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("start-batch missing body", func(t *testing.T) {
		if err := backupCmds[1].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}

func TestCleanValidate_MissingRequiredFields(t *testing.T) {
	cleanCmds := clean.Commands()

	t.Run("start missing object-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte(`{"force":true}`)
		if err := cleanCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("start missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "obj-1")
		if err := cleanCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}

func TestRestoreValidate_MissingBody(t *testing.T) {
	restoreCmds := restore.Commands()
	if err := restoreCmds[0].Validate(meta.NewRuntime()); err == nil {
		t.Fatal("Validate() error = nil, want missing body")
	}
}

func TestJobValidate_MissingRequiredFields(t *testing.T) {
	jobCmds := job.Commands()
	backupDetail := commandByPath(t, jobCmds, "job", "backup-detail")
	childList := commandByPath(t, jobCmds, "job", "child", "list")
	detail := commandByPath(t, jobCmds, "job", "detail")
	syncDetail := commandByPath(t, jobCmds, "job", "sync-detail")
	logs := commandByPath(t, jobCmds, "job", "logs")
	list := commandByPath(t, jobCmds, "job", "list")

	t.Run("backup-detail missing job-id", func(t *testing.T) {
		if err := backupDetail.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --job-id")
		}
	})

	t.Run("child list missing job-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetInt("index", 0)
		rt.SetInt("count", 30)
		if err := childList.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --job-id")
		}
	})

	t.Run("child list negative index", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("job-id", "job-1")
		rt.SetInt("index", -1)
		rt.SetInt("count", 30)
		if err := childList.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want negative index error")
		}
	})

	t.Run("child list zero count", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("job-id", "job-1")
		rt.SetInt("index", 0)
		rt.SetInt("count", 0)
		if err := childList.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want zero count error")
		}
	})

	t.Run("child list valid paging", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("job-id", "job-1")
		rt.SetInt("index", 0)
		rt.SetInt("count", 30)
		if err := childList.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	})

	t.Run("detail missing job-id", func(t *testing.T) {
		if err := detail.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --job-id")
		}
	})

	t.Run("sync-detail missing job-id", func(t *testing.T) {
		if err := syncDetail.Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --job-id")
		}
	})

	t.Run("logs missing job-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetInt("index", 0)
		rt.SetInt("count", 30)
		if err := logs.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --job-id")
		}
	})

	t.Run("list negative index", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetInt("index", -1)
		rt.SetInt("count", 10)
		if err := list.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want negative index error")
		}
	})

	t.Run("list zero count", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetInt("index", 0)
		rt.SetInt("count", 0)
		if err := list.Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want zero count error")
		}
	})

	t.Run("list valid paging", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetInt("index", 0)
		rt.SetInt("count", 10)
		if err := list.Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	})
}

func TestHostValidate_MissingBody(t *testing.T) {
	hostCmds := host.Commands()

	t.Run("object create missing body", func(t *testing.T) {
		if err := hostCmds[1].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("backup-config missing body", func(t *testing.T) {
		if err := hostCmds[3].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}

func TestClientValidate_MissingBody(t *testing.T) {
	clientCmds := client.Commands()
	if err := clientCmds[0].Validate(meta.NewRuntime()); err == nil {
		t.Fatal("Validate() error = nil, want missing body")
	}
}

func TestVirtualValidate(t *testing.T) {
	virtualCmds := virtual.Commands()

	t.Run("list invalid exec-status", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("exec-status", "0")
		if err := virtualCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid exec-status error")
		}
	})

	t.Run("list invalid sort-field", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("sort-field", "invalid")
		if err := virtualCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want invalid sort-field error")
		}
	})

	t.Run("list negative index", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("exec-status", "1")
		rt.SetString("sort-field", "name")
		rt.SetInt("index", -1)
		rt.SetInt("count", 10)
		if err := virtualCmds[0].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want negative index error")
		}
	})

	t.Run("list valid combination", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("exec-status", "1")
		rt.SetString("sort-field", "name")
		rt.SetBool("is-config", true)
		rt.SetInt("index", 0)
		rt.SetInt("count", 10)
		if err := virtualCmds[0].Validate(rt); err != nil {
			t.Fatalf("Validate() error = %v", err)
		}
	})

	t.Run("vmware sub-object-list missing production-system-id", func(t *testing.T) {
		if err := virtualCmds[1].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --production-system-id")
		}
	})

	t.Run("vmware vm-list missing production-system-id", func(t *testing.T) {
		if err := virtualCmds[2].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing --production-system-id")
		}
	})

	t.Run("vmware backup-config missing object-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte(`{"key":"value"}`)
		if err := virtualCmds[3].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("vmware backup-config missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "vm-1")
		if err := virtualCmds[3].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("vmware recover missing object-id", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.Body = []byte(`{"target":"host"}`)
		if err := virtualCmds[4].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing --object-id")
		}
	})

	t.Run("vmware recover missing body", func(t *testing.T) {
		rt := meta.NewRuntime()
		rt.SetString("object-id", "vm-1")
		if err := virtualCmds[4].Validate(rt); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}

func TestMySQLValidate_MissingBody(t *testing.T) {
	mysqlCmds := mysql.Commands()

	t.Run("backup-config missing body", func(t *testing.T) {
		if err := mysqlCmds[1].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})

	t.Run("recover missing body", func(t *testing.T) {
		if err := mysqlCmds[2].Validate(meta.NewRuntime()); err == nil {
			t.Fatal("Validate() error = nil, want missing body")
		}
	})
}
