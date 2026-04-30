package protect

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newStoragePoolAutoSelectCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "storage-pool", "auto-select"},
		Use:           "auto-select",
		Description:   "Auto select storage pool for backup destination",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringSliceFlag(cmd.Flags(), "exclude-id", nil, "Exclude storage service ID (repeatable)")
		},
		Validate: func(rt *meta.Runtime) error {
			for _, value := range rt.Strings("exclude-id") {
				if len(value) != 32 {
					return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", "exclude-id"))
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if values := rt.Strings("exclude-id"); len(values) > 0 {
				qb.AddAll("excludeIds", values)
			}

			path := "/backupmgm/v2/auto_select/storage_pool"
			if encoded := qb.Encode(); encoded != "" {
				path += "?" + encoded
			}

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     path,
				ReadOnly: true,
			}, nil
		},
	}
}
