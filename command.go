package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
)

var (
	mu     sync.Mutex
	output io.Writer
)

type Runner interface {
	Init() Command
	Run([]string) error
}

func Run(runner Runner) (Command, error) {
	if output == nil {
		SetOutput(os.Stdout)
	}

	cmd := runner.Init()

	if cmd.commands == nil {
		cmd.commands = make([]*Command, 0)
	}

	if cmd.runners == nil {
		cmd.runners = make(map[string]Runner)
	}

	if cmd.Version == "" {
		cmd.Version = "dev"
	}

	cmd.runner = runner
	cmd.runners[cmd.Name] = runner

	if err := cmd.parseCommands(cmd.Name, os.Args[1:]); err != nil {
		if errors.Is(err, PrintHelp) {
			return cmd, nil
		}

		return cmd, err
	}

	return cmd, nil
}

func SetOutput(out io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	output = out
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
	Flags []Flag

	// actual is the full set of flags passed to the command.
	actual map[Option]Flag

	// formal is the full set of flags represented as options.
	formal []Option

	// args is the full set of arguments passed to the command.
	// args []Option

	// commands is a list of commands.
	commands []*Command

	// runners is a list of commands that satisfies the Runner interface.
	runners map[string]Runner

	// runner is the runner for the current command.
	runner Runner

	// parent is the parent command for this command.
	parent *Command

	// flagUsage is the combined usage of all flags for the command.
	flagUsage string

	// commandUsage is the combined usage of all commands for the command.
	commandUsage string

	// fullName is the combined names of the root command and it's subcommands, if applicable.
	fullName string
}

func (c *Command) HasAvailableFlags() bool {
	return len(c.Flags) != 0
}

func (c *Command) HasAvailableCommands() bool {
	return len(c.commands) != 0
}

func (c *Command) HasParent() bool {
	return c.parent != nil
}

func (c *Command) CommandHelp() string {
	return c.commandUsage
}

func (c *Command) FlagHelp() string {
	return c.flagUsage
}

func (c *Command) Output() io.Writer {
	return output
}

// FullName returns the full name of the command (e.g. test server start).
func (c *Command) FullName() string {
	if c.fullName == "" {
		return c.Name
	}
	return c.fullName
}

func (c *Command) AddCommands(runners ...Runner) {
	if c.commands == nil {
		c.commands = make([]*Command, 0)
	}

	if c.runners == nil {
		c.runners = make(map[string]Runner)
	}

	for _, runner := range runners {
		cmd := runner.Init()
		c.commands = append(c.commands, &cmd)
		c.runners[cmd.Name] = runner
		cmd.parent = c
	}
}

// PrintHelp prints the command's help.
func (c *Command) PrintHelp() {
	if err := renderTemplate(output, UsageTemplate, c); err != nil {
		panic(err)
	}
}

// PrintVersion prints the version of the command.
func (c *Command) PrintVersion() {
	_, _ = fmt.Fprintln(output, c.Version)
}

func (c *Command) parseCommands(name string, args []string) error {
	if len(c.commands) > 0 {
		c.sortCommands()
	}

	for _, arg := range args {
		if cmd, ok := c.lookupCommand(arg); ok {
			fullName := fmt.Sprintf("%s %s", name, cmd.Name)

			if err := cmd.addParentFlags(c); err != nil {
				return err
			}

			if cmd.Version == "" {
				cmd.Version = cmd.parent.Version
			}

			cmd.runner = c.runners[cmd.Name]

			return cmd.parseCommands(fullName, args[1:])
		}
	}

	c.fullName = name

	if err := c.parseFlags(args); err != nil {
		if errors.Is(err, PrintHelp) {
			c.PrintHelp()
			return PrintHelp
		}

		return err
	}

	if err := c.runner.Run(args); err != nil {
		if errors.Is(err, PrintHelp) {
			c.PrintHelp()
			return PrintHelp
		}

		return err
	}

	return nil
}

func (c *Command) parseFlags(args []string) error {
	if err := c.addFlags(); err != nil {
		return err
	}

	c.sortFlags()
	c.parseUsage()

	var flagArgs []string

	for i, arg := range args {
		// If the current argument isn't a flag, then skip it.
		if !c.isFlag(arg) {
			continue
		}

		// Trim the dashes before we check if we've seen this flag
		// before.
		arg = strings.TrimLeft(arg, "-")

		flag, ok := c.lookupFlag(arg)
		if !ok {
			// This flag is not present in the commands list of flags
			// therefore it is invalid.
			err := ErrOptionNotDefined{arg: arg}
			fmt.Fprintln(output, err.Error())
			return PrintHelp
		}

		option, err := flag.Option()
		if err != nil {
			return err
		}

		if option.Shorthand == HelpFlag.Shorthand || option.Name == HelpFlag.Name {
			return PrintHelp
		}

		if option.Shorthand == VersionFlag.Shorthand || option.Name == VersionFlag.Name {
			c.PrintVersion()
			return nil
		}

		flagArgs = args[i:]

		switch x := flag.(type) {
		case *BoolFlag:
			if err := x.Set("true"); err != nil {
				return err
			}
		case *StringSliceFlag:
			x.Clear()

			for _, v := range flagArgs[1:] {
				if err := x.Set(v); err != nil {
					return err
				}
			}
		default:
			flagArgs = flagArgs[1:]

			if len(flagArgs) > 0 {
				if err := x.Set(flagArgs[0]); err != nil {
					return err
				}
			}
		}
	}

	if err := c.checkRequiredOptions(); err != nil {
		return err
	}

	return nil
}

func (c *Command) parseUsage() {
	for _, flag := range c.formal {
		name := colorize(ColorYellow, "--%s", flag.Name)
		shorthand := colorize(ColorYellow, "-%s", flag.Shorthand)
		usage := flag.Desc

		if flag.Shorthand != "" {
			c.flagUsage += fmt.Sprintf("    %s, %s\n\t    %s\n\n", shorthand, name, usage)
		} else {
			c.flagUsage += fmt.Sprintf("        %s\n\t    %s\n\n", name, usage)
		}
	}

	maxLen := findMaxLength(c.commands)

	for _, cmd := range c.commands {
		c.commandUsage += fmt.Sprintf("    %s%s\n", rpad(colorize(ColorYellow, cmd.Name), computePadding(maxLen, cmd.Name)), cmd.Desc)
	}
}

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

func (c *Command) addFlags() error {
	if c.Flags == nil {
		c.Flags = make([]Flag, 0)
	}

	if c.actual == nil {
		c.actual = make(map[Option]Flag)
	}

	if c.formal == nil {
		c.formal = make([]Option, 0)
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

func (c *Command) addFlag(flag Flag) error {
	opt, err := flag.Option()
	if err != nil {
		return err
	}

	c.actual[opt] = flag
	c.formal = append(c.formal, opt)
	return nil
}

func (c *Command) checkRequiredOptions() error {
	result := &multierror.Error{}

	for _, option := range c.formal {
		if !option.HasBeenSet() && option.Required {
			var err error

			if option.Shorthand != "" {
				err = fmt.Errorf(colorize(ColorRed, "-%s, --%s is required", option.Shorthand, option.Name))
			} else {
				err = fmt.Errorf(colorize(ColorRed, "--%s is required", option.Name))
			}

			_ = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func (c *Command) sortCommands() {
	sort.Sort(SortCommandsByName(c.commands))
}

func (c *Command) sortFlags() {
	sort.Sort(SortOptionsByName(c.formal))
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
	for option, flag := range c.actual {
		if option.Shorthand == arg || option.Name == arg {
			return flag, true
		}
	}
	return nil, false
}

func (c *Command) isFlag(arg string) bool {
	if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
		return true
	}
	return false
}

type SortCommandsByName []*Command

func (n SortCommandsByName) Len() int           { return len(n) }
func (n SortCommandsByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n SortCommandsByName) Less(i, j int) bool { return n[i].Name < n[j].Name }
