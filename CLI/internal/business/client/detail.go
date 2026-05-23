package client

import (
	"net/http"
	"net/url"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "detail"},
		Use:           "detail",
		Description:   "Get client detail by client ID",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "client-id", "", "Client ID")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("client-id", rt.String("client-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/commons/client/" + url.PathEscape(rt.String("client-id")),
				ReadOnly: true,
			}, nil
		},
	}
}

