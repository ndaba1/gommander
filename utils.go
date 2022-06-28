package gommander

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	leadingNewline    = regexp.MustCompile("^[\n]")
	whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
	leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
)

func suggestSubCmd(c *Command, val string) []string {
	var minMatchSize = 3
	var matches []string

	cmdMap := make(map[string]int, 0)

	for _, v := range c.subCommands {
		cmdMap[v.name] = 0
	}

	for _, sc := range c.subCommands {
		for i, v := range strings.Split(val, "") {
			if len(sc.name) > i {
				var next string
				current := string(sc.name[i])

				if len(sc.name) > i+1 {
					next = string(sc.name[i+1])
				}

				if next == v || current == v {
					cmdMap[sc.name] += 1
				}
			}
		}
	}

	for k, v := range cmdMap {
		if v >= minMatchSize {
			matches = append(matches, k)
		}
	}

	return matches
}

/********************************** Text wrap and formatting **************************************/

func dedent(text string) string {
	var margin string

	text = leadingNewline.ReplaceAllString(text, "")
	text = whitespaceOnly.ReplaceAllString(text, "")
	indents := leadingWhitespace.FindAllStringSubmatch(text, -1)

	for i, indent := range indents {
		if i == 0 {
			margin = indent[1]
		} else if strings.HasPrefix(indent[1], margin) {
			continue
		} else if strings.HasPrefix(margin, indent[1]) {
			margin = indent[1]
		} else {
			margin = ""
			break
		}
	}

	if margin != "" {
		text = regexp.MustCompile("(?m)^"+margin).ReplaceAllString(text, "")
	}
	return text
}

func indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	prefixed := []string{}

	for _, line := range lines {
		var temp strings.Builder
		temp.WriteString(prefix)
		temp.WriteString(line)
		prefixed = append(prefixed, temp.String())
	}

	return strings.Join(prefixed, "\n")
}

func wrapContent(text string, width int) []string {
	buff := make([]string, 0)
	line := ""
	for _, word := range regexp.MustCompile(" ").Split(text, -1) {
		if len(line+word) < width {
			line += word + " "
		} else {
			line = strings.TrimSpace(line)
			if line != "" {
				buff = append(buff, strings.TrimSpace(line))
			}
			line = word + " "
		}
	}
	line = strings.TrimSpace(line)
	if line != "" {
		buff = append(buff, strings.TrimSpace(line))
	}
	return buff
}

func fillContent(text string, width int) string {
	return strings.Join(wrapContent(text, width), "\n")
}

/********************************** Testing and Debug utilities **************************************/

func isTestMode() bool {
	val, exists := os.LookupEnv("GOMMANDER_TEST_MODE")

	return exists && val == "true"
}

func setGommanderTestMode() {
	os.Setenv("GOMMANDER_TEST_MODE", "true")
}

func assert(t *testing.T, val interface{}, msg ...interface{}) {
	if val != true {
		if len(msg) > 0 {
			t.Error(msg...)
		} else {
			t.Error("Assertion failed. Value is not truthy. ")
		}
		t.Errorf("Expected: %v to be truthy", val)
	}
}

func assertEq(t *testing.T, first interface{}, second interface{}, msg ...interface{}) {
	if first != second {
		if len(msg) > 0 {
			t.Error(msg...)
		} else {
			t.Error("Assertion failed. Expected values to be equal. ")
		}
		t.Errorf("Left hand side is: %v", first)
		t.Errorf("Right hand side is: %v", second)
	}
}

func assertNe(t *testing.T, first interface{}, second interface{}, msg ...interface{}) {
	if first == second {
		if len(msg) > 0 {
			t.Error(msg...)
		} else {
			t.Error("Assertion failed. Did not expect values to be equal. ")
		}
		t.Errorf("Left hand side is: %v", first)
		t.Errorf("Right hand side is: %v", second)
	}
}

type structComp[model structTypes] interface {
	compare(model) bool
}

type structTypes interface {
	*Flag | *Command | *Option | *Argument
}

func assertStructEq[model structTypes](t *testing.T, first structComp[model], second model, msg ...interface{}) {
	if !first.compare(second) {
		if len(msg) > 0 {
			t.Error(msg...)
		} else {
			t.Error("Assertion failed. Struct values not equal. ")
		}
		t.Errorf("Left hand side is: %+v", first)
		t.Errorf("Right hand side is: %+v", second)
	}
}
