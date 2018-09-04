package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tmp, err := ioutil.TempDir(wd, "helm-run")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmp)

	runCmd := runCmd{
		dos2unix: true,
		command:  path.Join(strings.Replace(tmp, wd+"/", "", -1), "hello"),
		args:     []string{"Matt"},
		pwd:      wd,
		local:    true,
	}

	ioutil.WriteFile(path.Join(tmp, "hello"), []byte(`
	echo "hello $@"
	`), defaultDirectoryPermission)

	err = runCmd.run()
	if err != nil {
		t.Error(err)
	}
}

func TestRun_Make(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tmp, err := ioutil.TempDir(wd, "helm-run")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmp)

	runCmd := runCmd{
		dos2unix: true,
		command:  path.Join(strings.Replace(tmp, wd+"/", "", -1), "hello"),
		args:     []string{"NAME=Matt"},
		pwd:      wd,
		local:    true,
		make:     true,
	}

	ioutil.WriteFile(path.Join(tmp, "hello"), []byte(`
NAME := anonymous
all: hello my_name
hello:
	echo "hello"
my_name:
	echo "$(NAME)"
	`), defaultDirectoryPermission)

	err = runCmd.run()
	if err != nil {
		t.Error(err)
	}
}

func TestGetCommandContents(t *testing.T) {
	//runCmd := runCmd{
	//	command:         "package",
	//	commandOwner:    commandOwner,
	//	commandRepo:     commandRepo,
	//	commandPathBase: commandPathBase,
	//}
	//c, err := runCmd.getCommandContents()
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(c)
	//if !strings.Contains(c, "#!/usr/bin/env bash") {
	//	t.Error("package should be a bash command")
	//}
}
