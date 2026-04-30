package host

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
	allowedDatasourceType := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}, "4": {}, "5": {}}
	allowedIntelligentArchive := map[string]struct{}{"0": {}, "1": {}, "2": {}}
	allowedExecStatus := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}, "4": {}}
	allowedHasTp := map[string]struct{}{"0": {}, "1": {}, "2": {}}
	allowedIsBackupConfig := map[string]struct{}{"1": {}, "2": {}}

	return meta.CommandMeta{
		Domain:        "host",
		CanonicalPath: []string{"host", "object", "list"},
		Use:           "list",
		Description:   "List host objects",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "production-system-id", "", "Production system ID")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 10, "Page size")
			rt.BindStringFlag(cmd.Flags(), "name", "", "Protect object name")
			rt.BindStringFlag(cmd.Flags(), "datasource-type", "", "Datasource type")
			rt.BindStringFlag(cmd.Flags(), "intelligent-archive", "", "Intelligent archive")
			rt.BindStringFlag(cmd.Flags(), "exec-status", "", "Execution status")
			rt.BindStringFlag(cmd.Flags(), "is-include-tenant-id", "", "Whether to include tenant ID")
			rt.BindStringFlag(cmd.Flags(), "query-tenant-id", "", "Tenant ID filter")
			rt.BindStringFlag(cmd.Flags(), "group-id", "", "Group ID")
			rt.BindStringFlag(cmd.Flags(), "has-tp", "", "Whether host object has timepoints")
			rt.BindStringFlag(cmd.Flags(), "is-backup-config", "", "Backup config filter")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Protect object ID")
		},
		Validate: func(rt *meta.Runtime) error {
			if index := rt.Int("index"); index < 0 || index > 20000 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --index")
			}
			if count := rt.Int("count"); count < 1 || count > 100 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --count")
			}
			if err := validateOptionalExactLength("production-system-id", rt.String("production-system-id"), 32); err != nil {
				return err
			}
			if v := rt.String("datasource-type"); v != "" {
				if err := inputs.RequireOneOfString("datasource-type", v, allowedDatasourceType); err != nil {
					return err
				}
			}
			if v := rt.String("intelligent-archive"); v != "" {
				if err := inputs.RequireOneOfString("intelligent-archive", v, allowedIntelligentArchive); err != nil {
					return err
				}
			}
			if v := rt.String("exec-status"); v != "" {
				if err := inputs.RequireOneOfString("exec-status", v, allowedExecStatus); err != nil {
					return err
				}
			}
			if v := rt.String("is-include-tenant-id"); v != "" {
				if _, err := strconv.ParseBool(v); err != nil {
					return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --is-include-tenant-id")
				}
			}
			if v := rt.String("has-tp"); v != "" {
				if err := inputs.RequireOneOfString("has-tp", v, allowedHasTp); err != nil {
					return err
				}
			}
			if v := rt.String("is-backup-config"); v != "" {
				if err := inputs.RequireOneOfString("is-backup-config", v, allowedIsBackupConfig); err != nil {
					return err
				}
			}
			if err := validateOptionalExactLength("object-id", rt.String("object-id"), 32); err != nil {
				return err
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			index := rt.Int("index")
			count := rt.Int("count")
			if count == 0 {
				count = 10
			}
			qb := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", index)).
				Add("count", fmt.Sprintf("%d", count))
			if v := rt.String("production-system-id"); v != "" {
				qb.Add("productionSystemId", v)
			}
			if v := rt.String("name"); v != "" {
				qb.Add("name", v)
			}
			if v := rt.String("datasource-type"); v != "" {
				qb.Add("datasourceType", v)
			}
			if v := rt.String("intelligent-archive"); v != "" {
				qb.Add("intelligentArchive", v)
			}
			if v := rt.String("exec-status"); v != "" {
				qb.Add("execStatus", v)
			}
			if v := rt.String("is-include-tenant-id"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("isIncludeTenantId", strconv.FormatBool(parsed))
			}
			if v := rt.String("query-tenant-id"); v != "" {
				qb.Add("tenantId", v)
			}
			if v := rt.String("group-id"); v != "" {
				qb.Add("groupId", v)
			}
			if v := rt.String("has-tp"); v != "" {
				qb.Add("hasTp", v)
			}
			if v := rt.String("is-backup-config"); v != "" {
				qb.Add("isBackupConfig", v)
			}
			if v := rt.String("object-id"); v != "" {
				qb.Add("objectId", v)
			}

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/backupmgm/v1/file/object_list?" + qb.Encode(),
				ReadOnly: true,
			}, nil
		},
	}
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
