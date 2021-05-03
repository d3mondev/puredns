package progressbar

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// ProgressBar is an asynchronous progress bar that can perform polling updates.
type ProgressBar struct {
	template string
	updateCB Update
	style    Style

	updateInterval time.Duration
	writer         io.Writer

	doneChan     chan bool
	finishedChan chan bool

	vars    map[interface{}]interface{}
	current int64
	total   int64

	rate       *MovingRate
	startTime  time.Time
	finishTime time.Time
}

// Update is a function that updates the progress bar.
type Update func(bar *ProgressBar)

// New creates a new progress bar from a template string and update function.
func New(update Update, total int64, opts ...Option) *ProgressBar {
	// Default options
	config := options{
		updateInterval: 200 * time.Millisecond,
		template:       "[ETA {{ eta }}] {{ bar }} {{ current }}/{{ total }} {{ rate }}/s (time: {{ time }})",
		writer:         os.Stderr,
		style:          DefaultStyle(),
	}

	// Apply options
	for _, o := range opts {
		o.apply(&config)
	}

	return &ProgressBar{
		template: config.template,
		updateCB: update,
		style:    config.style,

		updateInterval: config.updateInterval,
		writer:         config.writer,

		vars:  make(map[interface{}]interface{}),
		total: total,

		rate: NewMovingRate(time.Second, 10),
	}
}

// Start starts a progress bar.
func (p *ProgressBar) Start() {
	// Start the bar
	p.rate.Start()
	p.startTime = time.Now()

	// Create the channels
	p.doneChan = make(chan bool)
	p.finishedChan = make(chan bool)

	// Main loop for the progress bar
	go func() {
		var last bool

		// Update the progress bar every interval
		ticker := time.NewTicker(p.updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.updateCB(p)
				p.update()
				fmt.Fprintf(p.writer, "\r%s", p.Render())

				if last {
					fmt.Fprintln(p.writer, "")
					p.finishedChan <- true
					return
				}
			case <-p.doneChan:
				last = true
			}
		}
	}()
}

// Stop signals the progress bar to perform one last update and waits for it to finish.
func (p *ProgressBar) Stop() {
	// Signal the progress bar to perform one last update
	p.doneChan <- true

	// Wait for the progress bar to terminate
	<-p.finishedChan

	// Stop the rate calculation
	p.rate.Stop()

	// Record the finish time if it hasn't already been set
	if p.finishTime.IsZero() {
		p.finishTime = time.Now()
	}
}

// Set sets or updates a key value pair on the progress bar.
func (p *ProgressBar) Set(key interface{}, value interface{}) {
	p.vars[key] = value
}

// Get returns a key value if it exists, otherwise it returns nil.
func (p *ProgressBar) Get(key interface{}) interface{} {
	if val, ok := p.vars[key]; ok {
		return val
	}

	return nil
}

// SetCurrent sets the current progress bar value.
func (p *ProgressBar) SetCurrent(current int64) {
	diff := current - p.current

	if diff <= 0 {
		return
	}

	p.Increment(diff)
}

// Increment increments the progress bar by the specified value.
func (p *ProgressBar) Increment(val int64) {
	p.rate.Sample(float64(val))

	p.current = p.current + val

	if p.current == p.total {
		if p.finishTime.IsZero() {
			p.finishTime = time.Now()
		}
	}
}

// Current returns the current progress bar value.
func (p *ProgressBar) Current() int64 {
	return p.current
}

// Total returns the progress bar total value.
func (p *ProgressBar) Total() int64 {
	return p.total
}

// Rate returns the current rate.
func (p *ProgressBar) Rate() float64 {
	rate, _ := p.rate.Current()
	return rate
}

// ETA returns the current estimated time before the task finishes.
func (p *ProgressBar) ETA() (hours, minutes, seconds int) {
	current := p.Current()
	total := p.Total()
	remaining := total - current
	rate := p.Rate()

	if total == 0 {
		return 0, 0, 0
	}

	if remaining <= 0 || current == total {
		return 0, 0, 0
	}

	if rate == 0.0 {
		return 99, 59, 59
	}

	left := float64(remaining) / rate

	return convertTime(left)
}

// Time returns the time since the progress bar has been started until either it is stopped,
// or the current counter reaches the total.
func (p *ProgressBar) Time() (hours, minutes, seconds int) {
	// If the progress bar hasn't been started, the time is zero
	if p.startTime.IsZero() {
		return 0, 0, 0
	}

	var duration time.Duration
	if p.finishTime.IsZero() {
		// Still running, take the time since start
		duration = time.Since(p.startTime)
	} else {
		// Finished, take the time since finish
		duration = p.finishTime.Sub(p.startTime)
	}

	return convertTime(duration.Seconds())
}

// Render returns a string containing the progress bar.
func (p *ProgressBar) Render() string {
	render := parseVariables(p.template, p.vars)

	return render
}

// update updates the default internal values
func (p *ProgressBar) update() {
	// Update time
	th, tm, ts := p.Time()
	p.Set("time", fmt.Sprintf("%02d:%02d:%02d", th, tm, ts))

	// Update rate
	rate := p.Rate()
	p.Set("rate", fmt.Sprintf("%.f", rate))

	// Update ETA
	eh, em, es := p.ETA()
	p.Set("eta", fmt.Sprintf("%02d:%02d:%02d", eh, em, es))

	// Update current
	current := p.Current()
	p.Set("current", fmt.Sprintf("%d", current))

	// Update total
	total := p.Total()
	p.Set("total", fmt.Sprintf("%d", total))

	// Update percentage
	perc := float64(current) / float64(total) * 100.0
	p.Set("percent", fmt.Sprintf("%.f", perc))

	// Update bar
	p.Set("bar", p.drawBar(perc))
}

// drawBar renders the progress bar
func (p *ProgressBar) drawBar(percent float64) string {
	width := 40
	barWidth := width - 2

	var builder strings.Builder
	builder.WriteString(string(p.style.BarPrefixColor))
	builder.WriteRune('|')
	builder.WriteString(string(p.style.BarFullColor))

	setEmptyColor := true
	for i := 0; i < barWidth; i++ {
		cur := float64(i) / float64(barWidth) * 100.0
		if cur < percent {
			builder.WriteRune('█')
		} else {
			if setEmptyColor {
				builder.WriteString(string(p.style.BarEmptyColor))
				setEmptyColor = false
			}

			builder.WriteRune('░')
		}
	}

	builder.WriteString(string(p.style.BarSuffixColor))
	builder.WriteRune('|')
	builder.WriteString(string(ColorReset))

	return builder.String()
}

var rxVariableMatch = regexp.MustCompile(`{{\s*([a-zA-Z0-9-_.]+)\s*}}`)

// parseVariables replaces {{ variables }} in a string by their value from the variable map
func parseVariables(str string, vars map[interface{}]interface{}) string {
	matches := rxVariableMatch.FindAllStringSubmatch(str, -1)

	for _, group := range matches {
		expression := group[0]
		key := group[1]

		if value, ok := vars[key]; ok {
			str = strings.Replace(str, expression, value.(string), 1)
		}
	}

	return str
}

// converTime takes a time in seconds and converts it to hours, minutes, seconds.
func convertTime(t float64) (hours, minutes, seconds int) {
	hours = int(t / 3600.0)
	t -= float64(hours * 3600)

	minutes = int(t / 60.0)
	t -= float64(minutes * 60)

	seconds = int(t)

	return hours, minutes, seconds
}
