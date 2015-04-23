package main

import (
	"github.com/boombuler/knowledgedb/index"
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/storage"
	"html/template"
	"net/http"
	tt "text/template"
)

type SearchResult struct {
	Id    string
	Title template.HTML
	Texts []template.HTML
	Tags  []string
}

func (s *httpServer) serveSearchEntries(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res, err := index.Search(r.Form["query"][0])
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Warning(err)
		return
	}
	searchRes := make([]*SearchResult, 0, len(res.Hits))
	for _, hit := range res.Hits {
		entry := storage.LoadEntry(hit.ID)
		if entry != nil {
			title := template.HTML(tt.HTMLEscapeString(entry.Title))
			if thit, ok := hit.Fragments["Title"]; ok {
				title = template.HTML(thit[0])
			}

			searchHits, ok := hit.Fragments["Body"]
			result := &SearchResult{
				Id:    hit.ID,
				Title: title,
				Tags:  entry.Tags,
			}
			if ok {
				hitTexts := make([]template.HTML, 0, len(searchHits))
				for _, hitText := range searchHits {
					hitTexts = append(hitTexts, template.HTML(hitText))
				}
				result.Texts = hitTexts
			} else {
				result.Texts = []template.HTML{
					template.HTML(tt.HTMLEscapeString(entry.GetBodySnippet())),
				}
			}

			searchRes = append(searchRes, result)
		}
	}
	renderTemplate(w, "search_result", searchRes)
}
