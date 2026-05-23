package mysql

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newDatasourceBackupCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "datasource", "backup"},
		Use:           "backup",
		Description:   "Get MySQL backup datasource",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 20, "Page size")
			rt.BindStringFlag(cmd.Flags(), "system-id", "", "System ID")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			rt.BindStringFlag(cmd.Flags(), "full-path", "", "Full path")
			rt.BindStringFlag(cmd.Flags(), "request-id", "", "Request ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.Int("index") < 0 {
				return validateCountRange("index", -1, 0, 0)
			}
			if err := validateCountRange("count", rt.Int("count"), 1, 100); err != nil {
				return err
			}
			return inputs.RequireNonEmpty("object-id", rt.String("object-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			count := rt.Int("count")
			if count == 0 {
				count = 20
			}
			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", rt.Int("index"))).
				Add("count", fmt.Sprintf("%d", count)).
				Add("objectId", rt.String("object-id"))
			if v := rt.String("system-id"); v != "" {
				qb.Add("systemId", v)
			}
			if v := rt.String("full-path"); v != "" {
				qb.Add("fullPath", v)
			}
			if v := rt.String("request-id"); v != "" {
				qb.Add("requestId", v)
			}
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/backup_datasource?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
