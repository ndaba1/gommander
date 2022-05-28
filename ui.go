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
	Error
	Other
)

const (
	Colorful PredefinedTheme = iota
	Plain
)

var DESIGNATION_SLICE = []Designation{
	Keyword,
	Headline,
	Description,
	Error,
	Other,
}

type Formatter struct {
	theme       Theme
	buffer      bytes.Buffer
	prev_offset int
}

type Theme = map[Designation]color.Color

type FormatGenerator interface {
	generate() (string, string)
}

func NewFormatter(theme Theme) Formatter {
	return Formatter{
		theme: theme,
	}
}

func GetPredefinedTheme(val PredefinedTheme) Theme {
	switch val {
	case Colorful:
		return NewTheme(color.FgGreen, color.FgMagenta, color.FgHiBlue, color.FgHiRed, color.FgHiWhite)
	case Plain:
		return NewTheme(color.FgWhite, color.FgWhite, color.FgWhite, color.FgWhite, color.FgWhite)
	default:
		return DefaultTheme()
	}
}

// Creates as many colors as there are designations
func NewVariadicTheme(values ...color.Attribute) Theme {
	theme := make(map[Designation]color.Color)
	for i, v := range DESIGNATION_SLICE {
		theme[v] = *color.New(values[i])
	}

	return theme
}

func NewTheme(keyword, headline, description, errors, others color.Attribute) Theme {
	theme := make(map[Designation]color.Color)

	kw_color := color.New(keyword)
	hd_color := color.New(headline)
	ds_color := color.New(description)
	er_color := color.New(errors)
	ot_color := color.New(others)

	theme[Keyword] = *kw_color
	theme[Headline] = *hd_color
	theme[Description] = *ds_color
	theme[Error] = *er_color
	theme[Other] = *ot_color

	return theme
}

func DefaultTheme() Theme {
	return NewTheme(color.FgCyan, color.FgGreen, color.FgWhite, color.FgRed, color.FgWhite)
}

func (f *Formatter) section(val string) {
	f.add(Headline, fmt.Sprintf("\n%v: \n", val))
}

func (f *Formatter) close() {
	f.add(Other, "\n")
}

func (f *Formatter) add(dsgn Designation, val string) {
	c := f.theme[dsgn]

	colored_val := c.Sprintf(val)
	f.buffer.WriteString(colored_val)
}

func (f *Formatter) print() {
	color.New().Printf(f.buffer.String())
}

func (f *Formatter) format(items []FormatGenerator) {
	values := []([2]string){}

	// TODO: check for sort alphabetically setting
	sort.Slice(items, func(i, j int) bool {
		second, _ := items[j].generate()
		first, _ := items[i].generate()

		return second > first
	})

	for _, i := range items {
		leading, floating := i.generate()
		temp := [2]string{leading, floating}
		values = append(values, temp)
	}

	max_offset := 0
	current_offset := 0

	// Finds the longest value, adds some padding to it and sets it as the max offset
	for _, v := range values {
		capacity := len([]byte(v[0]))
		if capacity > current_offset {
			current_offset = capacity + 8 // Padding
		}
		if capacity > f.prev_offset {
			f.prev_offset = current_offset
		}
	}

	// If different sections have almost similar max_offsets, use equal values
	diff := f.prev_offset - current_offset
	if diff < 8 && diff > 0 {
		max_offset = f.prev_offset
	} else {
		max_offset = current_offset
	}

	for _, v := range values {
		leading := v[0]
		floating := v[1]

		f.print_output(leading, floating, max_offset)
	}

}

func (f *Formatter) print_output(leading string, floating string, offset int) {
	// TODO: Add support for sentence wrap
	buffer := make([]byte, offset)
	reader := strings.NewReader(leading)
	var temp_str strings.Builder

	numBytes, _ := reader.Read(buffer)
	temp_str.Write(buffer[:numBytes])
	diff := len(buffer) - numBytes

	for i := 0; i < diff; i++ {
		temp_str.Write([]byte(" "))
	}

	f.add(Keyword, fmt.Sprintf("    %v", temp_str.String()))
	f.add(Description, fmt.Sprintf("%v\n", floating))
}
