package find

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

// Files searches for files matching the path patterns,
// where '**' is any number of folders and '*' is any substring.
// A file only has to match one of the patterns. If no patterns
// are given then all files under the current folder that contain
// the filter substring are taken.
//
// Note: filter matching uses 'git grep' so only files that are added to git
// can be found by that method.
func Files(filter string, pathPatterns []string) []string {
	if len(pathPatterns) == 0 {
		if filter == "" {
			// All files in current folder.
			out, _ := exec.Command("ls", "-1").Output()
			return strings.Split(strings.Trim(string(out), "\n"), "\n")
		}

		if _, err := exec.LookPath("git"); err != nil {
			fmt.Fprintf(os.Stderr, "Missing git (error: %v)\n", err)
			return nil
		}
		// Do git-grep for files containing filter
		cmd := exec.Command("git", "grep", "-F", "-I", "--name-only", filter)
		stdErr := new(bytes.Buffer)
		cmd.Stderr = stdErr
		out, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
				// Not a git repository.
				/*
					cmd = exec.Command("git", "grep", "-F", "-I", "--name-only", "--no-index", filter)
					stdErr.Reset()
					out, err = cmd.Output()
				*/
				fmt.Fprintln(os.Stderr, stdErr.String())
				return nil
			}
		}
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				// Normal exit, just nothing found
			} else {
				fmt.Fprintf(os.Stderr, "git grep failure. StdOut: %v\nStdErr: %v\nError: %v\n", string(out), string(stdErr.Bytes()), err)
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
			fmt.Fprintf(os.Stderr, "bad pattern '%v', error: %v\n", pathPattern, err)
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
	sort.Strings(files)
	return files
}
