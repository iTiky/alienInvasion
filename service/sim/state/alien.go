package state

import (
	"context"
	"math/rand"
	"time"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/sim/types"
	"github.com/rs/zerolog"
)

// alienWorldNotifierExpected notifies the World simulation engine about Alien's intentions.
// Since Alien has no idea about what is happening to the World, it asks a sim engine to do stuff.
type alienWorldNotifierExpected interface {
	// MoveAlien sends Alien's move to a next city request.
	MoveAlien(r types.AlienMoveRequest)

	// EvacuateAlien sends Alien evacuation request if Alien has no steps left.
	EvacuateAlien(r types.AlienEvacuateRequest)
}

// Alien keeps an Alien runner state.
type Alien struct {
	model.Alien

	// State
	curLocation model.City
	curSteps    uint

	// Params
	worldNotifier alienWorldNotifierExpected

	// Input event channels
	worldEventsCh chan types.AlienEvent
}

// NewAlien creates a new Alien state.
// Contract: inputs are valid.
func NewAlien(alien model.Alien, startLocation model.City, worldNotifier alienWorldNotifierExpected) *Alien {
	return &Alien{
		Alien:         alien,
		curLocation:   startLocation,
		worldNotifier: worldNotifier,
		worldEventsCh: make(chan types.AlienEvent, 1),
	}
}

// Run is an Alien lifecycle worker which reacts to World events and sends requests to it.
func (a *Alien) Run(ctx context.Context) {
	stepTicker := time.NewTicker(a.Speed)
	defer stepTicker.Stop()

	for working := true; working; {
		select {
		case <-ctx.Done():
			working = false
		case eBz := <-a.worldEventsCh:
			if targetID := eBz.TargetID(); a.Name != targetID {
				a.log(ctx).Warn().Msgf("Event (%T) skipped: targetID mismatch (%s / %s)", eBz, targetID, a.Name)
				break
			}

			switch e := eBz.(type) {
			case types.AlienRelocatedEvent:
				a.handleRelocatedEvent(ctx, e)
			case types.AlienDismissedEvent:
				a.handleDismissedEvent(ctx, e)
				working = false
			default:
				a.log(ctx).Warn().Msgf("Event (%T) skipped: unknown type", eBz)
			}
		case <-stepTicker.C:
			a.handleNextStepEvent(ctx)
		}
	}
}

// Relocate notifies an Alien about a confirmed move.
func (a *Alien) Relocate(e types.AlienRelocatedEvent) {
	a.worldEventsCh <- e
}

// Dismiss notifies an Alien's runner to stop.
func (a *Alien) Dismiss(e types.AlienDismissedEvent) {
	a.worldEventsCh <- e
}

// handleRelocatedEvent handles a received type.AlienRelocatedEvent event.
func (a *Alien) handleRelocatedEvent(ctx context.Context, e types.AlienRelocatedEvent) {
	a.curLocation = e.NewLocation
}

// handleRelocatedEvent handles a received type.AlienRelocatedEvent event.
func (a *Alien) handleDismissedEvent(ctx context.Context, e types.AlienDismissedEvent) {
	a.log(ctx).Info().Msgf("Alien dismissed (%s)", e.Reason)
}

// handleNextStepEvent notifies a World engine about Alien's next move intention.
func (a *Alien) handleNextStepEvent(ctx context.Context) {
	if a.curSteps >= a.MaxSteps {
		// Out of steps
		r := types.NewAlienEvacuateRequest(a.Name)
		a.worldNotifier.EvacuateAlien(r)
		return
	}
	a.curSteps++

	availableRoads := a.curLocation.AvailableRoads()
	if len(availableRoads) == 0 {
		// Nowhere to move
		return
	}

	roadIdx := rand.Intn(len(availableRoads)) //nolint:gosec
	r := types.NewAlienMoveRequest(a.Name, availableRoads[roadIdx])
	a.worldNotifier.MoveAlien(r)
}

// log returns logger with object related fields set.
func (a *Alien) log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.AlienNameKey, a.Name).Logger()

	return &logger
}
