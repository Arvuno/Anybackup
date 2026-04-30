package policy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newListCommand() meta.CommandMeta {
	allowedValidateStatus := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}}
	allowedDisableMark := map[string]struct{}{"-1": {}, "0": {}, "1": {}, "2": {}}
	allowedCopyMode := map[string]struct{}{"1": {}, "2": {}}
	allowedBackupMode := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}}

	return meta.CommandMeta{
		Domain:        "policy",
		CanonicalPath: []string{"policy", "list"},
		Use:           "list",
		Description:   "List policy templates",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 10, "Page size")
			rt.BindStringFlag(cmd.Flags(), "filter", "", "Filter text")
			rt.BindStringFlag(cmd.Flags(), "validate-status", "", "Validate status")
			rt.BindStringFlag(cmd.Flags(), "disable-mark", "", "Disable mark")
			rt.BindStringFlag(cmd.Flags(), "type", "", "Business type")
			rt.BindStringSliceFlag(cmd.Flags(), "types", nil, "Business type list")
			rt.BindStringFlag(cmd.Flags(), "copy-mode", "", "Copy mode")
			rt.BindStringFlag(cmd.Flags(), "backup-mode", "", "Backup mode")
		},
		Validate: func(rt *meta.Runtime) error {
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
			if v := rt.String("copy-mode"); v != "" {
				if err := inputs.RequireOneOfString("copy-mode", v, allowedCopyMode); err != nil {
					return err
				}
			}
			if v := rt.String("backup-mode"); v != "" {
				if err := inputs.RequireOneOfString("backup-mode", v, allowedBackupMode); err != nil {
					return err
				}
			}

			if err := validateOptionalType("type", rt.String("type")); err != nil {
				return err
			}
			if len(rt.Strings("types")) > 10 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "too many --types values")
			}
			for _, v := range rt.Strings("types") {
				if err := validateOptionalType("types", v); err != nil {
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
			if values := rt.Strings("types"); len(values) > 0 {
				qb.AddAll("types", values)
			}
			if v := rt.String("copy-mode"); v != "" {
				qb.Add("copyMode", v)
			}
			if v := rt.String("backup-mode"); v != "" {
				qb.Add("backupMode", v)
			}

			path := "/api/sla/v1/templates"
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

func validateOptionalType(flag, value string) error {
	if value == "" {
		return nil
	}
	if n := len(value); n < 1 || n > 64 {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}
