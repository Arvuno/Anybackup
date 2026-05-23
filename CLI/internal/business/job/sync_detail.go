package job

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newSyncDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "sync-detail"},
		Use:           "sync-detail",
		Description:   "Get sync task detail by job id",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "job-id", "", "job id")
		},
		Validate: func(rt *meta.Runtime) error {
			return inputs.RequireNonEmpty("job-id", rt.String("job-id"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/sync_task/%s/detail", rt.String("job-id")),
				ReadOnly: true,
			}, nil
		},
	}
}
