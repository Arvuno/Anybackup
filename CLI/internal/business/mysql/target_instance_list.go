package mysql

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newTargetInstanceListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "target-instance", "list"},
		Use:           "list",
		Description:   "List target MySQL instances",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := validateCountRange("index", rt.Int("index"), 0, 1<<30); err != nil {
				return err
			}
			if err := validateCountRange("count", rt.Int("count"), 1, 100); err != nil {
				return err
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			count := rt.Int("count")
			if count == 0 {
				count = 30
			}
			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", rt.Int("index"))).
				Add("count", fmt.Sprintf("%d", count)).
				Add("type", "202").
				Add("customs", "clusterType:202")
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/resource_center/v1/databasealone/instance_list?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}
