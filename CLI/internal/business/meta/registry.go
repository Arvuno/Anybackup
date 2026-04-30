package meta

type CommandFactory func() CommandMeta

func CollectCommands(factories []CommandFactory) []CommandMeta {
	defs := make([]CommandMeta, 0, len(factories))
	for _, factory := range factories {
		defs = append(defs, factory())
	}
	return defs
}

func NewDomainCommands(factories ...CommandFactory) []CommandMeta {
	return CollectCommands(factories)
}
