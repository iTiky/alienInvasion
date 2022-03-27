package types

import "github.com/itiky/alienInvasion/model"

// World to Alien events.
type (
	// AlienEvent defines a common event interface.
	AlienEvent interface {
		TargetID() string
	}

	// AlienRelocatedEvent defines a confirmed Alien move from old to new City.
	AlienRelocatedEvent struct {
		AlienID     string
		NewLocation model.City
	}

	// AlienDismissedEvent defines a confirmed Alien dismiss event with a reason comment.
	AlienDismissedEvent struct {
		AlienID string
		Reason  string
	}
)

// TargetID implements the AlienEvent interface.
func (e AlienRelocatedEvent) TargetID() string {
	return e.AlienID
}

// TargetID implements the AlienEvent interface.
func (e AlienDismissedEvent) TargetID() string {
	return e.AlienID
}

// NewAlienRelocatedEvent creates a new AlienRelocatedEvent object.
func NewAlienRelocatedEvent(alienID string, newLocation model.City) AlienRelocatedEvent {
	return AlienRelocatedEvent{
		AlienID:     alienID,
		NewLocation: newLocation,
	}
}

// NewAlienDismissedEvent creates a new AlienDismissedEvent object.
func NewAlienDismissedEvent(alienID, reason string) AlienDismissedEvent {
	return AlienDismissedEvent{
		AlienID: alienID,
		Reason:  reason,
	}
}
