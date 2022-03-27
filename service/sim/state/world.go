package state

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/config"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/monitor"
	"github.com/itiky/alienInvasion/service/sim/types"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	_ alienWorldNotifierExpected = (*World)(nil)
	_ cityWorldNotifierExpected  = (*World)(nil)
)

// World keeps the World simulator engine runner state.
type World struct {
	// State
	cities       map[string]*City  // Cities state (key: CityID)
	alienCityMap map[string]string // AlienID-CityID matching map (key: AlienID, value: CityID)

	// Notifiers
	stateNotifier monitor.WorldEventsListener

	// Input request channels
	alienRequestsCh chan types.AlienRequest
	worldRequestsCh chan types.WorldRequest
}

// NewWorld creates a new World state.
// Contract: inputs are valid.
func NewWorld(cityMap model.CityMap, stateNotifier monitor.WorldEventsListener) *World {
	const inputChSize = 100

	w := World{
		cities:          make(map[string]*City, len(cityMap)),
		stateNotifier:   stateNotifier,
		alienRequestsCh: make(chan types.AlienRequest, inputChSize),
		worldRequestsCh: make(chan types.WorldRequest, inputChSize),
	}

	for _, city := range cityMap {
		w.cities[city.Name] = NewCity(city, &w)
	}

	return &w
}

// Run is the World lifecycle worker which reacts to input events from Aliens / Cities and notifies an external service.
func (w *World) Run(ctx context.Context, aliens []model.Alien, simStopCh chan struct{}) {
	// Aliens disembark

	// List of all cities should be done here, as it might change during the operation
	cityIDs := make([]string, 0, len(w.cities))
	for _, city := range w.cities {
		cityIDs = append(cityIDs, city.Name)
	}

	w.alienCityMap = make(map[string]string, len(aliens))
	go w.disembarkAliens(ctx, cityIDs, aliens)

	// Worker
	stopCheckTicker := time.NewTicker(viper.GetDuration(config.AppSimStopCheckRate))
	defer stopCheckTicker.Stop()

	for working := true; working; {
		select {
		case <-ctx.Done():
			working = false
		case <-stopCheckTicker.C:
			if w.checkStopConditions(ctx) {
				working = false
				close(simStopCh)
			}
		case rBz := <-w.alienRequestsCh:
			switch r := rBz.(type) {
			case types.AlienMoveRequest:
				w.handleAlienMoveRequest(ctx, r)
			case types.AlienEvacuateRequest:
				w.handleAlienEvacuateRequest(ctx, r)
			default:
				w.log(ctx).Warn().Msgf("Alien request (%T) skipped: unknown type", rBz)
			}
		case rBz := <-w.worldRequestsCh:
			switch r := rBz.(type) {
			case types.CityDestroyRequest:
				w.handleCityDestroyRequest(ctx, r)
			case types.AlienDisembarkRequest:
				w.handleAlienDisembarkRequest(ctx, r)
			default:
				w.log(ctx).Warn().Msgf("World request (%T) skipped: unknown type", rBz)
			}
		}
	}
}

// CityDestroyed implements the cityWorldNotifierExpected interface.
func (w *World) CityDestroyed(r types.CityDestroyRequest) {
	w.worldRequestsCh <- r
}

// MoveAlien implements the alienWorldNotifierExpected interface.
func (w *World) MoveAlien(r types.AlienMoveRequest) {
	w.alienRequestsCh <- r
}

// EvacuateAlien implements the alienWorldNotifierExpected interface.
func (w *World) EvacuateAlien(r types.AlienEvacuateRequest) {
	w.alienRequestsCh <- r
}

// handleStopCheckEvent checks if simulation should be stopped.
//   * no Aliens left;
//   * no Cities left;
func (w *World) checkStopConditions(ctx context.Context) (retStop bool) {
	aliens, cities := len(w.alienCityMap), len(w.cities)

	defer func() {
		w.stateNotifier.SimStatus(aliens, cities, retStop)
		if retStop {
			w.log(ctx).
				Info().
				Msgf("Simulation stopped (aliens / cities left: %d / %d)", aliens, cities)
		}
	}()

	if aliens <= 1 {
		retStop = true
	}
	if cities == 01 {
		retStop = true
	}

	return
}

// handleAlienDisembarkRequest handles Alien's request to disembark (be created).
func (w *World) handleAlienDisembarkRequest(ctx context.Context, r types.AlienDisembarkRequest) {
	// Check city exists
	city, ok := w.cities[r.CityID]
	if !ok {
		w.log(ctx).Warn().Msgf("Alien disembark failed: city (%s) not found", r.CityID)
		return
	}

	// Disembark
	alienState := NewAlien(r.Alien, city.City, w)
	w.moveAlienTo(ctx, alienState, nil, city)
	go alienState.Run(ctx)
}

// handleAlienMoveRequest handles Alien's request to move.
func (w *World) handleAlienMoveRequest(ctx context.Context, r types.AlienMoveRequest) {
	// Find all related objects
	oldCityID, ok := w.alienCityMap[r.AlienID]
	if !ok {
		return
	}

	oldCity, ok := w.cities[oldCityID]
	if !ok {
		return
	}

	alien, ok := oldCity.aliens[r.AlienID]
	if !ok {
		return
	}

	newCity, ok := w.cities[r.NewCityID]
	if !ok {
		return
	}

	// Relocate
	w.moveAlienTo(ctx, alien, oldCity, newCity)
}

// handleAlienEvacuateRequest handles Alien's request to evacuate which dismisses Alien from the map.
func (w *World) handleAlienEvacuateRequest(ctx context.Context, r types.AlienEvacuateRequest) {
	w.dismissAlien(ctx, r.AlienID, "evacuated")
}

// handleCityDestroyRequest handles City's request to be destroyed dismissing all Aliens, removing City from the map and notifying an external service.
func (w *World) handleCityDestroyRequest(ctx context.Context, r types.CityDestroyRequest) {
	// Find all related objects
	city, ok := w.cities[r.CityID]
	if !ok {
		return
	}

	// Remove connections
	removeConnection := func(connectedCityID string) {
		connectedCity, ok := w.cities[connectedCityID]
		if !ok {
			return
		}

		connectedCity.RemoveRoadsTo(r.CityID)
		w.stateNotifier.CityUpdated(connectedCity.City)
	}

	removeConnection(city.NorthRoad)
	removeConnection(city.EastRoad)
	removeConnection(city.SouthRoad)
	removeConnection(city.WestRoad)

	// Dismiss aliens
	aliensInvolved := city.AlienIDs()
	for _, alien := range city.aliens {
		w.dismissAlien(ctx, alien.Name, "destroyed")
	}

	// Remove city
	delete(w.cities, city.Name)

	// Log
	city.Log(ctx).Info().Msgf("Destroyed by: %s", strings.Join(aliensInvolved, ", "))

	// Notify
	w.stateNotifier.CityDestroyed(city.Name, aliensInvolved)
}

// dismissAlien dismisses a single Alien removing it from a City.
func (w *World) dismissAlien(ctx context.Context, alienID, reason string) {
	// Find all related objects
	cityID, ok := w.alienCityMap[alienID]
	if !ok {
		return
	}

	city, ok := w.cities[cityID]
	if !ok {
		return
	}

	alien, ok := city.aliens[alienID]
	if !ok {
		return
	}

	// Send dismiss event to the Alien
	r := types.NewAlienDismissedEvent(alien.Name, reason)
	alien.Dismiss(r)

	// Remove bindings
	delete(city.aliens, alien.Name)
	delete(w.alienCityMap, alien.Name)

	// Notify
	w.stateNotifier.AlienDismissed(alien.Name, reason)
}

// moveAlienTo moves Alien from an old location (optional) to a new one.
func (w *World) moveAlienTo(ctx context.Context, alien *Alien, oldCity, newCity *City) {
	if alien == nil || newCity == nil {
		return
	}

	// Check if Alien can be moved
	if oldCity != nil && oldCity.AtFight() {
		// Alien can't escape the fight
		return
	}

	// Update the map
	w.alienCityMap[alien.Name] = newCity.Name
	if oldCity != nil {
		oldCity.RemoveAlien(alien)
	}
	newCityFightStarted := newCity.AddAlien(alien)

	// Send relocate event to the Alien
	e := types.NewAlienRelocatedEvent(alien.Name, newCity.City)
	alien.Relocate(e)

	// Notify
	w.stateNotifier.AlienRelocated(alien.Name, newCity.Name)
	if newCityFightStarted {
		w.stateNotifier.CityFightStarted(newCity.Name)
	}
}

// disembarkAliens drops Aliens to a random City.
// Not all Aliens can land, since a target City might be already destroyed (it happens).
func (w *World) disembarkAliens(ctx context.Context, cityIDs []string, aliens []model.Alien) {
	disembarkMinRate, disembarkMaxRate := viper.GetDuration(config.AppAliensDisembarkMinRate), viper.GetDuration(config.AppAliensDisembarkMaxRate)
	disembarkDiff := int64(disembarkMaxRate - disembarkMinRate)

	for _, alien := range aliens {
		// Delay
		disembarkDelay := disembarkMinRate
		if disembarkMaxRate != disembarkMinRate {
			disembarkDelay += time.Duration(rand.Int63n(disembarkDiff)) //nolint:gosec
		}
		time.Sleep(disembarkDelay)

		// Pick a target location
		cityID := cityIDs[rand.Intn(len(cityIDs))] //nolint:gosec

		// Disembark request
		r := types.NewAlienDisembarkRequest(alien, cityID)
		w.worldRequestsCh <- r
	}
}

// log returns logger with object related fields set.
func (w *World) log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)

	return &logger
}
