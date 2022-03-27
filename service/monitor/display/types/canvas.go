package types

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png" // PNG image format registration
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/itiky/alienInvasion/model"
	"github.com/itiky/alienInvasion/service/monitor/display/resource"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var _ ebiten.Game = (*Canvas)(nil)

// Canvas keeps all Aliens and City sprites data and updates on external events.
type Canvas struct {
	citiesLock sync.RWMutex // cities map lock
	cities     citySprites  // Cities set

	aliensLock sync.RWMutex // aliens map lock
	aliens     alienSprites // Aliens set

	status *statusSprite // Status window

	screenWidth, screenHeight int // Window size
}

func NewCanvas(cityMap model.CityMap, aliens []model.Alien) (*Canvas, error) {
	// Sprite default params
	const (
		defCanvasWidth, defCanvasHeight = 800, 600

		cityWidth, cityHeight = 100, 100
		cityOffsetXY          = 50
		cityNameOffsetY       = 25

		roadWidth, roadHeight = 50, 50

		battleWidth, battleHeight = 100, 100

		alienWidth, alienHeight = 50, 50
		alienMoveSpeed          = 30

		cityNameFontSize = 24
		statusFontSize   = 18
		fontDPI          = 72

		statusMsgOffsetXY = 5
		statusMsgBufSize  = 5
		statusMsgMaxLen   = 50
	)

	// Build
	c := Canvas{
		screenWidth:  defCanvasWidth,
		screenHeight: defCanvasHeight,
	}

	// Crate fonts
	fontData, err := opentype.Parse(resource.DefaultFont)
	if err != nil {
		return nil, fmt.Errorf("decoding DefaultFont: %w", err)
	}

	cityNameFontFace, err := opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    cityNameFontSize,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("creating CityName font face: %w", err)
	}

	statusFontFace, err := opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    cityNameFontSize,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("creating status font face: %w", err)
	}

	// Decode images
	cityImage, _, err := image.Decode(bytes.NewReader(resource.CityImage))
	if err != nil {
		return nil, fmt.Errorf("decoding CityImage: %w", err)
	}
	cityEbitenImage := ebiten.NewImageFromImage(cityImage)

	roadImage, _, err := image.Decode(bytes.NewReader(resource.RoadImage))
	if err != nil {
		return nil, fmt.Errorf("decoding RoadImage: %w", err)
	}
	roadEbitenImage := ebiten.NewImageFromImage(roadImage)

	battleImage, _, err := image.Decode(bytes.NewReader(resource.BattleImage))
	if err != nil {
		return nil, fmt.Errorf("decoding RoadImage: %w", err)
	}
	battleEbitenImage := ebiten.NewImageFromImage(battleImage)

	alienImage, _, err := image.Decode(bytes.NewReader(resource.AlienImage))
	if err != nil {
		return nil, fmt.Errorf("decoding AlienImage: %w", err)
	}
	alienEbitenImage := ebiten.NewImageFromImage(alienImage)

	// Create sprites
	citySprites, err := newCitySprites(cityMap,
		[]citySpriteOption{
			withCityImage(cityEbitenImage),
			withCityScaling(cityWidth, cityHeight),
			withCityTopLeftOffset(cityOffsetXY),
			withCityNameFont(cityNameFontFace, color.White, cityNameOffsetY, cityNameFontSize),
			withRoadImage(roadEbitenImage),
			withRoadScaling(roadWidth, roadHeight),
			withBattleImage(battleEbitenImage),
			withBattleScaling(battleWidth, battleHeight),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating cities sprite map: %w", err)
	}
	c.cities = citySprites
	citiesWidth, citiesHeight := c.cities.Size()

	alienSprites, err := newAlienSprites(aliens,
		[]alienSpriteOption{
			withAlienImage(alienEbitenImage),
			withAlienScaling(alienWidth, alienHeight),
			withCitySize(cityOffsetXY, cityWidth, cityHeight),
			withAlienMoveSpeed(alienMoveSpeed),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating aliens sprite map: %w", err)
	}
	c.aliens = alienSprites

	statusSprite, err := newStatusSprite(
		withStatusWindowLocation(0, citiesHeight),
		withStatusMsgBuffer(statusMsgBufSize, statusMsgMaxLen),
		withStatusNameFont(statusFontFace, color.White, statusFontSize, statusMsgOffsetXY),
	)
	if err != nil {
		return nil, fmt.Errorf("creating status sprite: %w", err)
	}
	c.status = statusSprite

	// Adjust the screen size
	c.screenWidth, c.screenHeight = citiesWidth, citiesHeight+c.status.Height()

	return &c, nil
}

// Update implements the ebiten.Game interface.
func (c *Canvas) Update() error {
	if ebiten.IsWindowBeingClosed() {
		return ErrWindowClosed
	}

	return nil
}

// Draw implements the ebiten.Game interface.
func (c *Canvas) Draw(screen *ebiten.Image) {
	c.citiesLock.RLock()
	defer c.citiesLock.RUnlock()

	c.aliensLock.RLock()
	defer c.aliensLock.RUnlock()

	for _, sprite := range c.cities {
		sprite.Draw(screen)
	}

	for _, sprite := range c.aliens {
		sprite.Draw(screen)
	}

	c.status.Draw(screen)
}

// Layout implements the ebiten.Game interface.
func (c *Canvas) Layout(outsideWidth, outsideHeight int) (int, int) {
	return c.screenWidth, c.screenHeight
}

// UpdateCity updates a City (road connection have changed).
func (c *Canvas) UpdateCity(city model.City) {
	c.citiesLock.Lock()
	defer c.citiesLock.Unlock()

	sprite, ok := c.cities[city.Name]
	if !ok {
		return
	}
	sprite.UpdateCityData(city)
}

// SetCityOnFight sets a flag to display "City on fight" sprite.
func (c *Canvas) SetCityOnFight(cityID string) {
	c.citiesLock.Lock()
	defer c.citiesLock.Unlock()

	sprite, ok := c.cities[cityID]
	if !ok {
		return
	}
	sprite.SetOnFight()
}

// DestroyCity removes a City from the Canvas.
func (c *Canvas) DestroyCity(cityID string, alienIDs []string) {
	c.citiesLock.Lock()
	defer c.citiesLock.Unlock()

	delete(c.cities, cityID)

	msg := fmt.Sprintf("City %s destroyed by [%s]", cityID, strings.Join(alienIDs, ","))
	c.status.AddMsg(msg)
}

// RelocateAlien sets a new movement animation target for an Alien.
func (c *Canvas) RelocateAlien(alienID, cityID string) {
	c.aliensLock.Lock()
	defer c.aliensLock.Unlock()

	alienSprite, ok := c.aliens[alienID]
	if !ok {
		return
	}

	citySprite, ok := c.cities[cityID]
	if !ok {
		return
	}

	alienSprite.SetMoveTarget(citySprite.xIdx, citySprite.yIdx)
}

// DestroyAlien removes an Alien from the Canvas.
func (c *Canvas) DestroyAlien(alienID, reason string) {
	c.aliensLock.Lock()
	defer c.aliensLock.Unlock()

	delete(c.aliens, alienID)

	msg := fmt.Sprintf("Alien %s removed (%s)", alienID, reason)
	c.status.AddMsg(msg)
}

// PrintMsg adds a message to the Status sprite.
func (c *Canvas) PrintMsg(msg string) {
	c.status.AddMsg(msg)
}
