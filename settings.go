package gommander

type AppSettings = map[Setting]bool

type Setting byte

const (
	ShowCommandAliases Setting = iota
	ShowHelpOnAllErrors
	IncludeHelpSubcommand
	OverrideAllDefaultListeners
	SortCommandsAlphabetically
)
