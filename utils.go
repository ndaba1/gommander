package gommander

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var (
	leadingNewline    = regexp.MustCompile("^[\n]")
	whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
	leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
)

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

func _throwAssertionError(t *testing.T, errMsg string, first, second interface{}, msg ...interface{}) {
	if len(msg) > 0 {
		t.Error(msg...)
	} else {
		t.Errorf("Assertion failed. %s", strings.ToUpper(errMsg))
	}
	err := fmt.Sprintf("\n *********LEFT HAND SIDE IS*********:\n `%v` \n\n *********RIGHT HAND SIDE IS*********:\n `%v` \n", first, second)
	t.Error(err)
}

func assert(t *testing.T, val interface{}, msg ...interface{}) {
	if val != true {
		if len(msg) > 0 {
			t.Error(msg...)
		} else {
			t.Error("Assertion failed. Value is not truthy. ")
		}
		t.Errorf("Expected: `%v` to be truthy", val)
	}
}

func assertEq(t *testing.T, first, second interface{}, msg ...interface{}) {
	if first != second {
		_throwAssertionError(t, "Expected values to be equal. ", first, second, msg...)
	}
}

func assertDeepEq(t *testing.T, first, second interface{}, msg ...interface{}) {
	if !reflect.DeepEqual(first, second) {
		_throwAssertionError(t, "Expected values to be deeply equal. ", first, second, msg...)
	}
}

func assertNe(t *testing.T, first, second interface{}, msg ...interface{}) {
	if first == second {
		_throwAssertionError(t, "Did not expect values to be equal.", first, second, msg...)
	}
}

func assertStdOut(t *testing.T, expected string, exec func(), msg ...interface{}) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exec()
	_ = w.Close()
	res, _ := io.ReadAll(r)
	output := string(res)

	os.Stdout = stdOut

	// println(len(expected), "vs", len(output))
	if output != expected {
		_throwAssertionError(t, "Expected output was different from actual output", expected, output, msg...)
	}
}
