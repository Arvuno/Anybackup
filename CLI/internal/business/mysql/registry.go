package mysql

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newObjectListCommand,
		newTargetInstanceListCommand,
		newObjectGetCommand,
		newBackupDetailCommand,
		newAppBackupConfigDetailCommand,
		newAppBackupConfigSetCommand,
		newDatasourceBackupCommand,
		newRecoveryRangeCommand,
		newRestoreConfigCreateCommand,
		newRecoveryConfigDetailCommand,
		newDatasourceRecoveryCommand,
		newTimepointListCommand,
		newRecoveryDetailCommand,
		newAuthorizeCommand,
	)
}
