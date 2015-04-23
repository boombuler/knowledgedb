package main

import (
	"encoding/json"
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/storage"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"sort"
)

func (s *httpServer) serveAllTags(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(storage.GetAllTags())
}

func (s *httpServer) serveTagDetails(w http.ResponseWriter, r *http.Request) {
	tagId, ok := mux.Vars(r)["tag"]
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	tagId, err := url.QueryUnescape(tagId)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Warning(err)
		return
	}

	items := storage.EntriesByDateTime(storage.GetAllEntriesWithTag(tagId))
	sort.Sort(items)

	for _, item := range items {
		item.Body = item.GetBodySnippet()
	}

	renderTemplate(w, "tag", struct {
		Name  string
		Items []*storage.EntryWithId
	}{tagId, items})
}
