package gommander

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Make values to be more explicit, i.e. positional arg matches, matched_cmd_args etc.
type ParserMatches struct {
	arg_count       int
	raw_args        []string
	positional_args []string
	matched_cmd     *Command
	matched_cmd_idx int
	root_cmd        *Command
	flag_matches    []flag_matches
	option_matches  []option_matches
	arg_matches     []arg_matches
}

type flag_matches struct {
	matched_flag Flag
	// cursor_index int
}

type option_matches struct {
	matched_opt    Option
	instance_count int
	passed_args    []arg_matches
	// cursor_index   int
}

type arg_matches struct {
	raw_value   string
	instance_of Argument
	// cursor_index int
}

func (pm *ParserMatches) GetRawArgCount() int {
	return pm.arg_count
}

func (pm *ParserMatches) GetRawArgs() []string {
	return pm.raw_args
}

func (pm *ParserMatches) GetPositionalArgs() []string {
	return pm.positional_args
}

func (pm *ParserMatches) GetAppRef() *Command {
	return pm.root_cmd
}

func (pm *ParserMatches) GetMatchedCommand() *Command {
	return pm.matched_cmd
}

func (pm *ParserMatches) GetMatchedCommandIndex() int {
	return pm.matched_cmd_idx
}

// Returns whether or not a flag was passed to the program args
func (pm *ParserMatches) ContainsFlag(val string) bool {
	for _, v := range pm.flag_matches {
		flag := v.matched_flag
		if flag.short == val || flag.long == val || flag.name == val {
			return true
		}
	}
	return false
}

// Returns whether or not an option was passed to the program args
func (pm *ParserMatches) ContainsOption(val string) bool {
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			return true
		}
	}
	return false
}

// A method used to get the value of an argument passed to the program. Takes as input the name of the argument or the raw value of the argument. If no value is found, or the argument is misspelled, an error is returned. If no value was passed to the argument but it is required, the default value is used if one exists, otherwise an error is thrown.
func (pm *ParserMatches) GetArgValue(val string) (string, error) {
	for _, v := range pm.arg_matches {
		arg := v.instance_of
		if arg.name == val || arg.get_raw_value() == val {
			return v.raw_value, nil
		}
	}

	return "", errors.New("no value found for provided argument")
}

func (pm *ParserMatches) GetOptionValue(val string) (string, error) {
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			// TODO: Probably check if slice is empty
			return v.passed_args[0].raw_value, nil
		}
	}
	return "", errors.New("no value found for the provided option")
}

func (pm *ParserMatches) GetAllOptionInstances(val string) []string {
	instances := []string{}
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			for _, a := range v.passed_args {
				instances = append(instances, a.raw_value)
			}
		}
	}
	return instances
}

/*************************** Parser Functionality ************************/

type Parser struct {
	cursor        int
	root_cmd      *Command
	current_cmd   *Command
	matches       ParserMatches
	eaten         []string
	cmd_idx       int
	current_token string
}

func NewParser(entry *Command) Parser {
	return Parser{
		cursor:      0,
		root_cmd:    entry,
		current_cmd: entry,
		matches: ParserMatches{
			arg_count:   0,
			root_cmd:    entry,
			matched_cmd: entry,
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
	for _, f := range p.current_cmd.flags {
		if f.short == val || f.long == val {
			return f, nil
		}
	}
	return NewFlag(""), errors.New("flag not found")
}

func (p *Parser) getOption(val string) (*Option, error) {
	for _, o := range p.current_cmd.options {
		if o.short == val || o.long == val {
			return o, nil
		}
	}
	return NewOption(""), errors.New("no option found")
}

func (p *Parser) getSubCommand(val string) (*Command, error) {
	for _, s := range p.current_cmd.sub_commands {
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

func (p *Parser) parse(raw_args []string) (ParserMatches, GommanderError) {
	p.matches.raw_args = raw_args
	p.matches.arg_count = len(raw_args)

	allow_positional_args := false

	for index, arg := range raw_args {
		p.cursor = index
		p.current_token = arg

		if p.isFlagLike(arg) {
			if flag, err := p.getFlag(arg); err == nil {
				// handle is flag
				p._eat(arg)
				if !allow_positional_args {
					flag_cfg := flag_matches{
						matched_flag: *flag,
					}
					if !p.matches.ContainsFlag(flag.long) {
						p.matches.flag_matches = append(p.matches.flag_matches, flag_cfg)
					}
				}
			} else if opt, err := p.getOption(arg); err == nil {
				// Handle is option
				p._eat(arg)
				err := p.parse_option(opt, raw_args[(index+1):])
				if !err.is_nil {
					return p.matches, err
				}
			} else if arg == "--" {
				p._eat(arg)
				allow_positional_args = true
			} else if p.isSpecialOption(arg) && !allow_positional_args {
				// parse special option
				p._eat(arg)
				parts := strings.Split(arg, "=")

				opt, err := p.getOption(parts[0])
				if err != nil {
					msg := fmt.Sprintf("failed to resolve option: %v in value: %v", parts[0], arg)
					ctx := fmt.Sprintf("Found value: %v, with long option syntax but the option: %v is not valid in this context", arg, parts[0])

					return p.matches, throw_error(UnresolvedArgument, msg, ctx).set_args([]string{parts[0]})
				}

				temp := []string{parts[1]}
				temp = append(temp, raw_args[(index+1):]...)

				e := p.parse_option(opt, temp)
				if !e.is_nil {
					return p.matches, e
				}
			} else if allow_positional_args {
				p._eat(arg)
				p.matches.positional_args = append(p.matches.positional_args, arg)
			} else if !p._isEaten(arg) && !allow_positional_args {
				values := strings.Split(arg, "")

				// TODO: More validation
				if len(values) > 2 && values[0] == "-" {
					p._eat(arg)
					for _, v := range values[1:] {
						flag, err := p.getFlag(fmt.Sprintf("-%v", v))

						if err != nil {
							msg := fmt.Sprintf("unknown shorthand flag: `%v` in: `%v`", v, p.current_token)
							ctx := fmt.Sprintf("Expected to find valid flags values in: `%v`, but instead found: `-%v` , which could not be resolved as a flag", p.current_token, v)

							return p.matches, throw_error(UnknownOption, msg, ctx).set_args([]string{v, p.current_token})
						}

						flag_cfg := flag_matches{
							matched_flag: *flag,
						}
						if !p.matches.ContainsFlag(flag.long) {
							p.matches.flag_matches = append(p.matches.flag_matches, flag_cfg)
						}
					}
					continue
				}

				fmt.Printf("%v", strings.Split(arg, ""))
				msg := fmt.Sprintf("found unknown flag or option: `%v`", p.current_token)
				ctx := fmt.Sprintf("The value: `%v`, could not be resolved as a flag or option.", p.current_token)

				return p.matches, throw_error(UnresolvedArgument, msg, ctx).set_args([]string{p.current_token})
			}
		} else if sc, err := p.getSubCommand(arg); err == nil {
			// handle subcmd
			p._eat(arg)
			p.current_cmd = sc
			p.cmd_idx = index

			continue
		} else if allow_positional_args {
			// TODO: More conditionals
			p._eat(arg)
			p.matches.positional_args = append(p.matches.positional_args, arg)
		}
	}

	p.matches.matched_cmd = p.current_cmd
	p.matches.matched_cmd_idx = p.cmd_idx

	cmd_args := []string{}
	if len(raw_args) > p.cmd_idx+1 {
		cmd_args = append(cmd_args, raw_args[p.cmd_idx+1:]...)
	}

	err := p.parse_cmd(cmd_args)
	if !err.is_nil {
		return p.matches, err
	}

	if !p.matches.ContainsFlag("help") {
		for _, o := range p.current_cmd.options {
			if o.required && !p.matches.ContainsOption(o.long) {
				var arg_vals []string
				for _, a := range o.args {
					if len(a.default_value) == 0 {
						// No default value and value is required
						msg := fmt.Sprintf("missing required option: `%v`", o.long)
						ctx := fmt.Sprintf("The option: `%v` is marked as required but no value was provided and it is not configured with a default value", o.long)

						return p.matches, throw_error(MissingRequiredOption, msg, ctx).set_args([]string{o.name})
					} else {
						// Generate opt match with default value
						arg_vals = append(arg_vals, a.default_value)
					}
				}

				err := p.parse_option(o, arg_vals)
				if !err.is_nil {
					return p.matches, err
				}
			}
		}

	}

	return p.matches, nil_error()
}

func (p *Parser) parse_option(opt *Option, raw_args []string) GommanderError {
	args, err := p.get_arg_matches(opt.args, raw_args)
	if !err.is_nil {
		return err
	}

	if p.matches.ContainsOption(opt.long) {
		for i, cfg := range p.matches.option_matches {
			if cfg.matched_opt.long == opt.long {
				cfg.passed_args = append(cfg.passed_args, args...)
				cfg.instance_count += 1

				p.matches.option_matches[i] = cfg
			}
		}
	} else {
		opt_cfg := option_matches{
			matched_opt:    *opt,
			instance_count: 1,
			passed_args:    args,
		}

		p.matches.option_matches = append(p.matches.option_matches, opt_cfg)
	}

	return nil_error()
}

func (p *Parser) parse_cmd(raw_args []string) GommanderError {
	arg_cfg_vals, err := p.get_arg_matches(p.current_cmd.arguments, raw_args)
	if !err.is_nil {
		return err
	}

	if len(arg_cfg_vals) > 0 {
		p.matches.arg_matches = append(p.matches.arg_matches, arg_cfg_vals...)
	} else if len(raw_args) > 0 {
		if len(p.current_cmd.sub_commands) > 0 && !p._isEaten(p.current_token) {
			msg := fmt.Sprintf("no such subcommand found: `%v`", p.current_token)
			suggestions := suggest_sub_cmd(p.current_cmd, p.current_token)

			var ctx strings.Builder
			ctx.WriteString(fmt.Sprintf("The value: `%v`, could not be resolved as a subcommand. ", p.current_token))
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

			return throw_error(UnknownCommand, msg, ctx.String()).set_args([]string{p.current_token})
		} else if !p._isEaten(p.current_token) {
			msg := fmt.Sprintf("failed to resolve argument: `%v`", p.current_token)
			ctx := fmt.Sprintf("Found value: `%v`, which was unexpected or is invalid in this context", p.current_token)

			return throw_error(UnresolvedArgument, msg, ctx).set_args([]string{p.current_token})
		}
	}

	return nil_error()
}

func (p *Parser) get_arg_matches(list []*Argument, args []string) ([]arg_matches, GommanderError) {
	max_len := len(list)
	matches := []arg_matches{}

	for arg_idx, arg_val := range list {
		var builder strings.Builder

		if arg_val.is_variadic {
			for _, v := range args {
				if !p.isFlagLike(v) && !p._isEaten(v) {
					p._eat(v)
					builder.WriteString(v)
					builder.WriteRune(' ')
				}
			}
		} else {
			for i := arg_idx; i < max_len; i++ {
				if len(args) == 0 && arg_val.is_required {
					if !arg_val.has_default_value() {
						args := []string{arg_val.get_raw_value()}
						msg := fmt.Sprintf("missing required argument: `%v`", args[0])
						ctx := fmt.Sprintf("Expected a required value corresponding to: `%v` but none was provided", arg_val.get_raw_value())

						return matches, throw_error(MissingRequiredArgument, msg, ctx).set_args(args)
					} else {
						builder.WriteString(arg_val.default_value)
					}
				} else {
					v := args[i]
					if p.isSpecialValue(v) {
						break
					} else if !p.isFlagLike(v) && !p._isEaten(v) {
						p._eat(v)
						builder.WriteString(v)
					} else if arg_val.has_default_value() {
						builder.WriteString(arg_val.default_value)
					} else if arg_val.is_required {
						args := []string{arg_val.get_raw_value()}
						msg := fmt.Sprintf("missing required argument: `%v`", args[0])
						ctx := fmt.Sprintf("Expected a value for argument: `%v`, but instead found: `%v`", arg_val.name, v)

						return matches, throw_error(MissingRequiredArgument, msg, ctx).set_args(args)
					} else {
						continue
					}
				}
			}
		}

		// test the value against default values if any
		if len(arg_val.valid_values) > 0 && !arg_val.test_value(builder.String()) {
			args := []string{builder.String()}
			msg := fmt.Sprintf("the passed value: `%v`, is not a valid argument", args[0])
			ctx := fmt.Sprintf("Expected one of: `%v`, but instead found: `%v`, which is not a valid value", arg_val.valid_values, builder.String())

			return matches, throw_error(InvalidArgumentValue, msg, ctx).set_args(args)
		}

		// test the value against the validator func if any
		if arg_val.validator_fn != nil {
			if err := arg_val.validator_fn(builder.String()); err != nil {
				args := []string{builder.String()}
				msg := fmt.Sprintf("the passed value: `%v`, is not a valid argument", args[0])
				ctx := fmt.Sprintf("The validator function threw the following error: \"%v\" when checking the value: `%v`", err.Error(), builder.String())

				return matches, throw_error(InvalidArgumentValue, msg, ctx).set_args(args)
			}
		}

		arg_cfg := arg_matches{
			raw_value:   builder.String(),
			instance_of: *arg_val,
		}

		matches = append(matches, arg_cfg)
	}

	return matches, nil_error()
}
