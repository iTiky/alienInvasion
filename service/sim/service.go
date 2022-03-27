package sim

import (
	"context"
	"fmt"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/monitor"
	"github.com/itiky/alienInvasion/service/monitor/noop"
	"github.com/itiky/alienInvasion/service/sim/state"
)

type (
	// Processor implements the World simulation engine.
	Processor struct {
		// Params
		cityMap model.CityMap
		aliens  []model.Alien
		monitor monitor.WorldEventsListener

		// State
		worldState *state.World
	}

	Option func(p *Processor) error
)

// WithCityMap is the Processor constructor option that sets the CityMap param.
func WithCityMap(cm model.CityMap) Option {
	return func(p *Processor) error {
		if err := cm.Validate(); err != nil {
			return fmt.Errorf("validating city map: %w", err)
		}
		p.cityMap = cm

		return nil
	}
}

// WithAliens is the Processor constructor option that sets the Aliens param.
func WithAliens(aliens []model.Alien) Option {
	return func(p *Processor) error {
		p.aliens = aliens

		return nil
	}
}

// WithMonitor is the Processor constructor option that sets the Monitor param.
func WithMonitor(monitor monitor.WorldEventsListener) Option {
	return func(p *Processor) error {
		if monitor == nil {
			return fmt.Errorf("monitor service: nil")
		}
		p.monitor = monitor

		return nil
	}
}

// New creates a new Processor instance and performs basic dependencies validation.
func New(opts ...Option) (*Processor, error) {
	// Construction
	p := Processor{
		monitor: noop.New(),
	}
	for _, opt := range opts {
		if err := opt(&p); err != nil {
			return nil, err
		}
	}

	// Validation
	if len(p.cityMap) == 0 {
		return nil, fmt.Errorf("city map is not defined (empty)")
	}
	if len(p.aliens) == 0 {
		return nil, fmt.Errorf("aliens are not defined (empty)")
	}

	return &p, nil
}

// Start starts the simulation engine and returns simulation stopped channel (close channel)
// Contract: config is valid.
func (p *Processor) Start(ctx context.Context) chan struct{} {
	// Enrich logger context
	ctx, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Simulator").Logger()
	ctx = logging.SetCtxLogger(ctx, logger)

	// Start the engine worker
	simStopCh := make(chan struct{})

	worldState := state.NewWorld(p.cityMap, p.monitor)
	go worldState.Run(ctx, p.aliens, simStopCh)
	p.worldState = worldState

	return simStopCh
}
