package tempura

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

// DegToRad converts degrees to radians
func DegToRad(deg float64) (rad float64) {
	return deg * math.Pi / 180
}

// RadToDeg converts radians to degrees
func RadToDeg(rad float64) (deg float64) {
	return rad * 180 / math.Pi
}

// Fit returns the Matrix that will transform a source Rect
// into the dest Rect
func Fit(source, dest Rect) ebiten.GeoM {
	scaleX := dest.W() / source.W()
	scaleY := dest.H() / source.H()

	mat := ebiten.GeoM{}
	mat.Translate(-source.Min.X, -source.Min.Y)
	mat.Scale(scaleX, scaleY)
	mat.Translate(dest.Min.X, dest.Min.Y)

	return mat
}

// FitGeoM returns the Matrix that will transform a source Rect
// into the dest Rect
func FitRotated(rot float64, source, dest Rect) ebiten.GeoM {
	scaleX := dest.W() / source.W()
	scaleY := dest.H() / source.H()

	mat := ebiten.GeoM{}

	// rotate about center of source
	mat.Translate(-source.W()/2, -source.W()/2)
	mat.Rotate(rot)
	mat.Translate(source.W()/2, source.W()/2)

	// scale
	mat.Scale(scaleX, scaleY)

	// move to destination
	mat.Translate(dest.Min.X, dest.Min.Y)

	return mat
}

// Collision returns if two rectangles intersect
func Collision(r1, r2 Rect) bool {
	if r1.Min.X > r2.Max.X || r2.Min.X > r1.Max.X {
		return false
	}
	if r1.Min.Y > r2.Max.Y || r2.Min.Y > r1.Max.Y {
		return false
	}
	return true
}
