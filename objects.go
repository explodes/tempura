package tempura

import (
	"github.com/cevaris/ordered_map"
	"github.com/hajimehoshi/ebiten"
)

// ObjectContainer is a simple interface for Object containers
type ObjectContainer interface {
	Len() int
	Contains(obj *Object) bool
	Iterator() ObjectIterator
}

// ObjectContainer is a simple interface for Object containers that support tags
type TaggedObjectContainer interface {
	ObjectContainer
	TagIterator(tags ...string) ObjectIterator
}

// ObjectIterator is a function used for iterating over collections of Objects.
// It is used as follows:
//
// iter := set.Iterator()
// for obj, ok := iter(); ok; obj, ok = iter() {
//   ..use obj..
// }
// Removing an object during iteration is undefined.
type ObjectIterator func() (next *Object, ok bool)

// Layers is a container for multiple Objects collections such that
// a particular drawing order can be preserved. Updates and Draws will
// happen from the lowest layer to the highest layer.
type Layers []*Objects

// NewLayers creates a new container of Objects with a given amount of layers.
func NewLayers(n int) Layers {
	layers := make(Layers, 0, n)
	for i := 0; i < n; i++ {
		layers = append(layers, NewObjects())
	}
	return layers
}

// Update updates all Objects. Updates happen in the first layer forward.
func (ly Layers) Update(dt float64) {
	for _, layer := range ly {
		layer.Update(dt)
	}
}

// Draw draws all Objects Draws happen in the first layer forward.
func (ly Layers) Draw(image *ebiten.Image) {
	for _, layer := range ly {
		layer.Draw(image)
	}
}

// Len returns the total number of Objects in all layers
func (ly Layers) Len() int {
	sum := 0
	for _, layer := range ly {
		sum += layer.Len()
	}
	return sum
}

// Contains checks to see if an Object is in any layer
func (ly Layers) Contains(obj *Object) bool {
	for _, layer := range ly {
		if layer.Contains(obj) {
			return true
		}
	}
	return false
}

// Iterator returns an ObjectIterator for all objects in all layers
// from the lowest layer to highest
func (ly Layers) Iterator() ObjectIterator {
	iters := make([]ObjectIterator, len(ly))
	for index, layer := range ly {
		iters[index] = layer.All().Iterator()
	}
	return chainIterators(iters)
}

// Iterator returns an ObjectIterator for all objects in all layers
// from the highest layer to lowest
func (ly Layers) IteratorTop() ObjectIterator {
	iters := make([]ObjectIterator, len(ly))
	for index := len(ly) - 1; index >= 0; index-- {
		iters[index] = ly[index].All().Iterator()
	}
	return chainIterators(iters)
}

// Iterator returns an ObjectIterator for all objects
// with the given tags in all layers from the lowest layer to highest
func (ly Layers) TagIterator(tags ...string) ObjectIterator {
	if len(tags) == 0 {
		return emptyObjectIterator
	}
	iters := make([]ObjectIterator, len(ly)*len(tags))
	index := 0
	for _, layer := range ly {
		for _, tag := range tags {
			iters[index] = layer.Tagged(tag).Iterator()
			index++
		}
	}
	return chainIterators(iters)
}

// Iterator returns an ObjectIterator for all objects
// with the given tags in all layers from the highest layer to lowest
func (ly Layers) TagIteratorTop(tags ...string) ObjectIterator {
	if len(tags) == 0 {
		return emptyObjectIterator
	}
	iters := make([]ObjectIterator, len(ly)*len(tags))
	index := 0
	for layerIndex := len(ly) - 1; layerIndex >= 0; layerIndex-- {
		for _, tag := range tags {
			iters[index] = ly[layerIndex].Tagged(tag).Iterator()
			index++
		}
	}
	return chainIterators(iters)
}

// Objects is a container of Object so that Objects can be quickly added
// and removed from a single source. Objects are also retrievable by tag
// allowing for quick access for a particular subset of Object.
//
// The Tag of an Object should not be modified after being added to this
// container.
type Objects struct {
	all    *ObjectSet
	tagged objectTagMap
}

// NewObjects makes a new Objects container.
func NewObjects() *Objects {
	return &Objects{
		all:    NewObjectSet(),
		tagged: make(objectTagMap),
	}
}

// Len returns the amount of Objects in this container
func (o *Objects) Len() int {
	return o.all.Len()
}

// All returns the ObjectSet containing all Objects in this container
func (o *Objects) All() *ObjectSet {
	return o.all
}

// Tagged returns an ObjectSet containing all Objects in this container
// that have a particular tag. Tags with empty strings are not recorded
// and Objects whose tags were modified after being added are not
// considered.
func (o *Objects) Tagged(tag string) *ObjectSet {
	return o.tagged[tag]
}

// Add adds an object to this container. If the Object has a Tag, that
// tag is used to quickly access a particular subset of Object.
func (o *Objects) Add(obj *Object) {
	o.all.Add(obj)
	if obj.Tag != "" {
		o.tagged.add(obj.Tag, obj)
	}
}

// Remove removes an object from this container.
func (o *Objects) Remove(obj *Object) {
	o.all.Remove(obj)
	if obj.Tag != "" {
		o.tagged.remove(obj.Tag, obj)
	}
}

// Contains tests to see if an object is contained in this container.
func (o *Objects) Contains(obj *Object) bool {
	return o.all.Contains(obj)
}

// Update performs all PreSteps, then all Steps, then all PostSteps
// of Object in this container.
func (o *Objects) Update(dt float64) {
	iter := o.all.Iterator()
	for object, ok := iter(); ok; object, ok = iter() {
		object.PreSteps.Execute(object, dt)
	}
	iter = o.all.Iterator()
	for object, ok := iter(); ok; object, ok = iter() {
		object.Steps.Execute(object, dt)
	}
	iter = o.all.Iterator()
	for object, ok := iter(); ok; object, ok = iter() {
		object.PostSteps.Execute(object, dt)
	}
}

// Draw draws all Object in this container.
func (o *Objects) Draw(image *ebiten.Image) {
	iter := o.all.Iterator()
	for object, ok := iter(); ok; object, ok = iter() {
		object.Draw(image)
	}
}

// Iterator gets an ObjectIterator for all Object in this container
func (o *Objects) Iterator() ObjectIterator {
	return o.All().Iterator()
}

// TagIterator gets an ObjectIterator for all Object in this
// container with the given tags
func (o *Objects) TagIterator(tags ...string) ObjectIterator {
	if len(tags) == 0 {
		return emptyObjectIterator
	}
	iters := make([]ObjectIterator, len(tags))
	index := 0
	for _, tag := range tags {
		iters[index] = o.Tagged(tag).Iterator()
		index++
	}
	return chainIterators(iters)
}

// ObjectSet is an ordered set of Object.
type ObjectSet struct {
	set *ordered_map.OrderedMap
}

// NewObjectSet creates a new empty set
func NewObjectSet() *ObjectSet {
	return &ObjectSet{
		set: ordered_map.NewOrderedMap(),
	}
}

// emptyObjectIterator is an ObjectIterator that always
// returns the ultimate result.
func emptyObjectIterator() (*Object, bool) {
	return nil, false
}

// Iterator returns an iterator function that can be used
// to iterate over all objects in this set.
func (os *ObjectSet) Iterator() ObjectIterator {
	if os == nil {
		return emptyObjectIterator
	}
	iter := os.set.IterFunc()
	return func() (*Object, bool) {
		next, ok := iter()
		if ok {
			return next.Key.(*Object), true
		}
		return nil, false
	}
}

// Len returns the size of this set
func (os *ObjectSet) Len() int {
	if os == nil {
		return 0
	}
	return os.set.Len()
}

// Contains tests if an Object is contained in this set
func (os *ObjectSet) Contains(obj *Object) bool {
	if os == nil {
		return false
	}
	_, ok := os.set.Get(obj)
	return ok
}

// Add adds objects to this set
func (os *ObjectSet) Add(obj *Object) {
	os.set.Set(obj, struct{}{})
}

// Remove removes objects from this set
func (os *ObjectSet) Remove(obj *Object) {
	os.set.Delete(obj)
}

// objectTagMap is a defaultdict-like map for adding and removing
// objects from an ObjectSet by tag
type objectTagMap map[string]*ObjectSet

func (m objectTagMap) add(tag string, obj *Object) {
	set := m[tag]
	if set == nil {
		set = NewObjectSet()
		m[tag] = set
	}
	set.Add(obj)
}

func (m objectTagMap) remove(tag string, obj *Object) {
	set := m[tag]
	if set != nil {
		set.Remove(obj)
	}
}

// chainIterators will iterate through a slice of ObjectIterator consecutively
func chainIterators(iters []ObjectIterator) ObjectIterator {
	index := 0
	return func() (*Object, bool) {
		for index < len(iters) {
			next, ok := iters[index]()
			if !ok {
				index++
				continue
			}
			return next, ok
		}
		return nil, false
	}
}
