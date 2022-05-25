package gommander

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Designation byte

const (
	Keyword Designation = iota
	Headline
	Description
	Error
	Other
)

type Formatter struct {
	theme  Theme
	buffer bytes.Buffer
}

type Theme struct {
	values map[Designation]color.Color
}

type FormatGenerator interface {
	generate() (string, string)
}

func NewFormatter() Formatter {
	return Formatter{
		theme: DefaultTheme(),
	}
}

func DefaultTheme() Theme {
	theme := make(map[Designation]color.Color)

	kw_color := color.New(color.FgCyan)
	hd_color := color.New(color.FgGreen)
	ds_color := color.New(color.FgWhite)
	er_color := color.New(color.FgRed)
	ot_color := color.New(color.FgWhite)

	theme[Keyword] = *kw_color
	theme[Headline] = *hd_color
	theme[Description] = *ds_color
	theme[Error] = *er_color
	theme[Other] = *ot_color

	return Theme{values: theme}
}

func (f *Formatter) section(val string) {
	f.add(Headline, fmt.Sprintf("\n%v: \n", val))
}

func (f *Formatter) add(dsgn Designation, val string) {
	c := f.theme.values[dsgn]

	colored_val := c.Sprintf(val)
	f.buffer.WriteString(colored_val)
}

func (f *Formatter) print() {
	color.New().Printf(f.buffer.String())
}

func (f *Formatter) format(items []FormatGenerator) {
	values := []([2]string){}

	for _, i := range items {
		leading, floating := i.generate()
		temp := [2]string{leading, floating}
		values = append(values, temp)
	}

	max_offset := 0
	for _, v := range values {
		capacity := len([]byte(v[0]))
		if capacity > max_offset {
			max_offset = capacity + 5 // Padding
		}
	}

	for _, v := range values {
		leading := v[0]
		floating := v[1]

		f.print_output(leading, floating, max_offset)
	}
}

func (f *Formatter) print_output(leading string, floating string, offset int) {
	buffer := make([]byte, offset)
	reader := strings.NewReader(leading)

	numBytes, _ := reader.Read(buffer)
	diff := len(buffer) - numBytes

	for i := 0; i < diff; i++ {
		strings.NewReader(" ").Read(buffer)
	}

	f.add(Keyword, string(buffer))
	f.add(Description, floating)
}
