package progressbar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubClock struct {
	now time.Time
}

func newStubClock() *stubClock {
	return &stubClock{
		now: time.Now(),
	}
}

func (c *stubClock) advance(d time.Duration) {
	c.now = c.now.Add(d)
}

func (c *stubClock) Now() time.Time {
	return c.now
}

func (c *stubClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func newWithClock(d time.Duration, samples int) (*MovingRate, *stubClock) {
	clock := newStubClock()
	rate := NewMovingRate(d, samples)
	rate.now = clock.Now
	rate.since = clock.Since

	return rate, clock
}

func TestNewMovingRate(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)
	assert.NotNil(t, rate)
}

func TestMovingRateStart(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)

	err := rate.Start()
	assert.Nil(t, err)

	err = rate.Start()
	assert.ErrorIs(t, err, ErrAlreadyStarted)
}

func TestMovingRateStop(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)

	err := rate.Stop()
	assert.ErrorIs(t, err, ErrNotStarted)

	require.Nil(t, rate.Start())
	err = rate.Stop()
	assert.Nil(t, err)

	err = rate.Stop()
	assert.ErrorIs(t, err, ErrAlreadyStopped)
}

func TestMovingRateSample_NotStarted(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)
	got := rate.Sample(1)
	assert.ErrorIs(t, got, ErrNotStarted)
}

func TestMovingRateSample_Stopped(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)
	require.Nil(t, rate.Start())
	require.Nil(t, rate.Stop())

	got := rate.Sample(1)

	assert.ErrorIs(t, got, ErrStopped)
}

func TestMovingRateSample(t *testing.T) {
	rate, clock := newWithClock(time.Second, 2)
	require.Nil(t, rate.Start())

	rate.Sample(10)
	got, _ := rate.Current()
	assert.Equal(t, 10.0, got, "first sample taken")

	clock.advance(time.Second)
	rate.Sample(20)
	got, _ = rate.Current()
	assert.Equal(t, 15.0, got, "second sample taken")

	clock.advance(time.Second)
	rate.Sample(2)
	got, _ = rate.Current()
	assert.Equal(t, 11.0, got, "third sample taken, first discarded")

	clock.advance(500 * time.Millisecond)
	rate.Sample(4)
	got, _ = rate.Current()
	assert.Equal(t, 11.0, got, "fourth sample sent to accumulator since it's smaller than our interval")

	clock.advance(500 * time.Millisecond)
	rate.Sample(4)
	got, _ = rate.Current()
	assert.Equal(t, 5.0, got, "fourth + fifth sample taken together")
}

func TestMovingRateCurrent_NotStarted(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)

	_, gotErr := rate.Current()

	assert.ErrorIs(t, gotErr, ErrNotStarted)
}

func TestMovingRateCurrent_Initial(t *testing.T) {
	rate := NewMovingRate(time.Second, 10)

	rate.Start()

	got, gotErr := rate.Current()
	assert.Nil(t, gotErr)
	assert.Equal(t, 0.0, got)

	rate.Sample(1)

	got, gotErr = rate.Current()

	assert.Nil(t, gotErr)
	assert.Equal(t, 1.0, got)
}

func TestMovingRateCurrent_GlobalRate(t *testing.T) {
	rate, clock := newWithClock(time.Second, 10)

	rate.Start()
	rate.Sample(1)
	clock.advance(time.Second)
	rate.Sample(1)
	clock.advance(time.Second)
	rate.Stop()

	gotCurrent, gotErr := rate.Current()

	assert.Nil(t, gotErr)
	assert.Equal(t, 1.0, gotCurrent)
}

func TestMovingRateCurrent_WhileGathering(t *testing.T) {
	rate, clock := newWithClock(time.Second, 10)

	rate.Start()
	rate.Sample(1)
	clock.advance(time.Second)
	rate.Sample(1)
	clock.advance(time.Second)

	gotCurrent, gotErr := rate.Current()

	assert.Nil(t, gotErr)
	assert.Equal(t, 1.0, gotCurrent)
}

func TestMovingRateCurrent_Decay(t *testing.T) {
	rate, clock := newWithClock(time.Second, 1)

	rate.Start()
	rate.Sample(10)
	clock.advance(time.Second)
	clock.advance(time.Second)
	clock.advance(time.Second)
	clock.advance(time.Second)
	clock.advance(time.Second)
	rate.Stop()

	got, _ := rate.Current()

	assert.Equal(t, 2.0, got)
}
