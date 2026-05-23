package mysql

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newRecoveryRangeCommand() meta.CommandMeta {
	allowedGran := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {}}

	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "recovery", "range"},
		Use:           "range",
		Description:   "Get MySQL recovery range",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "restore-object-type", "", "Restore object type")
			rt.BindStringFlag(cmd.Flags(), "restore-gran", "", "Restore granularity")
			rt.BindStringFlag(cmd.Flags(), "timestamp", "", "Timestamp")
			rt.BindStringFlag(cmd.Flags(), "backup-task-type", "", "Backup task type")
			rt.BindStringFlag(cmd.Flags(), "storage-pool-id", "", "Storage pool ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if v := rt.String("restore-gran"); v != "" {
				if err := inputs.RequireOneOfString("restore-gran", v, allowedGran); err != nil {
					return err
				}
			}
			return validateOptionalNonNegativeInt64("timestamp", rt.String("timestamp"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if v := rt.String("object-id"); v != "" {
				qb.Add("objectId", v)
			}
			if v := rt.String("storage-service-id"); v != "" {
				qb.Add("storageServiceId", v)
			}
			if v := rt.String("restore-object-type"); v != "" {
				qb.Add("restoreObjectType", v)
			}
			if v := rt.String("restore-gran"); v != "" {
				qb.Add("restoreGran", v)
			}
			if v := rt.String("timestamp"); v != "" {
				qb.Add("timestamp", v)
			}
			if v := rt.String("backup-task-type"); v != "" {
				qb.Add("backupTaskType", v)
			}
			if v := rt.String("storage-pool-id"); v != "" {
				qb.Add("storagePoolId", v)
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/get_recovery_range?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
