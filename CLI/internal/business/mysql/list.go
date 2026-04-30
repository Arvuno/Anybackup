package mysql

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

func newObjectListCommand() meta.CommandMeta {
	allowedSortField := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}}
	allowedIsConfig := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}
	allowedObjectStatus := map[string]struct{}{"1": {}, "2": {}}
	allowedObjectMode := map[string]struct{}{"1": {}, "2": {}}
	allowedCanBackup := map[string]struct{}{"1": {}, "2": {}}
	allowedBindSlaStatus := map[string]struct{}{
		"1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {}, "8": {}, "9": {}, "10": {},
	}
	allowedBindSlaPresence := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}}
	allowedProtectStatus := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}
	allowedLastTaskStatus := map[string]struct{}{"600": {}, "700": {}, "800": {}, "900": {}, "1000": {}, "1100": {}}
	allowedExecStatus := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}

	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: []string{"mysql", "object", "list"},
		Use:           "list",
		Description:   "List MySQL objects",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "is-asc-sort", "", "Whether to sort ascending")
			rt.BindStringFlag(cmd.Flags(), "sort-field", "", "Sort field")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 20, "Page size")
			rt.BindStringFlag(cmd.Flags(), "filter", "", "Filter text")
			rt.BindStringFlag(cmd.Flags(), "included-path", "", "Whether filter should match path")
			rt.BindStringSliceFlag(cmd.Flags(), "production-system-ids", nil, "Production system IDs")
			rt.BindStringSliceFlag(cmd.Flags(), "object-ids", nil, "Object IDs")
			rt.BindStringSliceFlag(cmd.Flags(), "object-types", nil, "Object types")
			rt.BindStringSliceFlag(cmd.Flags(), "app-types", nil, "App types")
			rt.BindStringFlag(cmd.Flags(), "group-id", "", "Group ID")
			rt.BindStringFlag(cmd.Flags(), "host-id", "", "Host ID")
			rt.BindStringFlag(cmd.Flags(), "empty-host-id", "", "Whether to filter empty host ID")
			rt.BindStringFlag(cmd.Flags(), "host-tag", "", "Host tag")
			rt.BindStringSliceFlag(cmd.Flags(), "is-config", nil, "Config filters")
			rt.BindStringFlag(cmd.Flags(), "object-status", "", "Object status")
			rt.BindStringFlag(cmd.Flags(), "object-mode", "", "Object mode")
			rt.BindStringFlag(cmd.Flags(), "can-backup", "", "Whether object can be backed up")
			rt.BindStringSliceFlag(cmd.Flags(), "bind-sla-status", nil, "SLA binding status filters")
			rt.BindStringSliceFlag(cmd.Flags(), "included-bind-sla", nil, "Included SLA binding filters")
			rt.BindStringSliceFlag(cmd.Flags(), "exclude-bind-sla", nil, "Excluded SLA binding filters")
			rt.BindStringSliceFlag(cmd.Flags(), "bind-sla-ids", nil, "SLA IDs")
			rt.BindStringSliceFlag(cmd.Flags(), "protect-status", nil, "Protect status filters")
			rt.BindStringSliceFlag(cmd.Flags(), "last-backup-status", nil, "Last backup status filters")
			rt.BindStringSliceFlag(cmd.Flags(), "last-snapshot-status", nil, "Last snapshot status filters")
			rt.BindStringFlag(cmd.Flags(), "exec-status", "", "Execution status")
			rt.BindStringSliceFlag(cmd.Flags(), "parent-ids", nil, "Parent IDs")
			rt.BindStringFlag(cmd.Flags(), "all-child", "", "Whether to query all child objects")
			rt.BindStringFlag(cmd.Flags(), "name-filter", "", "Name filter")
			rt.BindStringFlag(cmd.Flags(), "app-type", "", "App type")
			rt.BindStringFlag(cmd.Flags(), "isolation", "", "Whether to query isolated objects")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if index := rt.Int("index"); index < 0 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --index")
			}
			if count := rt.Int("count"); count < 1 || count > 100 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --count")
			}
			if v := rt.String("is-asc-sort"); v != "" {
				if err := validateOptionalBool("is-asc-sort", v); err != nil {
					return err
				}
			}
			if v := rt.String("sort-field"); v != "" {
				if err := inputs.RequireOneOfString("sort-field", v, allowedSortField); err != nil {
					return err
				}
			}
			if err := validateOptionalMaxLength("filter", rt.String("filter"), 256); err != nil {
				return err
			}
			if v := rt.String("included-path"); v != "" {
				if err := validateOptionalBool("included-path", v); err != nil {
					return err
				}
			}
			if err := validateMaxItems("production-system-ids", rt.Strings("production-system-ids"), 100); err != nil {
				return err
			}
			if err := validateSliceExactLength("production-system-ids", rt.Strings("production-system-ids"), 32); err != nil {
				return err
			}
			if err := validateMaxItems("object-ids", rt.Strings("object-ids"), 100); err != nil {
				return err
			}
			if err := validateSliceExactLength("object-ids", rt.Strings("object-ids"), 32); err != nil {
				return err
			}
			if err := validateOptionalExactLength("host-id", rt.String("host-id"), 32); err != nil {
				return err
			}
			if v := rt.String("empty-host-id"); v != "" {
				if err := validateOptionalBool("empty-host-id", v); err != nil {
					return err
				}
			}
			if err := validateOptionalExactLength("host-tag", rt.String("host-tag"), 32); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("is-config", rt.Strings("is-config"), allowedIsConfig); err != nil {
				return err
			}
			if v := rt.String("object-status"); v != "" {
				if err := inputs.RequireOneOfString("object-status", v, allowedObjectStatus); err != nil {
					return err
				}
			}
			if v := rt.String("object-mode"); v != "" {
				if err := inputs.RequireOneOfString("object-mode", v, allowedObjectMode); err != nil {
					return err
				}
			}
			if v := rt.String("can-backup"); v != "" {
				if err := inputs.RequireOneOfString("can-backup", v, allowedCanBackup); err != nil {
					return err
				}
			}
			if err := inputs.RequireAllInSet("bind-sla-status", rt.Strings("bind-sla-status"), allowedBindSlaStatus); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("included-bind-sla", rt.Strings("included-bind-sla"), allowedBindSlaPresence); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("exclude-bind-sla", rt.Strings("exclude-bind-sla"), allowedBindSlaPresence); err != nil {
				return err
			}
			if err := validateMaxItems("bind-sla-ids", rt.Strings("bind-sla-ids"), 10); err != nil {
				return err
			}
			if err := validateSliceExactLength("bind-sla-ids", rt.Strings("bind-sla-ids"), 32); err != nil {
				return err
			}
			if err := validateMaxItems("protect-status", rt.Strings("protect-status"), 4); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("protect-status", rt.Strings("protect-status"), allowedProtectStatus); err != nil {
				return err
			}
			if err := validateMaxItems("last-backup-status", rt.Strings("last-backup-status"), 6); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("last-backup-status", rt.Strings("last-backup-status"), allowedLastTaskStatus); err != nil {
				return err
			}
			if err := validateMaxItems("last-snapshot-status", rt.Strings("last-snapshot-status"), 6); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("last-snapshot-status", rt.Strings("last-snapshot-status"), allowedLastTaskStatus); err != nil {
				return err
			}
			if v := rt.String("exec-status"); v != "" {
				if err := inputs.RequireOneOfString("exec-status", v, allowedExecStatus); err != nil {
					return err
				}
			}
			if err := validateMaxItems("parent-ids", rt.Strings("parent-ids"), 100); err != nil {
				return err
			}
			if err := validateSliceMaxLength("parent-ids", rt.Strings("parent-ids"), 36); err != nil {
				return err
			}
			if v := rt.String("all-child"); v != "" {
				if err := validateOptionalBool("all-child", v); err != nil {
					return err
				}
			}
			if err := validateOptionalMaxLength("name-filter", rt.String("name-filter"), 1024); err != nil {
				return err
			}
			if err := validateOptionalInt32("app-type", rt.String("app-type")); err != nil {
				return err
			}
			if v := rt.String("isolation"); v != "" {
				if err := validateOptionalBool("isolation", v); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			index := rt.Int("index")
			count := rt.Int("count")
			if count == 0 {
				count = 20
			}

			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", index)).
				Add("count", fmt.Sprintf("%d", count))

			if v := rt.String("is-asc-sort"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("isAscSort", strconv.FormatBool(parsed))
			}
			if v := rt.String("sort-field"); v != "" {
				qb.Add("sortField", v)
			}
			if v := rt.String("filter"); v != "" {
				qb.Add("filter", v)
			}
			if v := rt.String("included-path"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("includedPath", strconv.FormatBool(parsed))
			}
			if values := rt.Strings("production-system-ids"); len(values) > 0 {
				qb.AddAll("productionSystemIds", values)
			}
			if values := rt.Strings("object-ids"); len(values) > 0 {
				qb.AddAll("objectIds", values)
			}
			if values := rt.Strings("object-types"); len(values) > 0 {
				qb.AddAll("objectTypes", values)
			}
			if values := rt.Strings("app-types"); len(values) > 0 {
				qb.AddAll("appTypes", values)
			}
			if v := rt.String("group-id"); v != "" {
				qb.Add("groupId", v)
			}
			if v := rt.String("host-id"); v != "" {
				qb.Add("hostId", v)
			}
			if v := rt.String("empty-host-id"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("emptyHostId", strconv.FormatBool(parsed))
			}
			if v := rt.String("host-tag"); v != "" {
				qb.Add("hostTag", v)
			}
			if values := rt.Strings("is-config"); len(values) > 0 {
				qb.AddAll("isConfig", values)
			}
			if v := rt.String("object-status"); v != "" {
				qb.Add("objectStatus", v)
			}
			if v := rt.String("object-mode"); v != "" {
				qb.Add("objectMode", v)
			}
			if v := rt.String("can-backup"); v != "" {
				qb.Add("canBackup", v)
			}
			if values := rt.Strings("bind-sla-status"); len(values) > 0 {
				qb.AddAll("bindSlaStatus", values)
			}
			if values := rt.Strings("included-bind-sla"); len(values) > 0 {
				qb.AddAll("includedBindSla", values)
			}
			if values := rt.Strings("exclude-bind-sla"); len(values) > 0 {
				qb.AddAll("excludeBindSla", values)
			}
			if values := rt.Strings("bind-sla-ids"); len(values) > 0 {
				qb.AddAll("bindSlaIds", values)
			}
			if values := rt.Strings("protect-status"); len(values) > 0 {
				qb.AddAll("protectStatus", values)
			}
			if values := rt.Strings("last-backup-status"); len(values) > 0 {
				qb.AddAll("lastBackupStatus", values)
			}
			if values := rt.Strings("last-snapshot-status"); len(values) > 0 {
				qb.AddAll("lastSnapshotStatus", values)
			}
			if v := rt.String("exec-status"); v != "" {
				qb.Add("execStatus", v)
			}
			if values := rt.Strings("parent-ids"); len(values) > 0 {
				qb.AddAll("parentIds", values)
			}
			if v := rt.String("all-child"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("allChild", strconv.FormatBool(parsed))
			}
			if v := rt.String("name-filter"); v != "" {
				qb.Add("nameFilter", v)
			}
			if v := rt.String("app-type"); v != "" {
				qb.Add("appType", v)
			} else {
				qb.Add("appType", "202")
			}
			if v := rt.String("isolation"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("isolation", strconv.FormatBool(parsed))
			}
			if v := rt.String("object-id"); v != "" {
				qb.Add("objectId", v)
			}

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/database/object_list?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
}

func validateOptionalBool(flag, value string) error {
	if _, err := strconv.ParseBool(value); err != nil {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateOptionalMaxLength(flag, value string, max int) error {
	if value == "" {
		return nil
	}
	if len(value) > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateOptionalExactLength(flag, value string, length int) error {
	if value == "" {
		return nil
	}
	if len(value) != length {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateOptionalInt32(flag, value string) error {
	if value == "" {
		return nil
	}
	if _, err := strconv.ParseInt(value, 10, 32); err != nil {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateMaxItems(flag string, values []string, max int) error {
	if len(values) > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("too many --%s values", flag))
	}
	return nil
}

func validateSliceExactLength(flag string, values []string, length int) error {
	for _, value := range values {
		if len(value) != length {
			return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
		}
	}
	return nil
}

func validateSliceMaxLength(flag string, values []string, max int) error {
	for _, value := range values {
		if len(value) > max {
			return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
		}
	}
	return nil
}
