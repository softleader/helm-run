package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const desc = `
Run bash command in container-based environment.

To run commands that store on https://github.com/softleader/dockerfile/tree/master/helm
	
	$ helm run package

To run command on local, use '-l':

	$ helm run -l path/to/script

To execute via GNU make utility, not bash, use '-m':

	$ helm run -lm path/to/Makefile
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
	f.BoolVarP(&runCmd.updateImage, "update-image", "U", false, "update image before running command")
	f.BoolVarP(&runCmd.rm, "rm", "", true, "automatically remove the container when it exits")
	f.BoolVarP(&runCmd.local, "local", "l", false, "command store on local storage, not on github")
	f.BoolVarP(&runCmd.make, "make", "m", false, "executed via GNU make utility, not bash")
	f.BoolVarP(&runCmd.dos2unix, "dos2unix", "", true, "convert FILE from DOS to Unix format")
	f.StringVarP(&runCmd.owner, "owner", "", owner, "github owner of command")
	f.StringVarP(&runCmd.repo, "repo", "", repo, "github repository of command")
	f.StringVarP(&runCmd.pathBase, "path-base", "", pathBase, "github path base of command")
	f.StringVarP(&runCmd.token, "token", "", "", "github access token of command for private repositories")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
