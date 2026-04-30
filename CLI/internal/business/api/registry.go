package api

import "github.com/anybackup-ai/Anybackup/CLI/internal/business/meta"

type commandFactory func() meta.CommandMeta

var commandFactories = []commandFactory{
	newAPICommand,
}

func Commands() []meta.CommandMeta {
	defs := make([]meta.CommandMeta, 0, len(commandFactories))
	for _, factory := range commandFactories {
		defs = append(defs, factory())
	}
	return defs
}
