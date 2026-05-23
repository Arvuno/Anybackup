package host

import (
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newListCommand() meta.CommandMeta {
	allowedClientOSFilter := map[string]struct{}{
		"0": {}, "1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {},
	}

	return meta.CommandMeta{
		Domain:        "host",
		CanonicalPath: []string{"host", "list"},
		Use:           "list",
		Description:   "List hosts",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
			rt.BindStringSliceFlag(cmd.Flags(), "runner-types", nil, "Runner types")
			rt.BindStringFlag(cmd.Flags(), "group-id", "", "Group ID")
			rt.BindStringSliceFlag(cmd.Flags(), "client-os-filter", nil, "Client OS filter")
			rt.BindStringFlag(cmd.Flags(), "filter", "", "Filter text")
			rt.BindStringFlag(cmd.Flags(), "is-child", "", "Whether to query child hosts")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.Changed("index") || rt.Changed("count") {
				if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
					return err
				}
			}
			if err := inputs.RequireAllInSet("client-os-filter", rt.Strings("client-os-filter"), allowedClientOSFilter); err != nil {
				return err
			}
			if v := rt.String("is-child"); v != "" {
				if _, err := strconv.ParseBool(v); err != nil {
					return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --is-child")
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if values := rt.Strings("runner-types"); len(values) > 0 {
				qb.AddAll("runnerTypes", values)
			}
			if v := rt.String("group-id"); v != "" {
				qb.Add("groupId", v)
			}
			if values := rt.Strings("client-os-filter"); len(values) > 0 {
				qb.AddAll("clientOsFilter", values)
			}
			if v := rt.String("filter"); v != "" {
				qb.Add("filter", v)
			}
			if v := rt.String("is-child"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("isChild", strconv.FormatBool(parsed))
			}
			if rt.Changed("index") {
				qb.Add("index", strconv.Itoa(rt.Int("index")))
			}
			if rt.Changed("count") {
				qb.Add("count", strconv.Itoa(rt.Int("count")))
			}

			path := "/resource_center/v1/host/group/host_list"
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
