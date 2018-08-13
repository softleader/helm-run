package main

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"context"
	"io/ioutil"
	"os"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/google/go-github/github"
	"github.com/docker/docker/api/types/mount"
	"path"
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	defaultDirectoryPermission = 0755
	commandOwner               = "softleader"
	commandRepo                = "dockerfile"
	commandPathRoot            = "helm"
	workDir                    = "/data"
)

type runCmd struct {
	pwd             string
	command         string
	args            []string
	image           string
	alwaysPullImage bool
	rm              bool
	entryPoint      []string
	local           bool
}

func (cmd *runCmd) run() error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if cmd.alwaysPullImage {
		_, err = cli.ImagePull(ctx, cmd.image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
	}

	if cmd.local {
		// verify local command file exists
		_, err := os.Stat(cmd.command)
		if os.IsNotExist(err) {
			return fmt.Errorf("command '%s' does not exist", cmd.command)
		}
	} else {
		shell, err := getCommandContents(cmd.command)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(cmd.pwd, cmd.command), []byte(shell), defaultDirectoryPermission)
		if err != nil {
			return err
		}
		defer os.Remove(path.Join(cmd.pwd, cmd.command))
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:      cmd.image,
			Entrypoint: cmd.strSlice(),
			WorkingDir: workDir,
			Cmd:        append([]string{"./" + cmd.command}, cmd.args...),
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: cmd.pwd,
					Target: workDir,
				},
			},
		}, nil, "")
	if err != nil {
		return err
	}
	if cmd.rm {
		defer func(containerID string) {
			cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
				RemoveLinks:   true,
				RemoveVolumes: true,
				Force:         true,
			})
		}(resp.ID)
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

func (cmd *runCmd) strSlice() strslice.StrSlice {
	slice := make(strslice.StrSlice, len(cmd.entryPoint))
	for i, v := range cmd.entryPoint {
		slice[i] = v
	}
	return slice
}

func getCommandContents(command string) (contents string, err error) {
	gc := github.NewClient(nil)
	fileContent, _, _, err := gc.Repositories.GetContents(context.Background(), commandOwner, commandRepo, commandPathRoot+"/"+command, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get command: %s", err.Error())
	}
	return fileContent.GetContent()
}
