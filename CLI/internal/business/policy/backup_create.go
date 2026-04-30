package policy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newBackupCreateCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "policy",
		CanonicalPath: []string{"policy", "backup", "create"},
		Use:           "create",
		Description:   "Create policy backup info",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := ensureNormalizedBackupCreateBody(rt); err != nil {
				return err
			}
			if len(rt.Body) == 0 {
				return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
			}
			if _, err := policyNameFromBody(rt.Body); err != nil {
				return err
			}
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			if err := ensureNormalizedBackupCreateBody(rt); err != nil {
				return err
			}
			name, err := policyNameFromBody(rt.Body)
			if err != nil {
				return err
			}

			precheckResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/api/sla/v1/common/name/exists?name=" + url.QueryEscape(name),
				ReadOnly: true,
			})
			if err != nil {
				return err
			}
			exists, err := parsePolicyNameExists(precheckResp)
			if err != nil {
				return err
			}
			if exists {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "policy backup name already exists")
			}

			body, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/api/sla/v1/group/backup_info",
				Body:     rt.Body,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			if err := ensureNormalizedBackupCreateBody(rt); err != nil {
				return console.RequestSpec{}, err
			}
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/api/sla/v1/group/backup_info",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

func ensureNormalizedBackupCreateBody(rt *meta.Runtime) error {
	if len(rt.Body) == 0 {
		return nil
	}
	body, err := normalizeBackupCreateBody(rt.Body)
	if err != nil {
		return err
	}
	rt.Body = body
	return nil
}

func normalizeBackupCreateBody(body []byte) ([]byte, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	if _, ok := req["validatePeriod"]; !ok || req["validatePeriod"] == nil {
		req["validatePeriod"] = 1
	}
	if _, ok := req["effectiveType"]; !ok || req["effectiveType"] == nil {
		req["effectiveType"] = 1
	}
	normalized, err := json.Marshal(req)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	return normalized, nil
}

type policyCreateBody struct {
	Name string `json:"name"`
}

type policyNameExistsResponse struct {
	Status       string `json:"status"`
	ResponseData struct {
		IsExists bool `json:"isExists"`
	} `json:"responseData"`
}

func policyNameFromBody(body []byte) (string, error) {
	var req policyCreateBody
	if err := json.Unmarshal(body, &req); err != nil {
		return "", clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field name")
	}
	return name, nil
}

func parsePolicyNameExists(body []byte) (bool, error) {
	var resp policyNameExistsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid policy name exists response")
	}
	if strings.TrimSpace(resp.Status) == "" {
		return false, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid policy name exists response")
	}
	return resp.ResponseData.IsExists, nil
}
