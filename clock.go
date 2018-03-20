package tempura

import (
	"time"
)

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
	Sleep(time.Duration)
}

type SystemClock struct{}

func (c *SystemClock) Now() time.Time                  { return time.Now() }
func (c *SystemClock) Since(t time.Time) time.Duration { return time.Since(t) }
func (c *SystemClock) Sleep(d time.Duration)           { time.Sleep(d) }

type FakeClock struct {
	time uint64
}

func (c *FakeClock) Now() time.Time                  { return time.Unix(0, int64(c.time)) }
func (c *FakeClock) Since(t time.Time) time.Duration { return c.Now().Sub(t) }
func (c *FakeClock) Sleep(d time.Duration)           { c.time += uint64(d) }
func (c *FakeClock) Advance(d time.Duration)         { c.time += uint64(d) }

type Stopwatch interface {
	TimeDelta() float64
	Pause()
	Resume()
}

func NewStopwatch() Stopwatch {
	var clock *SystemClock
	return NewStopwatchClock(clock)
}

func NewStopwatchClock(clock Clock) Stopwatch {
	return &stopwatchImpl{
		clock:      clock,
		paused:     false,
		updateTime: clock.Now(),
	}
}

type stopwatchImpl struct {
	clock      Clock
	paused     bool
	skip       float64
	updateTime time.Time
}

func (p *stopwatchImpl) TimeDelta() float64 {
	if p == nil || p.paused {
		return 0
	}

	now := p.clock.Now()
	dt := now.Sub(p.updateTime).Seconds() + p.skip

	p.updateTime = now
	p.skip = 0

	return dt
}

func (p *stopwatchImpl) Pause() {
	p.skip = p.TimeDelta()
	p.paused = true
}

func (p *stopwatchImpl) Resume() {
	p.updateTime = p.clock.Now()
	p.paused = false
}
