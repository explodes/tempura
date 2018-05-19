package tinge

import (
	"image/color"
	"image/draw"
	"math"
)

func DrawCircle(canvas draw.Image, cx, cy, radius int, border color.Color) {
	// Algorithm taken from
	// http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	x0, y0 := cx, cy
	f := 1 - radius
	ddF_x, ddF_y := 1, -2*radius
	cx, cy = 0, radius

	canvas.Set(x0, y0+radius, border)
	canvas.Set(x0, y0-radius, border)
	canvas.Set(x0+radius, y0, border)
	canvas.Set(x0-radius, y0, border)

	for cx < cy {
		if f >= 0 {
			cy--
			ddF_y += 2
			f += ddF_y
		}
		cx++
		ddF_x += 2
		f += ddF_x
		canvas.Set(x0+cx, y0+cy, border)
		canvas.Set(x0-cx, y0+cy, border)
		canvas.Set(x0+cx, y0-cy, border)
		canvas.Set(x0-cx, y0-cy, border)
		canvas.Set(x0+cy, y0+cx, border)
		canvas.Set(x0-cy, y0+cx, border)
		canvas.Set(x0+cy, y0-cx, border)
		canvas.Set(x0-cy, y0-cx, border)
	}
}

func DrawLine(canvas draw.Image, x0, y0, x1, y1 int, stroke color.Color) {
	// Algorithm taken from
	// https://en.wikipedia.org/wiki/Digital_differential_analyzer_(graphics_algorithm)

	dx := float64(x1 - x0)
	dy := float64(y1 - y0)

	adx := math.Abs(float64(dx))
	ady := math.Abs(float64(dy))

	var steps float64
	if adx > ady {
		steps = adx
	} else {
		steps = ady
	}

	if steps == 0 {
		return
	}

	xIncrement := dx / steps
	yIncrement := dy / steps

	x := float64(x0)
	y := float64(y0)
	canvas.Set(x0, y0, stroke)

	for i := 0; i < int(steps); i++ {
		x += xIncrement
		y += yIncrement
		canvas.Set(int(x), int(y), stroke)
	}

}
