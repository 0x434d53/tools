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

func SubstituteFileAndBackup(oldPath, newPath string) error {
	bakDir := path.Dir(oldPath)
	temp_number := fmt.Sprintf("%v", rand.Uint32())
	bakFileName := path.Base(oldPath) + temp_number + ".bak"
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

func SearchAndReplaceInFileWithBackup(p string, re string, newValue []byte) error {
	regex_compiled := regexp.MustCompile(re)

	sourceFile, err := os.Open(p)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	r := bufio.NewReader(sourceFile)
	buf := make([]byte, 2048)
	result := make([][]byte, 0)

	for {

		buf, _, err = r.ReadLine()
		if err != nil {
			break
		}

		result_line := regex_compiled.ReplaceAll(buf, newValue)
		result_line = append(result_line, '\n') //Append newlines!
		result = append(result, result_line)
	}

	temp_number := fmt.Sprintf("%v", rand.Uint32())
	tempFilePath := path.Join(path.Dir(p), "config"+temp_number+".temp")
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

	return SubstituteFileAndBackup(p, tempFilePath)
}

func IsGitDirectory(p string, info os.FileInfo) bool {
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
	if IsGitDirectory(p, info) {
		configFile := path.Join(p, "config")

		re := `https://github.com/`
		newValue := []byte(`git@github.com:`)

		if err = SearchAndReplaceInFileWithBackup(configFile, re, newValue); err != nil {
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
