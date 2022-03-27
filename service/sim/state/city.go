package state

import (
	"context"
	"time"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/config"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/sim/types"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// alienWorldNotifierExpected notifies the World simulation engine about Alien's intentions.
// Since Alien has no idea about what is happening to the World, it asks a sim engine to do stuff.
type cityWorldNotifierExpected interface {
	// CityDestroyed sends City destroy request when the fight is over.
	CityDestroyed(r types.CityDestroyRequest)
}

// City keeps City's state with all Aliens on tile.
type City struct {
	model.City

	// State
	fightTimer *time.Timer
	aliens     map[string]*Alien // key: AlienID

	// Params
	worldNotifier cityWorldNotifierExpected
}

// NewCity creates a new City state.
// Contract: inputs are valid.
func NewCity(location model.City, worldNotifier cityWorldNotifierExpected) *City {
	return &City{
		City:          location,
		aliens:        make(map[string]*Alien),
		worldNotifier: worldNotifier,
	}
}

// AtFight checks that Aliens have a fight on that City tile.
func (c *City) AtFight() bool {
	return len(c.aliens) > 1
}

// AlienIDs returns all Alien IDs on that City tile.
func (c *City) AlienIDs() []string {
	ids := make([]string, 0, len(c.aliens))
	for id := range c.aliens {
		ids = append(ids, id)
	}

	return ids
}

// RemoveRoadsTo removes connection roads leading to a target cityID.
func (c *City) RemoveRoadsTo(cityID string) {
	if c.NorthRoad == cityID {
		c.NorthRoad = ""
	}
	if c.EastRoad == cityID {
		c.EastRoad = ""
	}
	if (c.SouthRoad) == cityID {
		c.SouthRoad = ""
	}
	if (c.WestRoad) == cityID {
		c.WestRoad = ""
	}
}

// AddAlien adds Alien on that City tile and informs that a new fight has started.
func (c *City) AddAlien(alien *Alien) bool {
	if alien == nil {
		return false
	}

	c.aliens[alien.Name] = alien

	// Check if a fight has started
	if len(c.aliens) == 1 {
		return false
	}

	// Estimated fight duration
	totalAlienPower := uint(0)
	for _, alien := range c.aliens {
		totalAlienPower += alien.Power
	}
	fightDuration := viper.GetDuration(config.CityFightDurK) * time.Duration(totalAlienPower)

	// Reset fight timer (prolong the fight)
	if c.fightTimer != nil {
		c.fightTimer.Reset(fightDuration)
		return true
	}

	// Start the notification routine
	c.fightTimer = time.NewTimer(fightDuration)
	go func() {
		defer func() {
			c.fightTimer = nil
		}()

		<-c.fightTimer.C

		r := types.NewCityDestroyRequest(c.Name)
		c.worldNotifier.CityDestroyed(r)
	}()

	return true
}

// RemoveAlien removes Alien from that City tile.
func (c *City) RemoveAlien(alien *Alien) {
	if alien == nil {
		return
	}

	delete(c.aliens, alien.Name)
}

// Log returns logger with object related fields set.
func (c *City) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger.UpdateContext(c.GetLoggerContext)

	return &logger
}
