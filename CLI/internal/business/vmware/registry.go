package vmware

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newObjectListCommand,
		newObjectInfoCommand,
		newDatasourceGetCommand,
		newBackupDetailCommand,
		newBackupConfigCreateCommand,
		newBackupConfigDetailCommand,
		newRecoveryDetailCommand,
		newRestoreConfigCreateCommand,
		newTimepointMetadataCommand,
	)
}
