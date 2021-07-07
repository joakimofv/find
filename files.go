package find

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

// Files searches for files matching the path patterns,
// where '**' is any number of folders and '*' is any substring.
// A file only has to match one of the patterns. If no patterns
// are given then all files under the current folder that contain
// the filter substring are taken.
//
// Note: This makes use of 'git grep' so only files that are added to git
// can be matched.
func Files(filter string, pathPatterns []string) []string {
	if len(pathPatterns) == 0 {
		if filter == "" {
			// All files in current folder.
			out, _ := exec.Command("ls", "-1").Output()
			return strings.Split(strings.Trim(string(out), "\n"), "\n")
		}

		if _, err := exec.LookPath("git"); err != nil {
			fmt.Printf("Missing git (error: %v)\n", err)
			return nil
		}
		// Do git-grep for files containing filter
		cmd := exec.Command("git", "grep", "-F", "-I", "--name-only", filter)
		stdErr := new(bytes.Buffer)
		cmd.Stderr = stdErr
		out, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
				// Not a git repository. Redo the command but without using the git index
				cmd = exec.Command("git", "grep", "-F", "-I", "--name-only", "--no-index", filter)
				stdErr.Reset()
				out, err = cmd.Output()
			}
		}
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				// Normal exit, just nothing found
			} else {
				fmt.Printf("git grep failure. StdOut: %v\nStdErr: %v\nError: %v\n", string(out), string(stdErr.Bytes()), err)
			}
			return nil
		}
		return strings.Split(strings.Trim(string(out), "\n"), "\n")
	}
	// Match the file patterns
	uniqueMatches := make(map[string]struct{})
	for _, pathPattern := range pathPatterns {
		matches, err := zglob.Glob(pathPattern)
		if err != nil {
			fmt.Printf("bad pattern '%v', error: %v\n", pathPattern, err)
			continue
		}
		for _, match := range matches {
			if abs, err := filepath.Abs(match); err == nil {
				uniqueMatches[abs] = struct{}{}
			}
		}
	}
	files := make([]string, len(uniqueMatches))
	i := 0
	for uniqueMatch, _ := range uniqueMatches {
		files[i] = uniqueMatch
		i += 1
	}
	return files
}
