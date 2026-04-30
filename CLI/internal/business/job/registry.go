package job

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newBackupDetailCommand,
		newChildListCommand,
		newSyncDetailCommand,
		newLogsCommand,
		newBusinessTypesCommand,
		newAppTypesCommand,
		newStopCommand,
		newDeleteCommand,
		newListCommand,
	)
}
