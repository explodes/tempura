// +build android ios

package tux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type inputAdapter struct {
	touches []Touch
	// tested is the list of touch indices that have been tested during
	// an update loop. indices won't get tested if there are no events
	// for a given pointer ID, so their corresponding Touch should get
	// changed to TouchUp or TouchNone.
	tested []bool
}

func newInputAdapter() *inputAdapter {
	return &inputAdapter{
		touches: make([]Touch, 3),
		tested:  make([]bool, 3),
	}
}

func (a *inputAdapter) update(camera *ebiten.GeoM) []Touch {

	// we haven't tested any touches
	for i := 0; i < len(a.tested); i++ {
		a.tested[i] = false
	}

	// process down touch events
	touches := ebiten.Touches()
	for _, touch := range touches {
		a.downTouch(camera, touch)
	}

	// process up touch events
	a.upTouches()

	return a.touches
}

func (a *inputAdapter) ensurePointerIndex(index int) {
	if len(a.touches) <= index {
		next := make([]Touch, index+1)
		copy(next, a.touches)
		a.touches = next
	}
	if len(a.tested) <= index {
		next := make([]bool, index+1)
		copy(next, a.tested)
		a.tested = next
	}
}

func (a *inputAdapter) downTouch(camera *ebiten.GeoM, touch ebiten.Touch) {
	index := touch.ID()
	a.ensurePointerIndex(index)

	wasDown := isDownEvent(a.touches[index].Event)

	x, y := touch.Position()
	cx, cy := cameraXY(camera, x, y)

	// update event type
	switch {
	case !wasDown:
		if a.touches[index].Event != TouchDown {
			fmt.Println("TouchDown")
		}
		a.touches[index].Event = TouchDown
	case wasDown:
		if a.touches[index].Event == TouchDown && cx != a.touches[index].Position.X && cy != a.touches[index].Position.Y {
			if a.touches[index].Event != TouchDrag {
				fmt.Println("TouchDrag")
			}
			a.touches[index].Event = TouchDrag
		}
	}

	// update position
	a.touches[index].Position.X = float64(cx)
	a.touches[index].Position.Y = float64(cy)
	a.tested[index] = true
}

func (a *inputAdapter) upTouches() {
	for index, tested := range a.tested {
		if tested {
			continue
		}
		wasDown := isDownEvent(a.touches[index].Event)
		switch {
		case !wasDown:
			if a.touches[index].Event != TouchNone {
				fmt.Println("TouchNone")
			}
			a.touches[index].Event = TouchNone
		case wasDown:
			if a.touches[index].Event != TouchUp {
				fmt.Println("TouchUp")
			}
			a.touches[index].Event = TouchUp
		}
	}
}
