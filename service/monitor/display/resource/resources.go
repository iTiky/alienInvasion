package resource

import _ "embed"

var (
	//go:embed Planet.png
	PlanetImage []byte

	//go:embed Alien.png
	AlienImage []byte

	//go:embed City.png
	CityImage []byte

	//go:embed Road.png
	RoadImage []byte

	//go:embed Battle.png
	BattleImage []byte
)

var (
	//go:embed mplus-1p-regular.ttf
	DefaultFont []byte
)
