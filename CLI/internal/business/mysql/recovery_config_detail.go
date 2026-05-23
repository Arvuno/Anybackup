package mysql

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newRecoveryConfigDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "recovery-config", "detail"},
		Use:           "detail",
		Description:   "Get MySQL recovery config detail",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "system-id", "", "System ID")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("system-id", rt.String("system-id")); err != nil {
				return err
			}
			return inputs.RequireNonEmpty("object-id", rt.String("object-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().
				Add("systemId", rt.String("system-id")).
				Add("objectId", rt.String("object-id"))
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/recovery_config?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
