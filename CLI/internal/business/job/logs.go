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

func newLogsCommand() meta.CommandMeta {
	allowedLevel := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}

	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "logs"},
		Use:           "logs",
		Description:   "Get job activity logs",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "job-id", "", "job id")
			rt.BindIntFlag(cmd.Flags(), "index", 0, "Page index")
			rt.BindIntFlag(cmd.Flags(), "count", 30, "Page size")
			rt.BindStringFlag(cmd.Flags(), "start-time", "", "Start time in milliseconds")
			rt.BindStringFlag(cmd.Flags(), "end-time", "", "End time in milliseconds")
			rt.BindStringFlag(cmd.Flags(), "level", "", "Log level")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("job-id", rt.String("job-id")); err != nil {
				return err
			}
			if err := inputs.ValidatePaging(rt.Int("index"), rt.Int("count")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("start-time", rt.String("start-time")); err != nil {
				return err
			}
			if err := validateOptionalNonNegativeInt64("end-time", rt.String("end-time")); err != nil {
				return err
			}
			if v := rt.String("level"); v != "" {
				if err := inputs.RequireOneOfString("level", v, allowedLevel); err != nil {
					return err
				}
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			queryBuilder := inputs.NewQueryBuilder().
				Add("index", fmt.Sprintf("%d", rt.Int("index"))).
				Add("count", fmt.Sprintf("%d", rt.Int("count")))
			if v := rt.String("start-time"); v != "" {
				queryBuilder.Add("startTime", v)
			}
			if v := rt.String("end-time"); v != "" {
				queryBuilder.Add("endTime", v)
			}
			if v := rt.String("level"); v != "" {
				queryBuilder.Add("level", v)
			}
			query := queryBuilder.Encode()

			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/job_center/v1/activity/%s/logs?%s", rt.String("job-id"), query),
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
