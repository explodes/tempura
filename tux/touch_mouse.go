// +build !android
// +build !ios

package tux

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type inputAdapter struct {
	touches []Touch
}

func newInputAdapter() *inputAdapter {
	return &inputAdapter{
		touches: make([]Touch, 3),
	}
}

func (a *inputAdapter) update(camera *ebiten.GeoM) []Touch {
	a.updateTouch(camera, 0, ebiten.MouseButtonLeft)
	a.updateTouch(camera, 1, ebiten.MouseButtonRight)
	a.updateTouch(camera, 2, ebiten.MouseButtonMiddle)
	return a.touches
}

func (a *inputAdapter) updateTouch(camera *ebiten.GeoM, index int, button ebiten.MouseButton) {
	wasDown := IsDownEvent(a.touches[index].State)
	justPressed := inpututil.IsMouseButtonJustPressed(button)
	justReleased := inpututil.IsMouseButtonJustReleased(button)
	isDown := justPressed || (wasDown && !justReleased)

	x, y := ebiten.CursorPosition()
	cx, cy := cameraXY(camera, x, y)

	switch {
	case !wasDown && !isDown:
		a.touches[index].State = TouchNone
	case !wasDown && isDown:
		a.touches[index].State = TouchDown
	case wasDown && isDown:
		if a.touches[index].State == TouchDown && cx != a.touches[index].Position.X && cy != a.touches[index].Position.Y {
			a.touches[index].State = TouchDrag
		}
	case wasDown && !isDown:
		a.touches[index].State = TouchUp
	}

	a.touches[index].Position.X = float64(cx)
	a.touches[index].Position.Y = float64(cy)
}
