package tempura

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDegToRad(t *testing.T) {
	as := assert.New(t)

	as.Equal(float64(0), DegToRad(0))
	as.Equal(math.Pi, DegToRad(180))
	as.Equal(2*math.Pi, DegToRad(360))
}

func TestRadToDeg(t *testing.T) {
	as := assert.New(t)

	as.Equal(float64(0), RadToDeg(0))
	as.Equal(float64(180), RadToDeg(math.Pi))
	as.Equal(float64(360), RadToDeg(2*math.Pi))
}
func TestFit(t *testing.T) {

	cases := []struct {
		name     string
		src, dst Rect
	}{
		{name: "untranslated", src: R(0, 0, 10, 11), dst: R(0, 0, 100, 110)},
		{name: "pretranslated", src: R(5, 6, 10, 11), dst: R(0, 0, 100, 110)},
		{name: "posttranslated", src: R(0, 0, 10, 11), dst: R(10, 10, 100, 110)},
		{name: "translated", src: R(10, 11, 300, 450), dst: R(23, 67, 899, 9202)},
		{name: "large_to_small", src: R(0, 0, 1000, 1000), dst: R(10, 10, 20, 20)},
		{name: "small_to_large", src: R(10, 10, 20, 20), dst: R(0, 0, 1000, 1000)},
		{name: "realistic", src: R(0, 0, 2037, 768), dst: R(0, 0, 800, 600)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testFit(t, c.name, c.src, c.dst)
		})
		t.Run(c.name+"Restore", func(t *testing.T) {
			testFit(t, c.name+"Restore", c.dst, c.src)
		})
	}

}

func vectorWithin1(a, b Vec) bool {
	return math.Abs(a.X-b.X) < 1 && math.Abs(a.Y-b.Y) < 1
}

func testFit(t *testing.T, name string, src, dst Rect) {
	as := assert.New(t)

	mat := Fit(src, dst)

	srcMin := V(src.Min.X, src.Min.Y)
	srcMax := V(src.Max.X, src.Max.Y)

	dstMin := V(dst.Min.X, dst.Min.Y)
	dstMax := V(dst.Max.X, dst.Max.Y)

	resultMin := V(mat.Apply(src.Min.X, src.Min.Y))
	resultMax := V(mat.Apply(src.Max.X, src.Max.Y))

	as.True(vectorWithin1(dstMin, resultMin), "MIN name=%s from=%v to=%v result=%v", name, srcMin, dstMin, resultMin)
	as.True(vectorWithin1(dstMax, resultMax), "MAX name=%s from=%v to=%v result=%v", name, srcMax, dstMax, resultMax)
}
