package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
)

const (
	commonPrefix      = "├───"
	lastPrefix        = "└───"
	commonLevelPrefix = "│\t"
	lastLevelPrefix   = "\t"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDirTree(out, path, printFiles, "")
}

func printDirTree(out io.Writer, path string, printFiles bool, levelPrefix string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", path, err)
	}
	defer file.Close()

	dir, err := file.Readdir(-1)
	if err != nil {
		return fmt.Errorf("dir reading failed: %w", err)
	}

	if !printFiles {
		dir = excludeFiles(dir)
	}

	sort.SliceStable(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})

	for idx, f := range dir {
		prefix := commonPrefix
		levelPrefixAddition := commonLevelPrefix
		if idx == len(dir)-1 {
			prefix = lastPrefix
			levelPrefixAddition = lastLevelPrefix
		}

		if f.IsDir() {
			if _, err = fmt.Fprintf(out, "%s%s%s\n", levelPrefix, prefix, f.Name()); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
			newPath := path + "/" + f.Name()
			newLevelPrefix := levelPrefix + levelPrefixAddition
			if err = printDirTree(out, newPath, printFiles, newLevelPrefix); err != nil {
				return err
			}
			continue
		}

		size := fmt.Sprintf("%db", f.Size())
		if f.Size() == 0 {
			size = "empty"
		}

		if _, err = fmt.Fprintf(out, "%s%s%s (%s)\n", levelPrefix, prefix, f.Name(), size); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}
	return nil
}

func excludeFiles(dir []fs.FileInfo) []fs.FileInfo {
	filteredDir := []fs.FileInfo{}
	for _, f := range dir {
		if !f.IsDir() {
			continue
		}
		filteredDir = append(filteredDir, f)
	}
	return filteredDir
}
