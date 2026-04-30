package client

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newRunnerListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "runner", "list"},
		Use:           "list",
		Description:   "List client runners",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.Changed("index") || rt.Changed("count") {
				return inputs.ValidatePaging(rt.Int("index"), rt.Int("count"))
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if rt.Changed("index") {
				qb.Add("index", fmt.Sprintf("%d", rt.Int("index")))
			}
			if rt.Changed("count") {
				qb.Add("count", fmt.Sprintf("%d", rt.Int("count")))
			}

			path := "/commons/clientRunner"
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
