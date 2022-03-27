package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	appPrefix = "app."

	AppLogLevel     = appPrefix + "logLevel"     // Logging level [debug, info, warn, error, fatal]
	AppScreenWidth  = appPrefix + "screenWidth"  // Visualization screen width [int]
	AppScreenHeight = appPrefix + "screenHeight" // Visualization screen height [int]

	AppAliensDisembarkMinRate = appPrefix + "aliensDisembarkMinRate" // Minimum time offset to disembark an alien [duration]
	AppAliensDisembarkMaxRate = appPrefix + "aliensDisembarkMaxRate" // Maximum time offset to disembark an alien [duration]

	AppSimStopCheckRate = appPrefix + "simStopCheckRate" // Check if simulation should be stopped rate [duration]
)

const (
	cityPrefix = "city."

	CityFightDurK = cityPrefix + "fightDurationCoef" // Fight duration per Alien power (K * totalAliensPower = OverallFightDuration) [duration]
)

const (
	alienPrefix = "alien."

	AlienStepMinDur = alienPrefix + "stepMinDuration" // Minimum time offset to move from a City [duration]
	AlienStepMaxDur = alienPrefix + "stepMaxDuration" // Maximum time offset to move from a City [duration]

	AlienMaxSteps = alienPrefix + "maxSteps" // Maximum number of steps before Alien stops moving [uint]

	AlienMinPower = alienPrefix + "minPower" // Minimum fighting power [uint]
	AlienMaxPower = alienPrefix + "maxPower" // Maximum fighting power [uint]
)

func init() {
	// app. defaults
	viper.SetDefault(AppLogLevel, zerolog.LevelInfoValue)

	viper.SetDefault(AppAliensDisembarkMinRate, 250*time.Millisecond)
	viper.SetDefault(AppAliensDisembarkMaxRate, 500*time.Millisecond)

	viper.SetDefault(AppSimStopCheckRate, 1*time.Second)

	viper.SetDefault(AppScreenWidth, 1200)
	viper.SetDefault(AppScreenHeight, 1000)

	// city. defaults
	viper.SetDefault(CityFightDurK, 150*time.Millisecond)

	// alien. defaults
	viper.SetDefault(AlienStepMinDur, 500*time.Millisecond)
	viper.SetDefault(AlienStepMaxDur, 1*time.Second)

	viper.SetDefault(AlienMaxSteps, 25)

	viper.SetDefault(AlienMinPower, 0)
	viper.SetDefault(AlienMaxPower, 10)
}
