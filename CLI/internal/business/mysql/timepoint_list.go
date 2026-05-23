package mysql

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newTimepointListCommand() meta.CommandMeta {
	allowedGran := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}}

	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "recovery", "timepoint", "list"},
		Use:           "list",
		Description:   "List MySQL recovery timepoints",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 20, "Page size")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "storage-pool-id", "", "Storage pool ID")
			rt.BindStringFlag(cmd.Flags(), "start-time", "", "Start time")
			rt.BindStringFlag(cmd.Flags(), "end-time", "", "End time")
			rt.BindStringFlag(cmd.Flags(), "backup-task-type", "", "Backup task type")
			rt.BindStringFlag(cmd.Flags(), "restore-gran", "", "Restore granularity")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.Int("index") < 0 {
				return validateCountRange("index", -1, 0, 0)
			}
			if err := validateCountRange("count", rt.Int("count"), 1, 100); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("start-time", rt.String("start-time")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("end-time", rt.String("end-time")); err != nil {
				return err
			}
			if v := rt.String("restore-gran"); v != "" {
				if err := inputs.RequireOneOfString("restore-gran", v, allowedGran); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			count := rt.Int("count")
			if count == 0 {
				count = 20
			}
			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", rt.Int("index"))).
				Add("count", fmt.Sprintf("%d", count))
			if v := rt.String("object-id"); v != "" {
				qb.Add("objectId", v)
			}
			if v := rt.String("storage-service-id"); v != "" {
				qb.Add("storageServiceId", v)
			}
			if v := rt.String("storage-pool-id"); v != "" {
				qb.Add("storagePoolId", v)
			}
			if v := rt.String("start-time"); v != "" {
				qb.Add("startTime", v)
			}
			if v := rt.String("end-time"); v != "" {
				qb.Add("endTime", v)
			}
			if v := rt.String("backup-task-type"); v != "" {
				qb.Add("backupTaskType", v)
			}
			if v := rt.String("restore-gran"); v != "" {
				qb.Add("restoreGran", v)
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/time_points?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
