package massdns

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func newWithClock(r io.Reader, rate int) (*LineReader, *stubClock) {
	clock := newStubClock()
	reader := NewLineReader(r, rate)
	reader.now = clock.Now
	reader.since = clock.Since

	return reader, clock
}

func TestLineReaderNew(t *testing.T) {
	r := NewLineReader(strings.NewReader("test"), 1)
	assert.NotNil(t, r)
}

func TestLineReaderRead_Unlimited(t *testing.T) {
	testString := "line1\nline2\nline3\n"
	r, _ := newWithClock(strings.NewReader(testString), 0)

	buf := make([]byte, 1)
	var got string
	var gotErr error

	for {
		var n int
		if n, gotErr = r.Read(buf); gotErr != nil {
			break
		}

		got = got + string(buf[:n])
	}

	assert.ErrorIs(t, gotErr, io.EOF)
	assert.Equal(t, testString, got)
	assert.Equal(t, 3, r.Count())
}

func TestLineReaderRead_Limited(t *testing.T) {
	r, clock := newWithClock(strings.NewReader("line1\nline2\n"), 1)

	buf := make([]byte, 4096)
	var got string

	// First read
	n, gotErr := r.Read(buf)
	got = got + string(buf[:n])
	assert.Nil(t, gotErr)
	assert.Equal(t, "line1\n", got, "should return 1 line")

	n, gotErr = r.Read(buf)
	assert.Nil(t, gotErr)
	assert.Equal(t, 0, n, "should not return bytes before advancing clock")

	// Second read
	clock.advance(time.Second)

	n, gotErr = r.Read(buf)
	got = got + string(buf[:n])
	assert.Nil(t, gotErr)
	assert.Equal(t, "line1\nline2\n", got)

	// EOF
	clock.advance(time.Second)

	n, gotErr = r.Read(buf)
	assert.ErrorIs(t, gotErr, io.EOF)
	assert.Equal(t, 0, n)
}
