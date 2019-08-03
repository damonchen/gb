package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/damonchen/gb/build"
	"github.com/urfave/cli"
)

var (
	history *History
)

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
			fmt.Fprintf(os.Stderr,"open history error %s\n", err)
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
