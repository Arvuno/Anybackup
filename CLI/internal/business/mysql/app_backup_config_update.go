package mysql

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"

	"github.com/spf13/cobra"
)

func newAppBackupConfigSetCommand() meta.CommandMeta {
	return newAppBackupConfigWriteCommand([]string{"mysql", "backup-config", "set"}, "set", "Set MySQL app backup config")
}

func newAppBackupConfigWriteCommand(path []string, use, description string) meta.CommandMeta {
	return meta.CommandMeta{
		Domain:        "mysql",
		CanonicalPath: path,
		Use:           use,
		Description:   description,
		Risk:          meta.RiskWrite,
		AuthType:      meta.AuthAKSK,
		ReadOnly:      false,
		BindFlags: func(cmd *cobra.Command, rt *meta.Runtime) {
			cmd.Flags().StringVarP(&rt.Inline, "data", "d", "", "inline json body")
			cmd.Flags().StringVar(&rt.BodyFile, "body-file", "", "json body file")
		},
		Validate: func(rt *meta.Runtime) error {
			body, err := normalizeAppBackupConfigUpdateBody(rt.Body)
			if err != nil {
				return err
			}
			rt.Body = body
			return nil
		},
		Execute: func(ctx context.Context, rt *meta.Runtime, deps meta.Dependencies) error {
			if deps.Console == nil {
				return clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "console executor is not configured")
			}
			body, err := fillStoragePoolFromAutoSelectIfMissing(ctx, rt, deps.Console)
			if err != nil {
				return err
			}
			rt.Body = body
			resp, err := deps.Console.Execute(ctx, rt, console.RequestSpec{
				Method:   http.MethodPut,
				Path:     "/backupmgm/v1/mysql/app_backup_config",
				Body:     rt.Body,
				ReadOnly: false,
			})
			if err != nil {
				return err
			}
			return output.WriteRaw(deps.Streams, resp)
		},
		BuildRequest: func(rt *meta.Runtime) (console.RequestSpec, error) {
			return console.RequestSpec{
				Method:   http.MethodPut,
				Path:     "/backupmgm/v1/mysql/app_backup_config",
				Body:     rt.Body,
				ReadOnly: false,
			}, nil
		},
	}
}

func normalizeAppBackupConfigUpdateBody(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return nil, clierrors.New(clierrors.CodeMissingRequestBody, clierrors.ExitInvalidArgs, "request body is required")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	if isBlankString(payload["objectId"]) {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field objectId")
	}
	if isBlankString(payload["systemId"]) {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field systemId")
	}
	if payload["backupType"] == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field backupType")
	}
	if payload["backupGran"] == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field backupGran")
	}

	commonConfig, ok := payload["commonConfigParams"].(map[string]interface{})
	if !ok || commonConfig == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field commonConfigParams")
	}
	backupDestination, ok := commonConfig["backupDestination"].(map[string]interface{})
	if !ok || backupDestination == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field commonConfigParams.backupDestination")
	}
	regionParams, ok := backupDestination["regionParams"].([]interface{})
	if !ok || len(regionParams) == 0 {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "commonConfigParams.backupDestination.regionParams must contain at least one item")
	}
	firstRegion, ok := regionParams[0].(map[string]interface{})
	if !ok || firstRegion == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "commonConfigParams.backupDestination.regionParams[0] must be an object")
	}

	// Alias compatibility: fill canonical fields only when canonical is missing.
	fillAliasIfMissing(payload, "isMergeBackup", "syntheticBackup")
	fillAliasIfMissing(payload, "isTableRestore", "isTableLevelRecovery")

	// Successful payloads treat these as separate control fields rather than aliases.
	ensureDerivedDefaults(payload)

	// Keep nested conventionalConfig complete, but do not overwrite explicit top-level fields:
	// successful backend payloads may intentionally carry different top-level and nested values.
	if conventionalConfig, ok := commonConfig["conventionalConfig"].(map[string]interface{}); ok && conventionalConfig != nil {
		ensureNestedNumberDefault(conventionalConfig, "failureRetry", 0)
		ensureNestedNumberDefault(conventionalConfig, "encryptionTrans", 0)
	}

	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, clierrors.Wrap(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "failed to encode request body", err)
	}
	return normalized, nil
}

func fillAliasIfMissing(payload map[string]interface{}, canonicalKey, aliasKey string) {
	if payload[canonicalKey] == nil && payload[aliasKey] != nil {
		payload[canonicalKey] = payload[aliasKey]
	}
}

func ensureNestedNumberDefault(nested map[string]interface{}, key string, defaultValue int) {
	if _, exists := nested[key]; !exists {
		nested[key] = defaultValue
	}
}

func ensureDerivedDefaults(payload map[string]interface{}) {
	if payload["isMultiChannel"] == nil {
		if n, ok := asInt(payload["dataChannelNum"]); ok && n > 0 {
			payload["isMultiChannel"] = true
		}
	}
}

func asInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func isBlankString(value interface{}) bool {
	raw, ok := value.(string)
	return !ok || strings.TrimSpace(raw) == ""
}

func fillStoragePoolFromAutoSelectIfMissing(ctx context.Context, rt *meta.Runtime, executor meta.Executor) ([]byte, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(rt.Body, &payload); err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "request body is not valid json")
	}

	pool, err := getOrCreatePrimaryStoragePool(payload)
	if err != nil {
		return nil, err
	}
	if !isBlankString(pool["storageServiceId"]) && !isBlankString(pool["storagePoolId"]) {
		return rt.Body, nil
	}

	resp, err := executor.Execute(ctx, rt, console.RequestSpec{
		Method:   http.MethodGet,
		Path:     "/backupmgm/v2/auto_select/storage_pool",
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}

	selected, err := parseAutoSelectedStoragePool(resp)
	if err != nil {
		return nil, err
	}

	if isBlankString(pool["storageServiceId"]) {
		pool["storageServiceId"] = selected.StorageServiceID
	}
	if isBlankString(pool["storageServiceName"]) && selected.StorageServiceName != "" {
		pool["storageServiceName"] = selected.StorageServiceName
	}
	if isBlankString(pool["storagePoolId"]) {
		pool["storagePoolId"] = selected.StoragePoolID
	}
	if isBlankString(pool["storagePoolName"]) && selected.StoragePoolName != "" {
		pool["storagePoolName"] = selected.StoragePoolName
	}
	if pool["storagePoolType"] == nil && selected.StoragePoolType != nil {
		pool["storagePoolType"] = *selected.StoragePoolType
	}
	if pool["encryptionStorage"] == nil {
		pool["encryptionStorage"] = 0
	}
	if pool["compress"] == nil {
		pool["compress"] = 0
	}
	if pool["deduplication"] == nil {
		pool["deduplication"] = 0
	}
	if pool["dataConsistencyLogic"] == nil {
		pool["dataConsistencyLogic"] = 1
	}
	if pool["isMasterStorage"] == nil {
		pool["isMasterStorage"] = true
	}

	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, clierrors.Wrap(clierrors.CodeInvalidJSONBody, clierrors.ExitInvalidArgs, "failed to encode request body", err)
	}
	return normalized, nil
}

func getOrCreatePrimaryStoragePool(payload map[string]interface{}) (map[string]interface{}, error) {
	commonConfig, ok := payload["commonConfigParams"].(map[string]interface{})
	if !ok || commonConfig == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field commonConfigParams")
	}
	backupDestination, ok := commonConfig["backupDestination"].(map[string]interface{})
	if !ok || backupDestination == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "missing body field commonConfigParams.backupDestination")
	}

	var firstRegion map[string]interface{}
	regionParams, ok := backupDestination["regionParams"].([]interface{})
	if !ok || len(regionParams) == 0 {
		firstRegion = map[string]interface{}{}
		backupDestination["regionParams"] = []interface{}{firstRegion}
	} else {
		var regionOK bool
		firstRegion, regionOK = regionParams[0].(map[string]interface{})
		if !regionOK || firstRegion == nil {
			return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "commonConfigParams.backupDestination.regionParams[0] must be an object")
		}
	}

	storagePoolParams, ok := firstRegion["storagePoolParams"].([]interface{})
	if !ok || len(storagePoolParams) == 0 {
		pool := map[string]interface{}{}
		firstRegion["storagePoolParams"] = []interface{}{pool}
		return pool, nil
	}

	pool, ok := storagePoolParams[0].(map[string]interface{})
	if !ok || pool == nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "commonConfigParams.backupDestination.regionParams[0].storagePoolParams[0] must be an object")
	}
	return pool, nil
}

type autoSelectedStoragePool struct {
	StorageServiceID   string
	StorageServiceName string
	StoragePoolID      string
	StoragePoolName    string
	StoragePoolType    *int
}

func parseAutoSelectedStoragePool(body []byte) (autoSelectedStoragePool, error) {
	var envelope struct {
		ResponseData interface{} `json:"responseData"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return autoSelectedStoragePool{}, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "invalid auto-select response")
	}

	var data map[string]interface{}
	switch raw := envelope.ResponseData.(type) {
	case map[string]interface{}:
		data = raw
	case []interface{}:
		if len(raw) > 0 {
			if first, ok := raw[0].(map[string]interface{}); ok {
				data = first
			}
		}
	}
	if data == nil {
		return autoSelectedStoragePool{}, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "auto-select response missing responseData")
	}

	result := autoSelectedStoragePool{
		StorageServiceID:   strings.TrimSpace(stringFromAny(data["storageSvcId"])),
		StorageServiceName: strings.TrimSpace(stringFromAny(data["storageSvcName"])),
		StoragePoolID:      strings.TrimSpace(stringFromAny(data["storagePoolId"])),
		StoragePoolName:    strings.TrimSpace(stringFromAny(data["storagePoolName"])),
	}
	if v, ok := numberToInt(data["storagePoolType"]); ok {
		result.StoragePoolType = &v
	}

	if result.StorageServiceID == "" || result.StoragePoolID == "" {
		return autoSelectedStoragePool{}, clierrors.New(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "auto-select response missing storageSvcId or storagePoolId")
	}
	return result, nil
}

func stringFromAny(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func numberToInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}
