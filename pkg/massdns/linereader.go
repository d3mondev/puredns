package massdns

import (
	"bufio"
	"errors"
	"io"
	"math"
	"sync/atomic"
	"time"
)

// ErrNotStarted is an error happening when the LineReader hasn't been started.
var ErrNotStarted = errors.New("not started")

// LineReader is a line reader that limits the number of line per second read.
type LineReader struct {
	now   func() time.Time
	since func(time.Time) time.Duration

	reader       io.Reader
	readerBuffer *bufio.Reader

	rate      float64
	startTime time.Time
	lineCount int32
}

var _ io.Reader = (*LineReader)(nil)

// NewLineReader creates a new RateLimitLineReader.
func NewLineReader(r io.Reader, rate int) *LineReader {
	readerBuffer := bufio.NewReader(r)

	return &LineReader{
		now:   time.Now,
		since: time.Since,

		reader:       r,
		readerBuffer: readerBuffer,
		rate:         float64(rate),
	}
}

// Read reads from the reader, counting the number of lines read and applying rate limiting.
func (r *LineReader) Read(p []byte) (n int, err error) {
	const nl = byte('\n')
	var lines int

	canSend := r.canSend()

	for n = 0; n < len(p) && canSend > 0; n++ {
		var b byte
		if b, err = r.readerBuffer.ReadByte(); err != nil {
			break
		}

		if b == nl {
			lines++
			canSend--
		}

		p[n] = b
	}

	if r.rate > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	atomic.AddInt32(&r.lineCount, int32(lines))

	return n, err
}

// Count returns the number of lines read.
func (r *LineReader) Count() int {
	return int(atomic.LoadInt32(&r.lineCount))
}

// canSend calculates the number of lines that can be sent while respecting the rate limit.
func (r *LineReader) canSend() int {
	var canSend int
	if r.rate == 0 {
		canSend = math.MaxInt32
	} else {
		if r.startTime.IsZero() {
			r.startTime = r.now()
		}

		delta := r.since(r.startTime)
		canSend = int(r.rate*(delta.Seconds()+1)) - int(r.lineCount)
	}

	return canSend
}
