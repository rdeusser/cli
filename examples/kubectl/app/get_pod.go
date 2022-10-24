package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type GetPodCommand struct {
	Debug     bool
	Names     []string
	Namespace string
}

func (gpc *GetPodCommand) Init() *cli.Command {
	cmd := &cli.Command{
		Name: "pod",
		Desc: "Get a pod from Kubernetes",
		Args: cli.Args{
			&cli.Arg[[]string]{
				Name:     "names",
				Desc:     "Pod name",
				Value:    &gpc.Names,
				Required: true,
			},
		},
	}

	flag := cmd.Flags.Lookup("namespace")
	if flag != nil {
		gpc.Namespace = flag.String()
	}

	return cmd
}

// SetOptions is how we get flag values from parent commands. The root command
// has a debug flag. We don't want to copy and paste that flag to every command
// nor do we want to take that flag and put it in a different package so all the
// other commands can import it. We simply need the value. SetOptions lets us
// get that value.
func (gpc *GetPodCommand) SetOptions(flags cli.Flags) error {
	gpc.Debug = cli.ValueOf[bool](flags, "debug")
	gpc.Namespace = cli.ValueOf[string](flags, "namespace")

	return nil
}

func (gpc *GetPodCommand) Run() error {
	for _, name := range gpc.Names {
		fmt.Printf(`---
apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  containers:
  - name: %s
    image: %s:latest
    imagePullPolicy: Always
    ports:
    - containerPort: 8080
      name: http
`, name, gpc.Namespace, name, name, name)
	}

	return nil
}
