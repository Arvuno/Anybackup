package protect

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newConfigPolicyGetByAppTypeCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "config-policy", "get-by-app-type"},
		Use:           "get-by-app-type",
		Description:   "Get protect config policy by app type",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "app-type", "", "Application type")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("app-type", rt.String("app-type"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v2/app_type/%s/config_police", rt.String("app-type")),
				ReadOnly: true,
			}, nil
		},
	}
}
