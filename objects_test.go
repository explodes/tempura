package tempura

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_ TaggedObjectContainer = Layers{}
	_ TaggedObjectContainer = &Objects{}
	_ ObjectContainer       = &ObjectSet{}
)

func TestBehaviors_Execute(t *testing.T) {
	count := 0
	behavior := func(source *Object, dt float64) {
		count++
	}
	b := MakeBehaviors(behavior, behavior)

	b.Execute(nil, 0)

	assert.Equal(t, 2, count)
}

func TestMovement(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Pos = V(10, 10)
	obj.obj.Velocity = V(10, 10)

	Movement(obj.obj, 1)

	assert.Equal(t, V(20, 20), obj.obj.Pos)
}

func TestFaceDirection(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Velocity = V(1, 1)

	FaceDirection(obj.obj, 1)

	assert.Equal(t, V(1, 1).Angle(), obj.obj.Rot)
}

func TestFaceDirection_zero_velocity(t *testing.T) {
	obj := newTestObject("")
	obj.obj.Velocity = V(0, 0)

	FaceDirection(obj.obj, 1)

	assert.Equal(t, 0.0, obj.obj.Rot)
}

func TestNewLayers(t *testing.T) {
	layers := NewLayers(10)

	assert.Len(t, layers, 10)
}

func TestLayers_Draw(t *testing.T) {
	testObj := newTestObject("tag")
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	layers.Draw(nil, newTestImage(t))

	assert.Equal(t, 1, testObj.drawable.drawCount)
}

func TestLayers_Draw_no_drawable(t *testing.T) {
	testObj := newTestObject("tag")
	testObj.obj.Drawable = nil
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	layers.Draw(nil, newTestImage(t))

	assert.Equal(t, 0, testObj.drawable.drawCount)
}

func TestLayers_Update(t *testing.T) {
	testObj := newTestObject("tag")
	layers := NewLayers(1)
	layers[0].Add(testObj.obj)

	assert.Equal(t, 0, testObj.preCount)
	assert.Equal(t, 0, testObj.stepCount)
	assert.Equal(t, 0, testObj.postCount)

	layers.Update(1)

	assert.Equal(t, 1, testObj.preCount)
	assert.Equal(t, 1, testObj.stepCount)
	assert.Equal(t, 1, testObj.postCount)
}

func TestNewObjects(t *testing.T) {
	objects := NewObjects()

	assert.NotNil(t, objects)
}

func TestObjects_Add_tagged(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.False(t, objects.All().Contains(testObj.obj))
	assert.False(t, objects.Tagged("tag").Contains(testObj.obj))

	objects.Add(testObj.obj)

	assert.True(t, objects.All().Contains(testObj.obj))
	assert.True(t, objects.Tagged("tag").Contains(testObj.obj))
}

func TestObjects_Add_untagged(t *testing.T) {
	testObj := newTestObject("")
	objects := NewObjects()

	objects.Add(testObj.obj)

	assert.True(t, objects.All().Contains(testObj.obj))
	assert.Nil(t, objects.Tagged(""))
}

func TestObjects_Len(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Len())

	objects.Add(testObj.obj)

	assert.Equal(t, 1, objects.Len())
}

func TestObjects_Remove(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.Equal(t, 0, objects.Len())

	objects.Add(testObj.obj)
	objects.Remove(testObj.obj)

	assert.Equal(t, 0, objects.Len())
}

func TestObjects_Contains(t *testing.T) {
	testObj := newTestObject("tag")
	objects := NewObjects()

	assert.False(t, objects.Contains(testObj.obj))

	objects.Add(testObj.obj)

	assert.True(t, objects.Contains(testObj.obj))
}

func TestObjects_Iterator(t *testing.T) {
	testObj := newTestObject("tag")
	testObj2 := newTestObject("tag2")
	objects := NewObjects()
	objects.Add(testObj.obj)
	objects.Add(testObj2.obj)

	iter := objects.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 2, iterSize)
}

func TestObjects_TagIterator(t *testing.T) {
	testObj := newTestObject("tag")
	testObj2 := newTestObject("tag2")
	objects := NewObjects()
	objects.Add(testObj.obj)
	objects.Add(testObj2.obj)

	iter := objects.TagIterator("tag")
	iterSize := countIterator(iter)

	assert.Equal(t, 1, iterSize)
}

func TestObjects_TagIterator_emptyTags(t *testing.T) {
	testObj := newTestObject("tag")
	testObj2 := newTestObject("tag2")
	objects := NewObjects()
	objects.Add(testObj.obj)
	objects.Add(testObj2.obj)

	iter := objects.TagIterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 0, iterSize)
}

func TestObjectSet_Contains(t *testing.T) {
	testObj := newTestObject("tag")
	set := NewObjectSet()

	assert.False(t, set.Contains(testObj.obj))

	set.Add(testObj.obj)

	assert.True(t, set.Contains(testObj.obj))
}

func TestObjectSet_Iterator_nil(t *testing.T) {
	var set *ObjectSet = nil

	iter := set.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 0, iterSize)
}

func TestObjectSet_IteratorOrder(t *testing.T) {
	set := NewObjectSet()
	var expectedOrder []*Object
	for i := 0; i < 100; i++ {
		testObj := newTestObject("tag")
		expectedOrder = append(expectedOrder, testObj.obj)
		set.Add(testObj.obj)
	}

	count := 0
	iter := set.Iterator()
	for next, ok := iter(); ok; next, ok = iter() {
		assert.Equal(t, expectedOrder[count], next)
		count++
	}

	assert.Equal(t, len(expectedOrder), count)
}

func TestObjectSet_Len_nil(t *testing.T) {
	var set *ObjectSet = nil

	assert.Equal(t, 0, set.Len())
}

func TestObjectSet_Contains_nil(t *testing.T) {
	var set *ObjectSet = nil

	assert.False(t, set.Contains(nil))
}

func countIterator(iter ObjectIterator) int {
	count := 0
	for _, ok := iter(); ok; _, ok = iter() {
		count++
	}
	return count
}

func TestLayers_Iterator(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("b").obj)

	iter := layers.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 3, iterSize)
}

func TestLayers_Iterator_skipLayer(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	iter := layers.Iterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 3, iterSize)
}

func TestLayers_TagIterator_skipLayer(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	iter := layers.TagIterator("c")
	iterSize := countIterator(iter)

	assert.Equal(t, 1, iterSize)
}

func TestLayers_TagIterator_emptyTags(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	iter := layers.TagIterator()
	iterSize := countIterator(iter)

	assert.Equal(t, 0, iterSize)
}

func TestLayers_TagIterator(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[1].Add(newTestObject("d").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	iter := layers.TagIterator("b")
	iterSize := countIterator(iter)

	assert.Equal(t, 2, iterSize)
}

func TestLayers_TagIterator_multipleTags(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[1].Add(newTestObject("d").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	iter := layers.TagIterator("a", "b")
	iterSize := countIterator(iter)

	assert.Equal(t, 3, iterSize)
}

func TestLayers_Len(t *testing.T) {
	layers := NewLayers(3)
	layers[0].Add(newTestObject("a").obj)
	layers[1].Add(newTestObject("b").obj)
	layers[1].Add(newTestObject("d").obj)
	layers[2].Add(newTestObject("b").obj)
	layers[2].Add(newTestObject("c").obj)

	size := layers.Len()

	assert.Equal(t, 5, size)
}

func TestLayers_Contains(t *testing.T) {
	testObj := newTestObject("a")
	layers := NewLayers(3)

	assert.False(t, layers.Contains(testObj.obj))

	layers[0].Add(testObj.obj)

	assert.True(t, layers.Contains(testObj.obj))
}
