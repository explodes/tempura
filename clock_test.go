package tempura

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStopwatchImpl_TimeDelta(t *testing.T) {
	clock := &FakeClock{}
	timer := NewStopwatchClock(clock)

	assert.Equal(t, 0.0, timer.TimeDelta())

	clock.Advance(1 * time.Second)

	assert.Equal(t, 1.0, timer.TimeDelta())

	clock.Advance(2 * time.Second)

	assert.Equal(t, 2.0, timer.TimeDelta())
}

func TestStopwatchImpl_Pause(t *testing.T) {
	clock := &FakeClock{}
	timer := NewStopwatchClock(clock)

	clock.Advance(2 * time.Second)
	timer.Pause()
	clock.Advance(10 * time.Second)

	assert.Equal(t, 0.0, timer.TimeDelta())
}

func TestStopwatchImpl_SkipPause(t *testing.T) {
	clock := &FakeClock{}
	timer := NewStopwatchClock(clock)

	clock.Advance(2 * time.Second)
	timer.Pause()
	clock.Advance(10 * time.Second)
	timer.Resume()

	assert.Equal(t, 2.0, timer.TimeDelta())
}

func TestStopwatchImpl_PauseResume(t *testing.T) {
	clock := &FakeClock{}
	timer := NewStopwatchClock(clock)

	timer.Pause()
	clock.Advance(10 * time.Second)
	timer.Resume()

	assert.Equal(t, 0.0, timer.TimeDelta())

	timer.Pause()
	clock.Advance(10 * time.Second)
	timer.Resume()

	assert.Equal(t, 0.0, timer.TimeDelta())
}
