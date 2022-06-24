package gommander

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Make values to be more explicit, i.e. positional arg matches, matched_cmd_args etc.
type ParserMatches struct {
	argCount       int
	rawArgs        []string
	positionalArgs []string
	matchedCmd     *Command
	matchedCmdIdx  int
	rootCmd        *Command
	flagMatches    []flagMatches
	optionMatches  []optionMatches
	argMatches     []argMatches
}

type flagMatches struct {
	matchedFlag Flag
	// cursor_index int
}

type optionMatches struct {
	matchedOpt    Option
	instanceCount int
	passedArgs    []argMatches
	// cursor_index   int
}

type argMatches struct {
	rawValue   string
	instanceOf Argument
	// cursor_index int
}

func (pm *ParserMatches) GetRawArgCount() int {
	return pm.argCount
}

func (pm *ParserMatches) GetRawArgs() []string {
	return pm.rawArgs
}

func (pm *ParserMatches) GetPositionalArgs() []string {
	return pm.positionalArgs
}

func (pm *ParserMatches) GetAppRef() *Command {
	return pm.rootCmd
}

func (pm *ParserMatches) GetMatchedCommand() *Command {
	return pm.matchedCmd
}

func (pm *ParserMatches) GetMatchedCommandIndex() int {
	return pm.matchedCmdIdx
}

// Returns whether or not a flag was passed to the program args
func (pm *ParserMatches) ContainsFlag(val string) bool {
	for _, v := range pm.flagMatches {
		flag := v.matchedFlag
		if flag.short == val || flag.long == val || flag.name == val {
			return true
		}
	}
	return false
}

// Returns whether or not an option was passed to the program args
func (pm *ParserMatches) ContainsOption(val string) bool {
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.short == val || opt.long == val || opt.name == val {
			return true
		}
	}
	return false
}

// A method used to get the value of an argument passed to the program. Takes as input the name of the argument or the raw value of the argument. If no value is found, or the argument is misspelled, an error is returned. If no value was passed to the argument but it is required, the default value is used if one exists, otherwise an error is thrown.
func (pm *ParserMatches) GetArgValue(val string) (string, error) {
	for _, v := range pm.argMatches {
		arg := v.instanceOf
		if arg.name == val || arg.getRawValue() == val {
			return v.rawValue, nil
		}
	}

	return "", errors.New("no value found for provided argument")
}

func (pm *ParserMatches) GetOptionValue(val string) (string, error) {
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.short == val || opt.long == val || opt.name == val {
			// TODO: Probably check if slice is empty
			return v.passedArgs[0].rawValue, nil
		}
	}
	return "", errors.New("no value found for the provided option")
}

func (pm *ParserMatches) GetAllOptionInstances(val string) []string {
	instances := []string{}
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.short == val || opt.long == val || opt.name == val {
			for _, a := range v.passedArgs {
				instances = append(instances, a.rawValue)
			}
		}
	}
	return instances
}

/*************************** Parser Functionality ************************/

type Parser struct {
	cursor       int
	rootCmd      *Command
	currentCmd   *Command
	matches      ParserMatches
	eaten        []string
	cmdIdx       int
	currentToken string
}

func NewParser(entry *Command) Parser {
	return Parser{
		cursor:     0,
		rootCmd:    entry,
		currentCmd: entry,
		matches: ParserMatches{
			argCount:   0,
			rootCmd:    entry,
			matchedCmd: entry,
		},
	}
}

// Parser utilties
func (p *Parser) isFlagLike(val string) bool {
	return strings.HasPrefix(val, "-")
}

func (p *Parser) isSpecialOption(val string) bool {
	return strings.HasPrefix(val, "--") && strings.ContainsAny(val, "=")
}

func (p *Parser) isSpecialValue(val string) bool {
	return val == "-h" || val == "--help"
}

func (p *Parser) getFlag(val string) (*Flag, error) {
	for _, f := range p.currentCmd.flags {
		if f.short == val || f.long == val {
			return f, nil
		}
	}
	return NewFlag(""), errors.New("flag not found")
}

func (p *Parser) getOption(val string) (*Option, error) {
	for _, o := range p.currentCmd.options {
		if o.short == val || o.long == val {
			return o, nil
		}
	}
	return NewOption(""), errors.New("no option found")
}

func (p *Parser) getSubCommand(val string) (*Command, error) {
	for _, s := range p.currentCmd.subCommands {
		includes := func(val string) bool {
			for _, v := range s.aliases {
				if v == val {
					return true
				}
			}
			return false
		}

		if s.name == val || includes(val) {
			return s, nil
		}
	}
	return NewCommand(""), errors.New("no subcmd found")
}

func (p *Parser) _eat(val string) {
	p.eaten = append(p.eaten, val)
}

func (p *Parser) _isEaten(val string) bool {
	for _, v := range p.eaten {
		if v == val {
			return true
		}
	}

	return false
}

func (p *Parser) parse(rawArgs []string) (ParserMatches, Error) {
	p.matches.rawArgs = rawArgs
	p.matches.argCount = len(rawArgs)

	allowPositionalArgs := false

	for index, arg := range rawArgs {
		p.cursor = index
		p.currentToken = arg

		if p.isFlagLike(arg) {
			if flag, err := p.getFlag(arg); err == nil {
				// handle is flag
				p._eat(arg)
				if !allowPositionalArgs {
					flagCfg := flagMatches{
						matchedFlag: *flag,
					}
					if !p.matches.ContainsFlag(flag.long) {
						p.matches.flagMatches = append(p.matches.flagMatches, flagCfg)
					}
				}
			} else if opt, err := p.getOption(arg); err == nil {
				// Handle is option
				p._eat(arg)
				err := p.parseOption(opt, rawArgs[(index+1):])
				if !err.isNil {
					return p.matches, err
				}
			} else if arg == "--" {
				p._eat(arg)
				allowPositionalArgs = true
			} else if p.isSpecialOption(arg) && !allowPositionalArgs {
				// parse special option
				p._eat(arg)
				parts := strings.Split(arg, "=")

				opt, err := p.getOption(parts[0])
				if err != nil {
					msg := fmt.Sprintf("failed to resolve option: %v in value: %v", parts[0], arg)
					ctx := fmt.Sprintf("Found value: %v, with long option syntax but the option: %v is not valid in this context", arg, parts[0])

					return p.matches, throwError(UnresolvedArgument, msg, ctx).setArgs([]string{parts[0]})
				}

				temp := []string{parts[1]}
				temp = append(temp, rawArgs[(index+1):]...)

				e := p.parseOption(opt, temp)
				if !e.isNil {
					return p.matches, e
				}
			} else if allowPositionalArgs {
				p._eat(arg)
				p.matches.positionalArgs = append(p.matches.positionalArgs, arg)
			} else if !p._isEaten(arg) && !allowPositionalArgs {
				values := strings.Split(arg, "")

				// TODO: More validation
				if len(values) > 2 && values[0] == "-" {
					p._eat(arg)
					for _, v := range values[1:] {
						flag, err := p.getFlag(fmt.Sprintf("-%v", v))

						if err != nil {
							msg := fmt.Sprintf("unknown shorthand flag: `%v` in: `%v`", v, p.currentToken)
							ctx := fmt.Sprintf("Expected to find valid flags values in: `%v`, but instead found: `-%v` , which could not be resolved as a flag", p.currentToken, v)

							return p.matches, throwError(UnknownOption, msg, ctx).setArgs([]string{v, p.currentToken})
						}

						flagCfg := flagMatches{
							matchedFlag: *flag,
						}
						if !p.matches.ContainsFlag(flag.long) {
							p.matches.flagMatches = append(p.matches.flagMatches, flagCfg)
						}
					}
					continue
				}

				fmt.Printf("%v", strings.Split(arg, ""))
				msg := fmt.Sprintf("found unknown flag or option: `%v`", p.currentToken)
				ctx := fmt.Sprintf("The value: `%v`, could not be resolved as a flag or option.", p.currentToken)

				return p.matches, throwError(UnresolvedArgument, msg, ctx).setArgs([]string{p.currentToken})
			}
		} else if sc, err := p.getSubCommand(arg); err == nil {
			// handle subcmd
			p._eat(arg)
			p.currentCmd = sc
			p.cmdIdx = index

			continue
		} else if allowPositionalArgs {
			// TODO: More conditionals
			p._eat(arg)
			p.matches.positionalArgs = append(p.matches.positionalArgs, arg)
		}
	}

	p.matches.matchedCmd = p.currentCmd
	p.matches.matchedCmdIdx = p.cmdIdx

	cmdArgs := []string{}
	if len(rawArgs) > p.cmdIdx+1 {
		cmdArgs = append(cmdArgs, rawArgs[p.cmdIdx+1:]...)
	}

	err := p.parseCmd(cmdArgs)
	if !err.isNil {
		return p.matches, err
	}

	if !p.matches.ContainsFlag("help") {
		for _, o := range p.currentCmd.options {
			if o.required && !p.matches.ContainsOption(o.long) {
				var argVals []string
				for _, a := range o.args {
					if len(a.defaultValue) == 0 {
						// No default value and value is required
						msg := fmt.Sprintf("missing required option: `%v`", o.long)
						ctx := fmt.Sprintf("The option: `%v` is marked as required but no value was provided and it is not configured with a default value", o.long)

						return p.matches, throwError(MissingRequiredOption, msg, ctx).setArgs([]string{o.name})
					}
					// Generate opt match with default value
					argVals = append(argVals, a.defaultValue)
				}

				err := p.parseOption(o, argVals)
				if !err.isNil {
					return p.matches, err
				}
			}
		}

	}

	return p.matches, nilError()
}

func (p *Parser) parseOption(opt *Option, rawArgs []string) Error {
	args, err := p.getArgMatches(opt.args, rawArgs)
	if !err.isNil {
		return err
	}

	if p.matches.ContainsOption(opt.long) {
		for i, cfg := range p.matches.optionMatches {
			if cfg.matchedOpt.long == opt.long {
				cfg.passedArgs = append(cfg.passedArgs, args...)
				cfg.instanceCount++

				p.matches.optionMatches[i] = cfg
			}
		}
	} else {
		optCfg := optionMatches{
			matchedOpt:    *opt,
			instanceCount: 1,
			passedArgs:    args,
		}

		p.matches.optionMatches = append(p.matches.optionMatches, optCfg)
	}

	return nilError()
}

func (p *Parser) parseCmd(rawArgs []string) Error {
	argCfgVals, err := p.getArgMatches(p.currentCmd.arguments, rawArgs)
	if !err.isNil {
		return err
	}

	if len(argCfgVals) > 0 {
		p.matches.argMatches = append(p.matches.argMatches, argCfgVals...)
	} else if len(rawArgs) > 0 {
		if len(p.currentCmd.subCommands) > 0 && !p._isEaten(p.currentToken) {
			msg := fmt.Sprintf("no such subcommand found: `%v`", p.currentToken)
			suggestions := suggestSubCmd(p.currentCmd, p.currentToken)

			var ctx strings.Builder
			ctx.WriteString(fmt.Sprintf("The value: `%v`, could not be resolved as a subcommand. ", p.currentToken))
			if len(suggestions) > 0 {
				ctx.WriteString("Did you mean ")

				for i, s := range suggestions {
					if i > 0 {
						ctx.WriteString("or ")
					}
					ctx.WriteString(fmt.Sprintf("`%v` ", s))
				}

				ctx.WriteString("?")
			}

			return throwError(UnknownCommand, msg, ctx.String()).setArgs([]string{p.currentToken})
		} else if !p._isEaten(p.currentToken) {
			msg := fmt.Sprintf("failed to resolve argument: `%v`", p.currentToken)
			ctx := fmt.Sprintf("Found value: `%v`, which was unexpected or is invalid in this context", p.currentToken)

			return throwError(UnresolvedArgument, msg, ctx).setArgs([]string{p.currentToken})
		}
	}

	return nilError()
}

func (p *Parser) getArgMatches(list []*Argument, args []string) ([]argMatches, Error) {
	maxLen := len(list)
	matches := []argMatches{}

	for argIdx, argVal := range list {
		var builder strings.Builder

		if argVal.isVariadic {
			for _, v := range args {
				if !p.isFlagLike(v) && !p._isEaten(v) {
					p._eat(v)
					builder.WriteString(v)
					builder.WriteRune(' ')
				}
			}
		} else {
			for i := argIdx; i < maxLen; i++ {
				if len(args) == 0 && argVal.isRequired {
					if !argVal.hasDefaultValue() {
						args := []string{argVal.getRawValue()}
						msg := fmt.Sprintf("missing required argument: `%v`", args[0])
						ctx := fmt.Sprintf("Expected a required value corresponding to: `%v` but none was provided", argVal.getRawValue())

						return matches, throwError(MissingRequiredArgument, msg, ctx).setArgs(args)
					}
					builder.WriteString(argVal.defaultValue)
				} else if i < len(args) {
					v := args[i]
					if p.isSpecialValue(v) {
						break
					} else if !p.isFlagLike(v) && !p._isEaten(v) {
						p._eat(v)
						builder.WriteString(v)
					} else if argVal.hasDefaultValue() {
						builder.WriteString(argVal.defaultValue)
					} else if argVal.isRequired {
						args := []string{argVal.getRawValue()}
						msg := fmt.Sprintf("missing required argument: `%v`", args[0])
						ctx := fmt.Sprintf("Expected a value for argument: `%v`, but instead found: `%v`", argVal.name, v)

						return matches, throwError(MissingRequiredArgument, msg, ctx).setArgs(args)
					} else {
						continue
					}
				}
			}
		}

		// test the value against default values if any
		input := builder.String()
		if len(input) > 0 && len(argVal.validValues) > 0 && !argVal.testValue(input) {
			args := []string{input}
			msg := fmt.Sprintf("the passed value: `%v`, is not a valid argument", args[0])
			ctx := fmt.Sprintf("Expected one of: `%v`, but instead found: `%v`, which is not a valid value", argVal.validValues, input)

			return matches, throwError(InvalidArgumentValue, msg, ctx).setArgs(args)
		}

		// test the value against the validator func if any
		if argVal.validatorFn != nil {
			if err := argVal.validatorFn(input); err != nil {
				args := []string{input}
				msg := fmt.Sprintf("the passed value: `%v`, is not a valid argument", args[0])
				ctx := fmt.Sprintf("The validator function threw the following error: \"%v\" when checking the value: `%v`", err.Error(), input)

				return matches, throwError(InvalidArgumentValue, msg, ctx).setArgs(args)
			}
		}

		argCfg := argMatches{
			rawValue:   input,
			instanceOf: *argVal,
		}

		matches = append(matches, argCfg)
	}

	return matches, nilError()
}
