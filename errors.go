package gommander

import (
	"fmt"
	"io"
	"strings"
)

type GommanderError struct {
	kind    Event
	message string
	// args        []string
	context     string
	exit_code   int
	matched_cmd *Command
	is_nil      bool
}

func new_error(event Event) *GommanderError {
	return &GommanderError{
		kind:   event,
		is_nil: false,
	}
}

func nil_error() GommanderError {
	return GommanderError{is_nil: true}
}

func (e *GommanderError) msg(val string) *GommanderError {
	e.message = val
	return e
}

// func (e *GommanderError) values(vals []string) *GommanderError {
// 	e.args = vals
// 	return e
// }

func (e *GommanderError) ctx(val string) *GommanderError {
	e.context = val
	return e
}

func (e *GommanderError) exit(val int) *GommanderError {
	e.exit_code = val
	return e
}

func (e *GommanderError) cmd_ref(val *Command) *GommanderError {
	e.matched_cmd = val
	return e
}

func (e *GommanderError) Error() string {
	return e.message
}

func (e *GommanderError) Display() {
	fmter := NewFormatter()

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

func (e *GommanderError) into_event_cfg() EventConfig {
	return EventConfig{
		event:     e.kind,
		app_ref:   e.matched_cmd,
		exit_code: e.exit_code,
		err:       *e,
	}
}
