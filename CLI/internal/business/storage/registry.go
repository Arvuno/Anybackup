package storage

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newServiceListCommand,
		newPoolListCommand,
		newPoolCreateCommand,
		newPoolDeleteCommand,
		newPoolNodeListCommand,
		newPoolNodeDeviceListCommand,
	)
}
