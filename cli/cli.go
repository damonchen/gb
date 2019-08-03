package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/damonchen/gb/build"
	"github.com/urfave/cli"
)

var (
	history     *History
	projectPath string
)

func getGitProjectPath(path string) string {
	isGitDir := func(path string) bool {
		gitDir := filepath.Join(path, ".git")
		if stat, err := os.Stat(gitDir); err != nil {
			if os.IsExist(err) {
				return stat.IsDir()
			}
		}
		return false
	}

	testPath := path
	for {
		if isGitDir(testPath) {
			return testPath
		}

		if testPath == "/" {
			break
		}

		testPath = filepath.Dir(testPath)
	}
	return ""
}

// Run start command
func Run() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "gb"
	app.Usage = "Git Branch History Manager"
	app.Version = build.Version()
	app.Copyright = "Copyright (c) 2019, 2019, damonchen. All rights reserved."
	app.Authors = []cli.Author{
		{
			Name:  "damonchen",
			Email: "netubu@gmail.com",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		// search current working directory is in git project
		workDir, err := os.Getwd()
		if err != nil {
			return err
		}
		projectPath = getGitProjectPath(workDir)
		if len(projectPath) == 0 {
			return errors.New("current directory is not in git project")
		}

		// check git exists
		git := Git{}
		if !git.Exist() {
			return errors.New("git should exist in path")
		}

		rootDir := getHistoryRoot()
		if err := os.MkdirAll(rootDir, 0755); err != nil {
			return err
		}

		historyFile := getHistoryFileName()
		h, err := NewHistory(historyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open history error %s\n", err)
			return err
		}

		history = h
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		if history != nil {
			return history.Close()
		}
		return nil
	}

	app.Action = appRun
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[gb] %s\n", err.Error())
		os.Exit(1)
	}

}
