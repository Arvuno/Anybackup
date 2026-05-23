package job

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newStopCommand() meta.CommandMeta {
	var jobIDs []string

	return meta.CommandMeta{
		Domain:        "job",
		CanonicalPath: []string{"job", "stop"},
		Use:           "stop",
		Description:   "Stop jobs",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
			cmd.Flags().StringArrayVar(&jobIDs, "job-id", nil, "job id (repeatable)")
		},
		Validate: func(rt *meta.Runtime) error {
			if len(rt.Body) > 0 && len(jobIDs) > 0 {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "--job-id cannot be used with --data/--body-file")
			}
			if len(rt.Body) == 0 {
				normalized, err := normalizeJobIDs(jobIDs)
				if err != nil {
					return err
				}
				if len(normalized) == 0 {
					return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
				}
				body, err := json.Marshal(map[string][]string{"jobIds": normalized})
				if err != nil {
					return clierrors.Wrap(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "failed to build request body from --job-id", err)
				}
				rt.Body = body
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/job_center/v1/jobs/stop",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

func normalizeJobIDs(values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, nil
	}

	out := make([]string, 0, len(values))
	for _, value := range values {
		jobID := strings.TrimSpace(value)
		if jobID == "" {
			return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --job-id")
		}
		out = append(out, jobID)
	}
	return out, nil
}
