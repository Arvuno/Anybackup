package mysql

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newObjectGetCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "object", "get"},
		Use:           "get",
		Description:   "Get MySQL object by ID",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("object-id", rt.String("object-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().
				Add("objectId", rt.String("object-id"))
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/mysql/object?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
