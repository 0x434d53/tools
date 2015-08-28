package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func substituteFileAndBackup(oldPath, newPath string) error {
	bakDir := path.Dir(oldPath)
	tempNumber := fmt.Sprintf("%v", rand.Uint32())
	bakFileName := path.Base(oldPath) + tempNumber + ".bak"
	bakPath := path.Join(bakDir, bakFileName)

	err := os.Rename(oldPath, bakPath)

	if err != nil {
		return err
	}

	err = os.Rename(newPath, oldPath)

	if err != nil {
		err1 := os.Rename(bakPath, oldPath) // Try to revert the previous renaming
		if err1 != nil {
			panic(err1)
		}
		return err
	}

	return nil
}

func searchAndReplaceInFileWithBackup(p string, re string, newValue []byte) error {
	regex_compiled := regexp.MustCompile(re)

	sourceFile, err := os.Open(p)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	r := bufio.NewReader(sourceFile)
	buf := make([]byte, 2048)
	var result [][]byte

	for {

		buf, _, err = r.ReadLine()
		if err != nil {
			break
		}

		resultLine := regex_compiled.ReplaceAll(buf, newValue)
		resultLine = append(resultLine, '\n') //Append newlines!
		result = append(result, resultLine)
	}

	tempNumber := fmt.Sprintf("%v", rand.Uint32())
	tempFilePath := path.Join(path.Dir(p), "config"+tempNumber+".temp")
	tempFile, err := os.Create(tempFilePath)

	if err != nil {
		return err
	}

	defer tempFile.Close()

	sourceFileInfo, err := sourceFile.Stat()

	if err != nil {
		return err
	}

	tempFile.Chmod(sourceFileInfo.Mode())

	for _, line := range result {
		_, err := tempFile.Write(line)

		if err != nil {
			return err
		}
	}

	sourceFile.Close()

	return substituteFileAndBackup(p, tempFilePath)
}

func isGitDirectory(p string, info os.FileInfo) bool {
	b, err := path.Match(".git", info.Name())

	if err != nil {
		return false
	}

	if !b {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

func walker(p string, info os.FileInfo, err error) error {
	if isGitDirectory(p, info) {
		configFile := path.Join(p, "config")

		re := `https://github.com/`
		newValue := []byte(`git@github.com:`)

		if err = searchAndReplaceInFileWithBackup(configFile, re, newValue); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a root path as an argument")
		return
	}

	filepath.Walk(os.Args[1], walker)
}
