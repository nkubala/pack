package stack

import (
	"github.com/spf13/cobra"

	"github.com/buildpacks/pack/internal/commands"

	"github.com/buildpacks/pack/logging"
)

// Deprecated: Use `suggest` instead
func SuggestStacks(logger logging.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "suggest-stacks",
		Args:    cobra.NoArgs,
		Short:   "Display list of recommended stacks",
		Example: "pack suggest-stacks",
		Run: func(*cobra.Command, []string) {
			commands.DeprecationWarning(logger, "suggest-stacks", "stack suggest")
			Suggest(logger)
		},
		Hidden: true,
	}

	commands.AddHelpFlag(cmd, "suggest-stacks")
	return cmd
}
