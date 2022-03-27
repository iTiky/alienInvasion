package model

import (
	"github.com/itiky/alienInvasion/pkg/logging"
	"github.com/rs/zerolog"
)

// City keeps city connections.
type City struct {
	// Unique City ID
	Name string

	// Names of connected cities for each side (empty if none)
	NorthRoad string
	EastRoad  string
	SouthRoad string
	WestRoad  string
}

// AvailableRoads returns available roads to other cities.
func (c City) AvailableRoads() []string {
	var roads []string
	if c.NorthRoad != "" {
		roads = append(roads, c.NorthRoad)
	}
	if c.EastRoad != "" {
		roads = append(roads, c.EastRoad)
	}
	if c.SouthRoad != "" {
		roads = append(roads, c.SouthRoad)
	}
	if c.WestRoad != "" {
		roads = append(roads, c.WestRoad)
	}

	return roads
}

// HasNorthRoad checks if City has a side connection road.
func (c City) HasNorthRoad() bool {
	return c.NorthRoad != ""
}

// HasEastRoad checks if City has a side connection road.
func (c City) HasEastRoad() bool {
	return c.EastRoad != ""
}

// HasSouthRoad checks if City has a side connection road.
func (c City) HasSouthRoad() bool {
	return c.SouthRoad != ""
}

// HasWestRoad checks if City has a side connection road.
func (c City) HasWestRoad() bool {
	return c.WestRoad != ""
}

// GetLoggerContext enriches logger context with essential City fields.
func (c City) GetLoggerContext(logCtx zerolog.Context) zerolog.Context {
	return logCtx.
		Str(logging.CityNameKey, c.Name)
}
