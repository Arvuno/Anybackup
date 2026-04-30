package vmware

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newDatasourceGetCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "vmware",
		CanonicalPath: []string{"vmware", "datasource", "get"},
		Use:           "get",
		Description:   "Get VMware datasource sub-objects",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "production-system-id", "", "Production system ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("production-system-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --production-system-id")
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/virtual/vmware/%s/sub_objects", rt.String("production-system-id")),
				ReadOnly: true,
			}, nil
		},
	}
}
