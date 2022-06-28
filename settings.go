package gommander

type AppSettings = map[Setting]bool

type Setting byte

const (
	// Set whether or not to show the aliases of a subcommand, false by default
	ShowCommandAliases Setting = iota
	// Configure whether to print command help when an error occurs, false by default
	ShowHelpOnAllErrors
	// Configures whether to include the help subcommand, false by default
	IncludeHelpSubcommand
	// Removed all default event-listeners except the help event listener, cannot be overridden
	OverrideAllDefaultListeners
	// Removes the version flag when set to true
	DisableVersionFlag
	// Configure errors to be ignored when encountered
	IgnoreAllErrors
)
