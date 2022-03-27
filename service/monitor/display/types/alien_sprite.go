package types

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/itiky/alienInvasion/model"
)

const (
	alienSpriteStateIdle    = iota // skip draw
	alienSpriteStateLocated        // static draw
	alienSpriteStateMoving         // movement animation
)

type (
	// alienSprite keeps an Alien sprite data alongside movement animation state.
	alienSprite struct {
		// State
		name          string  // AlienID
		movementState int     // Movement state [idle, located, moving]
		x, y          float64 // Current abs coordinates

		// Moving state values
		xTarget, yTarget float64 // target abx coordinates
		xV, yV           float64 // movement speed
		moveSteps        int     // animation spped coef
		moveStepsLeft    int     // number of move steps left

		// Alien sprite params
		aImage           *ebiten.Image // image source
		aScaleX, aScaleY float64       // image scaling coefs
		aWidth, aHeight  float64       // image actual size (after scaling)
		aColor           color.RGBA    // image color adjustments

		// City sprite params
		cTileWidth     float64 // city sprite width
		cTileHeight    float64 // city sprite height
		cCenterOffsetX float64 // city center X coordinate offset
		cCenterOffsetY float64 // city center Y coordinate offset
	}

	// alienSpriteOption defines the newAlienSprite constructor option.
	alienSpriteOption func(s *alienSprite) error

	// alienSprites defines a set of alienSprite objects.
	alienSprites map[string]*alienSprite // key: AlienID
)

// withAlienImage sets Alien sprite image.
func withAlienImage(image *ebiten.Image) alienSpriteOption {
	return func(s *alienSprite) error {
		if image == nil {
			return fmt.Errorf("alien image: nil")
		}

		w, h := image.Size()
		s.aImage = image
		s.aWidth, s.aHeight = float64(w), float64(h)

		return nil
	}
}

// withAlienScaling adjusts Alien sprite image size coefs.
func withAlienScaling(targetWidth, targetHeight int) alienSpriteOption {
	return func(s *alienSprite) error {
		if targetWidth <= 0 {
			return fmt.Errorf("alien targetWidth must be GT 0")
		}
		if targetHeight <= 0 {
			return fmt.Errorf("alien targetHeight must be GT 0")
		}

		s.aScaleX = float64(targetWidth) / s.aWidth
		s.aScaleY = float64(targetHeight) / s.aHeight

		s.aWidth *= s.aScaleX
		s.aHeight *= s.aScaleY

		return nil
	}
}

// withCitySize sets City sprite offset coordinates to draw an Alien in the center of a City.
func withCitySize(offsetXY, width, height int) alienSpriteOption {
	return func(s *alienSprite) error {
		s.cTileWidth = float64(offsetXY) + float64(width)
		s.cTileHeight = float64(offsetXY) + float64(height)

		s.cCenterOffsetX = float64(offsetXY) + float64(width)/2.0
		s.cCenterOffsetY = float64(offsetXY) + float64(height)/2.0

		return nil
	}
}

// withAlienMoveSpeed overrides the default animation speed.
func withAlienMoveSpeed(speed int) alienSpriteOption {
	return func(s *alienSprite) error {
		if speed < 0 {
			return fmt.Errorf("animation speed: must be GTE 0")
		}

		s.moveSteps = speed

		return nil
	}
}

// newAlienSprite creates an alienSprite instance.
func newAlienSprite(alienID string, opts ...alienSpriteOption) (*alienSprite, error) {
	// Build
	s := alienSprite{
		name:      alienID,
		aScaleX:   1.0,
		aScaleY:   1.0,
		moveSteps: 30,
	}

	// Random Alien color
	randBytes := make([]byte, 3)
	if _, err := rand.Read(randBytes); err != nil { // nolint: gosec
		return nil, fmt.Errorf("reading random bytes: %w", err)
	}

	s.aColor = color.RGBA{
		R: randBytes[0],
		G: randBytes[1],
		B: randBytes[2],
		A: 0x0,
	}

	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	// Validate
	if s.aImage == nil {
		return nil, fmt.Errorf("alien image: nil")
	}

	return &s, nil
}

// SetMoveTarget sets a new movement animation target.
func (s *alienSprite) SetMoveTarget(xIdx, yIdx int) {
	//xIdxTarget, yIdxTarget := float64(xIdx), float64(yIdx)

	s.xTarget = float64(xIdx)*s.cTileWidth + s.cCenterOffsetX - s.aWidth/2.0
	s.yTarget = float64(yIdx)*s.cTileHeight + s.cCenterOffsetY - s.aHeight/2.0

	//s.xTarget = (xIdxTarget+1.0)*s.cCenterOffsetX + xIdxTarget*(s.aWidth/2.0)
	//s.yTarget = (yIdxTarget+1.0)*s.cCenterOffsetY + yIdxTarget*(s.aHeight/2.0)

	s.movementState = alienSpriteStateMoving

	s.moveStepsLeft = s.moveSteps
	s.xV = (s.xTarget - s.x) / float64(s.moveStepsLeft)
	s.yV = (s.yTarget - s.y) / float64(s.moveStepsLeft)
}

// Draw implements the ebiten.Game interface.
func (s *alienSprite) Draw(screen *ebiten.Image) {
	switch s.movementState {
	case alienSpriteStateIdle:
		return
	case alienSpriteStateLocated:
	case alienSpriteStateMoving:
		s.updateMovingState()
	}

	drawOpts := &ebiten.DrawImageOptions{
		Filter: ebiten.FilterLinear,
	}

	// Alien image
	{
		r := float64(s.aColor.R) / 0xFF
		g := float64(s.aColor.G) / 0xFF
		b := float64(s.aColor.B) / 0xFF
		a := float64(s.aColor.A) / 0xFF

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Scale(s.aScaleX, s.aScaleY)
		drawOpts.GeoM.Translate(s.x, s.y)
		drawOpts.ColorM.Reset()
		drawOpts.ColorM.Translate(r, g, b, a)
		screen.DrawImage(s.aImage, drawOpts)
	}
}

// updateMovingState updates current coordinates for moving animation state.
func (s *alienSprite) updateMovingState() {
	s.x += s.xV
	s.y += s.yV

	s.moveStepsLeft--
	if s.moveStepsLeft <= 0 {
		s.movementState = alienSpriteStateLocated
		s.x, s.y = s.xTarget, s.yTarget
	}
}

// newAlienSprites creates a new alienSprites sprites set.
func newAlienSprites(aliens []model.Alien, alienSpriteOpts []alienSpriteOption) (alienSprites, error) {
	sprites := make(alienSprites, len(aliens))

	for _, alien := range aliens {
		sprite, err := newAlienSprite(alien.Name, alienSpriteOpts...)
		if err != nil {
			return nil, fmt.Errorf("creating alienSprite (%s): %w", alien.Name, err)
		}
		sprites[alien.Name] = sprite
	}

	return sprites, nil
}
