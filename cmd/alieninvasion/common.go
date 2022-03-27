package alieninvasion

import (
	"fmt"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg"
	"github.com/itiky/alienInvasion/pkg/config"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagConfigPath      = "config"
	flagShortConfigPath = "c"

	flagMapPath      = "map"
	flagShortMapPath = "m"

	flagAliens      = "aliens"
	flagShortAliens = "a"
)

// loadConfig loads and validate a config file if file path is provided.
func loadConfig(cmd *cobra.Command) error {
	configPath, err := pkg.GetStringFlag(cmd, flagConfigPath, true)
	if err != nil {
		return err
	}

	if configPath != nil {
		if err := config.ReadConfigFile(*configPath); err != nil {
			return pkg.BuildParamErr(
				flagConfigPath, pkg.ParamTypeFlag,
				fmt.Errorf("reading config file: %w", err),
			)
		}
		if err := config.Validate(); err != nil {
			return pkg.BuildParamErr(
				flagConfigPath, pkg.ParamTypeFlag,
				fmt.Errorf("config invalid: %w", err),
			)
		}
	}

	return nil
}

// buildCityMap builds a CityMap by file path.
func buildCityMap(cmd *cobra.Command) (model.CityMap, error) {
	mapPath, err := pkg.GetStringFlag(cmd, flagMapPath, false)
	if err != nil {
		return nil, err
	}

	cityMap, err := model.NewCityMapFromFile(*mapPath)
	if err != nil {
		return nil, pkg.BuildParamErr(
			flagMapPath, pkg.ParamTypeFlag,
			fmt.Errorf("parsing map file: %w", err),
		)
	}

	return cityMap, nil
}

// buildAliens generates aliens slice by count provided.
func buildAliens(cmd *cobra.Command) ([]model.Alien, error) {
	aliensCount, err := pkg.GetUintFlag(cmd, flagAliens, false)
	if err != nil {
		return nil, err
	}

	return model.GenAliensFromConfig(*aliensCount), nil
}

// buildLogger builds a new logger using logLevel from config.
func buildLogger() (zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(viper.GetString(config.AppLogLevel))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parsing logLevel: %w", err)
	}

	logger := logging.NewLogger(
		logging.WithLogLevel(logLevel),
	)

	return logger, nil
}
