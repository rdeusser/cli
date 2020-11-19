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

func Run(runner Runner) error {
	if output == nil {
		SetOutput(os.Stdout)
	}

	cmd := runner.Init()

	if cmd.commands == nil {
		cmd.commands = make([]*Command, 0, 0)
	}

	if cmd.runners == nil {
		cmd.runners = make(map[string]Runner)
	}

	if cmd.Version == "" {
		cmd.Version = "dev"
	}

	cmd.runners[cmd.Name] = runner

	if err := cmd.parseCommands(cmd.Name, os.Args[1:]); err != nil {
		return err
	}

	return nil
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

	// flags is the full set of flags passed to the command.
	flags []Option

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
	return len(c.flags) != 0
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
		c.commands = make([]*Command, 0, 0)
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
			fullName := name + " " + cmd.Name

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
			return nil
		}

		return err
	}

	if c.helpRequested(args) || c.versionRequested(args) {
		c.PrintHelp()
		return nil
	}

	if err := c.runner.Run(args); err != nil {
		if errors.Is(err, PrintHelp) {
			c.PrintHelp()
			return nil
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

	seen := make(map[string]struct{})

	for _, flag := range c.flags {
		if flag.Name != "" {
			seen[flag.Name] = struct{}{}
		}

		if flag.Shorthand != "" {
			seen[flag.Shorthand] = struct{}{}
		}
	}

	for _, arg := range args {
		// If the current argument isn't a flag, then skip it.
		if !c.isFlag(arg) {
			continue
		}

		// Trim the dashes before we check if we've seen this flag
		// before.
		flag := strings.TrimLeft(arg, "-")

		// If we've seen this flag before, then skip it.
		if _, ok := seen[flag]; ok {
			continue
		}

		// This flag is not present in the commands list of flags
		// therefore it is invalid.
		err := ErrOptionNotDefined{arg: arg}
		fmt.Fprintln(output, err.Error())
		return PrintHelp
	}

	if err := c.checkRequiredFlags(); err != nil {
		return err
	}

	return nil
}

func (c *Command) parseUsage() {
	for _, flag := range c.flags {
		name := good(fmt.Sprintf("--%s", flag.Name))
		shorthand := good(fmt.Sprintf("-%s", flag.Shorthand))
		usage := flag.Desc

		if flag.Shorthand != "" {
			c.flagUsage += fmt.Sprintf("    %s, %s\n\t    %s\n\n", shorthand, name, usage)
		} else {
			c.flagUsage += fmt.Sprintf("        %s\n\t    %s\n\n", name, usage)
		}
	}

	maxLen := findMaxLength(c.commands)

	for _, cmd := range c.commands {
		c.commandUsage += fmt.Sprintf("    %s%s\n", rpad(good(cmd.Name), computePadding(maxLen, cmd.Name)), cmd.Desc)
	}
}

func (c *Command) addFlags() error {
	if c.Flags == nil {
		c.Flags = make([]Flag, 0)
	}

	if c.flags == nil {
		c.flags = make([]Option, 0, 0)
	}

	c.addFlag(HelpFlag)
	c.addFlag(VersionFlag)

	for _, flag := range c.Flags {
		if err := c.addFlag(flag); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) addFlag(flag Flag) error {
	opt, err := flag.GetOption()
	if err != nil {
		return err
	}

	c.flags = append(c.flags, opt)
	return nil
}

func (c *Command) checkRequiredFlags() error {
	result := &multierror.Error{}

	for name, flag := range c.flags {
		if !flag.HasBeenSet() && flag.Required {
			_ = multierror.Append(result, fmt.Errorf(bad("%s flag must be provided", name)))
		}
	}

	return result.ErrorOrNil()
}

func (c *Command) sortCommands() {
	sort.Sort(SortCommandsByName(c.commands))
}

func (c *Command) sortFlags() {
	sort.Sort(SortOptionsByName(c.flags))
}

func (c *Command) lookupCommand(name string) (*Command, bool) {
	for _, cmd := range c.commands {
		if cmd.Name == name {
			return cmd, true
		}
	}
	return nil, false
}

func (c *Command) helpRequested(args []string) bool {
	return c.isFlagSet(args, []string{HelpFlag.Name, HelpFlag.Shorthand})
}

func (c *Command) versionRequested(args []string) bool {
	return c.isFlagSet(args, []string{VersionFlag.Name, VersionFlag.Shorthand})
}

func (c *Command) isFlagSet(args, matches []string) bool {
	if len(args) == 0 || len(matches) == 0 {
		return false
	}

	for _, arg := range args {
		if !c.isFlag(arg) {
			continue
		}

		for _, match := range matches {
			if match == strings.TrimLeft(arg, "-") {
				return true
			}
		}
	}

	return false
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
