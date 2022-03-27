package noop

import (
	"strings"

	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/itiky/alienInvasion/service/monitor"
	"github.com/rs/zerolog"
)

const serviceName = "NoopMonitor"

var _ monitor.WorldEventsListener = (*Monitor)(nil)

type (
	Monitor struct {
		logsEnabled bool
		logger      zerolog.Logger
	}

	Option func(m *Monitor)
)

// WithLogs enables logs for all events.
func WithLogs() Option {
	return func(m *Monitor) {
		m.logsEnabled = true
	}
}

// New creates a new Monitor instance.
func New(opts ...Option) *Monitor {
	m := Monitor{
		logger: logging.NewLogger(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

// CityUpdated implements the WorldEventsListener interface.
func (m *Monitor) CityUpdated(city model.City) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "CityUpdated").
		Msgf("CityID = %s", city.Name)
}

// CityFightStarted implements the WorldEventsListener interface.
func (m *Monitor) CityFightStarted(cityID string) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "CityFightStarted").
		Msgf("CityID = %s", cityID)
}

// CityDestroyed implements the WorldEventsListener interface.
func (m *Monitor) CityDestroyed(cityID string, aliens []string) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "CityDestroyed").
		Msgf("CityID = %s, Aliens = [%s]", cityID, strings.Join(aliens, ","))
}

// AlienRelocated implements the WorldEventsListener interface.
func (m *Monitor) AlienRelocated(alienID, newCityID string) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "AlienRelocated").
		Msgf("AlienID = %s, NewCityID = %s", alienID, newCityID)
}

// AlienDismissed implements the WorldEventsListener interface.
func (m *Monitor) AlienDismissed(alienID, reason string) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "AlienDismissed").
		Msgf("AlienID = %s, Reason = %s", alienID, reason)
}

// SimStatus implements the WorldEventsListener interface.
func (m *Monitor) SimStatus(aliens, cities int, simStopped bool) {
	if !m.logsEnabled {
		return
	}

	m.logger.
		Debug().
		Str(logging.ServiceKey, serviceName).
		Str("event", "SimStatus").
		Msgf("Aliens = %d, Cities = %d, Stopped = %v", aliens, cities, simStopped)
}
