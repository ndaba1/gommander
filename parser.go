package gommander

import (
	"errors"
	"fmt"
	"os"
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
	cursor_index   int
}

type arg_matches struct {
	raw_value    string
	instance_of  Argument
	cursor_index int
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

func (pm *ParserMatches) GetMatchedCommand() (*Command, int) {
	return pm.matched_cmd, pm.matched_cmd_idx
}

func (pm *ParserMatches) ContainsFlag(val string) bool {
	for _, v := range pm.flag_matches {
		flag := v.matched_flag
		if flag.short == val || flag.long == val || flag.name == val {
			return true
		}
	}
	return false
}

func (pm *ParserMatches) ContainsOption(val string) bool {
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			return true
		}
	}
	return false
}

func (pm *ParserMatches) GetArgValue(val string) (string, int, error) {
	for _, v := range pm.arg_matches {
		arg := v.instance_of
		if arg.name == val || arg.get_raw_value() == val {
			return v.raw_value, v.cursor_index, nil
		}
	}

	return "", -1, errors.New("no value found for provided argument")
}

func (pm *ParserMatches) GetOptionArg(val string) (string, int, error) {
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			// TODO: Probably check if slice is empty
			return v.passed_args[0].raw_value, v.cursor_index, nil
		}
	}
	return "", -1, errors.New("no value found for the provided option")
}

func (pm *ParserMatches) GetOptionArgsCount(val string) int {
	for _, v := range pm.option_matches {
		opt := v.matched_opt
		if opt.short == val || opt.long == val || opt.name == val {
			return v.instance_count
		}
	}
	return 0
}

func (pm *ParserMatches) GetAllOptionArgs(val string) []string {
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
	cursor      int
	root_cmd    *Command
	current_cmd *Command
	matches     ParserMatches
	eaten       []string
	cmd_idx     int
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
		if s.name == val || s.alias == val {
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

func (p *Parser) parse(raw_args []string) (ParserMatches, error) {
	if len(raw_args) == 0 {
		p.current_cmd.PrintHelp()
		os.Exit(0)
	}

	p.matches.raw_args = raw_args
	p.matches.arg_count = len(raw_args)

	allow_positional_args := false

	for index, arg := range raw_args {
		p.cursor = index
		// Basic lookahead
		// next_token := raw_args[index+1]
		// prev_token := raw_args[index-1]

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
				if err != nil {
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
					return p.matches, err
				}

				temp := []string{parts[1]}
				temp = append(temp, raw_args[(index+1):]...)

				err = p.parse_option(opt, temp)
				if err != nil {
					return p.matches, err
				}
			} else if allow_positional_args {
				p._eat(arg)
				p.matches.positional_args = append(p.matches.positional_args, arg)
			} else if !p._isEaten(arg) && !allow_positional_args {
				// TODO: Throw unresolved option error
				return p.matches, fmt.Errorf("unresolved option found: %v", arg)
			}
		} else if sc, err := p.getSubCommand(arg); err == nil {
			// handle subcmd
			p._eat(arg)
			p.current_cmd = sc
			p.cmd_idx = index
			// p.parse_cmd(raw_args[p.cursor:])
			continue
		} else if allow_positional_args {
			// TODO: More conditionals
			p._eat(arg)
			p.matches.positional_args = append(p.matches.positional_args, arg)
		} else if !p._isEaten(arg) {
			p.parse_cmd(raw_args[p.cursor:])
		}

	}
	// sanity check incase its not set for some reason
	p.matches.matched_cmd = p.current_cmd
	p.matches.matched_cmd_idx = p.cmd_idx

	return p.matches, nil
}

func (p *Parser) parse_option(opt *Option, raw_args []string) error {
	args, err := p.get_arg_matches(opt.args, raw_args)
	if err != nil {
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

	return nil
}

func (p *Parser) parse_cmd(raw_args []string) error {
	arg_cfg_vals, err := p.get_arg_matches(p.current_cmd.arguments, raw_args[p.cmd_idx:])
	if err != nil {
		return err
	}

	if len(arg_cfg_vals) > 0 {
		p.matches.arg_matches = append(p.matches.arg_matches, arg_cfg_vals...)
	}

	if len(p.current_cmd.sub_commands) > 0 {
		return errors.New("no such subcmd found")
	} else {
		return errors.New("argument not resolved")
	}

	// return nil
}

func (p *Parser) get_arg_matches(list []*Argument, args []string) ([]arg_matches, error) {
	max_len := len(list)
	matches := []arg_matches{}

	for arg_idx, arg_val := range list {
		var builder strings.Builder

		if arg_val.variadic {
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
					// TODO: Throw option missing argument error
					return matches, errors.New("missing argument")
				}

				v := args[i]
				if p.isSpecialValue(v) {
					break
				} else if !p.isFlagLike(v) && !p._isEaten(v) {
					p._eat(v)
					builder.WriteString(v)
				} else if arg_val.is_required {
					// TODO: Throw option missing argument error
					return matches, errors.New("missing required argument")
				} else {
					continue
				}
			}

		}

		if len(arg_val.valid_values) > 0 && !arg_val.ValueIsValid(builder.String()) {
			// TODO: Throw invalid argument value
			return matches, errors.New("invalid argument value")
		}

		arg_cfg := arg_matches{
			raw_value:   builder.String(),
			instance_of: *arg_val,
		}

		matches = append(matches, arg_cfg)
	}

	return matches, nil
}
