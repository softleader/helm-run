package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	defaultDirectoryPermission = 0755
	owner                      = "softleader"
	repo                       = "dockerfile"
	pathBase                   = "helm"
	workDir                    = "/data"
	image                      = "softleader/helm"
	entrypoint                 = "/bin/bash"
)

type runCmd struct {
	pwd         string
	owner       string
	repo        string
	pathBase    string
	token       string
	command     string
	args        []string
	updateImage bool
	rm          bool
	local       bool
	dos2unix    bool
	make        bool
}

func (cmd *runCmd) run() error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if cmd.updateImage {
		out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer out.Close()
		termFd, isTerm := term.GetFdInfo(os.Stderr)
		jsonmessage.DisplayJSONMessagesStream(out, os.Stdout, termFd, isTerm, nil)
	}

	if cmd.local {
		_, err := os.Stat(cmd.command)
		if os.IsNotExist(err) {
			return fmt.Errorf("command '%s' does not exist", cmd.command)
		}
		cmd.command = filepath.ToSlash(cmd.command)
	} else {
		shell, err := cmd.getCommandContents()
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
			Image:      image,
			Entrypoint: []string{entrypoint},
			WorkingDir: workDir,
			Cmd:        cmd.cmd(),
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
	defer out.Close()
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}

func (cmd *runCmd) cmd() strslice.StrSlice {
	if cmd.make {
		c := []string{fmt.Sprintf("make -f ./%s", cmd.command)}
		c = append(c, cmd.args...)
		return []string{"-c", strings.Join(c, " ")}
	} else if cmd.dos2unix {
		c := []string{fmt.Sprintf("cat ./%s | dos2unix | bash", cmd.command)}
		if len(cmd.args) > 0 {
			c = append(c, "-s")
			c = append(c, cmd.args...)
		}
		return []string{"-c", strings.Join(c, " ")}
	} else {
		c := []string{fmt.Sprintf("./%s", cmd.command)}
		c = append(c, cmd.args...)
		return c
	}
}

func (cmd *runCmd) getCommandContents() (contents string, err error) {
	ctx := context.Background()
	var client *github.Client
	if cmd.token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cmd.token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}
	path := cmd.command
	if cmd.pathBase != "" {
		path = cmd.pathBase + "/" + cmd.command
	}
	fileContent, _, _, err := client.Repositories.GetContents(ctx, cmd.owner, cmd.repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get command: %s", err.Error())
	}
	return fileContent.GetContent()
}
