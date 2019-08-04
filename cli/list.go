package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

func listHistory(ctx *cli.Context) error {
	historyRecords, err := history.GetProjectBranchSwitchRecords()
	if err != nil {
		return err
	}

	for _, h := range historyRecords {
		fmt.Println(h.Normalize())
	}

	return nil
}


// remove history information
func removeHistory(ctx *cli.Context) error {
	return history.RemoveProjectBranchSwitchRecords()
}