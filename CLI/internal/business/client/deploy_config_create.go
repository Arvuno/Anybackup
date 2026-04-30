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

func newDeployConfigCreateCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "client",
		CanonicalPath: []string{"client", "deploy-config", "create"},
		Use:           "create",
		Description:   "Create deploy host config",
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
			if _, err := deployConfigNameFromBody(rt.Body); err != nil {
				return err
			}
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			name, err := deployConfigNameFromBody(rt.Body)
			if err != nil {
				return err
			}

			precheckResp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodGet,
				Path:     "/deploy/v1/hostConfig/nameExist?name=" + url.QueryEscape(name),
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
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "deploy config name already exists")
			}

			osPath, _ := deployConfigOSPathSegment(rt.String("os"))
			body, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/deploy/v1/hostConfig/" + osPath,
				Body:     rt.Body,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, body)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			osPath, _ := deployConfigOSPathSegment(rt.String("os"))
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/deploy/v1/hostConfig/" + osPath,
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

type deployConfigCreateBody struct {
	Name string `json:"name"`
}

type deployConfigNameExistResponse struct {
	Status       string `json:"status"`
	ResponseData struct {
		Exist bool `json:"exist"`
	} `json:"responseData"`
}

func deployConfigOSPathSegment(os string) (string, bool) {
	switch os {
	case "linux":
		return "Linux", true
	case "windows":
		return "Windows", true
	case "unix":
		return "Unix", true
	default:
		return "", false
	}
}

func deployConfigNameFromBody(body []byte) (string, error) {
	var req deployConfigCreateBody
	if err := json.Unmarshal(body, &req); err != nil {
		return "", clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field name")
	}
	return name, nil
}

func parseDeployConfigNameExists(body []byte) (bool, error) {
	var resp deployConfigNameExistResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid deploy config nameExist response")
	}
	if strings.TrimSpace(resp.Status) == "" {
		return false, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid deploy config nameExist response")
	}
	return resp.ResponseData.Exist, nil
}
