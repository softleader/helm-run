package main

import (
	"testing"
)

func TestRun(t *testing.T) {
//	wd, err := os.Getwd()
//	if err != nil {
//		t.Error(err)
//	}
//	tmp, err := ioutil.TempDir(wd, "helm-run")
//	if err != nil {
//		t.Error(err)
//	}
//	defer os.RemoveAll(tmp)
//
//	runCmd := runCmd{
//		command: path.Join(strings.Replace(tmp, wd+"/", "", -1), "hello"),
//		args:    []string{"Matt"},
//		image:   "softleader/helm",
//		pwd:     wd,
//		local:   true,
//	}
//
//	ioutil.WriteFile(path.Join(tmp, "hello"), []byte(`
//echo "hello $@"
//`), defaultDirectoryPermission)
//
//	err = runCmd.run()
//	if err != nil {
//		t.Error(err)
//	}
}

func TestGetCommandContents(t *testing.T) {
	//c, err := getCommandContents("package")
	//if err != nil {
	//	t.Error(err)
	//}
	//if !strings.Contains(c, "#!/usr/bin/env bash") {
	//	t.Error("package should be a bash command")
	//}
}
