package main

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

const desc = `
Run command in container-based environment, commands store on https://github.com/softleader/dockerfile/tree/master/helm 
	$ helm run package"
`

func main() {
	runCmd := runCmd{}

	cmd := &cobra.Command{
		Use:   "helm run [flags] COMMAND [ARGS]",
		Short: fmt.Sprintf("run command in container-based environment"),
		Long:  desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("command is required")
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			runCmd.pwd = wd
			runCmd.command = args[0]
			runCmd.args = args[1:]
			return runCmd.run()
		},
	}
	f := cmd.Flags()
	f.StringVarP(&runCmd.image, "image", "i", "softleader/helm", "image for running command")
	f.BoolVarP(&runCmd.alwaysPullImage, "always-pull-image", "", false, "always pull image before running command")
	f.BoolVarP(&runCmd.rm, "rm", "", true, "automatically remove the container when it exits")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
