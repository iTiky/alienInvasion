package monitor

import "github.com/itiky/alienInvasion/model"

// WorldEventsListener defines an external service that reacts to World / City / Alien events.
type WorldEventsListener interface {
	// CityUpdated is triggered when a City connection roads have been updated.
	CityUpdated(city model.City)

	// CityFightStarted is triggered when a City fight has started / prolonged.
	CityFightStarted(cityID string)

	// CityDestroyed is triggered when a City has been destroyed.
	CityDestroyed(cityID string, alienIDs []string)

	// AlienRelocated is triggered when an Alien has moved.
	AlienRelocated(alienID, newCityID string)

	// AlienDismissed is triggered when an Alien has been dismissed (evacuated / destroyed).
	AlienDismissed(alienID, reason string)

	// SimStatus is periodically triggered to inform about the current simulation state.
	SimStatus(aliens, cities int, stimStopped bool)
}
