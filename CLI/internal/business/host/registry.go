package host

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newListCommand,
		newObjectCreateCommand,
		newObjectListCommand,
		newObjectDetailCommand,
		newBackupConfigCommand,
		newBackupConfigDetailCommand,
		newRestoreStartCommand,
	)
}
