package tgModel

type Commands map[string]Command

func NewCommands() Commands {
	return make(Commands)
}

func (cs Commands) Merge(list Commands) Commands {
	merged := make(Commands)
	for key, value := range cs {
		merged[key] = value
	}
	for key, value := range list {
		merged[key] = value
	}
	return merged
}

func (cs Commands) SetBotData(botName, serviceName string) Commands {
	for key, value := range cs {
		value.BotName = botName
		value.Service = serviceName
		cs[key] = value
	}
	return cs
}

func (cs Commands) Add(name string, item Command) Commands {
	if cs == nil {
		cs = make(Commands)
	}
	cs[name] = item
	return cs
}

func (cs Commands) AddWIth(
	name, description, commandType string,
	synonyms, triggers, templates []string,
	listExclude bool,
	permissions CommandPermissions,
	handler HandlerFunc,
) Commands {
	if cs == nil {
		cs = make(Commands)
	}
	cs[name] = Command{
		Command:     "/" + name,
		Synonyms:    synonyms,
		Triggers:    triggers,
		Templates:   templates,
		Description: description,
		CommandType: commandType,
		ListExclude: listExclude,
		Permissions: permissions,
		Handler:     handler,
	}
	return cs
}

func (cs Commands) AddSimple(
	name, description string,
	handler HandlerFunc,
	synonyms ...string,
) Commands {
	if cs == nil {
		cs = make(Commands)
	}
	item := FreeCommand().Simple(name, description, handler, synonyms...)
	cs[name] = *item
	return cs
}

func (cs Commands) Exclude() Commands {
	for index, item := range cs {
		item.ListExclude = true
		cs[index] = item
	}
	return cs
}
