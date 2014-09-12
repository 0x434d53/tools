package lib

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type GitInfos struct {
	Path        string
	Username    string
	Projectname string
}

func ExtractUserAndProject(p string) (string, string, error) {
	regexSSH := `git@github.com:(?P<username>.+)/(?P<projectname>.+)`
	regexHTTP := `https?://github.com/(?P<username>.+)/(?P<projectname>.+)`
	rcSSH := regexp.MustCompile(regexSSH)
	rcHTTP := regexp.MustCompile(regexHTTP)

	f, err := os.Open(path.Join(p, "config"))
	if err != nil {
		return "", "", err
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
		if len(m) == 3 {
			return m[1], m[2], nil
		}

		m = rcHTTP.FindStringSubmatch(string(buf))
		if len(m) == 3 {
			return m[1], m[2], nil
		}
	}

	return "", "", fmt.Errorf("Not Found")
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

func CollectGitRepositories(p string) ([]GitInfos, error) {
	gitInfos := make([]GitInfos, 0)

	walker := func(p string, info os.FileInfo, err error) error {
		if IsGitDirectory(p, info) {
			username, projectname, err := ExtractUserAndProject(p)

			if err != nil {
				fmt.Printf("Error extracting Username and Project: %v", err)
				return nil
			}

			gi := GitInfos{Path: path.Dir(p), Username: username, Projectname: projectname}
			gitInfos = append(gitInfos, gi)
		}

		return nil
	}

	filepath.Walk(p, walker)

	return gitInfos, nil
}
