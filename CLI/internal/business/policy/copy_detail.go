package policy

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newCopyDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "policy",
		CanonicalPath: []string{"policy", "copy", "detail"},
		Use:           "detail",
		Description:   "Get copy policy detail",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "group-id", "", "Policy group ID")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("group-id", rt.String("group-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/api/sla/v1/group/%s/copy_detail", rt.String("group-id")),
				ReadOnly: true,
			}, nil
		},
	}
}
