package tempura

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type Drawable interface {
	DrawAbsolute(image *ebiten.Image, mat ebiten.GeoM)
	Bounds() Rect
}

var _ Drawable = (*ImageDrawable)(nil)

type ImageDrawable struct {
	src       *ebiten.Image
	frames    []Rect
	imgFrames []image.Rectangle
	frameNum  int
	opts      *ebiten.DrawImageOptions
}

func imageRectangleToRect(r image.Rectangle) Rect {
	return R(
		float64(r.Min.X),
		float64(r.Min.Y),
		float64(r.Max.X),
		float64(r.Max.Y),
	)
}

func NewImageDrawable(src *ebiten.Image) *ImageDrawable {
	frame := imageRectangleToRect(src.Bounds())
	return NewImageDrawableFrames(src, frame)
}

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

func (d *ImageDrawable) SetFrame(frameNum int) {
	d.frameNum = frameNum
}

func (d *ImageDrawable) NumFrames() int {
	return len(d.frames)
}

func (d *ImageDrawable) DrawAbsolute(image *ebiten.Image, mat ebiten.GeoM) {
	frame := d.imgFrames[d.frameNum]
	d.opts.SourceRect = &frame

	d.opts.GeoM = mat

	image.DrawImage(d.src, d.opts)
}

func (d *ImageDrawable) Bounds() Rect {
	return d.frames[d.frameNum]
}

func MakeFrames(wImage, hImage, cols, rows, totalFrames int) []Rect {
	fw := float64(wImage)
	fh := float64(hImage)
	wFrame, hFrame := fw/float64(cols), fh/float64(rows)
	frames := make([]Rect, 0, cols*rows)
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
