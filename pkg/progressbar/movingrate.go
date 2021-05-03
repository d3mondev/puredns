package progressbar

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrNotStarted is an error happening when the MovingRate hasn't been started.
	ErrNotStarted = errors.New("rate is not started")

	// ErrAlreadyStarted is an error happening when the MovingRate has already been started.
	ErrAlreadyStarted = errors.New("rate is already started")

	// ErrStopped is an error happening when the MovingRate has been stopped.
	ErrStopped = errors.New("rate has been stopped")

	// ErrAlreadyStopped is an error happening when the MovingRate has already been stopped.
	ErrAlreadyStopped = errors.New("rate is already stopped")
)

// MovingRate calculates the rate of elements sampled using a moving average.
type MovingRate struct {
	now   func() time.Time
	since func(time.Time) time.Duration

	mu sync.Mutex

	movingAvgSamplingRate time.Duration
	movingAvgMaxSamples   int
	movingAvgSamples      []float64

	accumCounter     float64
	accumCounterTime time.Time

	startTime    time.Time
	stopTime     time.Time
	totalCounter float64
}

// NewMovingRate creates a new MovingRate object with the specified sampling rate and number of samples
// to consider in the moving average.
func NewMovingRate(samplingRate time.Duration, samples int) *MovingRate {
	return &MovingRate{
		now:   time.Now,
		since: time.Since,

		movingAvgSamplingRate: samplingRate,
		movingAvgMaxSamples:   samples,
	}
}

// Start starts the MovingRate object.
func (r *MovingRate) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.startTime.IsZero() {
		return ErrAlreadyStarted
	}

	r.startTime = r.now()

	return nil
}

// Stop stops gathering
func (r *MovingRate) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.startTime.IsZero() {
		return ErrNotStarted
	}

	if !r.stopTime.IsZero() {
		return ErrAlreadyStopped
	}

	r.stopTime = r.now()

	return nil
}

// Sample records new data in the moving average. If there is not enough time elapsed between
// the previous call of Sample, the data is accumulated into a buffer until a proper rate can
// be calculated.
func (r *MovingRate) Sample(count float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Return an error if the sampler hasn't been started
	if r.startTime.IsZero() {
		return ErrNotStarted
	}

	// Return an error if the sampler has been stopped
	if !r.stopTime.IsZero() {
		return ErrStopped
	}

	// Set initial value
	if r.accumCounterTime.IsZero() {
		r.totalCounter += count
		r.movingAvgSamples = append(r.movingAvgSamples, count)
		r.accumCounterTime = r.now()
		return nil
	}

	// Accumulate values
	r.accumCounter += count
	r.totalCounter += count

	// Don't update the rates if we're below our sampling rate
	delta := r.since(r.accumCounterTime).Seconds()
	if delta < r.movingAvgSamplingRate.Seconds() {
		return nil
	}

	// Calculate the current rate and add it to the moving average
	curRate := r.accumCounter / delta
	r.movingAvgSamples = append(r.movingAvgSamples, curRate)
	r.accumCounter = 0

	// Trim moving average values if we have too many samples
	if len(r.movingAvgSamples) > r.movingAvgMaxSamples {
		r.movingAvgSamples = r.movingAvgSamples[1:]
	}

	r.accumCounterTime = r.now()

	return nil
}

// Current returns the current rate based on the moving average.
// If the MovingRate object has been stopped, return the global rate instead.
func (r *MovingRate) Current() (float64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If the counter hasn't been started, return an error
	if r.startTime.IsZero() {
		return 0, ErrNotStarted
	}

	// If the counter has been stopped, calculate the global rate
	if !r.stopTime.IsZero() {
		delta := r.stopTime.Sub(r.startTime).Seconds()
		return r.totalCounter / delta, nil
	}

	// If we don't have data yet, calculate the global rate
	// using the current time
	if len(r.movingAvgSamples) == 0 {
		delta := r.since(r.startTime).Seconds()
		return r.totalCounter / delta, nil
	}

	// Calculate the moving average
	var total float64
	for _, rate := range r.movingAvgSamples {
		total += rate
	}
	total = total / float64(len(r.movingAvgSamples))

	return total, nil
}
