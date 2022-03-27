package alieninvasion

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/itiky/alienInvasion/pkg"
	"github.com/itiky/alienInvasion/pkg/config"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/monitor"
	"github.com/itiky/alienInvasion/service/monitor/display"
	"github.com/itiky/alienInvasion/service/monitor/noop"
	"github.com/itiky/alienInvasion/service/sim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDisplay      = "display"
	flagShortDisplay = "d"
)

// NewStartCmd creates the /start command.
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the simulation engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Inputs build
			if err := loadConfig(cmd); err != nil {
				return err
			}

			cityMap, err := buildCityMap(cmd)
			if err != nil {
				return err
			}

			aliens, err := buildAliens(cmd)
			if err != nil {
				return err
			}

			logger, err := buildLogger()
			if err != nil {
				return err
			}
			ctx := logging.SetCtxLogger(context.Background(), logger)

			// Monitor
			visualizationEnabled, err := pkg.GetBoolFlag(cmd, flagDisplay, false)
			if err != nil {
				return err
			}

			var monitorSvc monitor.WorldEventsListener
			monitorStopCh := make(chan struct{})
			if *visualizationEnabled {
				m, err := display.New(
					cityMap, aliens,
					display.WithScreenSize(
						viper.GetInt(config.AppScreenWidth), viper.GetInt(config.AppScreenHeight),
					),
				)
				if err != nil {
					return fmt.Errorf("building visualization service: %w", err)
				}
				monitorSvc = m
			} else {
				monitorSvc = noop.New(
					noop.WithLogs(),
				)
			}

			// Simulation engine
			simSvc, err := sim.New(
				sim.WithCityMap(cityMap),
				sim.WithAliens(aliens),
				sim.WithMonitor(monitorSvc),
			)
			if err != nil {
				return fmt.Errorf("building simulation service: %w", err)
			}

			// Run
			ctx, ctxCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
			defer ctxCancel()

			simStopCh := simSvc.Start(ctx)
			if svc, ok := monitorSvc.(*display.Monitor); ok {
				svc.Run(ctx)
				close(monitorStopCh)
			}

			select {
			case <-ctx.Done():
				logger.Info().Msg("Closing app: signal received")
			case <-simStopCh:
				logger.Info().Msg("Closing app: simulation stopped")
			case <-monitorStopCh:
				logger.Info().Msg("Closing app: monitor stopped")
			}

			return nil
		},
	}

	cmd.Flags().StringP(flagConfigPath, flagShortConfigPath, "./config.toml", "Config file path (optional)")
	cmd.Flags().StringP(flagMapPath, flagShortMapPath, "./map.aimap", "Map file path")
	cmd.Flags().UintP(flagAliens, flagShortAliens, 25, "Number of Aliens to disembark")
	cmd.Flags().BoolP(flagDisplay, flagShortDisplay, false, "Enable visualization")

	return cmd
}
