package tinge

import (
	"image/color"
	"image/draw"
)

func DrawCircle(img draw.Image, cx, cy, radius int, border color.Color) {
	// Algorithm taken from
	// http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	x0, y0 := cx, cy
	f := 1 - radius
	ddF_x, ddF_y := 1, -2*radius
	cx, cy = 0, radius

	img.Set(x0, y0+radius, border)
	img.Set(x0, y0-radius, border)
	img.Set(x0+radius, y0, border)
	img.Set(x0-radius, y0, border)

	for cx < cy {
		if f >= 0 {
			cy--
			ddF_y += 2
			f += ddF_y
		}
		cx++
		ddF_x += 2
		f += ddF_x
		img.Set(x0+cx, y0+cy, border)
		img.Set(x0-cx, y0+cy, border)
		img.Set(x0+cx, y0-cy, border)
		img.Set(x0-cx, y0-cy, border)
		img.Set(x0+cy, y0+cx, border)
		img.Set(x0-cy, y0+cx, border)
		img.Set(x0+cy, y0-cx, border)
		img.Set(x0-cy, y0-cx, border)
	}
}
