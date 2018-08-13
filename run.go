package main

import (
	"errors"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"context"
	"io/ioutil"
	"os"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/google/go-github/github"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"path"
	"strings"
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
}

func (cmd *runCmd) run() (err error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		errors.New("docker is required: " + err.Error())
	}

	if cmd.alwaysPullImage {
		_, err = cli.ImagePull(ctx, cmd.image, types.ImagePullOptions{})
		if err != nil {
			return
		}
	}

	shell, err := getCommandContents(cmd.command)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(cmd.pwd, cmd.command), []byte(shell), defaultDirectoryPermission)
	if err != nil {
		return
	}
	defer os.Remove(path.Join(cmd.pwd, cmd.command))

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:      cmd.image,
			Entrypoint: strslice.StrSlice{"bash"},
			WorkingDir: workDir,
			Cmd:        []string{"-c", strings.Join(append([]string{"./" + cmd.command}, cmd.args...), " ")},
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
		return
	}
	defer func(containerID string) {
		cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			RemoveLinks:   true,
			RemoveVolumes: true,
			Force:         true,
		})
	}(resp.ID)

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return
}

func getCommandContents(command string) (contents string, err error) {
	gc := github.NewClient(nil)
	fileContent, _, _, err := gc.Repositories.GetContents(context.Background(), commandOwner, commandRepo, commandPathRoot+"/"+command, nil)
	if err != nil {
		return
	}
	return fileContent.GetContent()
}
