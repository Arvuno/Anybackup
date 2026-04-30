package storage

import (
	"encoding/json"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

func validateStorageServiceID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("storage-svc-id", rt.String("storage-svc-id"))
}

func validatePoolID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("pool-id", rt.String("pool-id"))
}

func validateNodeID(rt *meta.Runtime) error {
	return inputs.RequireNonEmpty("node-id", rt.String("node-id"))
}

type storagePoolCreatePayload struct {
	Name             string                     `json:"name"`
	Type             json.RawMessage            `json:"type"`
	RedundancyPolicy map[string]json.RawMessage `json:"redundancyPolicy"`
	Resource         []json.RawMessage          `json:"resource"`
	WarnThreshold    *int                       `json:"warnThreshold"`
	SafeThreshold    *int                       `json:"safeThreshold"`
}

func validatePoolCreateBody(body []byte) error {
	if len(body) == 0 {
		return clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	var req storagePoolCreatePayload
	if err := json.Unmarshal(body, &req); err != nil {
		return clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	if strings.TrimSpace(req.Name) == "" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field name")
	}
	if field, ok := raw["type"]; !ok || len(field) == 0 || string(field) == "null" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field type")
	}
	if _, ok := raw["redundancyPolicy"]; !ok || req.RedundancyPolicy == nil {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field redundancyPolicy")
	}
	if field, ok := req.RedundancyPolicy["policy"]; !ok || len(field) == 0 || string(field) == "null" {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field redundancyPolicy.policy")
	}
	if field, ok := raw["resource"]; !ok || string(field) == "null" || req.Resource == nil {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field resource")
	}
	if req.WarnThreshold != nil && (*req.WarnThreshold < 1 || *req.WarnThreshold > 98) {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid body field warnThreshold")
	}
	if req.SafeThreshold != nil && (*req.SafeThreshold < 2 || *req.SafeThreshold > 99) {
		return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid body field safeThreshold")
	}

	return nil
}
