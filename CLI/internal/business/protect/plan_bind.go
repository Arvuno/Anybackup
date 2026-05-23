package protect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"

	"github.com/spf13/cobra"
)

func newPolicyBindCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "policy", "bind"},
		Use:           "bind",
		Description:   "Bind policy to a protect object",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("object-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --object-id")
			}
			if len(rt.Body) == 0 {
				return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
			}
			return nil
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			body, err := normalizePolicyBindBody(rt.String("object-id"), rt.Body)
			if err != nil {
				return console.RequestSpec{}, err
			}
			return console.RequestSpec{
				Method:   http.MethodPost,
				Path:     fmt.Sprintf("/backupmgm/v1/protect_object/%s/slas", rt.String("object-id")),
				Body:     body,
				ReadOnly: false,
			}, nil
		},
	}
}

func normalizePolicyBindBody(objectID string, body []byte) ([]byte, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	if raw, ok := req["objectId"]; ok {
		if bodyObjectID, ok := raw.(string); ok && strings.TrimSpace(bodyObjectID) != "" {
			if strings.TrimSpace(bodyObjectID) != strings.TrimSpace(objectID) {
				return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "body objectId does not match --object-id")
			}
		}
	}
	req["objectId"] = objectID

	normalized, err := json.Marshal(req)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}
	return normalized, nil
}
