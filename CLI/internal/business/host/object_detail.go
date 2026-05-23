package host

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newObjectDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "host",
		CanonicalPath: []string{"host", "object", "detail"},
		Use:           "detail",
		Description:   "Get host protect object detail",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Protect object ID")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("object-id", rt.String("object-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/file/fileset/%s", rt.String("object-id")),
				ReadOnly: true,
			}, nil
		},
	}
}
