package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Git git commands
type Git struct {
}

func executeCommand(c string, args []string, stdin []byte) ([]byte, int, error) {
	cmd := exec.Command(c, args...)
	if len(stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(stdin)
	}

	output, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()

	return output, exitCode, err
}

func (g Git) do(args []string) ([]byte, int, error) {
	return executeCommand("git", args, nil)
}

// GitLog information
type GitLog struct {
	BranchName string
	Commit     string
	Subject    string
}

// GetGitLog get git log info
func (g Git) GetGitLog() (GitLog, error) {
	commit, err := g.CurrentBranchCommit()
	if err != nil {
		return GitLog{}, err
	}

	branchName, err := g.CurrentBranchName()
	if err != nil {
		return GitLog{}, err
	}

	subject, err := g.BranchSubject()
	if err != nil {
		return GitLog{}, err
	}

	return GitLog{
		BranchName: branchName,
		Commit:     commit,
		Subject:    subject,
	}, nil
}

func (g Git) Exist() bool {
	_, exitCode, err := executeCommand("which", []string{"git"}, nil)
	if err != nil {
		return false
	}

	return exitCode == 0
}

// CurrentBranchName current branch name, git version should >= 1.6.3
func (g Git) CurrentBranchName() (string, error) {
	output, _, err := g.do([]string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(output), "\r\n"), nil
}

// CurrentBranchCommit get current Branch Commit
func (g Git) CurrentBranchCommit() (string, error) {
	output, _, err := g.do([]string{"rev-parse", "HEAD"})
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(output), "\r\n"), nil
}

// BranchSubject branch subject
func (g Git) BranchSubject() (string, error) {
	output, _, err := g.do([]string{"show", "--format=\"%s\"", "--no-patch"})
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(output), "\r\n"), nil
}

// NewBranch new branch
func (g Git) NewBranch(name string) error {
	output, _, err := g.do([]string{"checkout", "-b", name})
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

// SwitchToBranch switch to Branch
func (g Git) SwitchToBranch(name string) error {
	output, _, err := g.do([]string{"checkout", name})
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

// DeleteBranch delete branch
func (g Git) DeleteBranch(name string) error {
	output, _, err := g.do([]string{"branch", "-d", name})
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

// ListBranches list all branch
func (g Git) ListBranches() error {
	output, _, err := g.do([]string{"branch"})
	if err != nil {
		return errors.New(string(output))
	}

	fmt.Printf("%s",string(output))
	return nil
}