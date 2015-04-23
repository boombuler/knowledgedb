package main

import (
	"github.com/boombuler/knowledgedb/config"
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/markdown"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

var templates map[string]*template.Template
var mdContent map[string]string = map[string]string{
	"help": "help.md",
}

func initTemplates() {
	templateDir := path.Join(config.ApplicationDir, "templates")

	html := []string{
		"markdown", // needed for "renderMarkdown"

		"edit_document",
		"view_document",
		"create_document",
		"search_result",
		"index",
		"tag",
		"licenses",
	}

	myFuncs := template.FuncMap{
		"fmtDateTime": func(t time.Time) string {
			return t.Format("02.01.2006 15:04:05")
		},
		"markdown": func(text string) template.HTML {
			return markdown.ToHtml(text)
		},
	}

	templates = make(map[string]*template.Template)

	addTemplate := func(name string, files ...string) {
		fList := make([]string, len(files))
		for i, f := range files {
			fList[i] = path.Join(templateDir, f)
		}
		templates[name] = template.Must(template.New(name).Funcs(myFuncs).ParseFiles(fList...))
	}

	for _, t := range html {
		addTemplate(t, "base.html", t+".html")
	}

	for k, v := range mdContent {
		f, err := os.Open(path.Join(templateDir, v))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		mdContent[k] = string(content)
	}

}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	tmplInstance, ok := templates[tmpl]
	if !ok {
		log.Fatal("Invalid Template:", tmpl)
	}
	err := tmplInstance.ExecuteTemplate(w, "base.html", p)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template Error", err)
	}
}

func renderMarkdown(w http.ResponseWriter, contentName string) {
	md, ok := mdContent[contentName]
	if !ok {
		log.Fatal("Invalid Pagename:", contentName)
	}
	renderTemplate(w, "markdown", md)
}
