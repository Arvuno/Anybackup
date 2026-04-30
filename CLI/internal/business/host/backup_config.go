package host

import (
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newBackupConfigCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "host",
		CanonicalPath: []string{"host", "backup-config"},
		Use:           "backup-config",
		Description:   "Start batch backup for host objects",
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
				Path:     "/backupmgm/v1/batch/backup_tasks/start",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}
