package vmware

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newBackupDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "vmware",
		CanonicalPath: []string{"vmware", "backup-detail"},
		Use:           "backup-detail",
		Description:   "Get VMware backup task detail",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "task-id", "", "Task ID")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("task-id", rt.String("task-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/virtual/vmware/backup_task/%s/detail", rt.String("task-id")),
				ReadOnly: true,
			}, nil
		},
	}
}
