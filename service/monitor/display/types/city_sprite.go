package types

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/itiky/alienInvasion/model"
	"golang.org/x/image/font"
)

type (
	// citySprite keeps a City sprite data.
	citySprite struct {
		// State
		model.City      // City data
		hasFight   bool // City is on a fight flag (draws Fight sprite on top)
		xIdx, yIdx int  // Sprite matrix coordinates (relative values that are converted to abx X and Y)

		// City sprite params
		cImage           *ebiten.Image // image source
		cX, cY           float64       // City abs coordinates
		cOffsetXY        float64       // sprite top-left offset (from x, y pair)
		cScaleX, cScaleY float64       // image scaling coefs
		cWidth, cHeight  float64       // image actual size (after scaling)

		// City name font params
		cFontFace    font.Face   // font
		cFontColor   color.Color // text color
		cFontOffsetY int         // text Y offset (depends on font size)
		cFontSize    int         // font size

		// Road sprite params
		rImage           *ebiten.Image // image source
		rScaleX, rScaleY float64       // image scaling coefs
		rWidth, rHeight  float64       // image actual size (after scaling)

		// Battle sprite params
		bImage           *ebiten.Image // image source
		bScaleX, bScaleY float64       // image scaling coefs
		bWidth, bHeight  float64       // image actual size (after scaling)
	}

	// citySpriteOption defines the newCitySprite constructor option.
	citySpriteOption func(s *citySprite) error

	// citySprites defines a set of citySprite objects.
	citySprites map[string]*citySprite // key: CityID
)

// withCityImage sets City sprite image.
func withCityImage(image *ebiten.Image) citySpriteOption {
	return func(s *citySprite) error {
		if image == nil {
			return fmt.Errorf("city image: nil")
		}

		w, h := image.Size()
		s.cImage = image
		s.cWidth, s.cHeight = float64(w), float64(h)

		return nil
	}
}

// withCityScaling adjusts City sprite image size coefs.
func withCityScaling(targetWidth, targetHeight int) citySpriteOption {
	return func(s *citySprite) error {
		if targetWidth <= 0 {
			return fmt.Errorf("city targetWidth must be GT 0")
		}
		if targetHeight <= 0 {
			return fmt.Errorf("city targetHeight must be GT 0")
		}

		s.cScaleX = float64(targetWidth) / s.cWidth
		s.cScaleY = float64(targetHeight) / s.cHeight

		s.cWidth *= s.cScaleX
		s.cHeight *= s.cScaleY

		return nil
	}
}

// withCityTopLeftOffset sets City sprite top-left offset.
func withCityTopLeftOffset(offset int) citySpriteOption {
	return func(s *citySprite) error {
		s.cOffsetXY = float64(offset)

		return nil
	}
}

// withCityNameFont sets City name text params.
func withCityNameFont(fFace font.Face, fColor color.Color, offsetY, fSize int) citySpriteOption {
	return func(s *citySprite) error {
		if fFace == nil {
			return fmt.Errorf("city name fontFace: nil")
		}
		if fColor == nil {
			return fmt.Errorf("city name fontColor: nil")
		}
		if fSize <= 0 {
			return fmt.Errorf("city name fontSize: must be GT 0")
		}

		s.cFontFace = fFace
		s.cFontColor = fColor
		s.cFontOffsetY = offsetY
		s.cFontSize = fSize

		return nil
	}
}

// withRoadImage sets Road sprite image.
func withRoadImage(image *ebiten.Image) citySpriteOption {
	return func(s *citySprite) error {
		if image == nil {
			return fmt.Errorf("road image: nil")
		}

		w, h := image.Size()
		s.rImage = image
		s.rWidth, s.rHeight = float64(w), float64(h)

		return nil
	}
}

// withRoadScaling adjusts Road sprite image size coefs.
func withRoadScaling(targetWidth, targetHeight int) citySpriteOption {
	return func(s *citySprite) error {
		if targetWidth <= 0 {
			return fmt.Errorf("road targetWidth must be GT 0")
		}
		if targetHeight <= 0 {
			return fmt.Errorf("road targetHeight must be GT 0")
		}

		s.rScaleX = float64(targetWidth) / s.rWidth
		s.rScaleY = float64(targetHeight) / s.rHeight

		s.rWidth *= s.rScaleX
		s.rHeight *= s.rScaleY

		return nil
	}
}

// withBattleImage sets Battle sprite image.
func withBattleImage(image *ebiten.Image) citySpriteOption {
	return func(s *citySprite) error {
		if image == nil {
			return fmt.Errorf("battle image: nil")
		}

		w, h := image.Size()
		s.bImage = image
		s.bWidth, s.bHeight = float64(w), float64(h)

		return nil
	}
}

// withBattleScaling adjusts Battle sprite image size coefs.
func withBattleScaling(targetWidth, targetHeight int) citySpriteOption {
	return func(s *citySprite) error {
		if targetWidth <= 0 {
			return fmt.Errorf("battle targetWidth must be GT 0")
		}
		if targetHeight <= 0 {
			return fmt.Errorf("battle targetHeight must be GT 0")
		}

		s.bScaleX = float64(targetWidth) / s.bWidth
		s.bScaleY = float64(targetHeight) / s.bHeight

		s.bWidth *= s.bScaleX
		s.bHeight *= s.bScaleY

		return nil
	}
}

// newCitySprite creates a citySprite instance without coordinates.
func newCitySprite(city model.City, opts ...citySpriteOption) (*citySprite, error) {
	// Build
	s := citySprite{
		City:    city,
		cScaleX: 1.0,
		cScaleY: 1.0,
	}

	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	// Validate
	if s.cImage == nil {
		return nil, fmt.Errorf("city image: nil")
	}
	if s.rImage == nil {
		return nil, fmt.Errorf("road image: nil")
	}
	if s.bImage == nil {
		return nil, fmt.Errorf("battle image: nil")
	}
	if s.cFontFace == nil {
		return nil, fmt.Errorf("city fontFace: nil")
	}

	return &s, nil
}

// SetLocation sets the City sprite abs coordinates based on relative sprite matrix values.
func (s *citySprite) SetLocation(xIdx, yIdx int) {
	s.xIdx, s.yIdx = xIdx, yIdx

	s.cX = float64(s.xIdx)*s.cWidth + float64(s.xIdx+1)*s.cOffsetXY
	s.cY = float64(s.yIdx)*s.cHeight + float64(s.yIdx+1)*s.cOffsetXY
}

// UpdateCityData updates the City data.
func (s *citySprite) UpdateCityData(city model.City) {
	s.City = city
}

// SetOnFight sets "City on fight" flag (enabled Fight sprite render).
func (s *citySprite) SetOnFight() {
	s.hasFight = true
}

// Draw implements the ebiten.Game interface.
func (s *citySprite) Draw(screen *ebiten.Image) {
	const (
		rotateCoefRads = 90.0 * math.Pi / 180.0
	)

	drawOpts := &ebiten.DrawImageOptions{
		Filter: ebiten.FilterLinear,
	}

	// Roads
	if s.HasNorthRoad() {
		x := s.cX + s.cWidth/2.0 + s.rHeight/2.0
		y := s.cY - s.rWidth

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.rScaleX, s.rScaleY)
		drawOpts.GeoM.Rotate(rotateCoefRads)
		drawOpts.GeoM.Translate(x, y)
		screen.DrawImage(s.rImage, drawOpts)
	}
	if s.HasEastRoad() {
		x := s.cX + s.cWidth
		y := s.cY + s.cHeight/2.0 - s.rHeight/2.0

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.rScaleX, s.rScaleY)
		drawOpts.GeoM.Translate(x, y)
		screen.DrawImage(s.rImage, drawOpts)
	}
	if s.HasSouthRoad() {
		x := s.cX + s.cWidth/2.0 + s.rHeight/2.0
		y := s.cY + s.cHeight

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.rScaleX, s.rScaleY)
		drawOpts.GeoM.Rotate(rotateCoefRads)
		drawOpts.GeoM.Translate(x, y)
		screen.DrawImage(s.rImage, drawOpts)
	}
	if s.HasWestRoad() {
		x := s.cX - s.rWidth
		y := s.cY + s.cHeight/2.0 - s.rHeight/2.0

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.rScaleX, s.rScaleY)
		drawOpts.GeoM.Translate(x, y)
		screen.DrawImage(s.rImage, drawOpts)
	}

	// City image
	{
		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.cScaleX, s.cScaleY)
		drawOpts.GeoM.Translate(s.cX, s.cY)
		screen.DrawImage(s.cImage, drawOpts)
	}

	// Battle image
	if s.hasFight {
		x := s.cX + s.cWidth/2.0 - s.bWidth/2.0
		y := s.cY + s.cHeight/2.0 - s.bHeight/2.0

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.bScaleX, s.bScaleY)
		drawOpts.GeoM.Translate(x, y)
		screen.DrawImage(s.bImage, drawOpts)
	}

	// City name
	{
		x := s.cX
		y := s.cY + s.cHeight
		text.Draw(screen, s.Name, s.cFontFace, int(x), int(y)+s.cFontOffsetY, s.cFontColor)
	}
}

// newCitySprites creates a new citySprites sprites set.
// Constructor places all cities onto a sprites matrix (relative to other sprites tile matrix).
func newCitySprites(cityMap model.CityMap, citySpriteOpts []citySpriteOption) (citySprites, error) {
	sprites := make(citySprites, len(cityMap))
	if len(cityMap) == 0 {
		return sprites, nil
	}

	// Add sprites without coordinates for now
	var rootSprite *citySprite             // root sprite is used for starting xIdx and yIdx (0, 0)
	spriteSet := make(map[string]struct{}) // set is used to track unprocessed sprites
	for _, city := range cityMap {
		sprite, err := newCitySprite(city, citySpriteOpts...)
		if err != nil {
			return nil, fmt.Errorf("creating citySprite (%s): %w", city.Name, err)
		}

		sprites[city.Name] = sprite
		spriteSet[city.Name] = struct{}{}

		if rootSprite == nil {
			rootSprite = sprite
		}
	}

	// Distribute raw sprite matrix coordinates (raws can be negative)
	xIdxMin, yIdxMin := math.MaxInt, math.MaxInt

	var setCitySpriteNeighboursXYIdx func(sprite *citySprite) // anon func prototype
	setCitySpriteNeighboursXYIdx = func(sprite *citySprite) {
		// Skip if already handled
		if _, ok := spriteSet[sprite.Name]; !ok {
			return
		}

		// Adjust XY mins
		if sprite.xIdx < xIdxMin {
			xIdxMin = sprite.xIdx
		}
		if sprite.yIdx < yIdxMin {
			yIdxMin = sprite.yIdx
		}

		// Remove from set (here to avoid stack overflow)
		delete(spriteSet, sprite.Name)

		// Recursive calls for neighbours
		if sprite.HasNorthRoad() {
			sideSprite := sprites[sprite.NorthRoad]
			sideSprite.SetLocation(sprite.xIdx, sprite.yIdx-1)
			setCitySpriteNeighboursXYIdx(sideSprite)
		}
		if sprite.HasEastRoad() {
			sideSprite := sprites[sprite.EastRoad]
			sideSprite.SetLocation(sprite.xIdx+1, sprite.yIdx)
			setCitySpriteNeighboursXYIdx(sideSprite)
		}
		if sprite.HasSouthRoad() {
			sideSprite := sprites[sprite.SouthRoad]
			sideSprite.SetLocation(sprite.xIdx, sprite.yIdx+1)
			setCitySpriteNeighboursXYIdx(sideSprite)
		}
		if sprite.HasWestRoad() {
			sideSprite := sprites[sprite.WestRoad]
			sideSprite.SetLocation(sprite.xIdx-1, sprite.yIdx)
			setCitySpriteNeighboursXYIdx(sideSprite)
		}
	}

	setCitySpriteNeighboursXYIdx(rootSprite)

	// Set sprite matrix coordinates for leftovers (without neighbours)
	// Extra column on the left is created
	yIdxLeftover := yIdxMin
	if len(spriteSet) > 0 {
		xIdxMin--
	}
	for cityID := range spriteSet {
		sprite := sprites[cityID]
		sprite.SetLocation(xIdxMin, yIdxLeftover)
		yIdxLeftover++
	}

	// Adjust sprite matrix coordinates to remove negative ones
	for _, sprite := range sprites {
		sprite.SetLocation(sprite.xIdx-xIdxMin, sprite.yIdx-yIdxMin)
	}

	return sprites, nil
}

// Size returns sprites matrix actual size.
func (s citySprites) Size() (width int, height int) {
	if len(s) == 0 {
		return 0, 0
	}

	xIdxMax, yIdxMax := 0, 0
	cityWidth, cityHeight := 0, 0
	for _, s := range s {
		if s.xIdx > xIdxMax {
			xIdxMax = s.xIdx
		}
		if s.yIdx > yIdxMax {
			yIdxMax = s.yIdx
		}

		if cityWidth == 0 {
			cityWidth = int(s.cOffsetXY + s.cWidth)
		}
		if cityHeight == 0 {
			cityHeight = int(s.cOffsetXY+s.cHeight) + s.cFontOffsetY + s.cFontSize
		}
	}

	citiesWidth := (xIdxMax + 1) * cityWidth
	citiesHeight := (yIdxMax + 1) * cityHeight

	return citiesWidth, citiesHeight
}
