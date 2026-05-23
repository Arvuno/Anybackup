package protect

import (
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPolicyBindListCommand() meta.CommandMeta {
	allowedValidateStatus := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}}
	allowedDisableMark := map[string]struct{}{"-1": {}, "0": {}, "1": {}, "2": {}}

	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "policy", "bind-list"},
		Use:           "bind-list",
		Description:   "List policy bindings for a protect object",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 10, "Page size")
			rt.BindStringFlag(cmd.Flags(), "filter", "", "Filter text")
			rt.BindStringFlag(cmd.Flags(), "validate-status", "", "Validate status")
			rt.BindStringFlag(cmd.Flags(), "disable-mark", "", "Disable mark")
			rt.BindStringFlag(cmd.Flags(), "type", "", "Business type")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("object-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --object-id")
			}
			if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
				return err
			}
			if v := rt.String("validate-status"); v != "" {
				if err := inputs.RequireOneOfString("validate-status", v, allowedValidateStatus); err != nil {
					return err
				}
			}
			if v := rt.String("disable-mark"); v != "" {
				if err := inputs.RequireOneOfString("disable-mark", v, allowedDisableMark); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if index := rt.Int("index"); index > 0 {
				qb.Add("index", strconv.Itoa(index))
			}
			if count := rt.Int("count"); count > 0 && count != 10 {
				qb.Add("count", strconv.Itoa(count))
			}
			if v := rt.String("filter"); v != "" {
				qb.Add("filter", v)
			}
			if v := rt.String("validate-status"); v != "" {
				qb.Add("validateStatus", v)
			}
			if v := rt.String("disable-mark"); v != "" {
				qb.Add("disableMark", v)
			}
			if v := rt.String("type"); v != "" {
				qb.Add("type", v)
			}

			path := "/api/sla/v1/group/object/" + rt.String("object-id") + "/templates"
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
