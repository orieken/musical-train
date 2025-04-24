package command

// RegisterCommands registers all commands with the registry
func RegisterCommands(registry *Registry) {
	// Register built-in commands
	RegisterListBinariesCommand(registry)
	
	// Register the subcommand executor
	RegisterSubcommandExecutor(registry)
}