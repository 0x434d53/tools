package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/0x434D53/tools/git/lib"
	"github.com/russross/blackfriday"
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

var projects []ProjectAndId

type TemplateData struct {
	Projects           []ProjectAndId
	CurrentReadme      template.HTML
	CurrentProjectname string
	CurrentLastUpdate  string
	CurrentUser        string
	CurrentURL         string
}

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

var repos []GitRepo
var tmpl *template.Template

func Serve() {
	var err error
	tmpl, err = template.ParseFiles("template.html")

	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", MainHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	data := TemplateData{
		Projects:           projects,
		CurrentProjectname: repos[0].Projectname,
		CurrentReadme:      template.HTML(repos[0].ReadMeRendered),
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
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

	repos, err = CollectRenderingInformation(gitInfos)

	projects = []ProjectAndId{}

	for i, gi := range repos {
		p := ProjectAndId{}
		p.Projectname = gi.Projectname
		p.Id = i

		projects = append(projects, p)
	}

	if err != nil {
		panic(err)
	}

	for _, r := range repos {
		fmt.Println(r.Projectname)
	}

	Serve()
}
