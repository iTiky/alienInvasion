package display

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/monitor/display/types"
	"github.com/rs/zerolog"
)

const serviceName = "CanvasMonitor"

type (
	// Monitor defines a service to visualize the simulation.
	Monitor struct {
		canvas *types.Canvas // Sprites storage

		screenWidth, screenHeight int // Screen size
	}

	// Option defines the New constructor options.
	Option func(m *Monitor) error
)

// WithScreenSize overrides default screen size.
func WithScreenSize(width, height int) Option {
	return func(m *Monitor) error {
		if width <= 0 {
			return fmt.Errorf("screen width: must be GT 0")
		}
		if height <= 0 {
			return fmt.Errorf("screen height: must be GT 0")
		}

		m.screenWidth, m.screenHeight = width, height

		return nil
	}
}

// New creates a new Monitor instance.
func New(cityMap model.CityMap, aliens []model.Alien, opts ...Option) (*Monitor, error) {
	m := Monitor{
		screenWidth:  1000,
		screenHeight: 800,
	}

	for _, opt := range opts {
		if err := opt(&m); err != nil {
			return nil, err
		}
	}

	canvas, err := types.NewCanvas(cityMap, aliens)
	if err != nil {
		return nil, fmt.Errorf("canvas build: %w", err)
	}
	m.canvas = canvas

	return &m, nil
}

// Run starts the ebiten.Game.
// Uses the runtime.LockOSThread call and must be started from the main routine.
// Contract: not canceled by the {ctx}, so used have to close the window.
func (m *Monitor) Run(ctx context.Context) {
	ebiten.SetWindowTitle("Alien invasion simulation")
	ebiten.SetWindowSize(1024, 800)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowClosingHandled(true)

	if err := ebiten.RunGame(m.canvas); err != nil {
		if errors.Is(err, types.ErrWindowClosed) {
			return
		}
		m.log(ctx).Warn().Err(err).Msg("RunGame failed")
	}
}

// log returns logger with object related fields set.
func (m *Monitor) log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
