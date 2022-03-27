package types

import (
	"container/ring"
	"fmt"
	"image/color"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type (
	statusSprite struct {
		sync.RWMutex

		sX, sY int

		tFontFace  font.Face   // font
		tFontColor color.Color // text color
		tFontSize  int         // font size

		msgBuf     *ring.Ring
		msgBufSize int
		msgLenMax  int
	}

	statusSpriteOption func(s *statusSprite) error
)

// withStatusWindowSize sets Status window size.
func withStatusWindowLocation(x, y int) statusSpriteOption {
	return func(s *statusSprite) error {
		if x < 0 {
			return fmt.Errorf("status window X: must be GT 0")
		}
		if y < 0 {
			return fmt.Errorf("status window Y: must be GT 0")
		}

		s.sX, s.sY = x, y

		return nil
	}
}

// withStatusWindowSize sets Status window size.
func withStatusMsgBuffer(size, msgMaxLen int) statusSpriteOption {
	return func(s *statusSprite) error {
		if size <= 0 {
			return fmt.Errorf("status msg buffer size: must be GT 0")
		}
		if msgMaxLen <= 0 {
			return fmt.Errorf("status msg max length: must be GT 0")
		}

		s.msgBufSize, s.msgLenMax = size, msgMaxLen

		return nil
	}
}

// withStatusNameFont sets Status msg text params.
func withStatusNameFont(fFace font.Face, fColor color.Color, fSize, fOffsetXY int) statusSpriteOption {
	return func(s *statusSprite) error {
		if fFace == nil {
			return fmt.Errorf("status fontFace: nil")
		}
		if fColor == nil {
			return fmt.Errorf("status fontColor: nil")
		}
		if fSize <= 0 {
			return fmt.Errorf("status fontSize: must be GT 0")
		}
		if fOffsetXY <= 0 {
			return fmt.Errorf("status fontOffsetXY: must be GT 0")
		}

		s.tFontFace = fFace
		s.tFontColor = fColor
		s.tFontSize = fSize

		return nil
	}
}

func newStatusSprite(opts ...statusSpriteOption) (*statusSprite, error) {
	// Build
	s := statusSprite{
		msgBufSize: 5,
		msgLenMax:  20,
	}

	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	s.msgBuf = ring.New(s.msgBufSize)

	// Validate
	if s.tFontFace == nil {
		return nil, fmt.Errorf("status fontFace: nil")
	}

	return &s, nil
}

func (s *statusSprite) AddMsg(msg string) {
	s.Lock()
	defer s.Unlock()

	msgRunes := []rune(msg)
	if len(msgRunes) > s.msgLenMax {
		msgRunes = msgRunes[:s.msgLenMax]
	}

	s.msgBuf.Value = string(msgRunes)
	s.msgBuf = s.msgBuf.Next()
}

func (s *statusSprite) Height() int {
	return s.msgBufSize * s.tFontSize * 2
}

// Draw implements the ebiten.Game interface.
func (s *statusSprite) Draw(screen *ebiten.Image) {
	s.RLock()
	defer s.RUnlock()

	statusText := strings.Builder{}
	s.msgBuf.Do(func(i interface{}) {
		msg, ok := i.(string)
		if !ok {
			return
		}
		statusText.WriteString(msg + "\n")
	})

	text.Draw(screen, statusText.String(), s.tFontFace, s.sX, s.sY, s.tFontColor)
}
