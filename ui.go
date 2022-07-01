package gommander

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type Designation byte
type PredefinedTheme byte

const (
	Keyword Designation = iota
	Headline
	Description
	ErrorMsg
	Other
)

const (
	ColorfulTheme PredefinedTheme = iota
	PlainTheme
)

type Formatter struct {
	theme      Theme
	buffer     bytes.Buffer
	prevOffset int
	appRef     *Command
}

type Theme = map[Designation]color.Attribute

type FormatGenerator interface {
	generate(*Command) (string, string)
}

func NewFormatter(cmd *Command) Formatter {
	return Formatter{
		theme:  cmd.theme,
		appRef: cmd,
	}
}

func GetPredefinedTheme(val PredefinedTheme) Theme {
	switch val {
	case ColorfulTheme:
		return NewTheme(color.FgGreen, color.FgMagenta, color.FgHiBlue, color.FgHiRed, color.FgHiWhite)
	case PlainTheme:
		return NewTheme(color.FgWhite, color.FgWhite, color.FgWhite, color.FgWhite, color.FgWhite)
	default:
		return DefaultTheme()
	}
}

// A constructor function that takes in color attributes in a specific order and creates a new theme from the provided color attributes from the `fatih/color` package
func NewTheme(keyword, headline, description, errors, others color.Attribute) Theme {
	theme := make(map[Designation]color.Attribute)

	theme[Keyword] = keyword
	theme[Headline] = headline
	theme[Description] = description
	theme[ErrorMsg] = errors
	theme[Other] = others

	return theme
}

// A simple function that returns the default package-defined theme
func DefaultTheme() Theme {
	return NewTheme(color.FgCyan, color.FgGreen, color.FgWhite, color.FgRed, color.FgWhite)
}

func (f *Formatter) section(val string) {
	f.Add(Headline, fmt.Sprintf("\n%v: \n", strings.ToUpper(val)))
}

func (f *Formatter) discussion(val string) {
	// Dedent, indent then word-wrap
	text := dedent(val)
	// text = fillContent(text, 80)
	text = dedent(text)
	text = indent(text, "    ")

	f.Add(Description, text)
	f.close()
}

func (f *Formatter) close() {
	f.Add(Other, "\n")
}

func (f *Formatter) Add(dsgn Designation, val string) {
	c := color.New(f.theme[dsgn])

	coloredVal := c.Sprintf(val)
	f.buffer.WriteString(coloredVal)
}

func (f *Formatter) AddAndPrint(dsgn Designation, val string) {
	f.Add(dsgn, val)
	f.Print()
}

func (f *Formatter) Color(color color.Color, val string) {
	colored := color.Sprintf(val)
	f.buffer.WriteString(colored)
}

func (f *Formatter) ColorAndPrint(color color.Color, val string) {
	f.Color(color, val)
	f.Print()
}

func (f *Formatter) Print() {
	if isTestMode() {
		fmt.Print(color.New().Sprintf(f.buffer.String()))
	} else {
		color.New().Printf(f.buffer.String())
	}
}

func (f *Formatter) GetString() string {
	return color.New().Sprintf(f.buffer.String())
}

func (f *Formatter) format(items []FormatGenerator) {
	values := []([2]string){}

	// TODO: check for sort alphabetically setting
	sort.Slice(items, func(i, j int) bool {
		second, _ := items[j].generate(f.appRef)
		first, _ := items[i].generate(f.appRef)

		return second > first
	})

	for _, i := range items {
		leading, floating := i.generate(f.appRef)
		temp := [2]string{leading, floating}
		values = append(values, temp)
	}

	maxOffset := 0
	currentOffset := 0

	// Finds the longest value, adds some padding to it and sets it as the max offset
	for _, v := range values {
		capacity := len([]byte(v[0]))
		if capacity > currentOffset {
			currentOffset = capacity + 8 // Padding
		}
		if capacity > f.prevOffset {
			f.prevOffset = currentOffset
		}
	}

	// If different sections have almost similar max_offsets, use equal values
	diff := f.prevOffset - currentOffset
	if diff < 8 && diff > -8 {
		maxOffset = f.prevOffset
	} else {
		maxOffset = currentOffset
	}

	for _, v := range values {
		leading := v[0]
		floating := v[1]

		f.printOutput(leading, floating, maxOffset)
	}

}

func (f *Formatter) printOutput(leading string, floating string, offset int) {
	// TODO: Add support for sentence wrap
	buffer := make([]byte, offset)
	reader := strings.NewReader(leading)
	var tempStr strings.Builder

	numBytes, _ := reader.Read(buffer)
	tempStr.Write(buffer[:numBytes])
	diff := len(buffer) - numBytes

	for i := 0; i < diff; i++ {
		tempStr.Write([]byte(" "))
	}

	f.Add(Keyword, fmt.Sprintf("    %v", tempStr.String()))
	f.Add(Description, fmt.Sprintf("%v\n", floating))
}
