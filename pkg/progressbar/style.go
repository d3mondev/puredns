package progressbar

// Color is a terminal escape string that defines a color.
type Color string

const (
	// ColorBlack is the black foreground color.
	ColorBlack Color = "\033[0;30m"

	// ColorGray is the gray foreground color.
	ColorGray Color = "\033[1;30m"

	// ColorRed is the red foreground color.
	ColorRed Color = "\033[0;31m"

	// ColorBrightRed is the bright red foreground color.
	ColorBrightRed Color = "\033[1;31m"

	// ColorGreen is the green foreground color.
	ColorGreen Color = "\033[0;32m"

	// ColorBrightGreen is the bright green foreground color.
	ColorBrightGreen Color = "\033[1;32m"

	// ColorYellow is the yellow foreground color.
	ColorYellow Color = "\033[0;33m"

	// ColorBrightYellow is the yellow foreground color.
	ColorBrightYellow Color = "\033[1;33m"

	// ColorBlue is the blue foreground color.
	ColorBlue Color = "\033[0;34m"

	// ColorBrightBlue is the bright blue foreground color.
	ColorBrightBlue Color = "\033[1;34m"

	// ColorMagenta is the magenta foreground color.
	ColorMagenta Color = "\033[0;35m"

	// ColorBrightMagenta is the bright magenta foreground color.
	ColorBrightMagenta Color = "\033[1;35m"

	// ColorCyan is the cyan foreground color.
	ColorCyan Color = "\033[0;36m"

	// ColorBrightCyan is the bright cyan foreground color.
	ColorBrightCyan Color = "\033[1;36m"

	// ColorWhite is the white foreground color.
	ColorWhite Color = "\033[0;37m"

	// ColorBrightWhite is the bright white foreground color.
	ColorBrightWhite Color = "\033[1;37m"

	// ColorReset is the code to reset all attributes.
	ColorReset Color = "\033[0m"
)

// Style defines a progress bar style.
type Style struct {
	BarPrefix rune
	BarSuffix rune
	BarFull   rune
	BarEmpty  rune

	BarPrefixColor Color
	BarSuffixColor Color
	BarFullColor   Color
	BarEmptyColor  Color
}

// DefaultStyle create a Style object with the default styling.
func DefaultStyle() Style {
	return Style{
		BarPrefix: '|',
		BarSuffix: '|',
		BarEmpty:  '░',
		BarFull:   '█',

		BarPrefixColor: ColorWhite,
		BarSuffixColor: ColorWhite,
		BarEmptyColor:  ColorGray,
		BarFullColor:   ColorWhite,
	}
}
