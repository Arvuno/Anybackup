package client

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newDatasourceListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "datasource", "list"},
		Use:           "list",
		Description:   "List client datasources",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "client-id", "", "Client ID")
			rt.BindStringFlag(cmd.Flags(), "full-path", "", "Datasource full path")
			rt.BindIntFlag(cmd.Flags(), "business-type", 0, "Business type")
			rt.BindStringFlag(cmd.Flags(), "runner-type", "", "Runner type")
			rt.BindStringFlag(cmd.Flags(), "runner-user", "", "Runner user")
			rt.BindStringFlag(cmd.Flags(), "request-id", "", "Request ID")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 100, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("client-id", rt.String("client-id")); err != nil {
				return err
			}
			if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
				return err
			}
			if rt.Int("count") > 100 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "--count must be <= 100")
			}
			if rt.Int("business-type") < 0 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "--business-type must be non-negative")
			}
			if businessType := rt.Int("business-type"); businessType != 0 && businessType != 1 && businessType != 2 && businessType != 3 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --business-type")
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().
				Add("clientId", rt.String("client-id")).
				Add("index", fmt.Sprintf("%d", rt.Int("index"))).
				Add("count", fmt.Sprintf("%d", rt.Int("count")))

			if rt.Changed("full-path") {
				qb.Add("fullPath", rt.String("full-path"))
			} else if v := rt.String("full-path"); v != "" {
				qb.Add("fullPath", v)
			}
			if rt.Changed("business-type") || rt.Int("business-type") != 0 {
				qb.Add("businessType", fmt.Sprintf("%d", rt.Int("business-type")))
			}
			if v := rt.String("runner-type"); v != "" {
				qb.Add("runnerType", v)
			}
			if v := rt.String("runner-user"); v != "" {
				qb.Add("runnerUser", v)
			}
			if rt.Changed("request-id") {
				qb.Add("requestId", rt.String("request-id"))
			} else if v := rt.String("request-id"); v != "" {
				qb.Add("requestId", v)
			}

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/commons/v1/datasources?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
