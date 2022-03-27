package model

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CityMap keeps city map data (key: City name).
type CityMap map[string]City

// Validate performs CityMap integrity validation:
//   * city name is valid;
//   * city name exists for road links;
//   * city road links are correct (if city A is connected to city B, B must be connected to A as well);
func (m CityMap) Validate() error {
	cityNameRegexp := regexp.MustCompile(`^[a-zA-Z]+(?:[\s\-)][a-zA-Z]+)*$`)

	cityNameSet := make(map[string]struct{}, len(m))
	for name := range m {
		cityNameSet[name] = struct{}{}
	}

	for k, city := range m {
		if k != city.Name {
			return fmt.Errorf("%s: key mismatch with .Name (%s)", k, city.Name)
		}

		if !cityNameRegexp.MatchString(k) {
			return fmt.Errorf("%s: invalid city name", k)
		}

		if city.NorthRoad != "" && m[city.NorthRoad].SouthRoad != k {
			return fmt.Errorf("%s: invalid North road link (%s)", k, city.NorthRoad)
		}
		if city.EastRoad != "" && m[city.EastRoad].WestRoad != k {
			return fmt.Errorf("%s: invalid Easth road link (%s)", k, city.EastRoad)
		}
		if city.SouthRoad != "" && m[city.SouthRoad].NorthRoad != k {
			return fmt.Errorf("%s: invalid South road link (%s)", k, city.SouthRoad)
		}
		if city.WestRoad != "" && m[city.WestRoad].EastRoad != k {
			return fmt.Errorf("%s: invalid West road link (%s)", k, city.WestRoad)
		}
	}

	return nil
}

// NewCityMapFromFile parses city map file.
// Format:
//   {CityName} [(north/east/south/west)={OtherCityName}]
// Example:
//   Foo north=Bar west=Baz south=Qu-ux
//   Bar south=Foo west=Bee
func NewCityMapFromFile(filePath string) (CityMap, error) {
	const (
		northRoadKey = "north"
		eastRoadKey  = "east"
		southRoadKey = "south"
		westRoadKey  = "west"
	)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	cityMap := make(CityMap)
	for lineN := 0; scanner.Scan(); lineN++ {
		lineParts := strings.Split(scanner.Text(), " ")
		if len(lineParts) == 0 {
			continue
		}

		city := City{
			Name: lineParts[0],
		}

		for _, roadBz := range lineParts[1:] {
			roadParts := strings.Split(roadBz, "=")
			if len(roadParts) != 2 {
				return nil, fmt.Errorf("parsing line (%d): road link (%s): invalid format (north=CityName is expected)", lineN, roadBz)
			}

			switch side, cityName := strings.ToLower(roadParts[0]), roadParts[1]; side {
			case northRoadKey:
				city.NorthRoad = cityName
			case eastRoadKey:
				city.EastRoad = cityName
			case southRoadKey:
				city.SouthRoad = cityName
			case westRoadKey:
				city.WestRoad = cityName
			default:
				return nil, fmt.Errorf("parsing line (%d): road link: invalid side (%s) (north/east/south/west is expected)", lineN, side)
			}
		}

		if _, ok := cityMap[city.Name]; ok {
			return nil, fmt.Errorf("parsing line (%d): duplicate found", lineN)
		}
		cityMap[city.Name] = city
	}

	return cityMap, nil
}
