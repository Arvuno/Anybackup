package vmware

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newRestoreConfigCreateCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "vmware",
		CanonicalPath: []string{"vmware", "restore-config", "create"},
		Use:           "create",
		Description:   "Create VMware restore config",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("object-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --object-id")
			}
			if len(rt.Body) == 0 {
				return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     fmt.Sprintf("/backupmgm/v1/virtual/vmware/%s/recovery", rt.String("object-id")),
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}
