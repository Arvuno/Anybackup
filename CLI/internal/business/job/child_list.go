package job

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newChildListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "child", "list"},
		Use:           "list",
		Description:   "List child jobs by job id",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "job-id", "", "job id")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("job-id", rt.String("job-id")); err != nil {
				return err
			}
			return inputs.ValidatePaging(rt.Int("index"), rt.Int("count"))
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			jobID := rt.String("job-id")
			index := rt.Int("index")
			count := rt.Int("count")
			if count <= 0 {
				count = 30
			}
			query := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", index)).
				Add("count", fmt.Sprintf("%d", count)).
				Add("taskId", jobID).
				Encode()

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/task/%s/child?%s", jobID, query),
				ReadOnly: true,
			}, nil
		},
	}
}
