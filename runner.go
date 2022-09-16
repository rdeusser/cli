package cli

// PersistentPreRunner is a runner that each command starting from the parent
// will run. This runner is always run and is the first to run.
type PersistentPreRunner interface {
	PersistentPreRun() error
}

// PreRunner is a runner that will run before the main runner for the command.
type PreRunner interface {
	PreRun() error
}

// Runner is the main runner that every command will implement.
type Runner interface {
	Init() *Command
	Run() error
}

// PostRunner is a runner that will run after the main runner for the command.
type PostRunner interface {
	PostRun() error
}

// PersistentPostRunner is a runner that each command starting from the parent
// will run. This runner is always run and is the last to run.
type PersistentPostRunner interface {
	PersistentPostRun() error
}
