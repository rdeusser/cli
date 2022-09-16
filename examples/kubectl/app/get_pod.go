package app

import (
	"fmt"

	"github.com/rdeusser/cli"
)

type GetPodCommand struct {
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

func (gpc *GetPodCommand) SetOptions(flags cli.Flags) error {
	namespace := flags.Lookup("namespace")
	if namespace != nil {
		gpc.Namespace = namespace.String()
	}

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
