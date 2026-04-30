package mysql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/output"
)

func TestNormalizeAppBackupConfigUpdateBody_RequiredFields(t *testing.T) {
	_, err := normalizeAppBackupConfigUpdateBody([]byte(`{"objectId":"x"}`))
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestNormalizeAppBackupConfigUpdateBody_AliasFillAndDerivedDefaults(t *testing.T) {
	input := []byte(`{
		"objectId":"428e4c05bf87deef20a86e670f3eb9a9",
		"systemId":"a63ba04bfb530fefa3b09d1f89d47432",
		"commonConfigParams":{
			"backupDestination":{
				"regionParams":[{"storagePoolParams":[{"storagePoolId":"0c6022822f1d11f19b0600163e33df73"}]}]
			},
			"conventionalConfig":{
				"failureRetry":0,
				"failureRetryCount":3,
				"failureRetryInterval":5,
				"forcedRetentionSwitch":true
			}
		},
		"backupType":2,
		"backupGran":1,
		"syntheticBackup":true,
		"isTableLevelRecovery":true,
		"dataChannelNum":2
	}`)

	normalized, err := normalizeAppBackupConfigUpdateBody(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(normalized, &out); err != nil {
		t.Fatalf("failed to parse normalized body: %v", err)
	}

	if v, ok := out["isMergeBackup"].(bool); !ok || !v {
		t.Fatalf("expected isMergeBackup=true from syntheticBackup alias, got %#v", out["isMergeBackup"])
	}
	if v, ok := out["isTableRestore"].(bool); !ok || !v {
		t.Fatalf("expected isTableRestore=true from isTableLevelRecovery alias, got %#v", out["isTableRestore"])
	}
	if v, ok := out["isMultiChannel"].(bool); !ok || !v {
		t.Fatalf("expected isMultiChannel=true derived from dataChannelNum>0, got %#v", out["isMultiChannel"])
	}
}

func TestNormalizeAppBackupConfigUpdateBody_DefaultConventionalConfigSwitchesToZero(t *testing.T) {
	input := []byte(`{
		"objectId":"428e4c05bf87deef20a86e670f3eb9a9",
		"systemId":"a63ba04bfb530fefa3b09d1f89d47432",
		"commonConfigParams":{
			"backupDestination":{
				"regionParams":[{"storagePoolParams":[{"storagePoolId":"0c6022822f1d11f19b0600163e33df73"}]}]
			},
			"conventionalConfig":{
				"failureRetryCount":3,
				"failureRetryInterval":5,
				"forcedRetentionSwitch":true
			}
		},
		"backupType":2,
		"backupGran":1
	}`)

	normalized, err := normalizeAppBackupConfigUpdateBody(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(normalized, &out); err != nil {
		t.Fatalf("failed to parse normalized body: %v", err)
	}

	commonConfig, ok := out["commonConfigParams"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected commonConfigParams object, got %#v", out["commonConfigParams"])
	}
	conventionalConfig, ok := commonConfig["conventionalConfig"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected conventionalConfig object, got %#v", commonConfig["conventionalConfig"])
	}

	if v, ok := conventionalConfig["failureRetry"].(float64); !ok || v != 0 {
		t.Fatalf("expected conventionalConfig.failureRetry=0 by default, got %#v", conventionalConfig["failureRetry"])
	}
	if v, ok := conventionalConfig["encryptionTrans"].(float64); !ok || v != 0 {
		t.Fatalf("expected conventionalConfig.encryptionTrans=0 by default, got %#v", conventionalConfig["encryptionTrans"])
	}
	if _, exists := out["failureRetry"]; exists {
		t.Fatalf("expected top-level failureRetry to remain unset when only conventionalConfig is defaulted, got %#v", out["failureRetry"])
	}
}

func TestNormalizeAppBackupConfigUpdateBody_DoesNotTreatControlFlagsAsAliases(t *testing.T) {
	input := []byte(`{
		"objectId":"428e4c05bf87deef20a86e670f3eb9a9",
		"systemId":"a63ba04bfb530fefa3b09d1f89d47432",
		"commonConfigParams":{
			"backupDestination":{
				"regionParams":[{"storagePoolParams":[{"storagePoolId":"0c6022822f1d11f19b0600163e33df73"}]}]
			},
			"conventionalConfig":{
				"encryptionTrans":1,
				"failureRetry":1,
				"failureRetryCount":2,
				"failureRetryInterval":5,
				"forcedRetentionSwitch":true,
				"forcedRetentionCycle":1
			}
		},
		"backupType":1,
		"backupGran":1,
		"backupMode":3,
		"dataChannelNum":2,
		"isDataChannelNum":false,
		"isMultiChannel":true,
		"nonAutomaticTransferComplete":true,
		"isNotAutoToFull":false,
		"failureRetry":false,
		"failureRetryCount":1,
		"failureRetryInterval":1,
		"forcedRetentionSwitch":false
	}`)

	normalized, err := normalizeAppBackupConfigUpdateBody(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(normalized, &out); err != nil {
		t.Fatalf("failed to parse normalized body: %v", err)
	}

	if v, ok := out["isMultiChannel"].(bool); !ok || !v {
		t.Fatalf("expected isMultiChannel to preserve explicit true, got %#v", out["isMultiChannel"])
	}
	if v, ok := out["isDataChannelNum"].(bool); !ok || v {
		t.Fatalf("expected isDataChannelNum to preserve explicit false, got %#v", out["isDataChannelNum"])
	}
	if v, ok := out["nonAutomaticTransferComplete"].(bool); !ok || !v {
		t.Fatalf("expected nonAutomaticTransferComplete to preserve explicit true, got %#v", out["nonAutomaticTransferComplete"])
	}
	if v, ok := out["isNotAutoToFull"].(bool); !ok || v {
		t.Fatalf("expected isNotAutoToFull to preserve explicit false, got %#v", out["isNotAutoToFull"])
	}
	if v, ok := out["failureRetry"].(bool); !ok || v {
		t.Fatalf("expected top-level failureRetry to preserve explicit false, got %#v", out["failureRetry"])
	}
	if v, ok := out["failureRetryCount"].(float64); !ok || v != 1 {
		t.Fatalf("expected top-level failureRetryCount to preserve explicit 1, got %#v", out["failureRetryCount"])
	}
	if v, ok := out["failureRetryInterval"].(float64); !ok || v != 1 {
		t.Fatalf("expected top-level failureRetryInterval to preserve explicit 1, got %#v", out["failureRetryInterval"])
	}
	if v, ok := out["forcedRetentionSwitch"].(bool); !ok || v {
		t.Fatalf("expected top-level forcedRetentionSwitch to preserve explicit false, got %#v", out["forcedRetentionSwitch"])
	}
}

type fakeExecutor struct {
	exec func(req console.RequestSpec) ([]byte, error)
}

func (f fakeExecutor) Execute(_ context.Context, _ *meta.Runtime, req console.RequestSpec) ([]byte, error) {
	return f.exec(req)
}

func TestAppBackupConfigSet_Execute_AutoFillStoragePool(t *testing.T) {
	def := newAppBackupConfigSetCommand()
	rt := meta.NewRuntime()
	rt.Body = []byte(`{
		"objectId":"obj-1",
		"systemId":"sys-1",
		"commonConfigParams":{
			"backupDestination":{
				"regionParams":[{"storagePoolParams":[{}]}]
			}
		},
		"backupType":1,
		"backupGran":1
	}`)
	if err := def.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	var stdout bytes.Buffer
	autoSelectCalled := false
	putCalled := false
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &stdout},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				switch {
				case req.Method == "GET" && req.Path == "/backupmgm/v2/auto_select/storage_pool":
					autoSelectCalled = true
					return []byte(`{"status":"success","responseData":{"storageSvcId":"svc-1","storageSvcName":"svc-name","storagePoolId":"pool-1","storagePoolName":"pool-name","storagePoolType":2}}`), nil
				case req.Method == "PUT" && req.Path == "/backupmgm/v1/mysql/app_backup_config":
					putCalled = true
					var payload map[string]interface{}
					if err := json.Unmarshal(req.Body, &payload); err != nil {
						t.Fatalf("PUT body json parse error: %v", err)
					}
					pool := payload["commonConfigParams"].(map[string]interface{})["backupDestination"].(map[string]interface{})["regionParams"].([]interface{})[0].(map[string]interface{})["storagePoolParams"].([]interface{})[0].(map[string]interface{})
					if pool["storageServiceId"] != "svc-1" || pool["storagePoolId"] != "pool-1" {
						t.Fatalf("storage pool auto-fill failed, pool=%#v", pool)
					}
					if pool["encryptionStorage"] != float64(0) || pool["compress"] != float64(0) || pool["deduplication"] != float64(0) {
						t.Fatalf("expected storage defaults encryption/compress/deduplication=0, pool=%#v", pool)
					}
					if pool["dataConsistencyLogic"] != float64(1) || pool["isMasterStorage"] != true {
						t.Fatalf("expected storage defaults dataConsistencyLogic=1 and isMasterStorage=true, pool=%#v", pool)
					}
					return []byte(`{"status":"success","error":null,"responseData":null}`), nil
				default:
					return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
				}
			},
		},
	}

	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !autoSelectCalled {
		t.Fatal("auto-select should be called when storage fields are missing")
	}
	if !putCalled {
		t.Fatal("PUT should be called")
	}
	if !strings.Contains(stdout.String(), `"status":"success"`) {
		t.Fatalf("stdout = %q, want success payload", stdout.String())
	}
}

func TestAppBackupConfigSet_Execute_KeepProvidedStoragePool(t *testing.T) {
	def := newAppBackupConfigSetCommand()
	rt := meta.NewRuntime()
	rt.Body = []byte(`{
		"objectId":"obj-1",
		"systemId":"sys-1",
		"commonConfigParams":{
			"backupDestination":{
				"regionParams":[{"storagePoolParams":[{"storageServiceId":"svc-fixed","storagePoolId":"pool-fixed"}]}]
			}
		},
		"backupType":1,
		"backupGran":1
	}`)
	if err := def.Validate(rt); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	autoSelectCalled := false
	deps := meta.Dependencies{
		Streams: output.Streams{Stdout: &bytes.Buffer{}},
		Console: fakeExecutor{
			exec: func(req console.RequestSpec) ([]byte, error) {
				if req.Method == "GET" && req.Path == "/backupmgm/v2/auto_select/storage_pool" {
					autoSelectCalled = true
					return nil, errors.New("auto-select should not be called")
				}
				if req.Method == "PUT" && req.Path == "/backupmgm/v1/mysql/app_backup_config" {
					return []byte(`{"status":"success","error":null,"responseData":null}`), nil
				}
				return nil, errors.New("unexpected request: " + req.Method + " " + req.Path)
			},
		},
	}

	if err := def.Execute(context.Background(), rt, deps); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if autoSelectCalled {
		t.Fatal("auto-select should not be called when storage fields already exist")
	}
}
