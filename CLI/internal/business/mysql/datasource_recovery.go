package mysql

import (
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newDatasourceRecoveryCommand() meta.CommandMeta {
	allowedGran := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}}

	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "datasource", "recovery"},
		Use:           "recovery",
		Description:   "Get MySQL recovery datasource",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "data-set-id", "", "Data set ID")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "timestamp", "", "Timestamp")
			rt.BindStringFlag(cmd.Flags(), "restore-gran", "", "Restore granularity")
			rt.BindStringFlag(cmd.Flags(), "request-id", "", "Request ID")
			rt.BindStringFlag(cmd.Flags(), "full-path", "/", "Full path")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 100, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("data-set-id", rt.String("data-set-id")); err != nil {
				return err
			}
			if err := inputs.RequireNonEmpty("storage-service-id", rt.String("storage-service-id")); err != nil {
				return err
			}
			if err := inputs.RequireNonEmpty("timestamp", rt.String("timestamp")); err != nil {
				return err
			}
			if v := rt.String("restore-gran"); v != "" {
				if err := inputs.RequireOneOfString("restore-gran", v, allowedGran); err != nil {
					return err
				}
			}
			if err := validateOptionalNonNegativeInt64("timestamp", rt.String("timestamp")); err != nil {
				return err
			}
			if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
				return err
			}
			return validateCountRange("count", rt.Int("count"), 1, 100)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().
				Add("dataSetId", rt.String("data-set-id")).
				Add("storageServiceId", rt.String("storage-service-id")).
				Add("timestamp", rt.String("timestamp")).
				Add("fullPath", rt.String("full-path")).
				Add("index", strconv.Itoa(rt.Int("index"))).
				Add("count", strconv.Itoa(rt.Int("count")))
			if v := rt.String("restore-gran"); v != "" {
				qb.Add("restoreGran", v)
			}
			if v := rt.String("request-id"); v != "" {
				qb.Add("requestId", v)
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/recovery_datasource?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
