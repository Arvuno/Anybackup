package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newDeployCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "deploy"},
		Use:           "deploy",
		Description:   "Deploy client configuration",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "os", "", "Target OS: linux/windows/unix")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if err := inputs.RequireNonEmpty("os", rt.String("os")); err != nil {
				return err
			}
			if _, ok := deployConfigOSPathSegment(rt.String("os")); !ok {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --os")
			}
			if len(rt.Body) == 0 {
				return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
			}
			if _, err := deployNameFromBody(rt.Body); err != nil {
				return err
			}
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			name, err := deployNameFromBody(rt.Body)
			if err != nil {
				return err
			}

			hostConfigPrecheckResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/deploy/v1/hostConfig/nameExist?name=" + url.QueryEscape(name),
				ReadOnly: true,
			})
			if err != nil {
				return err
			}
			hostConfigExists, err := parseDeployConfigNameExists(hostConfigPrecheckResp)
			if err != nil {
				return err
			}
			if hostConfigExists {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "deploy config name already exists")
			}

			osPath, _ := deployConfigOSPathSegment(rt.String("os"))
			hostConfigBody, err := ensureHostConfigCreateBody(rt.Body, name)
			if err != nil {
				return err
			}
			hostConfigCreateResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/deploy/v1/hostConfig/" + osPath,
				Body:     hostConfigBody,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			configIDs, err := extractConfigIDsFromCreateResponse(hostConfigCreateResp)
			if err != nil {
				return err
			}

			precheckResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/deploy/v1/job/config/nameExist?name=" + url.QueryEscape(name),
				ReadOnly: true,
			})
			if err != nil {
				return err
			}
			exists, err := parseDeployConfigNameExists(precheckResp)
			if err != nil {
				return err
			}
			if exists {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "deploy job config name already exists")
			}

			deployBody, err := mergeDeployHostListIDs(rt.Body, configIDs)
			if err != nil {
				return err
			}
			body, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/deploy/v1/job/config",
				Body:     deployBody,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/deploy/v1/job/config",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

type deployConfigCreateResponse struct {
	Status       string `json:"status"`
	ResponseData []any  `json:"responseData"`
}

func deployNameFromBody(body []byte) (string, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return "", clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	for _, key := range []string{"name", "jobName"} {
		if v, ok := req[key]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s), nil
			}
		}
	}
	return "", clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field name or jobName")
}

func ensureHostConfigCreateBody(body []byte, name string) ([]byte, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	if _, ok := payload["name"]; !ok {
		payload["name"] = name
	}
	out, err := json.Marshal(payload)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	return out, nil
}

func extractConfigIDsFromCreateResponse(body []byte) ([]string, error) {
	var resp deployConfigCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid deploy config create response")
	}
	if strings.TrimSpace(resp.Status) == "" {
		return nil, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid deploy config create response")
	}
	if len(resp.ResponseData) == 0 {
		return nil, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "deploy config create response missing id array")
	}

	ids := make([]string, 0, len(resp.ResponseData))
	for _, item := range resp.ResponseData {
		id, ok := item.(string)
		if !ok || strings.TrimSpace(id) == "" {
			return nil, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "deploy config create response contains invalid id")
		}
		ids = append(ids, strings.TrimSpace(id))
	}
	return ids, nil
}

func mergeDeployHostListIDs(body []byte, ids []string) ([]byte, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	idValues := make([]any, 0, len(ids))
	for _, id := range ids {
		idValues = append(idValues, id)
	}

	// /deploy/v1/job/config expects hostList as deploy target list.
	// In merged mode, hostList must carry IDs returned by /deploy/v1/hostConfig/*.
	payload["hostList"] = idValues

	out, err := json.Marshal(payload)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	return out, nil
}
