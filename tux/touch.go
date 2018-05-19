package tux

import (
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
)

// TouchState describes the state of a touch pointer
type TouchState uint8

const (
	// TouchNone is when there are no current touch states associated with a touch
	TouchNone TouchState = iota
	// TouchDown is when a pointer is currently down, but not yet dragged
	TouchDown
	// TouchDrag is when a pointer was down but has been moved while being down
	TouchDrag
	// TouchUp is when a pointer was down or dragged but has been lifted
	TouchUp
)

var noTouch = Touch{State: TouchNone}

// Touch describes a touch event at a given position
type Touch struct {
	// State is the last known state of the touch
	State TouchState
	// Position is the position at which the State changed to its current state
	Position tempura.Vec
}

// TouchInput is a mechanism for collecting all touch events and transforming
// them into meaningful touch states, such as TouchDown, TouchUp, TouchDrag, and TouchNone.
type TouchInput struct {
	touches      []Touch
	inputAdapter *inputAdapter
}

// NewTouchInput creates a new TouchInput
func NewTouchInput() *TouchInput {
	return &TouchInput{
		inputAdapter: newInputAdapter(),
	}
}

// Update will re-compute touch events and positions for a given camera
func (t *TouchInput) Update(camera *ebiten.GeoM) {
	t.touches = t.inputAdapter.update(camera)
}

// GetTouch returns a touch for a specific pointer index.
//
// For mouse events: 0=Left Mouse, 1=Right Mouse, 2=Middle Mouse
// For touch events, the index is equivalent to the pointer index.
func (t *TouchInput) GetTouch(index int) Touch {
	if index < len(t.touches) {
		return t.touches[index]
	}
	return noTouch
}

// cameraXY adjusts x and y for a possibly nil camera
func cameraXY(camera *ebiten.GeoM, x, y int) (cx, cy float64) {
	if camera == nil {
		cx = float64(x)
		cy = float64(y)
	} else {
		cx, cy = camera.Apply(float64(x), float64(y))
	}
	return
}

// IsDownEvent returns if the touch event is a pressed state
func IsDownEvent(event TouchState) bool {
	return event == TouchDown || event == TouchDrag
}
