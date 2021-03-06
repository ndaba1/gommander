package gommander

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	if p.rootCmd.settings[AllowNegativeNumbers] {
		if _, e := strconv.Atoi(val); e == nil {
			return false
		}
	}
	return strings.HasPrefix(val, "-")
}

func (p *Parser) isLongOptSyntax(val string) bool {
	return strings.HasPrefix(val, "--") && strings.ContainsAny(val, "=")
}

func (p *Parser) isSpecialValue(val string) bool {
	return val == "-h" || val == "--help"
}

func (p *Parser) isPosixFlagSyntax(vals []string) bool {
	return len(vals) > 2 && vals[0] == "-" && vals[1] != "-"
}

func (p *Parser) getFlag(val string) (*Flag, error) {
	for _, f := range p.currentCmd.flags {
		if f.ShortVal == val || f.LongVal == val {
			return f, nil
		}
	}
	return NewFlag(""), errors.New("flag not found")
}

func (p *Parser) getOption(val string) (*Option, error) {
	for _, o := range p.currentCmd.options {
		if o.ShortVal == val || o.LongVal == val {
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

func (p *Parser) reset() {
	p.eaten = []string{}
	p.cursor = 0
	p.cmdIdx = -1
}

func (p *Parser) parse(rawArgs []string) (*ParserMatches, *Error) {
	defer p.reset()

	p.matches.rawArgs = rawArgs
	p.matches.argCount = len(rawArgs)
	p.cmdIdx = -1

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
					if !p.matches.ContainsFlag(flag.LongVal) {
						p.matches.flagMatches = append(p.matches.flagMatches, flagCfg)
					}
				}
			} else if opt, err := p.getOption(arg); err == nil {
				// Handle is option
				p._eat(arg)
				err := p.parseOption(opt, rawArgs[(index+1):])
				if err != nil {
					return &p.matches, err
				}
			} else if arg == "--" {
				p._eat(arg)
				allowPositionalArgs = true
			} else if p.isLongOptSyntax(arg) && !allowPositionalArgs {
				// parse special option
				p._eat(arg)
				parts := strings.Split(arg, "=")

				opt, err := p.getOption(parts[0])
				if err != nil {
					err := generateError(p.currentCmd, UnknownOption, []string{parts[0], arg})
					return &p.matches, &err
				}

				temp := []string{parts[1]}
				temp = append(temp, rawArgs[(index+1):]...)

				e := p.parseOption(opt, temp)
				if e != nil {
					return &p.matches, e
				}
			} else if allowPositionalArgs {
				p._eat(arg)
				p.matches.positionalArgs = append(p.matches.positionalArgs, arg)
			} else if strings.ContainsRune(arg, '=') {
				err := generateError(p.currentCmd, UnknownOption, []string{})
				return &p.matches, &err
			} else if !p._isEaten(arg) && !allowPositionalArgs {
				values := strings.Split(arg, "")

				// TODO: More validation
				if p.isPosixFlagSyntax(values) {
					p._eat(arg)
					for _, v := range values[1:] {
						flag, err := p.getFlag(fmt.Sprintf("-%v", v))

						if err != nil {
							err := generateError(p.currentCmd, UnknownOption, []string{v, p.currentToken, ""})
							return &p.matches, &err
						}

						flagCfg := flagMatches{
							matchedFlag: *flag,
						}
						if !p.matches.ContainsFlag(flag.LongVal) {
							p.matches.flagMatches = append(p.matches.flagMatches, flagCfg)
						}
					}
					continue
				}

				err := generateError(p.currentCmd, UnknownOption, []string{p.currentToken})
				return &p.matches, &err
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
	if p.cmdIdx == -1 {
		// No subcommands matched
		cmdArgs = rawArgs
	} else if len(rawArgs) > p.cmdIdx+1 {
		cmdArgs = append(cmdArgs, rawArgs[p.cmdIdx+1:]...)
	}

	err := p.parseCmd(cmdArgs)
	if err != nil {
		return &p.matches, err
	}

	if !p.matches.ContainsFlag("help") {
		for _, o := range p.currentCmd.options {
			if o.IsRequired && !p.matches.ContainsOption(o.LongVal) {
				var argVals []string
				if o.Arg != nil {
					a := o.Arg
					if len(a.DefaultValue) == 0 {
						// No default value and value is required
						err := generateError(p.currentCmd, MissingRequiredOption, []string{o.LongVal})
						return &p.matches, &err
					}
					// Generate opt match with default value
					argVals = append(argVals, a.DefaultValue)
				}

				err := p.parseOption(o, argVals)
				if err != nil {
					return &p.matches, err
				}
			}
		}

	}

	return &p.matches, nil
}

func (p *Parser) parseOption(opt *Option, rawArgs []string) *Error {
	argList := []*Argument{}
	if opt.Arg != nil {
		argList = append(argList, opt.Arg)
	}

	args, err := p.getArgMatches(argList, rawArgs)
	if err != nil {
		return err
	}

	if p.matches.ContainsOption(opt.LongVal) {
		for i, cfg := range p.matches.optionMatches {
			if cfg.matchedOpt.LongVal == opt.LongVal {
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

	return nil
}

func (p *Parser) parseCmd(rawArgs []string) *Error {
	argCfgVals, err := p.getArgMatches(p.currentCmd.arguments, rawArgs)
	if err != nil {
		return err
	}

	// expected no args, probably a subcommand
	if len(rawArgs) > 0 && len(argCfgVals) == 0 {
		if p.currentCmd.hasSubcommands() && !p._isEaten(rawArgs[0]) {
			err := generateError(p.currentCmd, UnknownCommand, []string{rawArgs[0]})
			return &err
		}
	}

	// any unresolved arguments
	for _, a := range rawArgs {
		if !p._isEaten(a) {
			err := generateError(p.currentCmd, UnresolvedArgument, []string{a})
			return &err
		}
	}

	p.matches.argMatches = append(p.matches.argMatches, argCfgVals...)
	return nil
}

func (p *Parser) getArgMatches(list []*Argument, args []string) ([]argMatches, *Error) {
	// maxLen := len(list)
	matches := []argMatches{}

	for argIdx, argVal := range list {
		var builder strings.Builder

		if argVal.IsVariadic {
			for _, v := range args {
				if !p.isFlagLike(v) && !p._isEaten(v) {
					p._eat(v)
					builder.WriteString(v)
					builder.WriteRune(' ')
				}
			}
		} else if argIdx < len(args) {
			v := args[argIdx]

			if p.isSpecialValue(v) {
				break
			} else if !p.isFlagLike(v) && !p._isEaten(v) {
				p._eat(v)
				builder.WriteString(v)
			} else if argVal.hasDefaultValue() {
				builder.WriteString(argVal.DefaultValue)
			} else if argVal.IsRequired {
				args := []string{argVal.getRawValue(), v}
				err := generateError(p.currentCmd, MissingRequiredArgument, args)

				return matches, &err
			} else {
				continue
			}
		} else if argVal.IsRequired {
			args := []string{argVal.getRawValue()}
			err := generateError(p.currentCmd, MissingRequiredArgument, args)

			return matches, &err
		}

		// test the value against default values if any
		input := builder.String()
		if len(input) > 0 && len(argVal.ValidValues) > 0 && !argVal.testValue(input) {
			args := []string{input}
			args = append(args, argVal.ValidValues...)
			err := generateError(p.currentCmd, InvalidArgumentValue, args)

			return matches, &err
		}

		// test the value against the validator func if any
		for _, fn := range argVal.ValidatorFns {
			if err := fn(input); err != nil {
				args := []string{input, err.Error()}
				err := generateError(p.currentCmd, InvalidArgumentValue, args)

				return matches, &err
			}
		}

		// test against validator regex if any
		if argVal.ValidatorRe != nil && !argVal.ValidatorRe.MatchString(input) {
			args := []string{input, "failed to match value against validator regex"}
			err := generateError(p.currentCmd, InvalidArgumentValue, args)

			return matches, &err
		}

		argCfg := argMatches{
			rawValue:   input,
			instanceOf: *argVal,
		}

		matches = append(matches, argCfg)
	}

	return matches, nil
}
