package storage

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPoolDeleteCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "storage",
		CanonicalPath: []string{"storage", "pool", "delete"},
		Use:           "delete",
		Description:   "Delete a storage pool",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "storage-svc-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "pool-id", "", "Storage pool ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateStorageServiceID(rt); err != nil {
				return err
			}
			return validatePoolID(rt)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder().Add("id", rt.String("pool-id"))
			return console.RequestSpec{
				Method:   http.MethodDelete,
				Path:     fmt.Sprintf("/storageresmgm/v1/%s/pool?%s", rt.String("storage-svc-id"), qb.Encode()),
				Headers:  http.Header{"X-Foundation-Cli-Sign-Path": []string{"/storageresmgm/v1/pool"}},
				ReadOnly: false,
			}, nil
		},
	}
}
