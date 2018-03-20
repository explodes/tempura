package tempura

// Behavior is what happens when an object meets a condition for a given time delta
type Behavior func(source *Object, dt float64)

// Reaction is what happens when two objects meet a condition for a given time delta.
// Source is the Object performing a behavior "with" is the object the source is
// reacting with.
type Reaction func(source, with *Object, dt float64)

// Behaviors is a slice of Behavior that should happen in succession
type Behaviors []Behavior

// MakeBehaviors is a convenience function for turning a sequence
// of Behaviors or Behavior functions into Behaviors
func MakeBehaviors(behaviors ...Behavior) Behaviors {
	return Behaviors(behaviors)
}

// Execute executes all behaviors for an object with a time delta
func (b Behaviors) Execute(source *Object, dt float64) {
	for _, behavior := range b {
		behavior(source, dt)
	}
}

// Movement is a Behavior that will move a source an object
// by its velocity scaled by time delta
var Movement = Behavior(func(source *Object, dt float64) {
	v := source.Velocity.Scaled(dt)
	source.Pos = source.Pos.Add(v)
})

// FaceDirection is a behavior that adjusts an Object's
// Rot (rotation) to face the same angle as its Velocity.
func FaceDirection(source *Object, dt float64) {
	source.Rot = source.Velocity.Angle()
}
