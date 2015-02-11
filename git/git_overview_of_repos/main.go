package main

import (
	"fmt"
	"github.com/0x434D53/tools/git/lib"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type GitRepo struct {
	lib.GitInfos
	ReadMeRendered []byte
	LastUpdated    time.Time
	URL            string
}

type ProjectAndId struct {
	Id          int
	Projectname string
}

type templateData struct {
	Project            []ProjectAndId
	CurrentReadme      []byte
	CurrentProjectname string
	CurrentLastUpdate  string
	CurrentUser        string
	CurrentURL         string
}

var gitRepos []GitRepo
var servePath string

func CollectRenderingInformation(gi []lib.GitInfos) ([]GitRepo, error) {
	gitRepos := make([]GitRepo, 0)

	for _, g := range gi {
		gr := GitRepo{}

		gr.Path = g.Path
		gr.Projectname = g.Projectname
		gr.Username = g.Username

		gr.ReadMeRendered = RenderReadme(gr.Path)
		gr.LastUpdated = GetLastUpdate(gr.Path)
		gr.URL = GetURL(g)
		gitRepos = append(gitRepos, gr)
	}

	return gitRepos, nil
}

func GetLastUpdate(path string) time.Time {
	return time.Now()
}

func GetURL(gi lib.GitInfos) string {
	return ""
}

func RenderReadme(p string) []byte {
	file, err := findReadme(p)

	if err != nil {
		return []byte("This project has no README")
	}

	ext := path.Ext(file)

	f, err := os.Open(file)

	if err != nil {
		return []byte("This project has no valid README")
	}

	contents, err := ioutil.ReadAll(f)

	if ext == ".md" || ext == ".MD" {
		return blackfriday.MarkdownCommon(contents)
	}

	return contents
}

func findReadme(p string) (string, error) {
	filestoCheck := []string{"README.md", "readme.md", "Readme.md", "readme.txt", "README.markdown"}

	for _, f := range filestoCheck {
		_, err := os.Stat(path.Join(p, f))

		if err == nil {
			return path.Join(p, f), nil
		}
	}

	return "", fmt.Errorf("No Readme found")
}

func ServeWebSite(servePath string) {
	log.Panic(http.ListenAndServe(":8081", http.FileServer(http.Dir(servePath))))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a root path as an argument")
		return
	}

	gitInfos, err := lib.CollectGitRepositories(os.Args[1])
	if err != nil {
		panic(err)
	}

	repos, err := CollectRenderingInformation(gitInfos)

	if err != nil {
		panic(err)
	}

	for _, r := range repos {
		fmt.Println(r.Projectname)
	}

	//	ServeWebSite(servePath)
}
