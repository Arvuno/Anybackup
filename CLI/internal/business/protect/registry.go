package protect

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newPolicyCreateBindCommand,
		newPolicyBindCommand,
		newPolicyBindBatchCommand,
		newPolicyBindListCommand,
		newPolicyBatchUnbindCommand,
		newBackupStartCommand,
		newBackupBatchStartCommand,
		newConfigPolicyGetCommand,
		newConfigPolicyGetByAppTypeCommand,
		newStoragePoolAutoSelectCommand,
		newBackupConfigGetCommand,
		newClientListCommand,
	)
}
