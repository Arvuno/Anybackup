package timepoint

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
	allowedBusiness := map[string]struct{}{"1": {}, "4": {}, "5": {}}
	allowedBusinesses := map[string]struct{}{"1": {}, "4": {}, "5": {}, "6": {}}
	allowedUsable := map[string]struct{}{"1": {}, "2": {}}
	allowedTimePointType := map[string]struct{}{"0": {}, "1": {}, "2": {}, "3": {}}
	allowedBackupTypes := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {}}

	return meta.CommandMeta{
		Domain:        "timepoint",
		CanonicalPath: []string{"timepoint", "list"},
		Use:           "list",
		Description:   "List time points for a protect object",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "object id")
			rt.BindStringFlag(cmd.Flags(), "business", "", "Time point business")
			rt.BindStringFlag(cmd.Flags(), "start-time", "", "Start time")
			rt.BindStringFlag(cmd.Flags(), "end-time", "", "End time")
			rt.BindStringFlag(cmd.Flags(), "storage-pool-id", "", "Storage pool ID")
			rt.BindStringFlag(cmd.Flags(), "is-duplication", "", "Is duplication")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "data-set-id", "", "Data set ID")
			rt.BindStringSliceFlag(cmd.Flags(), "businesses", nil, "Time point businesses")
			rt.BindStringFlag(cmd.Flags(), "expiration-start-time", "", "Expiration start time")
			rt.BindStringFlag(cmd.Flags(), "expiration-end-time", "", "Expiration end time")
			rt.BindStringFlag(cmd.Flags(), "usable", "", "Usable state")
			rt.BindStringSliceFlag(cmd.Flags(), "backup-types", nil, "Backup types")
			rt.BindStringSliceFlag(cmd.Flags(), "include-storage-types", nil, "Include storage types")
			rt.BindStringSliceFlag(cmd.Flags(), "exclude-storage-types", nil, "Exclude storage types")
			rt.BindStringFlag(cmd.Flags(), "time-point-type", "", "Time point type")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("object-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --object-id")
			}
			if v := rt.String("business"); v != "" {
				if err := inputs.RequireOneOfString("business", v, allowedBusiness); err != nil {
					return err
				}
			}
			if err := validateOptionalNonNegativeInt64("start-time", rt.String("start-time")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("end-time", rt.String("end-time")); err != nil {
				return err
			}
			if err := validateOptionalLength("storage-pool-id", rt.String("storage-pool-id"), 32, 48); err != nil {
				return err
			}
			if v := rt.String("is-duplication"); v != "" {
				if _, err := strconv.ParseBool(v); err != nil {
					return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --is-duplication")
				}
			}
			if err := validateOptionalLength("storage-service-id", rt.String("storage-service-id"), 32, 48); err != nil {
				return err
			}
			if err := validateOptionalLength("data-set-id", rt.String("data-set-id"), 32, 48); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("businesses", rt.Strings("businesses"), allowedBusinesses); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("expiration-start-time", rt.String("expiration-start-time")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("expiration-end-time", rt.String("expiration-end-time")); err != nil {
				return err
			}
			if v := rt.String("usable"); v != "" {
				if err := inputs.RequireOneOfString("usable", v, allowedUsable); err != nil {
					return err
				}
			}
			if err := inputs.RequireAllInSet("backup-types", rt.Strings("backup-types"), allowedBackupTypes); err != nil {
				return err
			}
			if v := rt.String("time-point-type"); v != "" {
				if err := inputs.RequireOneOfString("time-point-type", v, allowedTimePointType); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			if v := rt.String("business"); v != "" {
				qb.Add("business", v)
			}
			if v := rt.String("start-time"); v != "" {
				qb.Add("startTime", v)
			}
			if v := rt.String("end-time"); v != "" {
				qb.Add("endTime", v)
			}
			if v := rt.String("storage-pool-id"); v != "" {
				qb.Add("storagePoolId", v)
			}
			if v := rt.String("is-duplication"); v != "" {
				parsed, _ := strconv.ParseBool(v)
				qb.Add("isDuplication", strconv.FormatBool(parsed))
			}
			if v := rt.String("storage-service-id"); v != "" {
				qb.Add("storageServiceId", v)
			}
			if v := rt.String("data-set-id"); v != "" {
				qb.Add("dataSetId", v)
			}
			if values := rt.Strings("businesses"); len(values) > 0 {
				qb.AddAll("businesses", values)
			}
			if v := rt.String("expiration-start-time"); v != "" {
				qb.Add("expirationStartTime", v)
			}
			if v := rt.String("expiration-end-time"); v != "" {
				qb.Add("expirationEndTime", v)
			}
			if v := rt.String("usable"); v != "" {
				qb.Add("usable", v)
			}
			if values := rt.Strings("backup-types"); len(values) > 0 {
				qb.AddAll("backupTypes", values)
			}
			if values := rt.Strings("include-storage-types"); len(values) > 0 {
				qb.AddAll("includeStorageTypes", values)
			}
			if values := rt.Strings("exclude-storage-types"); len(values) > 0 {
				qb.AddAll("excludeStorageTypes", values)
			}
			if v := rt.String("time-point-type"); v != "" {
				qb.Add("timePointType", v)
			}

			path := fmt.Sprintf("/backupmgm/v1/protect_object/%s/time_points", rt.String("object-id"))
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

func validateOptionalNonNegativeInt64(flag string, value string) error {
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateOptionalLength(flag string, value string, min int, max int) error {
	if value == "" {
		return nil
	}
	if n := len(value); n < min || n > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}
