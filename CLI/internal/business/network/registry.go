package network

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newSubnetListCommand,
		newSubnetNodeListCommand,
		newSubnetNodeIPListCommand,
		newSubnetNodeIPSetCommand,
		newSubnetNodeIPRemoveCommand,
	)
}
