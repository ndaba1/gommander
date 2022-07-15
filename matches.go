package gommander

import "errors"

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

// Returns the number of arguments that were passed to the program for parsing
func (pm *ParserMatches) GetRawArgCount() int {
	return pm.argCount
}

// This method returns the actual raw arg values passed to the program
func (pm *ParserMatches) GetRawArgs() []string {
	return pm.rawArgs
}

func (pm *ParserMatches) GetPositionalArgs() []string {
	return pm.positionalArgs
}

// Returns a reference to the app instance
func (pm *ParserMatches) GetAppRef() *Command {
	return pm.rootCmd
}

// Returns a reference to the command or subcommand that was matched by the parser
func (pm *ParserMatches) GetMatchedCommand() *Command {
	return pm.matchedCmd
}

// Returns the index of the matched command. An index of -1 means that the program itself was the matched command
func (pm *ParserMatches) GetMatchedCommandIndex() int {
	return pm.matchedCmdIdx
}

// Returns whether or not a flag was passed to the program args.
// Accepts the name of the flag, or the short or long version of the flag
func (pm *ParserMatches) ContainsFlag(val string) bool {
	for _, v := range pm.flagMatches {
		flag := v.matchedFlag
		if flag.ShortVal == val || flag.LongVal == val || flag.Name == val {
			return true
		}
	}
	return false
}

// Returns whether or not an option was passed to the program args
// Accepts as input the name of the option, or its short or long version
func (pm *ParserMatches) ContainsOption(val string) bool {
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.ShortVal == val || opt.LongVal == val || opt.Name == val {
			return true
		}
	}
	return false
}

// A method used to get the value of an argument passed to the program.
// Takes as input the name of the argument or the raw value of the argument.
// If no value is found, or the argument is misspelled, an error is returned.
// If no value was passed to the argument but it is required, the default value is used if one exists, otherwise an error is thrown.
func (pm *ParserMatches) GetArgValue(val string) (string, error) {
	for _, v := range pm.argMatches {
		arg := v.instanceOf
		if arg.Name == val || arg.getRawValue() == val {
			return v.rawValue, nil
		}
	}

	return "", errors.New("no value found for provided argument")
}

// This method returns the value passed to an option, if any.
// An error is thrown if no such option exists
// If an option has a default value and none was provided, the default value is used.
func (pm *ParserMatches) GetOptionValue(val string) (string, error) {
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.ShortVal == val || opt.LongVal == val || opt.Name == val {
			// TODO: Probably check if slice is empty
			return v.passedArgs[0].rawValue, nil
		}
	}
	return "", errors.New("no value found for the provided option")
}

// If option values are provided multiple times, all the instances can be acquired using this method
// For example, `-p 80 -p 90 -p 100`. All these instances are stored in a single slice to be acquired via this method
func (pm *ParserMatches) GetAllOptionInstances(val string) []string {
	instances := []string{}
	for _, v := range pm.optionMatches {
		opt := v.matchedOpt
		if opt.ShortVal == val || opt.LongVal == val || opt.Name == val {
			for _, a := range v.passedArgs {
				instances = append(instances, a.rawValue)
			}
		}
	}
	return instances
}
