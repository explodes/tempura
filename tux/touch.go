package tux

import (
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
)

type TouchEvent uint8

const (
	TouchNone TouchEvent = iota
	TouchDown
	TouchDrag
	TouchUp
)

var noTouch = Touch{Event: TouchNone}

type Touch struct {
	Event    TouchEvent
	Position tempura.Vec
}

type TouchInput struct {
	touches      []Touch
	inputAdapter *inputAdapter
}

func NewTouchInput() *TouchInput {
	return &TouchInput{
		inputAdapter: newInputAdapter(),
	}
}

func (t *TouchInput) Update(camera *ebiten.GeoM) {
	t.touches = t.inputAdapter.update(camera)
}

func (t *TouchInput) GetTouch(index int) Touch {
	if index < len(t.touches) {
		return t.touches[index]
	}
	return noTouch
}

func cameraXY(camera *ebiten.GeoM, x, y int) (cx, cy float64) {
	if camera == nil {
		cx = float64(x)
		cy = float64(y)
	} else {
		cx, cy = camera.Apply(float64(x), float64(y))
	}
	return
}

func isDownEvent(event TouchEvent) bool {
	return event == TouchDown || event == TouchDrag
}
