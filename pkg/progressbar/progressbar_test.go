package progressbar

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateFn(bar *ProgressBar) {}

func TestNew_Default(t *testing.T) {
	pb := New(updateFn, 100)
	assert.NotNil(t, pb)
}

func TestNew_Options(t *testing.T) {
	pb := New(
		updateFn,
		100,
		WithTemplate(""),
		WithWriter(io.Discard),
		WithInterval(100*time.Millisecond),
		WithStyle(DefaultStyle()),
	)
	assert.NotNil(t, pb)
}

func TestStart_OK(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
}

func TestStop(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
	pb.Stop()
}

func TestGetSet(t *testing.T) {
	pb := New(updateFn, 100)

	got := pb.Get("key")
	assert.Nil(t, got, "should not exist")

	pb.Set("key", "value")
	got = pb.Get("key")
	assert.Equal(t, "value", got, "should exist")
}

func TestIncrement(t *testing.T) {
	pb := New(updateFn, 100)

	got := pb.Current()
	assert.EqualValues(t, 0, got)

	pb.Increment(1)

	got = pb.Current()
	assert.EqualValues(t, 1, got)
}

func TestSetCurrent(t *testing.T) {
	pb := New(updateFn, 100)

	pb.SetCurrent(10)

	got := pb.Current()
	assert.EqualValues(t, 10, got)
}

func TestSetCurrent_Lower(t *testing.T) {
	pb := New(updateFn, 100)

	pb.SetCurrent(10)
	pb.SetCurrent(5)

	got := pb.Current()
	assert.EqualValues(t, 10, got)
}

func TestTotal(t *testing.T) {
	pb := New(updateFn, 100)
	got := pb.Total()
	assert.EqualValues(t, 100, got)
}

func TestRate_Initial(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()

	got := pb.Rate()
	assert.Equal(t, 0.0, got)
}

func TestRate_Increment(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))

	pb.Start()
	pb.Increment(1)

	got := pb.Rate()
	assert.Equal(t, 1.0, got)
}

func TestRate_SetCurrent(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))

	pb.Start()
	pb.SetCurrent(1)

	got := pb.Rate()
	assert.Equal(t, 1.0, got)
}

func TestETA_OK(t *testing.T) {
	pb := New(updateFn, 7777, WithWriter(io.Discard))
	pb.Start()
	pb.SetCurrent(1)

	require.Equal(t, 1.0, pb.Rate(), "rate should be 1/sec")

	gotH, gotM, gotS := pb.ETA()

	assert.Equal(t, 2, gotH)
	assert.Equal(t, 9, gotM)
	assert.Equal(t, 36, gotS)
}

func TestETA_NoTotal(t *testing.T) {
	pb := New(updateFn, 0, WithWriter(io.Discard))
	pb.Start()
	pb.SetCurrent(1)

	gotH, gotM, gotS := pb.ETA()

	assert.Equal(t, 0, gotH)
	assert.Equal(t, 0, gotM)
	assert.Equal(t, 0, gotS)
}

func TestETA_NoRate(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
	pb.SetCurrent(0)

	require.Equal(t, 0.0, pb.Rate(), "rate should be 0/sec")

	gotH, gotM, gotS := pb.ETA()

	assert.Equal(t, 99, gotH)
	assert.Equal(t, 59, gotM)
	assert.Equal(t, 59, gotS)
}

func TestETA_Done(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
	pb.SetCurrent(100)

	gotH, gotM, gotS := pb.ETA()

	assert.Equal(t, 0, gotH)
	assert.Equal(t, 0, gotM)
	assert.Equal(t, 0, gotS)
}

func TestTime_Initial(t *testing.T) {
	pb := New(updateFn, 100)

	gotH, gotM, gotS := pb.Time()

	assert.Equal(t, 0, gotH)
	assert.Equal(t, 0, gotM)
	assert.Equal(t, 0, gotS)
}

func TestTime_Running(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
	pb.startTime = time.Now().Add(-30 * time.Second)

	gotH, gotM, gotS := pb.Time()

	assert.Equal(t, 0, gotH)
	assert.Equal(t, 0, gotM)
	assert.Greater(t, gotS, 0)
}

func TestTime_Finished(t *testing.T) {
	pb := New(updateFn, 100, WithWriter(io.Discard))
	pb.Start()
	pb.startTime = time.Now().Add(-30 * time.Second)
	pb.Stop()

	gotH, gotM, gotS := pb.Time()

	assert.Equal(t, 0, gotH)
	assert.Equal(t, 0, gotM)
	assert.Greater(t, gotS, 0)
}
