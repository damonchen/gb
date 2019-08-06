package cli

import "github.com/urfave/cli"

var (
	commands = []cli.Command{
		{
			Name:      "h",
			Aliases:   []string{"history"},
			Usage:     "List history information",
			UsageText: "gb h",
			Action:    listHistory,
			Subcommands: []cli.Command{
				{
					Name:      "rm",
					Aliases:   []string{"remove"},
					Usage:     "remove all history information, dangerous",
					UsageText: "gb h remove",
					Action:    removeHistory,
				},
			},
		},
		{
			Name:      "n",
			Aliases:   []string{"new", "create"},
			Usage:     "new git branch",
			UsageText: "gb n branch-name",
			Action:    newBranch,
		},
		{
			Name:         "s",
			Aliases:      []string{"switch", "checkout"},
			Usage:        "switch to another branch",
			UsageText:    "gb s branch-name",
			Action:       switchBranch,
			BashComplete: switchBranchComplete,
		},
		{
			Name:      "d",
			Aliases:   []string{"delete", "del"},
			Usage:     "delete branch",
			UsageText: "gb d branch-name",
			Action:    deleteBranch,
		},
		{
			Name:      "l",
			Aliases:   []string{"list", "ls"},
			Usage:     "List local branch",
			UsageText: "gb ls",
			Action:    listBranches,
		},		{
			Name:      "p",
			Aliases:   []string{"push"},
			Usage:     "push local branch to remote branch",
			UsageText: "gb p",
			Action:    pushBranch,
		},
	}
)
