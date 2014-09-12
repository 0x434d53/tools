package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func ExtractURL(p string) (string, error) {
	regexSSH := `(?P<ssh>git@github.com:.+/.+)`
	regexHTTP := `(?P<http>Phttps?://github.com/.+/.+)`
	rcSSH := regexp.MustCompile(regexSSH)
	rcHTTP := regexp.MustCompile(regexHTTP)

	f, err := os.Open(path.Join(p, "config"))
	if err != nil {
		return "", err
	}

	defer f.Close()

	r := bufio.NewReader(f)

	buf := make([]byte, 2048)
	for {
		buf, _, err = r.ReadLine()
		if err != nil {
			break
		}
		m := rcSSH.FindStringSubmatch(string(buf))
		if len(m) == 2 {
			return m[1], nil
		}

		m = rcHTTP.FindStringSubmatch(string(buf))
		if len(m) == 2 {
			return m[1], nil
		}
	}

	return "", fmt.Errorf("%v", p)
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

func CollectGitRepositories(root string) ([]string, error) {
	dirs := make([]string, 0)

	walker := func(p string, info os.FileInfo, err error) error {
		if IsGitDirectory(p, info) {
			path, err := ExtractURL(p)

			if err != nil {
				fmt.Printf("Error extracting Url: %v", err)
				return nil
			}

			dirs = append(dirs, path)
		}

		return nil
	}

	filepath.Walk(root, walker)

	return dirs, nil
}

func WriteToFile(p string, repos []string) error {
	f, err := os.Create(path.Join(p, "repos.txt"))

	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, p := range repos {
		_, err = f.WriteString(p + "\n")

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a root path as an argument")
		return
	}

	dirs, err := CollectGitRepositories(os.Args[1])

	if err != nil {
		panic(err)
	}
	err = WriteToFile(os.Args[1], dirs)

	if err != nil {
		fmt.Println(err)
	}
}
