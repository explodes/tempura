package tinge

import (
	"image"
	"image/color"
)

type Transform func(base image.Image) (transformed image.Image, err error)

var _ color.Color = Color64(0)

type Color64 uint32

func (c Color64) RGBA() (r, g, b, a uint32) {
	r = uint32((c & 0x00ff0000) >> 16)
	r |= r << 8
	g = uint32((c & 0x0000ff00) >> 8)
	g |= g << 8
	b = uint32((c & 0x000000ff) >> 0)
	b |= b << 8
	a = uint32((c & 0xff000000) >> 24)
	a |= a << 8
	return
}

func SetAlpha(base color.Color, alpha uint8) color.Color {
	r, g, b, _ := base.RGBA()
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: alpha,
	}
}

func Colorize(newColor color.Color) Transform {
	return func(base image.Image) (image.Image, error) {
		img := &changeColorImage{
			Image:    base,
			newColor: newColor,
		}
		return img, nil
	}
}

type changeColorImage struct {
	image.Image
	newColor color.Color
}

func (i *changeColorImage) At(x, y int) color.Color {
	base := i.Image.At(x, y)
	_, _, _, ba := base.RGBA()
	return SetAlpha(i.newColor, uint8(ba))
}
