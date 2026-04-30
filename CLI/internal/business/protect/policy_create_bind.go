package protect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newPolicyCreateBindCommand() meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "protect",
		CanonicalPath: []string{"protect", "policy", "create-bind"},
		Use:           "create-bind",
		Description:   "Create backup policy then bind to a protect object",
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			rt.BindStringFlag(cmd.Flags(), "object-id", "", "Object ID")
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body for policy create")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file for policy create")
		},
		Validate: func(rt *meta.Runtime) error {
			if rt.String("object-id") == "" {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing --object-id")
			}
			if err := ensureNormalizedPolicyCreateBody(rt); err != nil {
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
			if deps.Console == nil {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "console executor is not configured")
			}
			if err := ensureNormalizedPolicyCreateBody(rt); err != nil {
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

			createBody, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     "/api/sla/v1/group/backup_info",
				Body:     rt.Body,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			policyID, err := extractCreatedPolicyID(createBody)
			if err != nil {
				return err
			}

			bindReqBody, err := json.Marshal(map[string]any{
				"objectId": rt.String("object-id"),
				"slaIds": []string{policyID},
			})
			if err != nil {
				return clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "failed to build bind request body")
			}
			bindBody, bindErr := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPost,
				Path:     fmt.Sprintf("/backupmgm/v1/protect_object/%s/slas", url.PathEscape(rt.String("object-id"))),
				Body:     bindReqBody,
				ReadOnly: false,
			})
			if bindErr != nil {
				rollbackBody, rbErr := deps.Console.Execute(ctx, rt, console.RequestSpec{
					Method:   http.MethodDelete,
					Path:     "/api/sla/v1/groups",
					Body:     []byte(fmt.Sprintf(`{"groupIds":["%s"]}`, policyID)),
					ReadOnly: false,
				})
				if rbErr != nil {
					return clierrors.New(
						clierrors.CodeBackendInvalidResp,
						clierrors.ExitTransport,
						fmt.Sprintf("policy bind failed and rollback delete failed (policyId=%s): bindErr=%v; rollbackErr=%v", policyID, bindErr, rbErr),
					)
				}
				rollbackData := parseJSONAny(rollbackBody)
				return clierrors.New(
					clierrors.CodeBackendInvalidResp,
					clierrors.ExitTransport,
					fmt.Sprintf("policy bind failed, rollback delete succeeded (policyId=%s, rollback=%s): bindErr=%v", policyID, oneLineJSON(rollbackData), bindErr),
				)
			}

			out := map[string]any{
				"status": "success",
				"error":  nil,
				"responseData": map[string]any{
					"policyId": policyID,
					"create":   parseJSONAny(createBody),
					"bind":     parseJSONAny(bindBody),
					"rollback": nil,
				},
			}
			return output.WriteJSON(deps.Streams, out)
		},
	}
}

func ensureNormalizedPolicyCreateBody(rt *meta.Runtime) error {
	if len(rt.Body) == 0 {
		return nil
	}
	body, err := normalizePolicyCreateBody(rt.Body)
	if err != nil {
		return err
	}
	rt.Body = body
	return nil
}

func normalizePolicyCreateBody(body []byte) ([]byte, error) {
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

func extractCreatedPolicyID(body []byte) (string, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid create policy response")
	}

	candidates := []string{
		stringFromPath(payload, "responseData"),
		stringFromPath(payload, "responseData", "id"),
		stringFromPath(payload, "responseData", "groupID"),
		stringFromPath(payload, "responseData", "groupId"),
		stringFromPath(payload, "responseData", "slaId"),
		stringFromPath(payload, "responseData", "data", "id"),
		stringFromPath(payload, "responseData", "data", "groupId"),
		stringFromPath(payload, "responseData", "result", "id"),
		stringFromPath(payload, "responseData", "result", "groupId"),
		stringFromFirstArrayItem(payload, "responseData", "data", "id"),
		stringFromFirstArrayItem(payload, "responseData", "data", "groupId"),
		stringFromFirstArrayItem(payload, "responseData", "groupIds"),
		stringFromFirstArrayItem(payload, "responseData", "slaIds"),
		stringFromPath(payload, "id"),
	}
	for _, id := range candidates {
		if strings.TrimSpace(id) != "" {
			return strings.TrimSpace(id), nil
		}
	}

	return "", clierrors.New(
		clierrors.CodeBackendInvalidResp,
		clierrors.ExitTransport,
		fmt.Sprintf("create policy response missing id; response=%s", truncateForErr(body, 512)),
	)
}

func stringFromPath(v any, path ...string) string {
	current := v
	for _, key := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return ""
		}
		next, ok := m[key]
		if !ok {
			return ""
		}
		current = next
	}
	return anyToString(current)
}

func stringFromFirstArrayItem(v any, first, second string, rest ...string) string {
	m, ok := v.(map[string]any)
	if !ok {
		return ""
	}
	level1, ok := m[first].(map[string]any)
	if !ok {
		return ""
	}
	arr, ok := level1[second].([]any)
	if !ok || len(arr) == 0 {
		return ""
	}
	if len(rest) == 0 {
		return anyToString(arr[0])
	}
	current := any(arr[0])
	for _, key := range rest {
		item, ok := current.(map[string]any)
		if !ok {
			return ""
		}
		next, ok := item[key]
		if !ok {
			return ""
		}
		current = next
	}
	return anyToString(current)
}

func anyToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	case float64:
		// JSON numbers are float64 when unmarshaled into interface{}.
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	default:
		return ""
	}
}

func truncateForErr(b []byte, n int) string {
	s := strings.TrimSpace(string(b))
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
}

func parseJSONAny(body []byte) any {
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return string(body)
	}
	return v
}

func oneLineJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}
