package client

import (
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newListCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "list"},
		Use:           "list",
		Description:   "List clients",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
			rt.BindStringFlag(cmd.Flags(), "type", "", "Client type")
			rt.BindStringFlag(cmd.Flags(), "status", "", "Client status")
			rt.BindStringFlag(cmd.Flags(), "client-type", "", "Client subtype")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.Changed("index") || rt.Changed("count") {
				return inputs.ValidatePaging(rt.Int("index"), rt.Int("count"))
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			index := 0
			if rt.Changed("index") {
				index = rt.Int("index")
			}
			count := 30
			if rt.Changed("count") {
				count = rt.Int("count")
			}
			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", index)).
				Add("count", fmt.Sprintf("%d", count))

			if rt.Changed("type") {
				qb.Add("type", rt.String("type"))
			} else {
				qb.Add("type", "0")
			}
			if rt.Changed("status") {
				qb.Add("status", rt.String("status"))
			} else {
				qb.Add("status", "2")
			}
			if rt.Changed("client-type") {
				qb.Add("clientType", rt.String("client-type"))
			} else {
				qb.Add("clientType", "0")
			}

			path := "/commons/all_clients"
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
