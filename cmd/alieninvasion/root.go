package alieninvasion

import "github.com/spf13/cobra"

// NewRootCmd creates the / command group.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "alieninvasion",
		Short:         "AlienInvasion emulates aliens invasion to Earth",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		NewStartCmd(),
		NewMapCmd(),
		NewVersionCmd(),
	)

	return cmd
}
