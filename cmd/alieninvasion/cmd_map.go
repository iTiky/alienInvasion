package alieninvasion

import (
	"context"
	"fmt"

	"github.com/itiky/alienInvasion/service/monitor/display"
	"github.com/spf13/cobra"
)

// NewMapCmd creates the /map command.
func NewMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map",
		Short: "Displays a map without simulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Inputs build
			if err := loadConfig(cmd); err != nil {
				return err
			}

			cityMap, err := buildCityMap(cmd)
			if err != nil {
				return err
			}

			monitor, err := display.New(cityMap, nil)
			if err != nil {
				return fmt.Errorf("building visualization service: %w", err)
			}

			// Run
			monitor.Run(context.Background())

			return nil
		},
	}

	cmd.Flags().StringP(flagConfigPath, flagShortConfigPath, "./config.toml", "Config file path (optional)")
	cmd.Flags().StringP(flagMapPath, flagShortMapPath, "./map.aimap", "Map file path")

	return cmd
}
