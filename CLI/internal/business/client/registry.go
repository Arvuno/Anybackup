package client

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

func Commands() []meta.CommandMeta {
	return meta.NewDomainCommands(
		newDeployCommand,
		newDeployHistoryCommand,
		newListCommand,
		newDetailCommand,
		newRunnerListCommand,
		newRunnerTypesCommand,
		newDatasourceListCommand,
		newDeployConfigListCommand,
		newDeployConfigCreateCommand,
		newDeployLogListCommand,
	)
}
