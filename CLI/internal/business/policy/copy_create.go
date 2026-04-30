package policy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newCopyCreateCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "policy",
		CanonicalPath: []string{"policy", "copy", "create"},
		Use:           "create",
		Description:   "Create policy copy info",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if len(rt.Body) == 0 {
				return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
			}
			if _, err := policyNameFromBody(rt.Body); err != nil {
				return err
			}
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			name, err := policyNameFromBody(rt.Body)
			if err != nil {
				return err
			}

			precheckResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/api/sla/v1/common/name/exists?name=" + url.QueryEscape(name),
				ReadOnly: true,
			})
			if err != nil {
				return err
			}
			exists, err := parsePolicyNameExists(precheckResp)
			if err != nil {
				return err
			}
			if exists {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "policy copy name already exists")
			}

			body, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/api/sla/v1/group/copy_info",
				Body:     rt.Body,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/api/sla/v1/group/copy_info",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}
