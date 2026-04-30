package protect

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPolicyBindBatchCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "policy", "bind-batch"},
		Use:           "bind-batch",
		Description:   "Batch bind policy to protect objects",
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
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/backupmgm/v1/protect_object/slas",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}
