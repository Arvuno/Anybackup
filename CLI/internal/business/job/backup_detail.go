package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newBackupDetailCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "backup-detail"},
		Use:           "backup-detail",
		Description:   "Get backup task detail by job id",
		Risk:          meta.RiskReadOnly,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      true,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "job-id", "", "job id")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("job-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --job-id")
			}
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			if deps.Console == nil {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "console executor is not configured")
			}
			jobID := rt.String("job-id")

			primaryBody, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/task/%s/detail", jobID),
				ReadOnly: true,
			})
			if err != nil {
				return err
			}

			jobBody, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/job_center/v1/job/%s", jobID),
				ReadOnly: true,
			})
			if err != nil {
				return output.WriteRaw(deps.Streams, primaryBody)
			}

			merged := injectObjectName(primaryBody, jobBody)
			return output.WriteRaw(deps.Streams, merged)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodGet,
				Path:     fmt.Sprintf("/backupmgm/v1/task/%s/detail", rt.String("job-id")),
				ReadOnly: true,
			}, nil
		},
	}
}

func injectObjectName(primaryBody, jobBody []byte) []byte {
	var primary map[string]interface{}
	if err := json.Unmarshal(primaryBody, &primary); err != nil {
		return primaryBody
	}

	objectName := extractObjectName(jobBody)
	if objectName == "" {
		return primaryBody
	}

	// Inject objectName into responseData when the payload uses the standard envelope.
	if responseData, ok := primary["responseData"].(map[string]interface{}); ok {
		responseData["objectName"] = objectName
	} else {
		// Fall back to a top-level objectName when the detail payload is flat.
		primary["objectName"] = objectName
	}

	out, err := json.Marshal(primary)
	if err != nil {
		return primaryBody
	}
	return out
}

func extractObjectName(jobBody []byte) string {
	var payload map[string]interface{}
	if err := json.Unmarshal(jobBody, &payload); err != nil {
		return ""
	}
	if responseData, ok := payload["responseData"].(map[string]interface{}); ok {
		if v, ok := responseData["objectName"].(string); ok {
			return v
		}
	}
	if v, ok := payload["objectName"].(string); ok {
		return v
	}
	return ""
}
