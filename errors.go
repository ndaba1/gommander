package gommander

import (
	"fmt"
	"io"
	"strings"
)

type GommanderError struct {
	kind      Event
	message   string
	args      []string
	context   string
	exit_code int
	is_nil    bool
}

func nil_error() GommanderError {
	return GommanderError{is_nil: true}
}

func throw_error(kind Event, msg string, ctx string) GommanderError {
	var exit_code int
	switch kind {
	case InvalidArgumentValue:
		exit_code = 10
	case MissingRequiredArgument:
		exit_code = 20
	case MissingRequiredOption:
		exit_code = 30
	case UnknownCommand:
		exit_code = 40
	case UnknownOption:
		exit_code = 50
	default:
		exit_code = 1
	}
	return GommanderError{
		kind:      kind,
		message:   msg,
		context:   ctx,
		exit_code: exit_code,
	}
}

func (e GommanderError) set_args(vals []string) GommanderError {
	e.args = vals
	return e
}

func (e *GommanderError) Error() string {
	return e.message
}

func (e *GommanderError) Display(theme Theme) {
	fmter := NewFormatter(theme)

	fmter.add(Error, "error:  ")
	fmter.add(Other, strings.ToLower(e.message))
	fmter.close()
	fmter.close()

	reader := strings.NewReader(e.context)
	// values := strings.Split(e.context, " ")
	buffer := make([]byte, 50)

	// TODO: Find a better way to word wrap
	for {
		bytes, err := reader.Read(buffer)
		chunk := buffer[:bytes]

		fmter.add(Description, fmt.Sprintf("   %v\n", string(chunk)))
		if err == io.EOF {
			break
		}
	}

	fmter.add(Other, "Run a COMMAND with --help for detailed usage information")
	fmter.close()
	fmter.print()
}
