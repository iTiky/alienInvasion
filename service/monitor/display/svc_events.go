package display

import (
	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/service/monitor"
)

var _ monitor.WorldEventsListener = (*Monitor)(nil)

// CityUpdated implements the WorldEventsListener interface.
func (m *Monitor) CityUpdated(city model.City) {
	m.canvas.UpdateCity(city)
}

// CityFightStarted implements the WorldEventsListener interface.
func (m *Monitor) CityFightStarted(cityID string) {
	m.canvas.SetCityOnFight(cityID)
}

// CityDestroyed implements the WorldEventsListener interface.
func (m *Monitor) CityDestroyed(cityID string, alienIDs []string) {
	m.canvas.DestroyCity(cityID, alienIDs)
}

// AlienRelocated implements the WorldEventsListener interface.
func (m *Monitor) AlienRelocated(alienID, newCityID string) {
	m.canvas.RelocateAlien(alienID, newCityID)
}

// AlienDismissed implements the WorldEventsListener interface.
func (m *Monitor) AlienDismissed(alienID, reason string) {
	m.canvas.DestroyAlien(alienID, reason)
}

// SimStatus implements the WorldEventsListener interface.
func (m *Monitor) SimStatus(aliens, cities int, simStopped bool) {
	if simStopped {
		m.canvas.PrintMsg("Simulation stopped")
	}
}
