package tempura

import (
	"testing"

	"github.com/hajimehoshi/ebiten"
)


var _ Drawable = (*testDrawable)(nil)

type testObj struct {
	obj *Object

	preCount  int
	stepCount int
	postCount int

	drawable *testDrawable
}

type testDrawable struct {
	drawCount int
}

func (t *testDrawable) DrawAbsolute(image *ebiten.Image, mat ebiten.GeoM) {
	t.drawCount++
}

func (t *testDrawable) Bounds() Rect {
	return R(0, 0, 10, 10)
}

func newTestObject(tag string) *testObj {
	drawable := &testDrawable{}

	testObject := &testObj{
		drawable: drawable,
	}

	testObject.obj = &Object{
		Tag:      tag,
		Drawable: drawable,
		PreSteps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.preCount++
		}),
		Steps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.stepCount++
		}),
		PostSteps: MakeBehaviors(func(source *Object, dt float64) {
			testObject.postCount++
		}),
	}

	return testObject
}

func newTestImage(t *testing.T) *ebiten.Image {
	t.Helper()
	img, err := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	if err != nil {
		t.Fatal(err)
	}
	return img
}
