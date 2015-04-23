package main

import (
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/storage"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"time"
)

func (s *httpServer) withDoc(fn func(id string, entry *storage.Entry, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if _, err := uuid.ParseHex(id); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Warningf("Bad Document ID Request: %v", id)
			return
		}

		entry := storage.LoadEntry(id)
		if entry == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Warningf("Document %v not found", id)
			return
		}

		fn(id, entry, w, r)
	}
}

func (s *httpServer) serveCreateEntry(w http.ResponseWriter, r *http.Request) {
	showPreview := false
	entry := new(storage.Entry)

	if r.Method == "POST" {
		r.ParseForm()
		entry.Body = r.Form["body"][0]
		entry.Title = r.Form["title"][0]
		entry.Tags = r.Form["tags"]
		entry.UnifyTags()
		if _, ok := r.Form["save"]; ok {
			id := storage.CreateEntry(entry)
			http.Redirect(w, r, "/view/"+id, http.StatusSeeOther)
		} else if _, ok := r.Form["preview"]; ok {
			showPreview = true
		}
	}

	renderTemplate(w, "create_document", struct {
		Title       string
		Content     string
		Tags        []string
		ShowPreview bool
	}{
		entry.Title,
		entry.Body,
		entry.Tags,
		showPreview,
	})
}

func (s *httpServer) serveDeleteEntry(entryId string, entry *storage.Entry, w http.ResponseWriter, r *http.Request) {
	storage.DeleteEntry(entryId, entry)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *httpServer) serveEditEntry(entryId string, entry *storage.Entry, w http.ResponseWriter, r *http.Request) {
	showPreview := false

	if r.Method == "POST" {
		r.ParseForm()
		entry.Tags = r.Form["tags"]
		entry.Body = r.Form["body"][0]
		entry.Title = r.Form["title"][0]
		entry.UnifyTags()
		if _, ok := r.Form["save"]; ok {
			storage.UpdateEntry(entryId, entry)
			http.Redirect(w, r, "/view/"+entryId, http.StatusSeeOther)
		} else if _, ok := r.Form["preview"]; ok {
			showPreview = true
		}
	}

	renderTemplate(w, "edit_document", struct {
		Title       string
		Id          string
		Content     string
		Tags        []string
		ShowPreview bool
	}{
		entry.Title,
		entryId,
		entry.Body,
		entry.Tags,
		showPreview,
	})
}

func (s *httpServer) serveViewEntry(entryId string, entry *storage.Entry, w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view_document", struct {
		Id           string
		Title        string
		Content      string
		Tags         []string
		LastModified time.Time
	}{
		entryId,
		entry.Title,
		entry.Body,
		entry.Tags,
		entry.LastModified,
	})
}
