package tempura

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Drawable interface {
	// DrawAbsolute draws this Drawable onto an image with the supplied transformation.
	// The transform has already had a camera applied to it.
	DrawAbsolute(image *ebiten.Image, mat ebiten.GeoM)

	// Bounds returns the dimension of the current frame of this Drawable.
	Bounds() Rect
}

var _ Drawable = (*ImageDrawable)(nil)

// ImageDrawable is a Drawable that is backed by an ebiten.Image.
// The backing image can be broken up into multiple frames, as in
// a sprite sheet or animation.
type ImageDrawable struct {
	src       *ebiten.Image
	frames    []Rect
	imgFrames []image.Rectangle
	frameNum  int
	opts      *ebiten.DrawImageOptions
}

// imageRectangleToRect convert an image.Rectangle into a Rect.
func imageRectangleToRect(r image.Rectangle) Rect {
	return R(
		float64(r.Min.X),
		float64(r.Min.Y),
		float64(r.Max.X),
		float64(r.Max.Y),
	)
}

// NewImageDrawable creates a new ImageDrawable consisting of a single frame.
func NewImageDrawable(src *ebiten.Image) *ImageDrawable {
	frame := imageRectangleToRect(src.Bounds())
	return NewImageDrawableFrames(src, frame)
}

// NewImageDrawableFrames creates a new ImageDrawable with the given frames.
func NewImageDrawableFrames(src *ebiten.Image, frames ...Rect) *ImageDrawable {
	imgFrames := make([]image.Rectangle, len(frames))
	for i, f := range frames {
		imgFrames[i] = image.Rect(
			int(f.Min.X),
			int(f.Min.Y),
			int(f.Max.X),
			int(f.Max.Y),
		)
	}
	return &ImageDrawable{
		src:       src,
		frameNum:  0,
		frames:    frames,
		imgFrames: imgFrames,
		opts:      &ebiten.DrawImageOptions{},
	}
}

// SetFrame sets the frame to draw
func (d *ImageDrawable) SetFrame(frameNum int) {
	d.frameNum = frameNum
}

// NumFrames returns the total number of frames in this drawable.
func (d *ImageDrawable) NumFrames() int {
	return len(d.frames)
}

// DrawAbsolute draws this ImageDrawable onto a canvas with the given transform.
func (d *ImageDrawable) DrawAbsolute(image *ebiten.Image, mat ebiten.GeoM) {
	frame := d.imgFrames[d.frameNum]
	d.opts.SourceRect = &frame

	d.opts.GeoM = mat

	image.DrawImage(d.src, d.opts)
}

// Bounds returns the bounds of the current frame.
func (d *ImageDrawable) Bounds() Rect {
	return d.frames[d.frameNum]
}

// MakeFrames is a utility for creating frames out of a sprite sheet.
// It does not account for strides between frames or padding, so it
// evenly divides an image according to the given spec.
func MakeFrames(imageWidth, imageHeight, imageColumns, imageRows, totalFrames int) []Rect {
	fw := float64(imageWidth)
	fh := float64(imageHeight)
	wFrame, hFrame := fw/float64(imageColumns), fh/float64(imageRows)
	frames := make([]Rect, 0, imageColumns*imageRows)
	frame := 0
	for y := 0.0; y < fh; y += hFrame {
		for x := 0.0; x < fw; x += wFrame {
			frames = append(frames, R(x, y, x+wFrame, y+hFrame))
			frame++
			if frame > totalFrames {
				return frames
			}
		}
	}
	return frames
}
