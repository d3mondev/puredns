package progressbar

import (
	"io"
	"time"
)

// Option configures a progress bar.
type Option interface {
	apply(*options)
}

type options struct {
	updateInterval time.Duration
	template       string
	writer         io.Writer
	style          Style
}

// WithTemplate provides a template string to the progress bar.
func WithTemplate(t string) Option {
	return templateOption(t)
}

type templateOption string

func (t templateOption) apply(opts *options) {
	opts.template = string(t)
}

// WithWriter provides a custom writer to the progress bar.
func WithWriter(w io.Writer) Option {
	return writerOption{w: w}
}

type writerOption struct {
	w io.Writer
}

func (w writerOption) apply(opts *options) {
	opts.writer = w.w
}

// WithInterval provides an update interval for the progress bar.
func WithInterval(d time.Duration) Option {
	return intervalOption(d)
}

type intervalOption time.Duration

func (i intervalOption) apply(opts *options) {
	opts.updateInterval = time.Duration(i)
}

// WithStyle provides a custom styling for the progress bar.
func WithStyle(style Style) Option {
	return styleOption(style)
}

type styleOption Style

func (s styleOption) apply(opts *options) {
	opts.style = Style(s)
}
