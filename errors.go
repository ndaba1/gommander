package gommander

import (
	"fmt"
	"io"
	"strings"
)

type Error struct {
	kind     Event
	message  string
	args     []string
	context  string
	exitCode int
	isNil    bool
}

func nilError() Error {
	return Error{isNil: true}
}

func throwError(kind Event, msg string, ctx string) Error {
	var exitCode int
	switch kind {
	case InvalidArgumentValue:
		exitCode = 10
	case MissingRequiredArgument:
		exitCode = 20
	case MissingRequiredOption:
		exitCode = 30
	case UnknownCommand:
		exitCode = 40
	case UnknownOption:
		exitCode = 50
	default:
		exitCode = 1
	}
	return Error{
		kind:     kind,
		message:  msg,
		context:  ctx,
		exitCode: exitCode,
	}
}

func (e Error) setArgs(vals []string) Error {
	e.args = vals
	return e
}

func (e *Error) compare(err *Error) bool {
	// TODO: Validate all fields
	return e.message == err.message && e.context == err.context && e.kind == err.kind
}

func (e *Error) ErrorMsg() string {
	return e.message
}

func (e *Error) Display(c *Command) {
	app := c._getAppRef()
	fmter := NewFormatter(app.theme)

	fmter.Add(ErrorMsg, "error:  ")
	fmter.Add(Other, strings.ToLower(e.message))
	fmter.close()
	fmter.close()

	reader := strings.NewReader(e.context)
	// values := strings.Split(e.context, " ")
	buffer := make([]byte, 50)

	// TODO: Find a better way to word wrap
	for {
		bytes, err := reader.Read(buffer)
		chunk := buffer[:bytes]

		fmter.Add(Description, fmt.Sprintf("   %v\n", string(chunk)))
		if err == io.EOF {
			break
		}
	}

	if app.settings[ShowHelpOnAllErrors] {
		c.PrintHelp()
		fmt.Println()
	}

	fmter.Add(Other, "Run a COMMAND with --help for detailed usage information")
	fmter.close()
	fmter.Print()
}
