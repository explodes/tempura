package tempura

import (
	"github.com/hajimehoshi/ebiten"
)

// Object is a game object that has basic physics, optional
// graphics, and associated Behaviors. It can be used standalone
// or managed (Updated and Drawn) by Objects.
type Object struct {
	// Tag is an optional identifier for this type of object.
	// It can be retrieved as an ObjectSet from an Objects by
	// this tag along with other Objects with the same tag.
	Tag string

	// Pos is the position of the Object. The Drawable, if any,
	// will be drawn with this as the origin.
	Pos Vec
	// Size is the size of the Object. The Drawable, if any,
	// will be scaled to fit.
	Size Vec
	// Velocity is the Vec describing the movement speed
	// and direction of this Object.
	Velocity Vec

	// Drawable is an optional Drawable to use to draw this
	// Object on a Target.
	Drawable Drawable
	// Rot is an amount in radians used to rotate the Drawable
	// where 0 degrees is right and 90 degrees is upwards.
	Rot float64
	// RotNormal is the amount that the drawable should be rotated
	// initially such that its default orientation is right-facing,
	// or 0 degrees.
	RotNormal float64

	// PreSteps is Behaviors to execute before Steps and
	// PostSteps during an Update performed by Objects.
	PreSteps Behaviors
	// Steps is Behaviors to execute before PostSteps and
	// after PreSteps during an Update performed by Objects.
	Steps Behaviors
	// PostSteps is Behaviors to execute after Steps during
	// an Update performed by Objects.
	PostSteps Behaviors
}

// Bounds gets the hitbox for this Object. Any Drawable will
// scaled and translated to fit this box. Collision detection
// can be performed using this Rect.
func (o *Object) Bounds() Rect {
	return R(o.Pos.X, o.Pos.Y, o.Pos.X+o.Size.X, o.Pos.Y+o.Size.Y)
}

// Draw will render this Object on a target if a Drawable is associated with
// this Object. The Object's Drawable will be scaled and translated to fit
// this Object's Bounds. It will also be rotated by Rot radians to
// This function does nothing if this Object has no Drawable.
func (o *Object) Draw(image *ebiten.Image) {
	if o.Drawable == nil {
		return
	}
	bounds := o.Bounds()
	mat := FitRotated(o.Rot+o.RotNormal, o.Drawable.Bounds(), bounds)
	o.Drawable.Draw(image, mat)
}
