package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type Runner interface {
	Init() Command
	Run() error
}

func Run(runner Runner) error {
	cmd := runner.Init()

	if cmd.commands == nil {
		cmd.commands = make(map[string]*Command)
	}

	if cmd.runners == nil {
		cmd.runners = make(map[string]Runner)
	}

	if cmd.Version == "" {
		cmd.Version = "dev"
	}

	cmd.runners[cmd.Name] = runner

	if err := cmd.parseCommands(os.Args[1:]); err != nil {
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
	Flags []Flag

	// flags is the full set of flags passed to the command.
	flags map[string]*Option

	// args is the full set of arguments passed to the command.
	// args map[string]*Option

	// commands is a list of commands.
	commands map[string]*Command

	// runners is a list of commands that satisfies the Runner interface.
	runners map[string]Runner

	// parent is the parent command for this command.
	parent *Command

	// output specifies where to write the output.
	output io.Writer // nil means stdout

	// flagUsage is the combined usage of all flags for the command.
	flagUsage string

	// commandUsage is the combined usage of all commands for the command.
	commandUsage string
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
	if c.output == nil {
		return os.Stdout
	}
	return c.output
}

func (c *Command) SetOutput(output io.Writer) {
	c.output = output
}

func (c *Command) AddCommands(runners ...Runner) {
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}

	if c.runners == nil {
		c.runners = make(map[string]Runner)
	}

	for _, runner := range runners {
		cmd := runner.Init()
		c.commands[cmd.Name] = &cmd
		c.commands[cmd.Name].parent = c
		c.runners[cmd.Name] = runner
	}
}

func (c *Command) PrintHelp() {
	if err := renderTemplate(c.Output(), UsageTemplate, c); err != nil {
		panic(err)
	}
}

func (c *Command) PrintVersion() {
	_, _ = fmt.Fprintln(c.Output(), c.Version)
}

func (c *Command) parseCommands(args []string) error {
	if len(c.commands) > 0 {
		c.sortCommands()
		c.setOutput()
	}

	for _, arg := range args {
		if cmd, ok := c.commands[arg]; ok {
			if err := cmd.parseFlags(args); err != nil {
				return err
			}

			if cmd.Version == "" {
				cmd.Version = cmd.parent.Version
			}

			if len(cmd.commands) > 0 {
				return cmd.parseCommands(args)
			}

			if cmd.helpRequested(args) || cmd.versionRequested(args) {
				return nil
			}

			// Subcommand
			err := c.runners[arg].Run()
			if errors.Is(err, PrintHelp) {
				c.commands[arg].PrintHelp()
				return nil
			}
			if err != nil {
				return err
			}

			return nil
		}
	}

	if err := c.parseFlags(args); err != nil {
		return err
	}

	if c.helpRequested(args) || c.versionRequested(args) {
		return nil
	}

	// Root command
	err := c.runners[c.Name].Run()
	if errors.Is(err, PrintHelp) {
		c.PrintHelp()
		return nil
	}
	if err != nil {
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
		err := ErrOptionNotDefined{opt: c.flags[flag], arg: arg}
		fmt.Println(err.Error())
		c.PrintHelp()
		os.Exit(0)
	}

	if c.helpRequested(args) {
		c.PrintHelp()
		return nil
	}

	if c.versionRequested(args) {
		c.PrintVersion()
		return nil
	}

	if err := c.checkRequiredFlags(); err != nil {
		return err
	}

	return nil
}

func (c *Command) parseUsage() {
	for name, flag := range c.flags {
		name := good(fmt.Sprintf("--%s", name))
		shorthand := good(fmt.Sprintf("-%s", flag.Shorthand))
		usage := flag.Desc

		if flag.Shorthand != "" {
			c.flagUsage += fmt.Sprintf("    %s, %s\n\t    %s\n\n", shorthand, name, usage)
		} else {
			c.flagUsage += fmt.Sprintf("        %s\n\t    %s\n\n", name, usage)
		}
	}

	maxLen := findMaxLength(c.commands)

	for name, sub := range c.commands {
		c.commandUsage += fmt.Sprintf("    %s%s\n", rpad(good(name), computePadding(maxLen, name)), sub.Desc)
	}
}

func (c *Command) addFlags() error {
	if c.Flags == nil {
		c.Flags = make([]Flag, 0)
	}

	if c.flags == nil {
		c.flags = make(map[string]*Option)
	}

	c.Flags = append(c.Flags, []Flag{HelpFlag, VersionFlag}...)

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

	c.flags[opt.Name] = &opt
	return nil
}

func (c *Command) checkRequiredFlags() error {
	result := &multierror.Error{}

	for name, flag := range c.flags {
		if !flag.hasBeenSet && flag.Required {
			_ = multierror.Append(result, fmt.Errorf(bad("%s flag must be provided", name)))
		}
	}

	return result.ErrorOrNil()
}

func (c *Command) sortCommands() {
	commands := make(map[string]*Command, len(c.commands))
	keys := make([]string, 0, len(commands))

	for cmd := range c.commands {
		keys = append(keys, cmd)
	}

	sort.Strings(keys)

	for _, k := range keys {
		commands[k] = c.commands[k]
	}

	c.commands = commands
}

func (c *Command) sortFlags() {
	flags := make(map[string]*Option, len(c.flags))
	keys := make([]string, 0, len(flags))

	for flag := range c.flags {
		keys = append(keys, flag)
	}

	sort.Strings(keys)

	for _, k := range keys {
		flags[k] = c.flags[k]
	}

	c.flags = flags
}

func (c *Command) setOutput() {
	for _, sub := range c.commands {
		if c.output != nil {
			sub.output = c.output
		}

		if len(sub.commands) > 0 {
			sub.setOutput()
		}
	}
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
