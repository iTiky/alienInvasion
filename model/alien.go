package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/itiky/alienInvasion/pkg/config"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Alien keeps alien params.
type Alien struct {
	// Unique Alien ID
	Name string

	// Alien power (higher the value, the longer City fight is gonna take)
	Power uint

	// Alien movement speed
	Speed time.Duration

	// Max number of movements
	MaxSteps uint
}

// GetLoggerContext enriches logger context with essential City fields.
func (a Alien) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.
		Str(logging.AlienNameKey, a.Name)
}

// GenAliensFromConfig generates Aliens with random stats according to config params.
// Contract: config is valid.
func GenAliensFromConfig(n uint) []Alien {
	stepMinDur, stepMaxDur := viper.GetDuration(config.AlienStepMinDur), viper.GetDuration(config.AlienStepMaxDur)
	pwrMin, pwrMax := viper.GetUint(config.AlienMinPower), viper.GetUint(config.AlienMaxPower)

	aliens := make([]Alien, 0, n)
	for id := uint(0); id < n; id++ {
		pwr := pwrMin
		if pwrMax != pwrMin {
			diff := int64(pwrMax - pwrMin)
			pwr += uint(rand.Int63n(diff)) //nolint:gosec
		}

		stepDur := stepMinDur
		if stepMaxDur != stepMinDur {
			diff := int64(stepMaxDur - stepMinDur)
			stepDur += time.Duration(rand.Int63n(diff)) //nolint:gosec
		}

		aliens = append(aliens, Alien{
			Name:     fmt.Sprintf("#%08d", id),
			Power:    pwr,
			Speed:    stepDur,
			MaxSteps: viper.GetUint(config.AlienMaxSteps),
		})
	}

	return aliens
}
