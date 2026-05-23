package storage

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPoolCreateCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "storage",
		CanonicalPath: []string{"storage", "pool", "create"},
		Use:           "create",
		Description:   "Create a storage pool",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			return validatePoolCreateBody(rt.Body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     fmt.Sprintf("/storageresmgm/v1/%s/pool", rt.String("storage-svc-id")),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/storageresmgm/v1/pool"}},
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}
