package mysql

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newRecoveryDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "recovery-detail"},
		Use:           "recovery-detail",
		Description:   "Get MySQL recovery task detail",
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
			qb := inputs.NewQueryBuilder().
				Add("taskId", rt.String("task-id"))
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/get_recovery_task_detail?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
