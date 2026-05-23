package job

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
	allowedStatuses := map[string]struct{}{
		"0": {}, "10": {}, "100": {}, "200": {}, "300": {}, "400": {}, "500": {}, "600": {},
		"700": {}, "800": {}, "900": {}, "1000": {}, "1100": {}, "1200": {},
	}
	allowedSort := map[string]struct{}{
		"startTime":  {},
		"-startTime": {},
		"endTime":    {},
		"-endTime":   {},
		"status":     {},
	}
	allowedTaskContinueType := map[string]struct{}{"0": {}, "1": {}}

	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "list"},
		Use:           "list",
		Description:   "List jobs",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 10, "Page size")
			rt.BindStringSliceFlag(cmd.Flags(), "statuses", nil, "Job statuses")
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			rt.BindStringFlag(cmd.Flags(), "object-name", "", "Object name")
			rt.BindStringSliceFlag(cmd.Flags(), "app-types", nil, "App types")
			rt.BindStringSliceFlag(cmd.Flags(), "operation-types", nil, "Operation types")
			rt.BindStringSliceFlag(cmd.Flags(), "business-types", nil, "Business types")
			rt.BindStringFlag(cmd.Flags(), "start-time", "", "Start time in milliseconds")
			rt.BindStringFlag(cmd.Flags(), "end-time", "", "End time in milliseconds")
			rt.BindStringFlag(cmd.Flags(), "sort", "", "Sort field")
			rt.BindStringFlag(cmd.Flags(), "remark", "", "Remark")
			rt.BindStringFlag(cmd.Flags(), "cluster-id", "", "Cluster ID")
			rt.BindStringFlag(cmd.Flags(), "storage-service-id", "", "Storage service ID")
			rt.BindStringFlag(cmd.Flags(), "storage-pool-id", "", "Storage pool ID")
			rt.BindStringFlag(cmd.Flags(), "strategy-name", "", "Strategy name")
			rt.BindStringFlag(cmd.Flags(), "client-name", "", "Client name")
			rt.BindStringFlag(cmd.Flags(), "task-continue-type", "", "Task continue type")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
				return err
			}
			if err := validateMaxItems("statuses", rt.Strings("statuses"), 13); err != nil {
				return err
			}
			if err := inputs.RequireAllInSet("statuses", rt.Strings("statuses"), allowedStatuses); err != nil {
				return err
			}
			if err := validateOptionalLengthRange("object-id", rt.String("object-id"), 1, 48); err != nil {
				return err
			}
			if err := validateOptionalLengthRange("object-name", rt.String("object-name"), 1, 256); err != nil {
				return err
			}
			if err := validateMaxItems("app-types", rt.Strings("app-types"), 100); err != nil {
				return err
			}
			if err := validateStringSliceLengthRange("app-types", rt.Strings("app-types"), 1, 255); err != nil {
				return err
			}
			if err := validateMaxItems("operation-types", rt.Strings("operation-types"), 10); err != nil {
				return err
			}
			if err := validateStringSliceNonNegativeInt64("operation-types", rt.Strings("operation-types")); err != nil {
				return err
			}
			if err := validateMaxItems("business-types", rt.Strings("business-types"), 20); err != nil {
				return err
			}
			if err := validateStringSliceNonNegativeInt64("business-types", rt.Strings("business-types")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("start-time", rt.String("start-time")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("end-time", rt.String("end-time")); err != nil {
				return err
			}
			if v := rt.String("sort"); v != "" {
				if err := inputs.RequireOneOfString("sort", v, allowedSort); err != nil {
					return err
				}
			}
			if v := rt.String("task-continue-type"); v != "" {
				if err := inputs.RequireOneOfString("task-continue-type", v, allowedTaskContinueType); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			qb := inputs.NewQueryBuilder()
			index := rt.Int("index")
			count := rt.Int("count")
			if index > 0 {
				qb.Add("index", fmt.Sprintf("%d", index))
			}
			if count > 0 && count != 10 {
				qb.Add("count", fmt.Sprintf("%d", count))
			}
			if values := rt.Strings("statuses"); len(values) > 0 {
				qb.AddAll("statuses", values)
			}
			if v := rt.String("object-id"); v != "" {
				qb.Add("objectId", v)
			}
			if v := rt.String("object-name"); v != "" {
				qb.Add("objectName", v)
			}
			if values := rt.Strings("app-types"); len(values) > 0 {
				qb.AddAll("appTypes", values)
			}
			if values := rt.Strings("operation-types"); len(values) > 0 {
				qb.AddAll("operationTypes", values)
			}
			if values := rt.Strings("business-types"); len(values) > 0 {
				qb.AddAll("businessTypes", values)
			}
			if v := rt.String("start-time"); v != "" {
				qb.Add("startTime", v)
			}
			if v := rt.String("end-time"); v != "" {
				qb.Add("endTime", v)
			}
			if v := rt.String("sort"); v != "" {
				qb.Add("sort", v)
			}
			if v := rt.String("remark"); v != "" {
				qb.Add("remark", v)
			}
			if v := rt.String("cluster-id"); v != "" {
				qb.Add("clusterId", v)
			}
			if v := rt.String("storage-service-id"); v != "" {
				qb.Add("storageServiceId", v)
			}
			if v := rt.String("storage-pool-id"); v != "" {
				qb.Add("storagePoolId", v)
			}
			if v := rt.String("strategy-name"); v != "" {
				qb.Add("strategyName", v)
			}
			if v := rt.String("client-name"); v != "" {
				qb.Add("clientName", v)
			}
			if v := rt.String("task-continue-type"); v != "" {
				qb.Add("taskContinueType", v)
			}
			encoded := qb.Encode()
			path := "/job_center/v1/jobs"
			if encoded != "" {
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

func validateMaxItems(flag string, values []string, max int) error {
	if len(values) > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("too many --%s values", flag))
	}
	return nil
}

func validateOptionalLengthRange(flag, value string, min, max int) error {
	if value == "" {
		return nil
	}
	if n := len(value); n < min || n > max {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
	}
	return nil
}

func validateStringSliceLengthRange(flag string, values []string, min, max int) error {
	for _, value := range values {
		if n := len(value); n < min || n > max {
			return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
		}
	}
	return nil
}

func validateStringSliceNonNegativeInt64(flag string, values []string) error {
	for _, value := range values {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 0 {
			return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, fmt.Sprintf("invalid --%s", flag))
		}
	}
	return nil
}
