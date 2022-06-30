package cli

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/termenv"
)

// Runner is the interface for initializing and running a command.
type Runner interface {
	Init() Command
	Run() error
}

// Run parses the commands, flags, and arguments and runs the command.
func Run(args []string, runner Runner) error {
	cmd := runner.Init()

	if cmd.Version == "" {
		cmd.Version = "dev"
	}

	if cmd.commands == nil {
		cmd.commands = make([]*Command, 0)
	}

	if cmd.runners == nil {
		cmd.runners = make(map[string]Runner)
	}

	if cmd.output == nil {
		cmd.output = os.Stdout
	}

	cmd.runner = runner
	cmd.runners[cmd.Name] = runner

	if err := cmd.parseCommands(cmd.Name, args[1:]); err != nil {
		if errors.Is(err, ErrPrintHelp) {
			return nil
		}

		return err
	}

	return nil
}

type Command struct {
	// Name is the name of the command.
	Name string

	// Desc is the short description the command.
	Desc string

	// LongDesc is the long description of the command.
	LongDesc string

	// Version is the version of the command.
	Version string

	// Flags is the full set of flags passed to the command.
	Flags Flags

	// Args is the arguments passed to the command after flags have been
	// processed.
	Args Args

	// flagOptions is the full set of flag options.
	flagOptions []FlagOption

	// globalFlagOptions are flags passed down by the parent command(s).
	globalFlagOptions []FlagOption

	// argOptions it he full set of arg options.
	argOptions []ArgOption

	// commands is a list of commands.
	commands []*Command

	// runners is a list of commands that satisfies the Runner interface.
	runners map[string]Runner

	// runner is the runner for the current command.
	runner Runner

	// parent is the parent command for this command.
	parent *Command

	// commandUsage is the combined usage of all commands for the command.
	commandUsage string

	// flagUsage is the combined usage of all flags for the command.
	flagUsage string

	// argUsage is the combined usage of all arguments for the command.
	argUsage string

	// output is where help and errors are written to.
	output io.Writer
}

// HasAvailableCommands returns true if there is at least one subcommand.
func (c *Command) HasAvailableCommands() bool {
	return len(c.commands) > 0
}

// HasAvailableFlags returns true if there is at least one flag.
func (c *Command) HasAvailableFlags() bool {
	return len(c.Flags) > 0
}

// HasAvailableArgs returns true if there is at least one argument.
func (c *Command) HasAvailableArgs() bool {
	return len(c.argOptions) > 0
}

// HasParent returns true if the command has a parent command (i.e the current
// command is a subcommand).
func (c *Command) HasParent() bool {
	return c.parent != nil
}

// CommandHelp returns the usage output for the command.
func (c *Command) CommandHelp() string {
	return c.commandUsage
}

// FlagHelp returns the usage output for the flags for this command.
func (c *Command) FlagHelp() string {
	return c.flagUsage
}

// ArgHelp returns the usage output for the arguments for this command.
func (c *Command) ArgHelp() string {
	return c.argUsage
}

// FullName returns the full name of the command (e.g. test server start).
func (c *Command) FullName() string {
	names := make([]string, 0)
	names = append(names, c.Name)

	parent := c.parent
	for parent != nil {
		names = append(names, parent.Name)
		parent = parent.parent
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	return strings.Join(names, " ")
}

// Output returns the io.Writer that the command uses to write output to.
func (c *Command) Output() io.Writer {
	return c.output
}

// Output sets the io.Writer that the command uses to write output to.
func (c *Command) SetOutput(w io.Writer) {
	c.output = w
}

// AddCommands adds a runner or runners to the command.
func (c *Command) AddCommands(runners ...Runner) {
	if c.commands == nil {
		c.commands = make([]*Command, 0)
	}

	if c.runners == nil {
		c.runners = make(map[string]Runner)
	}

	for _, runner := range runners {
		cmd := runner.Init()
		cmd.parent = c
		cmd.output = c.output
		c.commands = append(c.commands, &cmd)
		c.runners[cmd.Name] = runner
	}
}

// PrintHelp prints the command's help.
func (c *Command) PrintHelp() {
	description := c.LongDesc
	if description == "" {
		description = c.Desc
	}

	data := struct {
		Description          string
		Command              string
		HasAvailableCommands bool
		HasAvailableFlags    bool
		HasAvailableArgs     bool
		Args                 string
		CommandHelp          string
		FlagHelp             string
		ArgHelp              string
	}{
		Description:          description,
		Command:              c.FullName(),
		HasAvailableCommands: c.HasAvailableCommands(),
		HasAvailableFlags:    c.HasAvailableFlags(),
		HasAvailableArgs:     c.HasAvailableArgs(),
		Args:                 joinArgs(c.argOptions, " "),
		CommandHelp:          c.CommandHelp(),
		FlagHelp:             c.FlagHelp(),
		ArgHelp:              c.ArgHelp(),
	}

	if err := renderTemplate(c.output, UsageTemplate, data); err != nil {
		panic(err)
	}
}

// PrintVersion prints the version of the command.
func (c *Command) PrintVersion() {
	_, _ = fmt.Fprintln(c.output, c.Version)
}

func (c *Command) parseCommands(name string, args []string) error {
	if len(c.commands) > 0 {
		c.sortCommands()
	}

	for _, arg := range args {
		if cmd, ok := c.lookupCommand(arg); ok {
			if err := cmd.addParentFlags(c); err != nil {
				return err
			}

			if cmd.Version == "" {
				cmd.Version = cmd.parent.Version
			}

			cmd.runner = c.runners[cmd.Name]
			cmd.output = c.output

			return cmd.parseCommands(cmd.FullName(), args[1:])
		}
	}

	if err := c.addFlags(); err != nil {
		return err
	}

	if err := c.addArgs(); err != nil {
		return err
	}

	c.sortFlags()
	c.sortArgs()
	c.parseUsage()

	if err := c.parseFlags(args); err != nil {
		if errors.Is(err, ErrPrintHelp) {
			c.PrintHelp()
			return ErrPrintHelp
		}

		return err
	}

	if err := c.parseArgs(args); err != nil {
		return err
	}

	if err := c.runner.Run(); err != nil {
		if errors.Is(err, ErrPrintHelp) {
			c.PrintHelp()
			return ErrPrintHelp
		}

		return err
	}

	return nil
}

// parseFlags parses the flags passed to a command.
func (c *Command) parseFlags(args []string) error {
	var flags []string

	for i, arg := range args {
		// If the current argument isn't a flag, then skip it.
		if !c.isFlag(arg) {
			continue
		}

		flag, ok := c.lookupFlag(arg)
		if ok && c.isFlagLong(arg) && len(arg) == 3 {
			ok = false // edge case where passing a shorthand flag with an extra dash (so long form) is valid
		}
		if !ok {
			// This flag is not present in the commands list of flags
			// therefore it is invalid.
			err := ErrFlagNotDefined{flag: arg}
			fmt.Fprintln(c.output, err.Error())
			return ErrPrintHelp
		}

		option := flag.Option()

		if option.Shorthand == HelpFlag.Shorthand || option.Name == HelpFlag.Name {
			return ErrPrintHelp
		}

		if option.Shorthand == VersionFlag.Shorthand || option.Name == VersionFlag.Name {
			c.PrintVersion()
			return nil
		}

		flags = args[i:]

		switch x := flag.(type) {
		case *BoolFlag:
			if err := x.Set("true"); err != nil {
				return err
			}
		case *StringSliceFlag:
			flagArgs := strings.Split(arg, x.Separator)

			for _, flagArg := range flagArgs {
				if err := x.Set(flagArg); err != nil {
					return errors.Wrap(err, "setting flag value")
				}
			}
		default:
			flags = flags[1:]

			if len(flags) > 0 {
				if err := x.Set(flags[0]); err != nil {
					return err
				}
			}
		}

		if err := c.updateFlag(flag); err != nil {
			return errors.Wrap(err, "updating flag after setting value")
		}
	}

	return c.checkRequiredFlagOptions()
}

// parseArgs parses the arguments passed to a command.
func (c *Command) parseArgs(args []string) error {
	argArgs := args
	cmds := strings.Fields(c.FullName())

	if len(args) >= len(cmds) {
		for i, cmd := range cmds {
			// We should be okay to do this without any issues, but we'll
			// find out.
			if args[i] == cmd {
				argArgs = args[i+1:]
			}
		}
	}

	for i, arg := range argArgs {
		// If we've encountered a flag move past it.
		if c.isFlag(arg) {
			// We've already parsed flags by this point so we know
			// it's there and we know it's valid.
			flag, _ := c.lookupFlag(arg)

			switch x := flag.(type) {
			case *BoolFlag:
				argArgs = argArgs[i+1:]
			case *StringSliceFlag:
				// Move past the flag itself.
				argArgs = argArgs[i+1:]

				// And then move past all of it's arguments.
				flagArgs := strings.Split(arg, x.Separator)
				argArgs = argArgs[len(flagArgs):]

			default:
				argArgs = argArgs[i+2:]
			}
		}
	}

	// By this point we should have already moved past the command,
	// subcommands, and flags.
	for i, arg := range argArgs {
		actual, ok := c.lookupArg(i)
		if !ok {
			// This arg is not present in the commands list of args
			// therefore it is invalid.
			err := ErrArgNotDefined{
				arg:      arg,
				position: i,
			}

			fmt.Fprintln(c.output, err.Error())
			return ErrPrintHelp
		}

		if err := actual.Set(arg); err != nil {
			return err
		}

		switch x := actual.(type) {
		case *BoolArg:
			if err := x.Set(arg); err != nil {
				return nil
			}
		default:
			if err := x.Set(arg); err != nil {
				return nil
			}
		}

		// Flags have an actual name or shorthand we can find in the
		// input, but arguments don't so we have to do this a little bit
		// differently by setting the value of the argument and then
		// updating the value in the command.
		if err := c.updateArg(actual); err != nil {
			return err
		}
	}

	return c.checkRequiredArgOptions()
}

// parseUsage generates usage strings for commands, flags, and arguments.
func (c *Command) parseUsage() {
	for _, cmd := range c.commands {
		maxLen := findMaxCommandLength(c.commands)
		c.commandUsage += fmt.Sprintf("    %s%s\n", rpad(termenv.Colorize(termenv.ColorGreen, cmd.Name), computePadding(maxLen, cmd.Name)), cmd.Desc)
	}

	for _, flag := range c.flagOptions {
		maxLen := findMaxFlagLength(c.flagOptions)
		name := rpad(termenv.Colorize(termenv.ColorGreen, "--%s", flag.Name), computePadding(maxLen, flag.Name))
		shorthand := termenv.Colorize(termenv.ColorGreen, "-%s", flag.Shorthand)
		usage := flag.Desc

		if flag.Shorthand != "" {
			c.flagUsage += fmt.Sprintf("    %s, %s %s\n", shorthand, name, usage)
		} else {
			c.flagUsage += fmt.Sprintf("        %s %s\n", name, usage)
		}
	}

	for _, arg := range c.argOptions {
		maxLen := findMaxArgLength(c.argOptions)
		name := rpad(termenv.Colorize(termenv.ColorGreen, "<%s>", arg.Name), computePadding(maxLen, arg.Name))
		if arg.Name == "" {
			name = termenv.Colorize(termenv.ColorGreen, "<unknown argument name>")
		}

		usage := arg.Desc
		c.argUsage += fmt.Sprintf("    %s %s (type: %s)\n\n", name, usage, arg.Type)
	}
}

// addParentFlags adds the flags the parent currently has to this command.
func (c *Command) addParentFlags(parent *Command) error {
	seen := make(map[Flag]struct{})

	for _, flag := range c.Flags {
		seen[flag] = struct{}{}
	}

	for _, flag := range parent.Flags {
		if _, ok := seen[flag]; !ok {
			c.Flags = append(c.Flags, flag)
		}
	}

	return nil
}

// addFlags adds the flags passed to the command as flag options for parsing.
func (c *Command) addFlags() error {
	if c.Flags == nil {
		c.Flags = make([]Flag, 0)
	}

	if c.flagOptions == nil {
		c.flagOptions = make([]FlagOption, 0)
	}

	seen := make(map[Flag]struct{})

	for _, flag := range c.Flags {
		seen[flag] = struct{}{}
	}

	if _, ok := seen[&HelpFlag]; !ok {
		c.Flags = append(c.Flags, &HelpFlag)
	}

	if _, ok := seen[&VersionFlag]; !ok {
		c.Flags = append(c.Flags, &VersionFlag)
	}

	for _, flag := range c.Flags {
		if err := c.addFlag(flag); err != nil {
			return err
		}
	}

	return nil
}

// addFlag adds a flag option to the command.
func (c *Command) addFlag(flag Flag) error {
	if err := flag.Init(); err != nil {
		return errors.Wrap(err, "applying flag")
	}

	option := flag.Option()
	env, ok := os.LookupEnv(option.EnvVar)
	if ok {
		if err := flag.Set(env); err != nil {
			return errors.Wrap(err, "setting flag via environment variable")
		}
	}

	c.flagOptions = append(c.flagOptions, option)
	return nil
}

// updateFlag updates the flag option for the command.
func (c *Command) updateFlag(flag Flag) error {
	flagOption := flag.Option()

	for i, option := range c.flagOptions {
		// If the name and the shorthand are the same for the provided
		// flag and the one we already have, they're definitely the
		// same.
		if flagOption.Name == option.Name && flagOption.Shorthand == option.Shorthand {
			c.flagOptions[i] = flagOption
		}
	}

	return nil
}

// addArgs adds the arguments provided to the command as arg options to the
// command.
func (c *Command) addArgs() error {
	if c.Args == nil {
		c.Args = make([]Arg, 0)
	}

	if c.argOptions == nil {
		c.argOptions = make([]ArgOption, 0)
	}

	for _, arg := range c.Args {
		if err := c.addArg(arg); err != nil {
			return err
		}
	}

	if !c.consecutiveArgPositions() {
		return fmt.Errorf("positions in the args for %s command are not in consecutive order", c.Name)
	}

	return nil
}

// addArg adds an argument to the command.
func (c *Command) addArg(arg Arg) error {
	if err := arg.Init(); err != nil {
		return err
	}

	c.argOptions = append(c.argOptions, arg.Option())
	return nil
}

// updateArg updates the arg option for the command.
func (c *Command) updateArg(arg Arg) error {
	argOption := arg.Option()

	for i, option := range c.argOptions {
		// If the name and the position are the same for the provided
		// arg and the one we already have, they're definitely the same.
		if argOption.Name == option.Name && argOption.Position == option.Position {
			c.argOptions[i] = argOption
		}
	}

	return nil
}

// consecutiveArgPositions validates that the args added to the command have the
// correct positions.
//
// For example, if arg 1 has position 0 but arg 2 has position 3, that isn't
// valid.
func (c *Command) consecutiveArgPositions() bool {
	positions := make([]int, 0)

	for _, option := range c.argOptions {
		positions = append(positions, option.Position)
	}

	for i := 0; i < len(positions); i++ {
		if i+1 == len(positions) {
			break // we're at the end
		}

		if positions[i]+1 != positions[i+1] {
			return false
		}
	}

	return true
}

// checkRequiredFlagOptions ensures that flags that are required have been set.
func (c *Command) checkRequiredFlagOptions() error {
	result := &multierror.Error{}

	for _, flag := range c.flagOptions {
		if !flag.HasBeenSet && flag.Required {
			var err error

			if flag.Shorthand != "" {
				err = fmt.Errorf(termenv.Colorize(termenv.ColorRed, "-%s, --%s is required", flag.Shorthand, flag.Name))
			} else {
				err = fmt.Errorf(termenv.Colorize(termenv.ColorRed, "--%s is required", flag.Name))
			}

			_ = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// checkRequiredArgOptions ensures that args that are required have been set.
func (c *Command) checkRequiredArgOptions() error {
	result := &multierror.Error{}

	for _, arg := range c.argOptions {
		if !arg.HasBeenSet && arg.Required {
			var err error

			if arg.Name != "" {
				err = fmt.Errorf(termenv.Colorize(termenv.ColorRed, "%s is required", arg.Name))
			} else {
				err = fmt.Errorf(termenv.Colorize(termenv.ColorRed, "%s is required", reflect.TypeOf(arg).Name()))
			}

			_ = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// sortCommands sorts commands by name.
func (c *Command) sortCommands() {
	sort.Sort(SortCommandsByName(c.commands))
}

// sortFlags sorts flags by name.
func (c *Command) sortFlags() {
	sort.Sort(SortFlagOptionsByName(c.flagOptions))
}

// sortArgs sorts arguments by position.
func (c *Command) sortArgs() {
	sort.Sort(SortArgOptionsByPosition(c.argOptions))
}

func (c *Command) lookupCommand(arg string) (*Command, bool) {
	for _, cmd := range c.commands {
		if cmd.Name == arg {
			return cmd, true
		}
	}
	return nil, false
}

func (c *Command) lookupFlag(arg string) (Flag, bool) {
	// Trim the dashes before we check if we've seen this flag before.
	name := strings.TrimLeft(arg, "-")

	flag := c.Flags.Lookup(name, name)
	if flag == nil {
		return nil, false
	}

	return flag, true
}

func (c *Command) lookupArg(position int) (Arg, bool) {
	arg := c.Args.Lookup(position)
	if arg == nil {
		return nil, false
	}

	return arg, true
}

func (c *Command) isFlag(arg string) bool {
	if c.isFlagShorthand(arg) || c.isFlagLong(arg) {
		return true
	}

	return false
}

func (c *Command) isFlagShorthand(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func (c *Command) isFlagLong(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

// SortCommandOptionsByName sorts commands by name.
type SortCommandsByName []*Command

func (n SortCommandsByName) Len() int           { return len(n) }
func (n SortCommandsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortCommandsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

// rpad adds padding to the right side of a string.
func rpad(s string, count int) string {
	if count < 0 {
		count = 0
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", count))
}

// computePadding computes the padding needed for displaying usage text.
func computePadding(maxLen int, s string) int {
	return maxLen - len(s) + 4
}

// findMaxCommandLength sorts a map of commands by their length and returns the
// length of the longest command name.
func findMaxCommandLength(commands []*Command) int {
	if len(commands) == 0 {
		return 0
	}

	list := make([]int, 0, len(commands))

	for _, cmd := range commands {
		list = append(list, len(cmd.Name))
	}

	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(list)-1; i++ {
			if list[i+1] > list[i] {
				list[i+1], list[i] = list[i], list[i+1]
				swapped = true
			}
		}
	}

	return list[0]
}

// findMaxFlagLength finds the maximum length of a flag name.
func findMaxFlagLength(options []FlagOption) int {
	maxLength := 0

	for _, option := range options {
		length := len(option.Name)

		if maxLength == 0 {
			maxLength = length
		}

		if length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}

// findMaxArgLength finds the maximum length of an arg name.
func findMaxArgLength(options []ArgOption) int {
	maxLength := 0

	for _, option := range options {
		length := len(option.Name)

		if maxLength == 0 {
			maxLength = length
		}

		if length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}
