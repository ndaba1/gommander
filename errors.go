package gommander

import (
	"fmt"
	"strings"
)

type Error struct {
	kind     Event
	message  string
	args     []string
	context  string
	exitCode int
}

func generateError(cmd *Command, e Event, args []string) Error {
	var msg string
	var ctx string
	var code int

	switch e {
	case MissingRequiredArgument:
		{
			code = 20
			msg = fmt.Sprintf("missing required argument: `%v`", args[0])

			if len(args) == 1 {
				ctx = fmt.Sprintf("Expected a value corresponding to required argument: `%v` but none was provided", args[0])
			} else {
				ctx = fmt.Sprintf("Expected a value for argument: `%v`, but instead found option: `%v`", args[0], args[1])
			}

		}
	case MissingRequiredOption:
		{
			code = 30
			msg = fmt.Sprintf("missing required option: `%v`", args[0])
			ctx = fmt.Sprintf("The option: `%v` is marked as required but no value was provided and it is not configured with a default value", args[0])
		}
	case InvalidArgumentValue:
		{
			code = 10
			msg = fmt.Sprintf("the passed value: `%v`, is not a valid argument", args[0])

			switch len(args) {
			case 2:
				ctx = fmt.Sprintf("%v. Encountered this error when validating the value: `%v`", args[1], args[0])
			default:
				ctx = fmt.Sprintf("Expected one of: `[%v]`, but instead found: `%v`, which is not a valid value", strings.Join(args[1:], ", "), args[0])
			}
		}
	case UnknownOption:
		{
			code = 50

			switch len(args) {
			case 1:
				{
					msg = fmt.Sprintf("found unknown flag or option: `%v`", args[0])
					ctx = fmt.Sprintf("The value: `%v`, could not be resolved as a flag or option.", args[0])
				}
			case 2:
				{
					msg = fmt.Sprintf("failed to resolve option: %v in value: %v", args[0], args[1])
					ctx = fmt.Sprintf("Found value: %v, with long option syntax but the option: %v is not valid in this context", args[1], args[0])
				}
			case 3:
				{
					msg = fmt.Sprintf("unknown shorthand flag: `%v` in: `%v`", args[0], args[1])
					ctx = fmt.Sprintf("Expected to find valid flags values in: `%v`, but instead found: `-%v` , which could not be resolved as a flag", args[1], args[0])
				}
			default:
				{
					msg = "syntax not supported"
					ctx = "Passing option arguments using the `=` sign is only supported when using long options i.e. --port=8000 is valid but -p=8000 is not"
				}
			}

		}
	case UnresolvedArgument:
		{
			code = 60
			msg = fmt.Sprintf("failed to resolve argument: `%v`", args[0])
			ctx = fmt.Sprintf("Found value: `%v`, which was unexpected or is invalid in this context", args[0])
		}
	case UnknownCommand:
		{
			code = 40
			msg = fmt.Sprintf("no such subcommand found: `%v`", args[0])
			suggestions := cmd.suggestSubCmd(args[0])

			var context strings.Builder
			context.WriteString(fmt.Sprintf("The value: `%v`, could not be resolved as a subcommand. ", args[0]))
			if len(suggestions) > 0 {
				context.WriteString("Did you mean ")

				for i, s := range suggestions {
					if i > 0 {
						context.WriteString("or ")
					}
					context.WriteString(fmt.Sprintf("`%v` ", s))
				}

				context.WriteString("?")
			}
			ctx = context.String()
		}

	}

	return Error{
		kind:     e,
		message:  msg,
		context:  ctx,
		args:     args,
		exitCode: code,
	}
}

func (e *Error) ErrorMsg() string {
	return e.message
}

func (e *Error) GetErrorString(c *Command) string {
	fmter := e._writeError(c)
	return fmter.GetString()
}

func (e *Error) Display(c *Command) {
	fmter := e._writeError(c)
	fmter.Print()
}

func (e *Error) _writeError(c *Command) *Formatter {
	app := c._getAppRef()
	fmter := NewFormatter(app)

	fmter.Add(ErrorMsg, "error:  ")
	fmter.Add(Other, strings.ToLower(e.message))
	fmter.Add(Other, "\n\n")

	ctx := fillContent(e.context, 50)
	fmter.Add(Description, indent(ctx, "    "))
	fmter.Add(Other, "\n\n")

	if app.settings[ShowHelpOnAllErrors] {
		c.PrintHelp()
		fmt.Println()
	}

	fmter.Add(Other, "Run a COMMAND with --help for detailed usage information")
	fmter.close()

	return &fmter
}
