package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/rdeusser/cli/ast"
	"github.com/rdeusser/cli/help"
	"github.com/rdeusser/cli/internal/errors"
	"github.com/rdeusser/cli/internal/join"
	"github.com/rdeusser/cli/internal/multierror"
	"github.com/rdeusser/cli/internal/slice"
	"github.com/rdeusser/cli/parser"
	"github.com/rdeusser/cli/tablewriter"
)

type VisitOption int

const (
	VisitStartingAtChild VisitOption = iota
	VisitStartingAtParent
	VisitStartingAtChildReverse
)

type VisitFunc func(*Command) error

// Command is a command. How else are you supposed to describe this?
// e.g. `go run main.go`
type Command struct {
	// Name is the name of the command.
	Name string

	// Desc is the short description the command.
	Desc string

	// LongDesc is the long description of the command.
	LongDesc string

	// Flags is the full set of flags passed to the command.
	Flags Flags

	// Args is the arguments passed to the command after flags have been
	// processed.
	Args Args

	// parent of the current command.
	parent *Command

	// commands is a map of command names to the commands command.
	commands map[string]*Command

	// The order for the below setters and runners is as follows:
	// 1. OptionSetter
	// 2. PersistentPreRunner
	// 3. PreRunner
	// 4. Runner
	// 5. PostRunner
	// 6. PersistentPostRunner

	// optionSetter sets options from parent commands.
	optionSetter OptionSetter

	// persistentPreRunner is inherited and run by all children of this command
	// before all other runners.
	persistentPreRunner PersistentPreRunner

	// preRunner is run before the main runner.
	preRunner PreRunner

	// runner is the runner for the current command.
	runner Runner

	// postRunner is run after the main runner.
	postRunner PostRunner

	// persistentPostRunner is inherited and run by all children of this command
	// after all other runners.
	persistentPostRunner PersistentPostRunner

	// stmt is the parsed statement.
	stmt *ast.Statement

	// usage is the combined usage of commands, flags, and arguments.
	usage string

	// output is where help and errors are written to.
	output io.Writer
}

// AddCommands adds commands to the current command as children.
func (c *Command) AddCommands(runners ...Runner) {
	c.init()

	for _, runner := range runners {
		cmd := runner.Init()
		cmd.setRunners(runner)
		cmd.init()
		cmd.parent = c
		cmd.stmt = c.stmt
		cmd.output = c.Output()

		c.commands[cmd.Name] = cmd
	}
}

// Output returns the io.Writer that the command uses to write output to.
func (c *Command) Output() io.Writer {
	if c.output == nil {
		return os.Stdout
	}

	return c.output
}

// SetOutput sets the io.Writer that the command uses to write output to.
func (c *Command) SetOutput(w io.Writer) {
	c.output = w
}

// FullName returns the full name of the command starting from the root.
func (c *Command) FullName() string {
	commands := make([]string, 0)

	c.Visit(func(cmd *Command) error {
		commands = append(commands, cmd.Name)
		return nil
	}, VisitStartingAtParent)

	return strings.Join(commands, " ")
}

// HasFlag checks to see if the provided flag is already added to the command.
func (c *Command) HasFlag(name, shorthand string) bool {
	seen := make(map[string]struct{})

	for _, flag := range c.Flags {
		opt := flag.Options()

		seen[opt.Name] = struct{}{}
		seen[opt.Shorthand] = struct{}{}
	}

	_, ok := seen[name]
	if !ok {
		return false
	}

	if shorthand != "" {
		_, ok = seen[shorthand]
		if !ok {
			return false
		}
	}

	return true
}

// PrintHelp prints the command's help.
func (c *Command) PrintHelp() {
	c.Output().Write([]byte(c.usage))
}

// Visit runs fn for each command starting from the top-most parent.
func (c *Command) Visit(fn VisitFunc, option VisitOption) error {
	commands := make([]*Command, 0)
	commands = append(commands, c)

	switch option {
	case VisitStartingAtChild:
		for _, cmd := range c.commands {
			commands = append(commands, cmd)
		}

		return visit(fn, commands)
	default:
		parent := c.parent
		for parent != nil {
			commands = append(commands, parent)
			parent = parent.parent
		}

		switch option {
		case VisitStartingAtParent:
			// We're starting from the rightmost command, so we have to reverse the
			// slice to get the leftmost command in the first index.
			for i := len(commands)/2 - 1; i >= 0; i-- {
				j := len(commands) - i - 1
				commands[i], commands[j] = commands[j], commands[i]
			}
		case VisitStartingAtChildReverse:
			// `commands` is already in correct order for this visit option.
		}

		return visit(fn, commands)

	}
}

// parseCommands is the main method. It adds parent flags (if applicable), sorts
// commands and flags, generates the usage string, parses the args, sets
// options, runs the runners, and checks for unknown and required
// arguments/flags.
func (c *Command) parseCommands(args []string) error {
	c.init()

	if err := c.addParentFlags(); err != nil {
		return err
	}

	c.sortCommands()
	c.sortFlags()
	c.generateUsage()

	for _, arg := range args {
		if cmd, ok := c.commands[arg]; ok {
			return cmd.parseCommands(args[1:])
		}
	}

	if err := c.parseArgs(args); err != nil {
		return c.errOrPrintHelp(err)
	}

	if c.optionSetter != nil {
		if c.parent == nil {
			return ErrMustHaveParent
		}

		if err := c.optionSetter.SetOptions(c.parent.Flags); err != nil {
			return c.errOrPrintHelp(err)
		}
	}

	if err := c.Visit(func(cmd *Command) error {
		if cmd.persistentPreRunner != nil {
			if err := cmd.persistentPreRunner.PersistentPreRun(); err != nil {
				return cmd.errOrPrintHelp(err)
			}
		}

		return nil
	}, VisitStartingAtParent); err != nil {
		return c.errOrPrintHelp(err)
	}

	if c.preRunner != nil {
		if err := c.preRunner.PreRun(); err != nil {
			return c.errOrPrintHelp(err)
		}
	}

	if err := c.runner.Run(); err != nil {
		return c.errOrPrintHelp(err)
	}

	if c.postRunner != nil {
		if err := c.postRunner.PostRun(); err != nil {
			return c.errOrPrintHelp(err)
		}
	}

	if err := c.Visit(func(cmd *Command) error {
		if cmd.persistentPostRunner != nil {
			if err := cmd.persistentPostRunner.PersistentPostRun(); err != nil {
				return cmd.errOrPrintHelp(err)
			}
		}

		return nil
	}, VisitStartingAtChildReverse); err != nil {
		return c.errOrPrintHelp(err)
	}

	return nil
}

// init just makes sure slices and maps are initialized before use.
func (c *Command) init() {
	if c.Args == nil {
		c.Args = make(Args, 0)
	}

	if c.Flags == nil {
		c.Flags = make(Flags, 0)
	}

	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}

	if c.output == nil {
		c.SetOutput(os.Stdout)
	}
}

// addParentFlags adds the flags the parent currently has to this command.
func (c *Command) addParentFlags() error {
	var merr multierror.Error

	if c.parent == nil {
		return nil
	}

	for _, flag := range c.Flags {
		opt := flag.Options()
		if c.parent.HasFlag(opt.Name, opt.Shorthand) {
			merr.Append(ErrFlagAlreadyDefined{
				Name:      opt.Name,
				Shorthand: opt.Shorthand,
			})
		}
	}

	c.Flags = append(c.Flags, c.parent.Flags...)

	return merr.ErrorOrNil()
}

// setRunners sets thee
func (c *Command) setRunners(runner Runner) {
	if v, ok := runner.(OptionSetter); ok {
		c.optionSetter = v
	}

	if v, ok := runner.(PersistentPreRunner); ok {
		c.persistentPreRunner = v
	}

	if v, ok := runner.(PreRunner); ok {
		c.preRunner = v
	}

	c.runner = runner

	if v, ok := runner.(PostRunner); ok {
		c.postRunner = v
	}

	if v, ok := runner.(PersistentPostRunner); ok {
		c.persistentPostRunner = v
	}
}

// parseArgs sets the values for flags and arguments, checks for unknown
// arguments, and that all required flags and arguments have been set.
func (c *Command) parseArgs(args []string) error {
	if err := c.setValues(args); err != nil {
		return err
	}

	if err := c.checkUnknown(args); err != nil {
		return err
	}

	return c.checkRequired()
}

func (c *Command) setValues(args []string) error {
	for i, arg := range args {
		if matchesFlag(arg, HelpFlag) {
			return ErrPrintHelp
		}

		if isFlag(arg) {
			flag := c.Flags.Lookup(arg)
			if flag == nil {
				continue
			}

			if err := flag.Init(); err != nil {
				return err
			}

			opt := flag.Options()
			switch opt.Value.(type) {
			case *bool:
				if err := flag.Set("true"); err != nil {
					return err
				}

				args = slice.Remove(args, i, i)
			default:
				if opt.IsSlice && opt.Separator == 0 {
					return ErrFlagSliceMustHaveSeparator
				}

				if err := flag.Set(args[i+1]); err != nil {
					return err
				}

				args = slice.Remove(args, i, i+2)
			}
		}
	}

	for i := range args {
		arg := c.Args.Lookup(i)
		if arg == nil {
			continue
		}

		if err := arg.Init(); err != nil {
			return err
		}

		buf := slice.Reduce(args[i:], func(a string) bool {
			return !isFlag(a)
		})

		if err := arg.Set(join.Args(buf)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) checkUnknown(args []string) error {
	seen := make(map[string]struct{})

	for _, cmd := range c.getCommands() {
		seen[cmd.Name] = struct{}{}
	}

	for _, flag := range c.Flags {
		opt := flag.Options()
		shorthand := fmt.Sprintf("-%s", opt.Shorthand)
		name := fmt.Sprintf("--%s", opt.Name)

		if opt.Shorthand != "" {
			seen[shorthand] = struct{}{}
		}

		if opt.Name != "" {
			seen[name] = struct{}{}
		}

		seen[flag.String()] = struct{}{}
	}

	for _, arg := range c.Args {
		seen[arg.String()] = struct{}{}
	}

	for _, arg := range args {
		if _, ok := seen[arg]; !ok {
			start, end := c.stmt.Lookup(arg).Pos()

			return ErrUnknown{
				Input:    c.stmt.String(),
				Arg:      arg,
				StartPos: start,
				EndPos:   end,
			}
		}
	}

	return nil
}

func (c *Command) checkRequired() error {
	var merr multierror.Error

	for _, flag := range c.Flags {
		opt := flag.Options()
		if opt.Required && (!opt.HasBeenSet || trimBrackets(flag) == "") {
			merr.Append(ErrFlagRequired{
				Name:      opt.Name,
				Shorthand: opt.Shorthand,
			})
		}
	}

	for _, arg := range c.Args {
		opt := arg.Options()
		if opt.Required && !opt.HasBeenSet {
			merr.Append(ErrArgRequired{
				Name: opt.Name,
			})
		}
	}

	return merr.ErrorOrNil()
}

// errOrPrintHelp checks if the error returned is ErrPrintHelp. If so, then the
// user intends to print help text to the command's output and not actually
// return an error.
func (c *Command) errOrPrintHelp(err error) error {
	if errors.Is(err, ErrPrintHelp) {
		c.PrintHelp()
		return nil
	}

	return err
}

func (c *Command) getCommands() []*Command {
	commands := make([]*Command, 0, len(c.commands))
	for _, cmd := range c.commands {
		commands = append(commands, cmd)
	}

	sort.Sort(SortCommandsByName(commands))

	return commands
}

func (c *Command) sortCommands() {
	commands := c.getCommands()
	m := make(map[string]*Command)
	for _, cmd := range commands {
		m[cmd.Name] = cmd
	}
	c.commands = m
}

// sortFlags sorts flags by name.
func (c *Command) sortFlags() {
	sort.Sort(SortFlagsByName(c.Flags))
}

// generateUsage generates usage strings for commands, flags, and arguments.
func (c *Command) generateUsage() {
	builder := help.NewBuilder()
	indent := 4
	padding := 4
	description := c.LongDesc
	if description == "" {
		description = formatDesc(c.Desc)
	}

	builder.Text(description)
	builder.Newline()
	builder.Newline()
	builder.Header("USAGE:")
	builder.Newline()
	builder.Text(builder.WithIndent(c.FullName(), 4))

	if len(c.Flags) > 0 {
		builder.Text(" [flags]")
	}

	for _, arg := range c.Args {
		opt := arg.Options()
		if opt.IsSlice {
			builder.Text(" <%s>...", opt.Name)
		} else {
			builder.Text(" <%s>", opt.Name)
		}
	}

	if len(c.commands) > 0 {
		builder.Text(" [command]")
	}

	if len(c.commands) > 0 {
		commands := tablewriter.NewWriter()

		builder.Newline()
		builder.Newline()
		builder.Header("COMMANDS:")
		builder.Newline()

		for _, cmd := range c.getCommands() {
			commands.AddLine(
				tablewriter.Cell{
					Indent:  indent,
					Padding: padding,
					Text:    builder.Green(cmd.Name),
				},
				tablewriter.Cell{
					Padding: padding,
					Text:    formatDesc(cmd.Desc),
				},
			)
		}

		builder.Table(commands)
	}

	if len(c.Flags) > 0 {
		flags := tablewriter.NewWriter()

		builder.Newline()
		builder.Newline()
		builder.Header("FLAGS:")
		builder.Newline()

		for _, flag := range c.Flags {
			opt := flag.Options()

			flags.AddLine(
				tablewriter.Cell{
					Indent: indent,
					Text:   builder.Green("-%s", opt.Shorthand),
					Suffix: ", ",
				},
				tablewriter.Cell{
					Padding: padding,
					Text:    builder.Green("--%s", opt.Name),
				},
				tablewriter.Cell{
					Text: formatDesc(opt.Desc),
				},
			)
		}

		builder.Table(flags)
	}

	if len(c.Args) > 0 {
		args := tablewriter.NewWriter()

		builder.Newline()
		builder.Newline()
		builder.Header("ARGS:")
		builder.Newline()

		for _, arg := range c.Args {
			opt := arg.Options()
			text := builder.Green("<%s>", opt.Name)
			if opt.IsSlice {
				text += builder.Green("...")
			}

			args.AddLine(
				tablewriter.Cell{
					Indent:  indent,
					Padding: padding,
					Text:    text,
				},
				tablewriter.Cell{
					Text: formatDesc(opt.Desc),
				},
			)
		}

		builder.Table(args)
	}

	if len(c.commands) > 0 {
		builder.Newline()
		builder.Newline()
		builder.Text("Use \"%s [command] --help\" for more information about a command.", c.FullName())
	}

	c.usage = builder.String()
}

// SortCommandsByName sorts commands by name.
type SortCommandsByName []*Command

func (n SortCommandsByName) Len() int      { return len(n) }
func (n SortCommandsByName) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n SortCommandsByName) Less(i, j int) bool {
	return strings.Map(unicode.ToUpper, n[i].Name) < strings.Map(unicode.ToUpper, n[j].Name)
}

// Execute parses args and sets up the root command and it's children.
func Execute(runner Runner, args []string) error {
	p := parser.New(args)

	cmd := runner.Init()
	cmd.setRunners(runner)
	cmd.init()
	cmd.stmt = p.Parse()
	cmd.output = os.Stdout

	if !cmd.HasFlag(HelpFlag.Name, HelpFlag.Shorthand) {
		cmd.Flags = append(cmd.Flags, HelpFlag)
	}

	return cmd.parseCommands(args[1:])
}
