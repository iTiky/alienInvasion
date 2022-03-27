package alieninvasion

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitTag    string
	GitCommit string
	BuildDate string
)

const (
	flagFullVersion      = "full"
	flagFullVersionShort = "f"
)

// NewVersionCmd returns the /version cmd.
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints build version (use -f for more info)",
		Run: func(cmd *cobra.Command, args []string) {
			fullVersion, _ := cmd.Flags().GetBool(flagFullVersion)

			if fullVersion {
				fmt.Println("Tag:\t", GitTag)
				fmt.Println("Commit:\t", GitCommit)
				fmt.Println("Build:\t", BuildDate)
				return
			}

			fmt.Println(GitTag)
		},
	}

	cmd.Flags().BoolP(flagFullVersion, flagFullVersionShort, false, "Print full version info")

	return cmd
}
