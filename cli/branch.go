package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli"
)

func wrapperBranch(branchName string, realDo func(string) error) error {

	git := Git{}
	gitLog, err := git.GetGitLog()
	if err != nil {
		return err
	}

	err = realDo(branchName)
	if err != nil {
		return err
	}

	record := HistoryRecord{
		FromName:    gitLog.BranchName,
		ToName:      branchName,
		Occur:       time.Now(),
		FromCommit:  gitLog.Commit,
		FromSubject: gitLog.Subject,
	}

	_, err = history.AddNewHistoryRecord(record)
	return err

}

func switchBranch(ctx *cli.Context) error {
	args := ctx.Args()

	if len(args) != 1 {
		return errors.New("switch branch must given the branch name")
	}
	branchName := args[0]

	return wrapperBranch(branchName, func(branchName string) error {
		git := Git{}
		return git.SwitchToBranch(branchName)
	})
}

func switchBranchComplete(ctx *cli.Context) {
	historyRecords, err := history.GetAllHistoryRecord()
	if err != nil {
		return
	}

	for _, h := range historyRecords {
		fmt.Println(h.Normalize())
	}
}

func newBranch(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		fmt.Println("new branch should given name")
		return nil
	}

	branchName := args[0]

	return wrapperBranch(branchName, func(branchName string) error {
		git := Git{}
		return git.NewBranch(branchName)
	})

}

func appRun(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) == 0 {
		git := Git{}
		branch, err := git.CurrentBranchName()
		if err != nil {
			return err
		}
		fmt.Printf("current branch: %s\n", branch)
		return nil
	} else {
		return switchBranch(ctx)
	}
}

func deleteBranch(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		fmt.Println("branch name should be given")
		return nil
	}

	branchName := args[0]
	git := Git{}
	err := git.DeleteBranch(branchName)
	if err != nil {
		return err
	}

	return nil
	//return history.RemoveHistoryRecord(branchName)
}

func listBranches(ctx *cli.Context) error {
	git := Git{}
	return git.ListBranches()
}
