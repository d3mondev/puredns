package console

import (
	"fmt"
	"io"
	"os"
)

const (
	colorMeta        = ColorBrightWhite
	colorMessage     = ColorCyan
	colorMessageText = ColorWhite
	colorSuccess     = ColorBrightGreen
	colorSuccessText = ColorWhite
	colorWarning     = ColorBrightYellow
	colorWarningText = ColorWhite
	colorError       = ColorRed
	colorErrorText   = ColorRed
)

var (
	// Output is a writer where the messages are sent.
	Output io.Writer = os.Stderr

	// ExitHandler is a function called to exit the process.
	ExitHandler = os.Exit
)

// Message displays an informative message.
func Message(format string, a ...interface{}) {
	message := fmt.Sprintf("%s[%s*%s]%s %s%s\n", colorMeta, colorMessage, colorMeta, colorMessageText, format, ColorReset)
	fmt.Fprintf(Output, message, a...)
}

// Success displays a success message.
func Success(format string, a ...interface{}) {
	message := fmt.Sprintf("%s[%s+%s]%s %s%s\n", colorMeta, colorSuccess, colorMeta, colorSuccessText, format, ColorReset)
	fmt.Fprintf(Output, message, a...)
}

// Warning displays a warning message.
func Warning(format string, a ...interface{}) {
	message := fmt.Sprintf("%s[%s!%s]%s %s%s\n", colorMeta, colorWarning, colorMeta, colorWarningText, format, ColorReset)
	fmt.Fprintf(Output, message, a...)
}

// Error displays an error message.
func Error(format string, a ...interface{}) {
	message := fmt.Sprintf("%s[%sX%s]%s %s%s\n", colorMeta, colorError, colorMeta, colorErrorText, format, ColorReset)
	fmt.Fprintf(Output, message, a...)
}

// Printf prints a message without any formatting.
func Printf(format string, a ...interface{}) {
	fmt.Fprintf(Output, format, a...)
}

// Fatal prints a fatal error message and exists the process.
func Fatal(format string, a ...interface{}) {
	fmt.Fprintf(Output, ColorRed+format+ColorReset, a...)
	ExitHandler(-1)
}
