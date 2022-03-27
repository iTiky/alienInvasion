package types

import "github.com/itiky/alienInvasion/model"

// Alien to World requests.
type (
	// AlienRequest defines a common request interface.
	AlienRequest interface {
		SourceID() string
	}

	// AlienMoveRequest defines Alien's intention to move from old to new City location.
	AlienMoveRequest struct {
		AlienID   string
		NewCityID string
	}

	// AlienEvacuateRequest defines Alien's intention to evacuate do to "out of steps" event.
	AlienEvacuateRequest struct {
		AlienID string
	}
)

// SourceID implements the AlienRequest interface.
func (r AlienMoveRequest) SourceID() string {
	return r.AlienID
}

// SourceID implements the AlienRequest interface.
func (r AlienEvacuateRequest) SourceID() string {
	return r.AlienID
}

// NewAlienMoveRequest creates a new AlienMoveRequest object.
func NewAlienMoveRequest(alienID, newCityID string) AlienMoveRequest {
	return AlienMoveRequest{
		AlienID:   alienID,
		NewCityID: newCityID,
	}
}

// NewAlienEvacuateRequest creates a new AlienEvacuateRequest object.
func NewAlienEvacuateRequest(alienID string) AlienEvacuateRequest {
	return AlienEvacuateRequest{
		AlienID: alienID,
	}
}

// World to World requests.
type (
	// WorldRequest defines a common request interface.
	WorldRequest interface{}

	// CityDestroyRequest defines City has been destroyed event.
	CityDestroyRequest struct {
		CityID string
	}

	// AlienDisembarkRequest defines Alien created event.
	AlienDisembarkRequest struct {
		Alien  model.Alien
		CityID string
	}
)

// NewCityDestroyRequest creates a new CityDestroyRequest object.
func NewCityDestroyRequest(cityID string) CityDestroyRequest {
	return CityDestroyRequest{
		CityID: cityID,
	}
}

// NewAlienDisembarkRequest creates a new AlienDisembarkRequest object.
func NewAlienDisembarkRequest(alien model.Alien, cityID string) AlienDisembarkRequest {
	return AlienDisembarkRequest{
		Alien:  alien,
		CityID: cityID,
	}
}
