package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// ReadConfigFile reads config file and sets Viper values.
func ReadConfigFile(filePath string) error {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// Validate performs config values validation.
func Validate() error {
	// app
	{
		logLvl := viper.GetString(AppLogLevel)
		if _, err := zerolog.ParseLevel(logLvl); err != nil {
			return fmt.Errorf("%s: invalid", AppLogLevel)
		}

		minDrop, maxDrop := viper.GetDuration(AppAliensDisembarkMinRate), viper.GetDuration(AppAliensDisembarkMaxRate)
		if minDrop < 0 {
			return fmt.Errorf("%s key: must be GTE 0", AppAliensDisembarkMinRate)
		}
		if maxDrop < 0 {
			return fmt.Errorf("%s key: must be GTE 0", AppAliensDisembarkMaxRate)
		}
		if maxDrop < minDrop {
			return fmt.Errorf("%s key: must be GTE %s", AppAliensDisembarkMaxRate, AppAliensDisembarkMinRate)
		}

		stopCheckRate := viper.GetDuration(AppSimStopCheckRate)
		if stopCheckRate <= 0 {
			return fmt.Errorf("%s key: must be GT 0", AppSimStopCheckRate)
		}

		sWidth, sHeight := viper.GetInt(AppScreenWidth), viper.GetInt(AppScreenHeight)
		if sWidth <= 0 {
			return fmt.Errorf("%s key: must be GT 0", AppScreenWidth)
		}
		if sHeight <= 0 {
			return fmt.Errorf("%s key: must be GT 0", AppScreenHeight)
		}
	}

	// city
	{
		fightK := viper.GetDuration(CityFightDurK)
		if fightK < 0 {
			return fmt.Errorf("%s key: must be GTE 0", CityFightDurK)
		}
	}

	// alien
	{
		minStep, maxStep := viper.GetDuration(AlienStepMinDur), viper.GetDuration(AlienStepMaxDur)
		if minStep <= 0 {
			return fmt.Errorf("%s key: must be GT 0", AlienStepMinDur)
		}
		if maxStep <= 0 {
			return fmt.Errorf("%s key: must be GT 0", AlienStepMaxDur)
		}
		if maxStep < minStep {
			return fmt.Errorf("%s key: must be GTE %s", AlienStepMaxDur, AlienStepMinDur)
		}

		minPwr, maxPwr := viper.GetUint(AlienMinPower), viper.GetUint(AlienMaxPower)
		if maxPwr < minPwr {
			return fmt.Errorf("%s key: must be GTE %s", AlienMaxPower, AlienMinPower)
		}

		steps := viper.GetUint(AlienMaxSteps)
		if steps == 0 {
			return fmt.Errorf("%s key: must be GT 0", AlienMaxSteps)
		}
	}

	return nil
}

func init() {
	// Viper setup
	viper.SetEnvPrefix("AI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
