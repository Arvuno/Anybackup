package policy

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newListCommand,
		newBackupCreateCommand,
		newBackupDetailCommand,
		newCopyCreateCommand,
		newCopyDetailCommand,
		newBackupDeleteCommand,
	)
}
